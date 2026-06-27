package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			if origin == "" && containsWildcard(allowedOrigins) {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}

		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-ID,X-Restaurant-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if containsWildcard(allowedOrigins) {
		return true
	}

	for _, allowedOrigin := range allowedOrigins {
		if strings.EqualFold(strings.TrimSpace(allowedOrigin), origin) {
			return true
		}
	}

	return false
}

func containsWildcard(values []string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "*" {
			return true
		}
	}

	return false
}
