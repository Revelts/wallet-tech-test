package handler

import (
	"wallet-service/internal/app_error"
	"wallet-service/internal/infrastructure/jwt"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginUsecase *usecase.LoginUsecase
	jwtService   *jwt.JWTService
}

func NewLoginHandler(loginUsecase *usecase.LoginUsecase, jwtService *jwt.JWTService) *LoginHandler {
	return &LoginHandler{
		loginUsecase: loginUsecase,
		jwtService:   jwtService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var req LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		SendResponse(c, app_error.InvalidJsonRequest, nil)
		return
	}

	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.loginUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}

	token, err := h.jwtService.GenerateToken(output.UserID)
	if err != nil {
		SendResponse(c, app_error.FailedGenerateToken, nil)
		return
	}

	response := LoginResponse{
		AccessToken: token,
	}

	SendResponse(c, nil, response)
}
