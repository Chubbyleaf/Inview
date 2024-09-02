package route

import (
	"insense-local/config"
	"insense-local/controller"
	"insense-local/data"
	"insense-local/database"
	"insense-local/middleware"
	"insense-local/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(env *config.Env, timeout time.Duration, db database.Database, engine *gin.Engine) {
	publicRouter := engine.Group("/api")
	publicRouter.GET("/health-check", HealthCheck)
	publicRouter.GET("/fuck", CacheTest)

	//用户相关操作
	userService := services.UserService(data.UserData(db, "users"), timeout)
	userController := controller.UserController{
		UserService: userService,
	}
	publicRouter.POST("/user/register", func(c *gin.Context) {
		userController.Register(c)
	})

	publicRouter.POST("/user/login", func(c *gin.Context) {
		userController.Login(c)
	})

	publicRouter.POST("/user/editpsw", func(c *gin.Context) {
		userController.UpdatePassword(c)
	})
	publicRouter.DELETE("/user/delete/:id", func(c *gin.Context) {
		userController.DeleteUser(c)
	})

	publicRouter.POST("/user/editUserInfo", func(c *gin.Context) {
		userController.UpdateUserInfo(c)
	})
	publicRouter.GET("/user/info/:id", func(c *gin.Context) {
		userController.GetUserInfo(c)
	})

	// 创建 CameraController 实例
	cameraService := services.CameraService(data.CameraData(db, "cameras"), timeout)
	cameraController := controller.CameraController{
		CameraService: cameraService,
	}

	publicRouter.POST("/camera/add", func(c *gin.Context) {
		cameraController.CreateCamera(c)
	})

	//修改设备信息
	publicRouter.POST("/camera/editCameraInfo", func(c *gin.Context) {
		cameraController.UpdateCamera(c)
	})
	// 获取设备列表
	publicRouter.GET("/camera/getCameraList", func(c *gin.Context) {
		cameraController.GetCameras(c)
	})

	//任务相关操作
	taskService := services.TaskService(data.TaskData(db, "tasks"), timeout)
	taskController := controller.TaskController{
		TaskService: taskService,
	}

	//删除设备
	publicRouter.DELETE("/camera/delete/:deviceId", func(c *gin.Context) {
		cameraController.DeleteCameraByDeviceID(c, taskController)
	})
	//启动定时任务
	taskService.StartCronTask()

	publicRouter.POST("/algorithm/addTask", func(c *gin.Context) {
		taskController.CreateTask(c, env)
	})
	publicRouter.DELETE("/algorithm/deleteTask/:id", func(c *gin.Context) {
		taskController.DeleteTask(c)
	})

	publicRouter.POST("/algorithm/stopTask/:id", func(c *gin.Context) {
		taskController.StopTask(c)
	})

	publicRouter.POST("/algorithm/startTask/:id", func(c *gin.Context) {
		taskController.StartTask(c)
	})

	publicRouter.POST("/algorithm/editTask", func(c *gin.Context) {
		taskController.UpdateTask(c, env)
	})

	publicRouter.GET("/algorithm/getTaskList", func(c *gin.Context) {
		taskController.FindAllTasksWithCameras(c, cameraController)
	})

	publicRouter.POST("/algorithm/updateTaskStatus", func(c *gin.Context) {
		taskController.Call(c)
	})

	//任务结果相关操作
	taskResultService := services.TaskResultService(data.TaskResultData(db, "results"), timeout)
	taskResultController := controller.TaskResultController{
		TaskResultService: taskResultService,
	}

	publicRouter.POST("/algorithm/publishResults", func(c *gin.Context) {
		taskResultController.AddTaskResult(c, env, taskController)
	})

	publicRouter.POST("/algorithm/getTaskResultList", func(c *gin.Context) {
		taskResultController.GetTaskResultList(c)
	})

	publicRouter.GET("/algorithm/getTaskResult/:id", func(c *gin.Context) {
		taskResultController.GetTaskResultById(c)
	})

	publicRouter.GET("/algorithm/getTaskResultList", func(c *gin.Context) {
		taskResultController.FindAllTaskResults(c)
	})

	//任务结果相关操作
	systemService := services.SystemService(timeout)
	systemController := controller.SystemController{
		SystemService: systemService,
	}

	//获取系统信息
	publicRouter.GET("/system/getInfo", func(c *gin.Context) {
		systemController.GetSystemInfo(c)
	})

	//获取系统日志
	publicRouter.GET("/system/getLog", func(c *gin.Context) {
		systemController.GetSystemLog(c, env)
	})

	//获取网络信,IP地址等
	publicRouter.GET("/system/getIp", func(c *gin.Context) {
		systemController.GetNetworkInfo(c)
	})
	AuthRouter(env, timeout, db, publicRouter)

	protectedRouter := engine.Group("/api")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// DataRouter(env, timeout, db, protectedRouter)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func CacheTest(c *gin.Context) {
	// cacheConfig := config.SetupRedisCache("fuck", "fuck", "fuck")
	// cache.CacheByRequest
}
