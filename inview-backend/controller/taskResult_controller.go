package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"insense-local/config"
	"insense-local/middleware"
	"insense-local/models"
	"insense-local/services"
	"net/http"
)

type TaskResultController struct {
	TaskResultService services.TaskResultServiceInterface
}

type GetResultInfoRequest struct {
	DeviceID      int    `json:"deviceId" binding:"required"`
	Model         string `json:"model" binding:"required"`
	AlgorithmType string `json:"algorithmType" binding:"required"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
}

func (trc *TaskResultController) AddTaskResult(c *gin.Context, env *config.Env, tc TaskController) {
	var taskResult models.TaskResult
	if err := c.ShouldBindJSON(&taskResult); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	id, err := trc.TaskResultService.AddTaskResult(c, env, &taskResult, tc.TaskService)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "任务结果添加成功",
		Status:  http.StatusOK,
		Data:    id,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (trc *TaskResultController) GetTaskResultList(c *gin.Context) {
	var req GetResultInfoRequest
	if err := c.BindJSON(&req); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	taskResult, err := trc.TaskResultService.FindTaskResultList(c, req.DeviceID, req.Model, req.AlgorithmType, req.StartTime, req.EndTime)
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
		Data:    taskResult,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (trc *TaskResultController) GetTaskResultById(c *gin.Context) {
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

	res, err := trc.TaskResultService.FindTaskResultById(c, id)
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
		Data:    res,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (trc *TaskResultController) FindAllTaskResults(c *gin.Context) {
	results, err := trc.TaskResultService.GetResultList(c)
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
		Data:    results,
	}
	c.JSON(http.StatusOK, successResponse)

}
