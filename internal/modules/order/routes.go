package order

import (
	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/modules/auth"
	"restaurant-inventory-api/internal/platform/postgres"
)

func RegisterRoutes(rg *gin.RouterGroup, db *postgres.Client, authMiddleware gin.HandlerFunc) {
	handler := NewUnavailableHandler()
	if db != nil && db.Pool != nil {
		repo := NewRepository(db.Pool)
		service := NewService(repo)
		handler = NewHandler(&service)
	}

	routes := rg.Group("/orders", authMiddleware)
	{
		routes.POST("", handler.Create)
		routes.GET("", handler.List)
		routes.GET("/:id", handler.Get)
		routes.PATCH("/:id/status", auth.RequireRoles("admin", "store_manager"), handler.UpdateStatus)
	}
}
