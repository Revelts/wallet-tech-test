package handler

import (
	"net/http"
	"time"
	"wallet-service/internal/app_error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Code       int         `json:"code"`
	AccessTime string      `json:"accessTime"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	response := Response{
		Success:    true,
		Data:       data,
		Code:       200,
		AccessTime: time.Now().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, response)
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		appErr, ok := err.(*app_error.AppError)
		if !ok {
			appErr = app_error.InternalServerError
		}

		httpStatus, code, message, errorData := appErr.GetErrors()

		if errorData == nil {
			errorData = gin.H{"error": message}
		}

		response := Response{
			Success:    false,
			Data:       errorData,
			Code:       code,
			AccessTime: time.Now().Format(time.RFC3339),
		}
		c.JSON(httpStatus, response)
		return
	}

	SuccessResponse(c, data)
}
