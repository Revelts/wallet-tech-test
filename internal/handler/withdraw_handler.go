package handler

import (
	"wallet-service/internal/app_error"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type WithdrawHandler struct {
	withdrawUsecase *usecase.WithdrawUsecase
}

func NewWithdrawHandler(withdrawUsecase *usecase.WithdrawUsecase) *WithdrawHandler {
	return &WithdrawHandler{
		withdrawUsecase: withdrawUsecase,
	}
}

type WithdrawRequest struct {
	Amount int64  `json:"amount" binding:"required,gt=0"`
	Pin    string `json:"pin" binding:"required,len=6,numeric"`
}

type WithdrawResponse struct {
	TransactionID int64 `json:"transaction_id"`
	NewBalance    int64 `json:"new_balance"`
}

func (h *WithdrawHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		SendResponse(c, app_error.Unauthorized, nil)
		return
	}

	var req WithdrawRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		SendResponse(c, app_error.InvalidJsonRequest, nil)
		return
	}

	input := usecase.WithdrawInput{
		UserID: userID.(int64),
		Amount: req.Amount,
		Pin:    req.Pin,
	}

	output, err := h.withdrawUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}

	response := WithdrawResponse{
		TransactionID: output.TransactionID,
		NewBalance:    output.NewBalance,
	}

	SendResponse(c, nil, response)
}
