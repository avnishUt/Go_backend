package auth

import (
	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/config"
	"restaurant-inventory-api/internal/platform/postgres"
)

func RegisterRoutes(rg *gin.RouterGroup, db *postgres.Client, cfg config.AuthConfig) {
	handler := NewUnavailableHandler()
	if db != nil && db.Pool != nil {
		repo := NewRepository(db.Pool)
		service := NewService(repo, cfg)
		handler = NewHandler(&service)
	}

	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", handler.Register)
		authRoutes.POST("/login", handler.Login)
		authRoutes.GET("/me", Authenticate(cfg), handler.Me)
	}
}
