package controller

import (
	"github.com/gin-gonic/gin"
	"insense-local/config"
	"insense-local/middleware"
	"insense-local/services"
	"net/http"
)

type SystemController struct {
	SystemService services.SystemServiceInterface
}

func (sc *SystemController) GetSystemInfo(c *gin.Context) {

	info := sc.SystemService.SystemInfo()

	successResponse := middleware.SuccessResponse{
		Message: "成功获取本机信息",
		Status:  http.StatusOK,
		Data:    info,
	}

	c.JSON(http.StatusOK, successResponse)
}

func (sc *SystemController) GetSystemLog(c *gin.Context, env *config.Env) {

	info, err := sc.SystemService.SystemLog(env)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	successResponse := middleware.SuccessResponse{
		Message: "成功获取系统日志",
		Status:  http.StatusOK,
		Data:    info,
	}

	c.JSON(http.StatusOK, successResponse)
}

func (sc *SystemController) GetNetworkInfo(c *gin.Context) {

	info, err := sc.SystemService.NetworkInfo()
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	successResponse := middleware.SuccessResponse{
		Message: "成功获取网络信息",
		Status:  http.StatusOK,
		Data:    info,
	}

	c.JSON(http.StatusOK, successResponse)
}
