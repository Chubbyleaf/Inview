package controller

import (
	"insense-local/middleware"
	"insense-local/models"
	"insense-local/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CameraController struct {
	CameraService services.CameraServiceInterface
}

func (cc *CameraController) GetCameras(c *gin.Context) {

	cameras, err := cc.CameraService.FetchCameras(c)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "成功获取设备列表",
		Status:  http.StatusOK,
		Data:    cameras,
	}

	c.JSON(http.StatusOK, successResponse)
}
func (cc *CameraController) CreateCamera(c *gin.Context) {
	var camera models.Camera
	if err := c.ShouldBindJSON(&camera); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err := cc.CameraService.CreateCamera(c, &camera)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "设备创建成功",
		Status:  http.StatusOK,
		Data:    camera,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (cc *CameraController) DeleteCameraByDeviceID(c *gin.Context, tc TaskController) {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid device ID",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err = cc.CameraService.DeleteCamera(c, deviceID, tc.TaskService)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "设备删除成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (cc *CameraController) UpdateCamera(c *gin.Context) {
	var req middleware.UpdateCameraRequest
	if err := c.BindJSON(&req); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err := cc.CameraService.UpdateCamera(c, req.DeviceID, req.Name, req.Remark)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "设备信息更新成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}
