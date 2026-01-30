package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"
	"wallet-service/internal/domain/entity"
	"wallet-service/internal/domain/repository"
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	RequestID    string `json:"request_id"`
	Timestamp    string `json:"timestamp"`
	HTTPMethod   string `json:"http_method"`
	RequestPath  string `json:"request_path"`
	HTTPStatus   int    `json:"http_status"`
	LatencyMs    int64  `json:"latency_ms"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	UserID       int64  `json:"user_id,omitempty"`
	RequestBody  string `json:"request_body,omitempty"`
	ResponseBody string `json:"response_body,omitempty"`
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (n int, err error) {
	w.body.Write(b)
	n, err = w.ResponseWriter.Write(b)
	return
}

func LoggingMiddleware(apiLogRepo repository.APILogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		timestamp := time.Now()

		requestID := getRequestID(c)
		shouldLogBody := !isAuthEndpoint(c.Request.URL.Path)

		var requestBodyBytes []byte
		if shouldLogBody {
			requestBodyBytes = captureRequestBody(c)
		}

		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		c.Next()

		latency := time.Since(startTime)
		entry := buildLogEntry(c, requestID, timestamp, latency, requestBodyBytes, responseWriter.body.Bytes(), shouldLogBody)

		logJSON, err := json.Marshal(entry)
		if err != nil {
			return
		}

		println(string(logJSON))

		if apiLogRepo != nil {
			go saveLogToDB(apiLogRepo, entry, timestamp)
		}
	}
}

func saveLogToDB(apiLogRepo repository.APILogRepository, entry LogEntry, timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apiLog := &entity.APILog{
		RequestID:    entry.RequestID,
		Timestamp:    timestamp,
		HTTPMethod:   entry.HTTPMethod,
		RequestPath:  entry.RequestPath,
		HTTPStatus:   entry.HTTPStatus,
		LatencyMs:    entry.LatencyMs,
		IPAddress:    entry.IPAddress,
		UserAgent:    entry.UserAgent,
		UserID:       entry.UserID,
		RequestBody:  entry.RequestBody,
		ResponseBody: entry.ResponseBody,
		CreatedAt:    timestamp,
	}

	_ = apiLogRepo.Create(ctx, apiLog)
}

func getRequestID(c *gin.Context) (requestID string) {
	value, exists := c.Get(RequestIDKey)
	if exists {
		requestID, _ = value.(string)
	}
	if requestID == "" {
		requestID = "unknown"
	}
	return
}

func getUserID(c *gin.Context) (userID int64) {
	value, exists := c.Get("user_id")
	if exists {
		userID, _ = value.(int64)
	}
	return
}

func captureRequestBody(c *gin.Context) (bodyBytes []byte) {
	if c.Request.Body == nil {
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return
}

func buildLogEntry(c *gin.Context, requestID string, timestamp time.Time, latency time.Duration, requestBody []byte, responseBody []byte, shouldLogBody bool) (entry LogEntry) {
	entry = LogEntry{
		RequestID:   requestID,
		Timestamp:   timestamp.Format(time.RFC3339),
		HTTPMethod:  c.Request.Method,
		RequestPath: c.Request.URL.Path,
		HTTPStatus:  c.Writer.Status(),
		LatencyMs:   latency.Milliseconds(),
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		UserID:      getUserID(c),
	}

	if shouldLogBody {
		if len(requestBody) > 0 {
			redactedRequest := utils.RedactSensitiveData(requestBody)
			entry.RequestBody = string(redactedRequest)
		}

		if len(responseBody) > 0 {
			redactedResponse := utils.RedactSensitiveData(responseBody)
			entry.ResponseBody = string(redactedResponse)
		}
	}

	return
}

func isAuthEndpoint(path string) (result bool) {
	authPaths := []string{
		"/api/login",
		"/api/register",
	}

	for _, authPath := range authPaths {
		if strings.HasSuffix(path, authPath) {
			result = true
			return
		}
	}

	result = false
	return
}
