package data

import (
	"context"
	"insense-local/database"
	"insense-local/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CameraDataInterface interface {
	Create(c context.Context, camera *models.Camera) error
	Fetch(c context.Context) ([]models.Camera, error)
	Delete(c context.Context, id int) error
	Update(c context.Context, id int, camera *models.Camera) error
	FindByDeviceId(c context.Context, id int) (models.Camera, error)
	GetNextDeviceID(c context.Context) (int, error)
}

type cameraData struct {
	database   database.Database
	collection string
}

func CameraData(db database.Database, collection string) CameraDataInterface {
	return &cameraData{
		database:   db,
		collection: collection,
	}
}

func (cd *cameraData) Create(c context.Context, camera *models.Camera) error {
	collection := cd.database.Collection(cd.collection)
	_, err := collection.InsertOne(c, camera)
	return err
}

func (cd *cameraData) Fetch(c context.Context) ([]models.Camera, error) {
	collection := cd.database.Collection(cd.collection)

	// 创建排序选项
	findOptions := options.Find().SetSort(bson.D{{"deviceId", 1}})

	cursor, err := collection.Find(c, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	var cameras []models.Camera
	err = cursor.All(c, &cameras)
	if cameras == nil {
		return []models.Camera{}, err
	}

	return cameras, err
}

func (cd *cameraData) Delete(c context.Context, id int) error {
	collection := cd.database.Collection(cd.collection)

	_, err := collection.DeleteOne(c, bson.M{"deviceId": id})
	return err
}

func (cd *cameraData) Update(c context.Context, id int, camera *models.Camera) error {
	collection := cd.database.Collection(cd.collection)

	update := bson.M{
		"$set": bson.M{
			"name":   camera.Name,
			"remark": camera.Remark,
		},
	}

	_, err := collection.UpdateOne(c, bson.M{"deviceId": id}, update)
	return err
}

func (cd *cameraData) FindByDeviceId(c context.Context, id int) (models.Camera, error) {
	collection := cd.database.Collection(cd.collection)

	filter := bson.M{"deviceId": id}
	var camera models.Camera
	err := collection.FindOne(c, filter).Decode(&camera)
	if err != nil {
		return camera, err
	}
	return camera, nil
}

func (cd *cameraData) GetNextDeviceID(c context.Context) (int, error) {
	collection := cd.database.Collection(cd.collection)

	// 获取现有设备的所有 DeviceID
	var cameras []models.Camera
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var camera models.Camera
		if err = cursor.Decode(&camera); err != nil {
			return 0, err
		}
		cameras = append(cameras, camera)
	}

	// 查找最小的可用 DeviceID
	existingIDs := make(map[int]bool)
	for _, camera := range cameras {
		existingIDs[camera.DeviceID] = true
	}

	deviceID := 1
	for {
		if _, exists := existingIDs[deviceID]; !exists {
			break
		}
		deviceID++
	}

	return deviceID, nil
}
