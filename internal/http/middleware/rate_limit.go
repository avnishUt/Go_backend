package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/http/response"
)

type rateBucket struct {
	windowStart time.Time
	count       int
}

func RateLimit(enabled bool, requestsPerMinute int) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := map[string]rateBucket{}

	return func(c *gin.Context) {
		if !enabled || requestsPerMinute <= 0 {
			c.Next()
			return
		}

		key := c.ClientIP()
		now := time.Now()

		mu.Lock()
		bucket := buckets[key]
		if bucket.windowStart.IsZero() || now.Sub(bucket.windowStart) >= time.Minute {
			bucket = rateBucket{windowStart: now}
		}
		bucket.count++
		buckets[key] = bucket
		allowed := bucket.count <= requestsPerMinute
		mu.Unlock()

		if !allowed {
			response.Fail(c, http.StatusTooManyRequests, "RATE_LIMITED", "too many requests", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
