package app_error

import "net/http"

type AppError struct {
	HTTPStatus  int
	Code        int
	Message     string
	ErrorData   interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) GetErrors() (httpStatus int, code int, message string, data interface{}) {
	httpStatus = e.HTTPStatus
	code = e.Code
	message = e.Message
	data = e.ErrorData
	return
}

func New(httpStatus int, code int, message string) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		ErrorData:  nil,
	}
}

func NewWithData(httpStatus int, code int, message string, data interface{}) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		ErrorData:  data,
	}
}

var (
	InvalidJsonRequest      = New(http.StatusBadRequest, 400, "invalid json request")
	EmailAlreadyExists      = New(http.StatusBadRequest, 40001, "email already registered")
	InvalidCredentials      = New(http.StatusUnauthorized, 40002, "invalid email or password")
	InvalidToken            = New(http.StatusUnauthorized, 40003, "invalid or expired token")
	InvalidPIN              = New(http.StatusBadRequest, 40004, "invalid PIN")
	InsufficientBalance     = New(http.StatusBadRequest, 40005, "insufficient balance")
	InvalidAmount           = New(http.StatusBadRequest, 40006, "amount must be greater than zero")
	WalletNotFound          = New(http.StatusNotFound, 40007, "wallet not found")
	UserNotFound            = New(http.StatusNotFound, 40008, "user not found")
	Unauthorized            = New(http.StatusUnauthorized, 401, "unauthorized")
	AuthHeaderRequired      = New(http.StatusUnauthorized, 401, "authorization header required")
	InvalidAuthHeader       = New(http.StatusUnauthorized, 401, "invalid authorization header format")
	FailedGenerateToken     = New(http.StatusInternalServerError, 500, "failed to generate token")
	InternalServerError     = New(http.StatusInternalServerError, 500, "internal server error")
)
