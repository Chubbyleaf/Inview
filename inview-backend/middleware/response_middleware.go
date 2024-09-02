package middleware

type CommonRequest struct {
	UnionId string `json:"unionId"`
}

type LoginRequest struct {
	UserName string `form:"userName" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LogoutRequest struct {
	UserName string `form:"userName" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LogoutResponse struct {
	UserName string `form:"userName" binding:"required"`
}

type SignupRequest struct {
	Name     string `form:"name" binding:"required"`
	UserName string `form:"userName" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"` // 使用 interface{} 来支持任意数据类型
}

type RefreshTokenRequest struct {
	RefreshToken string `form:"refreshToken" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateCameraRequest struct {
	DeviceID int    `json:"deviceId"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
}
