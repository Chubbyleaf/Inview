package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"insense-local/config"
	"insense-local/data"
	"insense-local/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type TaskResultServiceInterface interface {
	AddTaskResult(ctx context.Context, env *config.Env, taskResult *models.TaskResult, ts TaskServiceInterface) (interface{}, error)
	FindTaskResultList(ctx context.Context, deviceID int, model, algorithmType, startTime, endTime string) ([]*models.TaskResult, error)
	FindTaskResultById(ctx context.Context, id primitive.ObjectID) (*models.TaskResult, error)
	GetResultList(ctx context.Context) ([]models.TaskResult, error)
}

type taskResultService struct {
	taskResultData data.TaskResultDataInterface
	contextTimeout time.Duration
}

func TaskResultService(taskResult data.TaskResultDataInterface, timeout time.Duration) TaskResultServiceInterface {
	return &taskResultService{
		taskResultData: taskResult,
		contextTimeout: timeout,
	}
}

func (trs *taskResultService) AddTaskResult(ctx context.Context, env *config.Env, taskResult *models.TaskResult, ts TaskServiceInterface) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, trs.contextTimeout)
	defer cancel()
	task, err := ts.FindTask(ctx, taskResult.DeviceID, taskResult.Model, taskResult.AlgorithmType)
	if err != nil {
		return nil, err
	}

	imagePath := filepath.Join(env.ImagePath, strconv.Itoa(taskResult.DeviceID), task.ID.Hex(), "image", taskResult.ImgName)
	imageBase64, err := trs.ImageToBase64(imagePath)
	if err != nil {
		log.Printf(err.Error())
	} else {
		//推到云端
		err = trs.PushResults(ctx, taskResult, imageBase64, models.DefaultCloudAPI)
		if err != nil {
			log.Printf(err.Error())
		}
		//推到目标API
		if task.TargetAPI != "" {
			err = trs.PushResults(ctx, taskResult, imageBase64, task.TargetAPI)
		}
		if err != nil {
			log.Printf(err.Error())
		}
	}
	return trs.taskResultData.Create(ctx, taskResult)
}
func (trs *taskResultService) FindTaskResultList(ctx context.Context, deviceID int, model, algorithmType, startTime, endTime string) ([]*models.TaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, trs.contextTimeout)
	defer cancel()
	return trs.taskResultData.FindTaskResultList(ctx, deviceID, model, algorithmType, startTime, endTime)
}

func (trs *taskResultService) FindTaskResultById(ctx context.Context, id primitive.ObjectID) (*models.TaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, trs.contextTimeout)
	defer cancel()
	return trs.taskResultData.FindByID(ctx, id)
}

func (trs *taskResultService) GetResultList(ctx context.Context) ([]models.TaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, trs.contextTimeout)
	defer cancel()
	return trs.taskResultData.Fetch(ctx)
}

func (trs *taskResultService) ImageToBase64(imagePath string) (string, error) {
	// 读取图片文件的内容
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}
	// 将图片内容编码为Base64
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	return imgBase64, nil
}

func (trs *taskResultService) PushResults(ctx context.Context, taskResult *models.TaskResult, imageBase64, apiPath string) error {
	// 动态构建请求体
	requestBodyMap := map[string]interface{}{
		"deviceId":      taskResult.DeviceID,
		"model":         taskResult.Model,
		"algorithmType": taskResult.AlgorithmType,
		"imgName":       taskResult.ImgName,
		"predResult":    taskResult.PredResult,
		"className":     taskResult.ClassName,
		"videoName":     taskResult.VideoName,
		"time":          taskResult.Time.Format(time.RFC3339), // 将 Time 转换为字符串
		"info":          taskResult.Info,
		"imageBase64":   imageBase64, // 添加 Base64 编码的图像
	}

	// 将 map 序列化为 JSON
	requestBody, err := json.Marshal(requestBodyMap)

	if err != nil {
		return fmt.Errorf("failed to marshal task result: %w", err)
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiPath, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	return nil
}
