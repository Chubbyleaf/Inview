package data

import (
	"context"
	"errors"
	"insense-local/database"
	"insense-local/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskDataInterface interface {
	Create(ctx context.Context, task *models.Task) error
	Fetch(c context.Context) ([]models.Task, error)
	FindOne(ctx context.Context, deviceId int, model string, algorithm string) (*models.Task, error)
	Update(ctx context.Context, id primitive.ObjectID, updateFields bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (models.Task, error)
	FindTasksByDeviceID(ctx context.Context, deviceID int) ([]models.Task, error)
}

type taskData struct {
	database   database.Database
	collection string
}

func TaskData(db database.Database, collection string) TaskDataInterface {
	return &taskData{
		database:   db,
		collection: collection,
	}
}

func (td *taskData) Create(ctx context.Context, task *models.Task) error {
	collection := td.database.Collection(td.collection)
	_, err := collection.InsertOne(ctx, task)
	return err
}

func (td *taskData) Fetch(c context.Context) ([]models.Task, error) {
	collection := td.database.Collection(td.collection)
	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	var tasks []models.Task
	err = cursor.All(c, &tasks)
	if tasks == nil {
		return []models.Task{}, err
	}
	return tasks, err
}

func (td *taskData) FindOne(ctx context.Context, deviceId int, model string, algorithm string) (*models.Task, error) {
	filter := bson.M{
		"deviceId":      deviceId,
		"model":         model,
		"algorithmType": algorithm,
	}

	collection := td.database.Collection(td.collection)
	var tmpTask models.Task
	err := collection.FindOne(ctx, filter).Decode(&tmpTask)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &tmpTask, err
}

func (td *taskData) Update(ctx context.Context, id primitive.ObjectID, updateFields bson.M) error {
	collection := td.database.Collection(td.collection)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (td *taskData) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := td.database.Collection(td.collection)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (td *taskData) FindByID(ctx context.Context, id primitive.ObjectID) (models.Task, error) {
	collection := td.database.Collection(td.collection)
	var task models.Task
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (td *taskData) FindTasksByDeviceID(ctx context.Context, deviceID int) ([]models.Task, error) {
	collection := td.database.Collection(td.collection)
	cursor, err := collection.Find(ctx, bson.M{"deviceId": deviceID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
