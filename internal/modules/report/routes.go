package report

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

	routes := rg.Group("/reports", authMiddleware, auth.RequireRoles("admin", "store_manager"))
	{
		routes.GET("/sales-summary", handler.SalesSummary)
		routes.GET("/low-stock.csv", handler.LowStockCSV)
	}
}
