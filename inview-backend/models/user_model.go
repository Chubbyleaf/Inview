package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserName string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Phone    string             `bson:"phone,omitempty" json:"phone"`
	Email    string             `bson:"email,omitempty" json:"email"`
}
