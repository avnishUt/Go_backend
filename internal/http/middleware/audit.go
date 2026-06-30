package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"restaurant-inventory-api/internal/modules/auth"
)

func Audit(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if db == nil || !isAuditableMethod(c.Request.Method) {
			return
		}

		requestID, _ := c.Get(RequestIDKey)
		var userID *string
		var restaurantID *string
		if claimsValue, exists := c.Get(auth.ContextClaimsKey); exists {
			claims := claimsValue.(auth.Claims)
			userID = &claims.UserID
			restaurantID = claims.RestaurantID
		}
		if restaurantID == nil {
			if headerValue := c.GetHeader("X-Restaurant-ID"); headerValue != "" {
				restaurantID = &headerValue
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := db.Exec(ctx, `
			INSERT INTO audit_logs (
				request_id, user_id, restaurant_id, method, path, status_code, client_ip, user_agent
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, requestID, userID, restaurantID, c.Request.Method, c.Request.URL.Path,
			c.Writer.Status(), c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			log.Printf("audit log write failed: %v", err)
		}
	}
}

func isAuditableMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}
