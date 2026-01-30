package middleware

import (
	"strings"
	"time"
	"wallet-service/internal/app_error"
	"wallet-service/internal/infrastructure/jwt"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Code       int         `json:"code"`
	AccessTime string      `json:"accessTime"`
}

func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			respondError(c, app_error.AuthHeaderRequired)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(c, app_error.InvalidAuthHeader)
			c.Abort()
			return
		}

		token := parts[1]

		userID, err := jwtService.ValidateToken(token)
		if err != nil {
			respondError(c, app_error.InvalidToken)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func respondError(c *gin.Context, appErr *app_error.AppError) {
	httpStatus, code, message, _ := appErr.GetErrors()

	response := ErrorResponse{
		Success:    false,
		Data:       gin.H{"error": message},
		Code:       code,
		AccessTime: time.Now().Format(time.RFC3339),
	}
	c.JSON(httpStatus, response)
}
