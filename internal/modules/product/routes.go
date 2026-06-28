package product

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

	categories := rg.Group("/product-categories", authMiddleware, auth.RequireRoles("admin", "store_manager"))
	{
		categories.POST("", handler.CreateCategory)
		categories.GET("", handler.ListCategories)
	}

	products := rg.Group("/products", authMiddleware)
	{
		products.POST("", auth.RequireRoles("admin", "store_manager"), handler.CreateProduct)
		products.GET("", handler.ListProducts)
		products.GET("/:id", handler.GetProduct)
		products.PATCH("/:id", auth.RequireRoles("admin", "store_manager"), handler.UpdateProduct)
	}
}
