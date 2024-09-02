package controller

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"insense-local/middleware"
	"insense-local/models"
	"insense-local/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserServiceInterface
}

type UpdatePasswordRequest struct {
	ID          string `json:"_id" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type UpdateUserInfoRequest struct {
	ID       string `json:"_id" binding:"required"`
	UserName string `json:"userName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err := uc.UserService.AddUser(c, &user)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "用户注册成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (uc *UserController) Login(c *gin.Context) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := uc.UserService.LoginUser(c, req.UserName, req.Password)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "登录成功",
		Status:  http.StatusOK,
		Data:    user,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (uc *UserController) UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest

	if err := c.BindJSON(&req); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid user ID",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err = uc.UserService.UpdatePassword(c, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "密码更新成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (uc *UserController) UpdateUserInfo(c *gin.Context) {
	var req UpdateUserInfoRequest
	if err := c.BindJSON(&req); err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid user ID",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user := &models.User{
		ID:       userID,
		UserName: req.UserName,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	err = uc.UserService.UpdateUserInfo(c, user)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "用户信息更新成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid user ID",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err = uc.UserService.DeleteUser(c, userID)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "用户删除成功",
		Status:  http.StatusOK,
	}
	c.JSON(http.StatusOK, successResponse)
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: "Invalid user ID",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := uc.UserService.GetUserByID(c, id)
	if err != nil {
		errorResponse := middleware.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := middleware.SuccessResponse{
		Message: "成功获取用户信息",
		Status:  http.StatusOK,
		Data:    user,
	}
	c.JSON(http.StatusOK, successResponse)
}
