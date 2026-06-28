package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/config"
	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/response"
)

const ContextClaimsKey = "auth_claims"

func Authenticate(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, apperrors.Unauthorized("bearer token required"))
			c.Abort()
			return
		}

		claims, err := ParseToken(strings.TrimPrefix(header, "Bearer "), cfg.JWTSecret)
		if err != nil {
			var appErr apperrors.AppError
			if errors.As(err, &appErr) {
				response.Error(c, appErr)
			} else {
				response.Error(c, apperrors.Unauthorized("invalid token"))
			}
			c.Abort()
			return
		}

		c.Set(ContextClaimsKey, claims)
		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, role := range roles {
		allowed[role] = true
	}

	return func(c *gin.Context) {
		claimsValue, exists := c.Get(ContextClaimsKey)
		if !exists {
			response.Error(c, apperrors.Unauthorized("authentication required"))
			c.Abort()
			return
		}

		claims := claimsValue.(Claims)
		for _, role := range claims.Roles {
			if role == "super_admin" || allowed[role] {
				c.Next()
				return
			}
		}

		response.Error(c, apperrors.Forbidden("permission denied"))
		c.Abort()
	}
}
