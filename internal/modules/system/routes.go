package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/http/middleware"
	"restaurant-inventory-api/internal/http/response"
)

func RegisterRoutes(r *gin.Engine, collector *middleware.MetricsCollector) {
	r.GET("/metrics", func(c *gin.Context) {
		response.OK(c, http.StatusOK, "metrics fetched", collector.Snapshot())
	})
}
