package handler

import (
	"wallet-service/internal/app_error"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
	getBalanceUsecase *usecase.GetBalanceUsecase
}

func NewBalanceHandler(getBalanceUsecase *usecase.GetBalanceUsecase) *BalanceHandler {
	return &BalanceHandler{
		getBalanceUsecase: getBalanceUsecase,
	}
}

type BalanceResponse struct {
	Balance int64 `json:"balance"`
}

func (h *BalanceHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		SendResponse(c, app_error.Unauthorized, nil)
		return
	}

	input := usecase.GetBalanceInput{
		UserID: userID.(int64),
	}

	output, err := h.getBalanceUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}

	response := BalanceResponse{
		Balance: output.Balance,
	}

	SendResponse(c, nil, response)
}
