package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionCameras = "cameras"
)

type Camera struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DeviceID int                `bson:"deviceId" json:"deviceId"`
	Name     string             `bson:"name" json:"name"`
	Remark   string             `bson:"remark,omitempty" json:"remark,omitempty"`
	Type     string             `bson:"type" json:"type"`
	URL      string             `bson:"url" json:"url"`
}
