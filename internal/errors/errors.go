package errors

import "net/http"

type AppError struct {
	Status  int
	Code    string
	Message string
	Details interface{}
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) IsZero() bool {
	return e.Status == 0 && e.Code == "" && e.Message == ""
}

func New(status int, code string, message string, details interface{}) AppError {
	return AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func Internal(message string) AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

func BadRequest(message string, details interface{}) AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(message string) AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(message string) AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message, nil)
}

func NotFound(message string) AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message, nil)
}

func MethodNotAllowed(message string) AppError {
	return New(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", message, nil)
}
