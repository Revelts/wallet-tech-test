package middleware

import (
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-Id"
	RequestIDKey    = "request_id"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)

		if requestID == "" {
			requestID = utils.GenerateRequestID()
		}

		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}
