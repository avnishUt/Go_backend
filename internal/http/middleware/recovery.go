package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/response"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get(RequestIDKey)
		log.Printf("panic recovered request_id=%v error=%v", requestID, recovered)

		response.Error(c, apperrors.New(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
			nil,
		))
		c.Abort()
	})
}
