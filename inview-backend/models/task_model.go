package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DefaultImgDir    = "/image"                                                      // 设置默认的 imgDir docker内部
	DefaultLogPath   = "/log"                                                        // 设置默认的 logPath
	DefaultVideoPath = "/video"                                                      // default video result path
	DefaultCloudAPI  = "http://116.62.47.112:8072/api/v1/data/uploadAlgorithmResult" // 设置默认的推送API
)

// Task represents a task document in MongoDB
type Task struct {
	ID                  primitive.ObjectID `bson:"_id" json:"_id"`
	DeviceID            int                `bson:"deviceId" json:"deviceId"`
	Model               string             `bson:"model" json:"model"`
	AlgorithmType       string             `bson:"algorithmType" json:"algorithmType"`
	LiveStreamInputURL  string             `bson:"liveStreamInputUrl" json:"liveStreamInputUrl"`
	LiveStreamOutputURL string             `bson:"liveStreamOutputUrl,omitempty" json:"liveStreamOutputUrl,omitempty"`
	Coordinate          []CoordinateGroup  `bson:"coordinate,omitempty" json:"coordinate,omitempty"`
	LogPath             string             `bson:"logPath" json:"logPath"`
	LogRotateDuration   int                `bson:"logRotateDuration" json:"logRotateDuration"`
	LogSize             int                `bson:"logSize" json:"logSize"`
	DataSize            int                `bson:"dataSize" json:"dataSize"`
	DataRotation        int                `bson:"dataRotation" json:"dataRotation"`
	WorkingTime         []WorkingPeriod    `bson:"workingTime" json:"workingTime"`
	TargetAPI           string             `bson:"targetApi" json:"targetApi"`
	CloudAPI            string             `bson:"cloudApi" json:"cloudApi"`
	ImgDir              string             `bson:"imgDir,omitempty" json:"imgDir,omitempty"`
	GPU                 string             `bson:"gpu,omitempty" json:"gpu,omitempty"`
	VideoPath           string             `bson:"videoPath" json:"videoPath"`
	Sound               bool               `bson:"sound" json:"sound"`
	SMS                 bool               `bson:"sms" json:"sms"`
	Tel                 bool               `bson:"tel" json:"tel"`
	InitTime            time.Time          `bson:"initTime" json:"initTime"`
	UpdateTime          time.Time          `bson:"updateTime" json:"updateTime"`
	LatestCallTime      time.Time          `bson:"latestCallTime" json:"latestCallTime"`
	Status              int                `bson:"status" json:"status"`     //-1表示运行失败 0表示空  1运行成功 2用户手动停止 3表示不在运行时间范围内
	DockerID            string             `bson:"dockerId" json:"dockerId"` // 算法运行对应的dockerID
	TaskResultsPath     string             `bson:"taskResultsPath" json:"taskResultsPath"`
}

type CoordinateGroup []Coordinate

// Coordinate represents x and y coordinates
type Coordinate struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}

// WorkingPeriod represents a working time period with start and end times
type WorkingPeriod struct {
	StartTime string `bson:"startTime" json:"startTime"`
	EndTime   string `bson:"endTime" json:"endTime"`
}
