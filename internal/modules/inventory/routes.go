package inventory

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

	items := rg.Group("/inventory/items", authMiddleware, auth.RequireRoles("admin", "store_manager"))
	{
		items.POST("", handler.CreateItem)
		items.GET("", handler.ListItems)
		items.GET("/:id", handler.GetItem)
		items.PATCH("/:id", handler.UpdateItem)
		items.POST("/:id/adjust", handler.AdjustStock)
	}
}
