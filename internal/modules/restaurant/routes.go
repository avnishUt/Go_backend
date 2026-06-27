package restaurant

import (
	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/platform/postgres"
)

func RegisterRoutes(rg *gin.RouterGroup, db *postgres.Client) {
	handler := NewUnavailableHandler()
	if db != nil && db.Pool != nil {
		repo := NewRepository(db.Pool)
		service := NewService(repo)
		handler = NewHandler(&service)
	}

	restaurants := rg.Group("/restaurants")
	{
		restaurants.POST("", handler.CreateRestaurant)
		restaurants.GET("", handler.ListRestaurants)
		restaurants.GET("/:id", handler.GetRestaurant)
		restaurants.PATCH("/:id", handler.UpdateRestaurant)
		restaurants.POST("/:id/branches", handler.CreateBranch)
		restaurants.GET("/:id/branches", handler.ListBranches)
		restaurants.GET("/:id/settings", handler.GetSettings)
		restaurants.PATCH("/:id/settings", handler.UpdateSettings)
	}
}
