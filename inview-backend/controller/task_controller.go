package controller

import (
	"insense-local/config"
	"insense-local/middleware"
	"insense-local/models"
	"insense-local/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	TaskService services.TaskServiceInterface
}

func (tc *TaskController) CreateTask(c *gin.Context, env *config.Env) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err := tc.TaskService.CreateTask(c, &task, env)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	// println(task.ID)
	successResponse := middleware.SuccessResponse{
		Message: "任务创建成功",
		Status:  http.StatusOK,
		Data:    task,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (tc *TaskController) UpdateTask(c *gin.Context, env *config.Env) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if task.ID.IsZero() {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid or missing ID",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err := tc.TaskService.UpdateTask(c, &task, env)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "任务更新成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (tc *TaskController) StopTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err = tc.TaskService.StopTaskByUser(c, id)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "算法成功停止!",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (tc *TaskController) StartTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err = tc.TaskService.StartTaskByUser(c, id)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "算法成功启动!",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err = tc.TaskService.DeleteTask(c, id)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "任务删除成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (tc *TaskController) FindAllTasksWithCameras(c *gin.Context, cc CameraController) {
	tasks, err := tc.TaskService.FindAllTasksWithCameras(c, cc.CameraService)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	successResponse := middleware.SuccessResponse{
		Message: "查询成功",
		Status:  http.StatusOK,
		Data:    tasks,
	}
	c.JSON(http.StatusOK, successResponse)

}

func (tc *TaskController) Call(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err := tc.TaskService.TaskAlive(c, &task)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}
