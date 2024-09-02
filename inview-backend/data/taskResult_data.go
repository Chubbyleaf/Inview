package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"insense-local/database"
	"insense-local/models"
	"time"
)

type TaskResultDataInterface interface {
	Create(ctx context.Context, taskResult *models.TaskResult) (interface{}, error)
	FindTaskResultList(ctx context.Context, deviceID int, model, algorithmType, startTime, endTime string) ([]*models.TaskResult, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.TaskResult, error)
	Fetch(c context.Context) ([]models.TaskResult, error)
}

type taskResultData struct {
	database   database.Database
	collection string
}

func TaskResultData(db database.Database, collection string) TaskResultDataInterface {
	return &taskResultData{
		database:   db,
		collection: collection,
	}
}

func (trd *taskResultData) Create(ctx context.Context, taskResult *models.TaskResult) (interface{}, error) {
	collection := trd.database.Collection(trd.collection)
	id, err := collection.InsertOne(ctx, taskResult)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (trd *taskResultData) FindByID(ctx context.Context, id primitive.ObjectID) (*models.TaskResult, error) {
	collection := trd.database.Collection(trd.collection)
	var taskResult models.TaskResult
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&taskResult)
	if err != nil {
		return nil, err
	}
	return &taskResult, nil
}

func (trd *taskResultData) Fetch(c context.Context) ([]models.TaskResult, error) {
	collection := trd.database.Collection(trd.collection)
	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	var taskResults []models.TaskResult
	err = cursor.All(c, &taskResults)
	if taskResults == nil {
		return []models.TaskResult{}, err
	}
	return taskResults, err
}

func (trd *taskResultData) FindTaskResultList(ctx context.Context, deviceID int, model, algorithmType, startTime, endTime string) ([]*models.TaskResult, error) {
	collection := trd.database.Collection(trd.collection)

	// 设置过滤条件
	filter := bson.M{
		"deviceId":      deviceID,
		"model":         model,
		"algorithmType": algorithmType,
	}

	// 处理日期范围过滤
	if startTime != "" && endTime != "" {
		start, err := time.Parse("20060102", startTime)
		if err != nil {
			return nil, err
		}
		end, err := time.Parse("20060102", endTime)
		if err != nil {
			return nil, err
		}
		filter["time"] = bson.M{
			"$gte": start,
			"$lt":  end.AddDate(0, 0, 1), // 包含结束日期的整个一天
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "time", Value: -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var taskResults []*models.TaskResult
	for cursor.Next(ctx) {
		var taskResult models.TaskResult
		if decodeErr := cursor.Decode(&taskResult); decodeErr != nil {
			return nil, decodeErr
		}
		taskResults = append(taskResults, &taskResult)
	}

	return taskResults, nil
}
