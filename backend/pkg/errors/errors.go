package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       int
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(message string, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewAppErrorWithErr creates an app error with an underlying error
func NewAppErrorWithErr(message string, statusCode int, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// BadRequest returns a 400 error
func BadRequest(message string) *AppError {
	return NewAppError(message, http.StatusBadRequest)
}

// Unauthorized returns a 401 error
func Unauthorized(message string) *AppError {
	return NewAppError(message, http.StatusUnauthorized)
}

// Forbidden returns a 403 error
func Forbidden(message string) *AppError {
	return NewAppError(message, http.StatusForbidden)
}

// NotFound returns a 404 error
func NotFound(message string) *AppError {
	return NewAppError(message, http.StatusNotFound)
}

// Conflict returns a 409 error
func Conflict(message string) *AppError {
	return NewAppError(message, http.StatusConflict)
}

// InternalServerError returns a 500 error
func InternalServerError(message string) *AppError {
	return NewAppError(message, http.StatusInternalServerError)
}

// InternalServerErrorWithErr returns a 500 error with underlying error
func InternalServerErrorWithErr(message string, err error) *AppError {
	return NewAppErrorWithErr(message, http.StatusInternalServerError, err)
}

// IsAppError checks if error is an AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// GetStatusCode extracts status code from error
func GetStatusCode(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// GetMessage extracts message from error
func GetMessage(err error) string {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Message
	}
	return "Internal server error"
}
