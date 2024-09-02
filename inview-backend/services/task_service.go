package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"insense-local/config"
	"insense-local/data"
	"insense-local/models"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, task *models.Task, env *config.Env) error
	UpdateTask(ctx context.Context, task *models.Task, env *config.Env) error
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
	FindAllTasksWithCameras(ctx context.Context, cs CameraServiceInterface) ([]models.TaskWithCamera, error)
	StartCronTask()
	TaskAlive(ctx context.Context, task *models.Task) error
	StopTaskByUser(ctx context.Context, id primitive.ObjectID) error
	StartTaskByUser(ctx context.Context, id primitive.ObjectID) error
	DeviceExistTask(ctx context.Context, id int) (bool, error)
	FindTask(c context.Context, deviceId int, model string, algorithm string) (*models.Task, error)
}

type taskService struct {
	taskData       data.TaskDataInterface
	contextTimeout time.Duration
}

func TaskService(taskData data.TaskDataInterface, timeout time.Duration) TaskServiceInterface {
	return &taskService{
		taskData:       taskData,
		contextTimeout: timeout,
	}
}

func (ts *taskService) CreateTask(ctx context.Context, task *models.Task, env *config.Env) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	existingTask, err := ts.taskData.FindOne(ctx, task.DeviceID, task.Model, task.AlgorithmType)
	if err == nil && existingTask != nil {
		return errors.New("task with the same deviceId, model, and algorithmType already exists")
	}

	// 设置默认值
	task.ImgDir = models.DefaultImgDir
	task.LogPath = models.DefaultLogPath
	task.VideoPath = models.DefaultVideoPath
	task.CloudAPI = models.DefaultCloudAPI
	task.ID = primitive.NewObjectID()
	println(task.ID.Hex())
	var zone = time.FixedZone("CST", 8*60*60)
	currentTime := time.Now().In(zone)
	task.InitTime = currentTime
	task.UpdateTime = currentTime
	task.LatestCallTime = currentTime

	resultsPath := createFileFolders(task, env)
	dockerID := ts.StartDocker(task.Model, task.AlgorithmType, resultsPath)
	if strings.HasPrefix(dockerID, "error") {
		task.Status = -1
		task.DockerID = "NA"
	} else {
		task.Status = 1
		task.DockerID = dockerID
	}
	res := ts.taskData.Create(ctx, task)
	return res
}

func createFileFolders(task *models.Task, env *config.Env) []string {
	taskPath := strconv.Itoa(task.DeviceID) + "/" + task.ID.Hex()
	configFileName := task.Model + "_" + task.AlgorithmType + ".json"
	fullTaskPath := filepath.Join(env.JsonPath, taskPath)
	task.TaskResultsPath = fullTaskPath
	configPath := filepath.Join(fullTaskPath, "config")
	logPath := filepath.Join(fullTaskPath, "log")
	videoPath := filepath.Join(fullTaskPath, "video")
	imagePath := filepath.Join(fullTaskPath, "image")
	jsonFile := filepath.Join(configPath, configFileName)

	if _, err := os.Stat(fullTaskPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullTaskPath, 0755)
		if err != nil {
			log.Printf("json path make failed", err)
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(configPath, 0755)
		if err != nil {
			log.Printf("config file save path make failed", err)
		}
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(logPath, 0755)
		if err != nil {
			log.Printf("log file save path make failed", err)
		}
	}
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		err := os.MkdirAll(videoPath, 0755)
		if err != nil {
			log.Printf("video file save path make failed", err)
		}
	}
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		err := os.MkdirAll(imagePath, 0755)
		if err != nil {
			log.Printf("image file save path make failed", err)
		}
	}

	file, err := os.Create(jsonFile)
	if err != nil {
		log.Printf("json file create failed")
	}
	defer file.Close()

	content, err := json.Marshal(task)
	if err != nil {
		log.Printf("json converted failed")
	}
	_, err = file.Write(content)
	if err != nil {
		log.Printf("file write failed")
	}
	var resultsPath = []string{
		logPath, imagePath, videoPath, configPath,
	}
	return resultsPath
}

func (ts *taskService) StartDocker(model string, algorithm string, paths []string) string {
	var cmd *exec.Cmd
	if model == "safety" && algorithm == "equipment" {
		var cmdArgs = []string{
			"run",
			"-d",
			"--gpus", "all",
			"--network", "host",
			"-v", paths[0] + ":/log",
			"-v", paths[1] + ":/image",
			"-v", paths[3] + ":/config",
			"safety_equipment:v1.0",
			"/config/safety_equipment.json",
		}
		cmd = exec.Command("docker", cmdArgs...)
	} else if model == "safety" && algorithm == "concreteSupport" {
		var cmdArgs = []string{
			"run",
			"-d",
			"--gpus", "all",
			"--network", "host",
			"-v", paths[0] + ":/log",
			"-v", paths[1] + ":/image",
			"-v", paths[3] + ":/config",
			"safety_concretesupport:v1.0",
			"/config/safety_concreteSupport.json",
		}
		cmd = exec.Command("docker", cmdArgs...)
	} else if model == "safety" && algorithm == "areaEdge" {
		var cmdArgs = []string{
			"run",
			"-d",
			"--gpus", "all",
			"--network", "host",
			"-v", paths[0] + ":/log",
			"-v", paths[1] + ":/image",
			"-v", paths[3] + ":/config",
			"safety_areaedge:v1",
			"/config/safety_areaEdge.json",
		}
		cmd = exec.Command("docker", cmdArgs...)
	} else if model == "safety" && algorithm == "fireSmoke" {
		var cmdArgs = []string{
			"run",
			"-d",
			"--runtime", "nvidia",
			"--network", "host",
			"-v", paths[0] + ":/log",
			"-v", paths[1] + ":/image",
			"-v", paths[2] + ":/video",
			"-v", paths[3] + ":/config",
			"safety_firesmoke:v1.0",
			"--json", "/config/safety_fireSmoke.json"}
		cmd = exec.Command("docker", cmdArgs...)
	} else {
		return "error"
	}

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Stderr:%s", stderr.String())
		log.Printf("Unable to start Docker container: %v", err)
		return "error: " + stderr.String()
	}
	outStr := out.String()
	fmt.Println(strconv.Quote("start docker out.Str: " + outStr))
	DockerID := strings.ReplaceAll(outStr, "\n", "")
	fmt.Println(strconv.Quote("start docker id" + DockerID))
	fmt.Println("results", DockerID)
	log.Println("Container started successfully!")
	return DockerID
}

func (ts *taskService) StopDocker(task *models.Task) error {
	var cmd *exec.Cmd
	cmd = exec.Command("docker", "stop", task.DockerID)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Printf("Stop docker id:%s", task.DockerID)
	if err != nil {
		log.Printf("Stderr:%s", stderr.String())
		log.Printf("Unable to stop task-%s Docker container: %v", task.ID, err)
		log.Printf("Stop docker id:%s", task.DockerID)
		return err
	}
	outStr := out.String()
	log.Println(strconv.Quote("stop docker out.Str: " + outStr))
	log.Printf("Task-%s docker container stop successfully!", task.ID)
	return nil
}

func (ts *taskService) RemoveDocker(task *models.Task) error {
	var cmd *exec.Cmd
	cmd = exec.Command("docker", "rm", "-f", task.DockerID)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Printf("Remove docker id:%s", task.DockerID)
	if err != nil {
		log.Printf("Stderr:%s", stderr.String())
		log.Printf("Unable to remove task-%s Docker container: %v", task.ID, err)
		return err
	}
	outStr := out.String()
	log.Println(strconv.Quote("remove docker out.Str: " + outStr))
	log.Printf("Task-%s docker container remove successfully!", task.ID)
	return nil
}

func (ts *taskService) StopTaskByUser(ctx context.Context, id primitive.ObjectID) error {
	err := ts.UpdateTaskStatus(ctx, id, 2)
	if err != nil {
		return err
	}
	task, _ := ts.taskData.FindByID(ctx, id)
	if task.DockerID == "NA" || task.DockerID == "" {
		return nil
	}
	return ts.StopDocker(&task)
}

// StartTaskByUser 用户手动启动算法
func (ts *taskService) StartTaskByUser(ctx context.Context, id primitive.ObjectID) error {
	task, err := ts.taskData.FindByID(ctx, id)
	if err != nil {
		return err
	} else {
		return ts.RestartDocker(ctx, task)
	}
}

func (ts *taskService) RestartDocker(ctx context.Context, task models.Task) error {
	var cmd *exec.Cmd

	if task.DockerID == "NA" || task.DockerID == "" {
		log.Println("dockerId,Na")
		if task.Model == "safety" && task.AlgorithmType == "equipment" {
			var args = []string{
				"run",
				"-d",
				"--gpus", "all",
				"--network", "host",
				"-v", task.TaskResultsPath + "/log:/log",
				"-v", task.TaskResultsPath + "/image:/image",
				"-v", task.TaskResultsPath + "/config:/config",
				"safety_equipment:v1.0",
				"/config/safety_equipment.json",
			}
			cmd = exec.Command("docker", args...)
		} else if task.Model == "safety" && task.AlgorithmType == "concreteSupport" {
			var args = []string{
				"run",
				"-d",
				"--gpus", "all",
				"--network", "host",
				"-v", task.TaskResultsPath + "/log:/log",
				"-v", task.TaskResultsPath + "/image:/image",
				"-v", task.TaskResultsPath + "/config:/config",
				"safety_concretesupport:v1.0",
				"/config/safety_concreteSupport.json",
			}
			cmd = exec.Command("docker", args...)
		} else if task.Model == "safety" && task.AlgorithmType == "areaEdge" {
			var args = []string{
				"run",
				"-d",
				"--gpus", "all",
				"--network", "host",
				"-v", task.TaskResultsPath + "/log:/log",
				"-v", task.TaskResultsPath + "/image:/image",
				"-v", task.TaskResultsPath + "/config:/config",
				"safety_areaedge:v1.0",
				"/config/safety_areaEdge.json",
			}
			cmd = exec.Command("docker", args...)
		} else if task.Model == "safety" && task.AlgorithmType == "fireSmoke" {
			var args = []string{
				"run",
				"-d",
				"--runtime", "nvidia",
				"--network", "host",
				"-v", task.TaskResultsPath + "/log:/log",
				"-v", task.TaskResultsPath + "/image:/image",
				"-v", task.TaskResultsPath + "/video:/video",
				"-v", task.TaskResultsPath + "/config:/config",
				"safety_fireSmoke:v1.0",
				"--json", "/config/safety_fireSmoke.json",
			}
			cmd = exec.Command("docker", args...)
		} else {
			return errors.New("task model and algorithm not supported")
		}
	} else {
		log.Println("dockerId", task.DockerID)
		task.DockerID = strings.ReplaceAll(task.DockerID, "\n", "")
		cmd = exec.Command("docker", "restart", task.DockerID)
	}
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	updateFields := bson.M{}
	err := cmd.Run()
	if err != nil {
		log.Printf("stderr:", stderr.String())
		log.Printf("Unable to start Docker container: %v", err)
		updateFields["status"] = -1
		ts.taskData.Update(ctx, task.ID, updateFields)
		return err
	}
	outStr := out.String()
	fmt.Println(strconv.Quote("restart docker out.Str: " + outStr))
	DockerID := strings.ReplaceAll(outStr, "\n", "")
	fmt.Println(strconv.Quote("restart docker id" + DockerID))
	updateFields["status"] = 1
	updateFields["dockerId"] = DockerID
	ts.taskData.Update(ctx, task.ID, updateFields)
	log.Println("Container restarted successfully!")
	return nil
}

// StartCronTask 定时任务每一分钟执行一次，检查算法是否正在运行
func (ts *taskService) StartCronTask() {
	ctx := context.Background()

	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() {
		ts.CheckAlgorithmStatus(ctx)
	})
	if err != nil {
		log.Println("添加定时任务失败:", err)
		return
	}
	c.Start()
	log.Println("定时任务启动成功！")
}

func (ts *taskService) CheckAlgorithmStatus(ctx context.Context) {
	var zone = time.FixedZone("CST", 8*60*60)
	currentTime := time.Now().In(zone)
	tmpCronTaskList, err := ts.FetchTasks(ctx)
	if err == nil {
		for _, cronTask := range tmpCronTaskList {
			//用户没有手动停止并且在时间范围内需要重新启动docker
			if cronTask.Status != 2 && isCurrentTimeWithinWorkingPeriod(currentTime, cronTask.WorkingTime) {
				if currentTime.Sub(cronTask.LatestCallTime).Minutes() > 2 {
					log.Println(cronTask.DeviceID, cronTask.Model, cronTask.AlgorithmType, "运行失败")
					cronTask.Status = -1
					ts.RestartDocker(ctx, cronTask)
				} else {
					cronTask.Status = 1
				}
				ts.UpdateTaskStatus(ctx, cronTask.ID, cronTask.Status)
			} else if !isCurrentTimeWithinWorkingPeriod(currentTime, cronTask.WorkingTime) {
				//不在运行时间内停止docker
				cronTask.Status = 3
				ts.StopDocker(&cronTask)
				ts.UpdateTaskStatus(ctx, cronTask.ID, cronTask.Status)
			}
		}
	}
}

// 检查现在的时间是不是工作时间
func isCurrentTimeWithinWorkingPeriod(currentTime time.Time, workingPeriods []models.WorkingPeriod) bool {
	currentTimeStr := currentTime.Format("150405")
	for _, period := range workingPeriods {
		if period.StartTime <= currentTimeStr && currentTimeStr <= period.EndTime {
			return true
		}
	}
	return false
}

func (ts *taskService) UpdateTask(ctx context.Context, task *models.Task, env *config.Env) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	tmpTask, err := ts.taskData.FindByID(ctx, task.ID)
	if err != nil {
		return err
	}

	updateFields := bson.M{}
	if task.Coordinate != nil {
		updateFields["coordinate"] = task.Coordinate
		tmpTask.Coordinate = task.Coordinate
	}
	if task.LogRotateDuration != 0 {
		updateFields["logRotateDuration"] = task.LogRotateDuration
		tmpTask.LogRotateDuration = task.LogRotateDuration
	}
	if task.LogSize != 0 {
		updateFields["logSize"] = task.LogSize
		tmpTask.LogSize = task.LogSize
	}
	if task.DataSize != 0 {
		updateFields["dataSize"] = task.DataSize
		tmpTask.DataSize = task.DataSize
	}
	if task.DataRotation != 0 {
		updateFields["dataRotation"] = task.DataRotation
		tmpTask.DataRotation = task.DataRotation
	}
	if task.WorkingTime != nil {
		updateFields["workingTime"] = task.WorkingTime
		tmpTask.WorkingTime = task.WorkingTime
	}
	if task.GPU != "" {
		updateFields["gpu"] = task.GPU
		tmpTask.GPU = task.GPU
	}
	updateFields["targetApi"] = task.TargetAPI
	updateFields["sound"] = task.Sound
	tmpTask.Sound = task.Sound
	updateFields["sms"] = task.SMS
	tmpTask.SMS = task.SMS
	updateFields["tel"] = task.Tel
	tmpTask.Tel = task.Tel
	updateFields["updateTime"] = time.Now().In(time.FixedZone("CST", 8*60*60))
	tmpTask.UpdateTime = task.UpdateTime
	createFileFolders(&tmpTask, env)
	ts.RestartDocker(ctx, tmpTask)
	return ts.taskData.Update(ctx, task.ID, updateFields)
}

func (ts *taskService) UpdateTaskStatus(ctx context.Context, id primitive.ObjectID, status int) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	updateFields := bson.M{}
	_, err := ts.taskData.FindByID(ctx, id)
	if err != nil {
		return err
	}
	updateFields["status"] = status
	return ts.taskData.Update(ctx, id, updateFields)
}

func (ts *taskService) DeviceExistTask(ctx context.Context, id int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	task, err := ts.taskData.FindTasksByDeviceID(ctx, id)
	if err != nil {
		return true, err
	}
	if task != nil {
		return true, err
	}
	return false, nil
}

// TaskAlive 检查算法是否正常运行
func (ts *taskService) TaskAlive(ctx context.Context, task *models.Task) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	task, err := ts.taskData.FindOne(ctx, task.DeviceID, task.Model, task.AlgorithmType)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found, please add using the api '/api/algorithm/addTask")
	}
	var zone = time.FixedZone("CST", 8*60*60)
	currentTime := time.Now().In(zone)
	updateFields := bson.M{}
	updateFields["status"] = 1
	updateFields["latestCallTime"] = currentTime
	return ts.taskData.Update(ctx, task.ID, updateFields)
}

func (ts *taskService) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	task, err := ts.taskData.FindByID(ctx, id)
	if err != nil {
		return err
	}
	err = ts.RemoveDocker(&task)
	if err != nil {
		return err
	}
	return ts.taskData.Delete(ctx, id)
}

func (ts *taskService) FetchTasks(c context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(c, ts.contextTimeout)
	defer cancel()
	return ts.taskData.Fetch(ctx)
}

func (ts *taskService) FindAllTasksWithCameras(ctx context.Context, cs CameraServiceInterface) ([]models.TaskWithCamera, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	cameras, err := cs.FetchCameras(ctx)
	if err != nil {
		return nil, err
	}

	var results []models.TaskWithCamera

	for _, camera := range cameras {
		tasks, err := ts.taskData.FindTasksByDeviceID(ctx, camera.DeviceID)
		if err != nil {
			return nil, err
		}

		if tasks == nil {
			results = append(results, models.TaskWithCamera{
				DeviceID:      camera.DeviceID,
				CameraID:      camera.ID,
				CameraIName:   camera.Name,
				CameraIRemark: camera.Remark,
				CameraIType:   camera.Type,
				CameraIURL:    camera.URL,
			})
		} else {
			for _, task := range tasks {
				results = append(results, models.TaskWithCamera{
					TaskID:              task.ID,
					DeviceID:            task.DeviceID,
					Model:               task.Model,
					AlgorithmType:       task.AlgorithmType,
					LiveStreamInputURL:  task.LiveStreamInputURL,
					LiveStreamOutputURL: task.LiveStreamOutputURL,
					Coordinate:          task.Coordinate,
					LogPath:             task.LogPath,
					LogRotateDuration:   task.LogRotateDuration,
					LogSize:             task.LogSize,
					DataSize:            task.DataSize,
					DataRotation:        task.DataRotation,
					WorkingTime:         task.WorkingTime,
					TargetAPI:           task.TargetAPI,
					ImgDir:              task.ImgDir,
					GPU:                 task.GPU,
					Sound:               task.Sound,
					SMS:                 task.SMS,
					Tel:                 task.Tel,
					InitTime:            task.InitTime,
					UpdateTime:          task.UpdateTime,
					Status:              task.Status,
					TaskResultsPath:     task.TaskResultsPath,
					CameraID:            camera.ID,
					CameraIName:         camera.Name,
					CameraIRemark:       camera.Remark,
					CameraIType:         camera.Type,
					CameraIURL:          camera.URL,
				})
			}

		}
	}

	return results, nil
}

func (ts *taskService) FindTask(c context.Context, deviceId int, model string, algorithm string) (*models.Task, error) {
	c, cancel := context.WithTimeout(c, ts.contextTimeout)
	defer cancel()
	return ts.taskData.FindOne(c, deviceId, model, algorithm)
}
