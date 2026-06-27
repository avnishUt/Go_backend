package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		requestID, _ := c.Get(RequestIDKey)
		log.Printf(
			"request_id=%v method=%s path=%s status=%d latency=%s client_ip=%s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(startedAt),
			c.ClientIP(),
		)
	}
}
