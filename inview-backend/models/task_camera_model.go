package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskWithCamera struct {
	TaskID              primitive.ObjectID `json:"taskId"`
	DeviceID            int                `json:"deviceId"`
	Model               string             `json:"model"`
	AlgorithmType       string             `json:"algorithmType"`
	LiveStreamInputURL  string             `json:"liveStreamInputUrl"`
	LiveStreamOutputURL string             `json:"liveStreamOutputUrl"`
	Coordinate          []CoordinateGroup  `json:"coordinate"`
	LogPath             string             `json:"logPath"`
	LogRotateDuration   int                `json:"logRotateDuration"`
	LogSize             int                `json:"logSize"`
	DataSize            int                `json:"dataSize"`
	DataRotation        int                `json:"dataRotation"`
	WorkingTime         []WorkingPeriod    `json:"workingTime"`
	TargetAPI           string             `json:"targetApi"`
	ImgDir              string             `json:"imgDir"`
	GPU                 string             `json:"gpu"`
	Sound               bool               `json:"sound"`
	SMS                 bool               `json:"sms"`
	Tel                 bool               `json:"tel"`
	InitTime            time.Time          `json:"initTime"`
	UpdateTime          time.Time          `json:"updateTime"`
	Status              int                `json:"status"`
	TaskResultsPath     string             `json:"taskResultsPath"`

	CameraID      primitive.ObjectID `json:"cameraId"`
	CameraIName   string             `json:"name"`
	CameraIRemark string             `json:"remark"`
	CameraIType   string             `json:"type"`
	CameraIURL    string             `json:"url"`
}
