package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskResult struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DeviceID      int                `bson:"deviceId" json:"deviceId"`
	Model         string             `bson:"model" json:"model" binding:"required"`
	AlgorithmType string             `bson:"algorithmType" json:"algorithmType" binding:"required"`
	ImgName       string             `bson:"imgName" json:"imgName" binding:"required"`
	PredResult    [][]float64        `bson:"predResult" json:"predResult" binding:"required""`
	ClassName     map[string]string  `bson:"className" json:"className" binding:"required"`
	VideoName     string             `bson:"videoName" json:"videoName" binding:"required"`
	Time          time.Time          `bson:"time" json:"time" binding:"required"`
	Info          string             `bson:"info" json:"info" binding:"required"`
}
