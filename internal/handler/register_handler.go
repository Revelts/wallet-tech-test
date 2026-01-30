package handler

import (
	"wallet-service/internal/app_error"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	registerUsecase *usecase.RegisterUsecase
}

func NewRegisterHandler(registerUsecase *usecase.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		registerUsecase: registerUsecase,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Pin      string `json:"pin" binding:"required,len=6,numeric"`
}

type RegisterResponse struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var req RegisterRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		SendResponse(c, app_error.InvalidJsonRequest, nil)
		return
	}

	input := usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Pin:      req.Pin,
	}

	output, err := h.registerUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}

	response := RegisterResponse{
		UserID: output.UserID,
		Email:  output.Email,
	}

	SendResponse(c, nil, response)
}
