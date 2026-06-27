package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Body struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func OK(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, status int, code string, message string, details interface{}) {
	c.JSON(status, Body{
		Success: false,
		Message: message,
		Error: &ErrorBody{
			Code:    code,
			Details: details,
		},
	})
}

func Error(c *gin.Context, err apperrors.AppError) {
	if err.IsZero() {
		err = apperrors.Internal("unexpected server error")
	}

	Fail(c, err.Status, err.Code, err.Message, err.Details)
}
