package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/config"
	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/health"
	"restaurant-inventory-api/internal/http/middleware"
	"restaurant-inventory-api/internal/http/response"
)

func New(cfg config.Config) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.CORS(cfg.AllowedOrigins),
	)

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, apperrors.NotFound("route not found"))
	})

	r.NoMethod(func(c *gin.Context) {
		response.Error(c, apperrors.MethodNotAllowed("method not allowed"))
	})

	healthHandler := health.NewHandler(cfg, time.Now())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "restaurant inventory api",
			"data": gin.H{
				"app":     cfg.AppName,
				"env":     cfg.Env,
				"version": "v1",
			},
		})
	})
	r.GET("/health", healthHandler.Service)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.API)
	}

	return r
}
