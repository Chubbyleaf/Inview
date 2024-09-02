package data

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"insense-local/database"
	"insense-local/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userData struct {
	database   database.Database
	collection string
}

func UserData(db database.Database, collection string) UserDataInterface {
	return &userData{
		database:   db,
		collection: collection,
	}
}

type UserDataInterface interface {
	Create(c context.Context, user *models.User) error
	Fetch(c context.Context) ([]models.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	FindByUserName(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

func (ud *userData) Create(c context.Context, user *models.User) error {
	collection := ud.database.Collection(ud.collection)
	_, err := collection.InsertOne(c, user)
	return err
}
func (ud *userData) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := ud.database.Collection(ud.collection)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (ud *userData) Fetch(c context.Context) ([]models.User, error) {
	collection := ud.database.Collection(ud.collection)
	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []models.User

	err = cursor.All(c, &users)
	if users == nil {
		return []models.User{}, err
	}

	return users, err
}

func (ud *userData) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	collection := ud.database.Collection(ud.collection)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ud *userData) FindByUserName(ctx context.Context, username string) (*models.User, error) {
	collection := ud.database.Collection(ud.collection)
	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ud *userData) Update(ctx context.Context, user *models.User) error {
	collection := ud.database.Collection(ud.collection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}
