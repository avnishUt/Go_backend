package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/config"
	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/health"
	"restaurant-inventory-api/internal/http/middleware"
	"restaurant-inventory-api/internal/http/response"
	"restaurant-inventory-api/internal/modules/auth"
	"restaurant-inventory-api/internal/modules/inventory"
	"restaurant-inventory-api/internal/modules/invtxn"
	"restaurant-inventory-api/internal/modules/notification"
	"restaurant-inventory-api/internal/modules/order"
	"restaurant-inventory-api/internal/modules/payment"
	"restaurant-inventory-api/internal/modules/product"
	"restaurant-inventory-api/internal/modules/report"
	restaurant "restaurant-inventory-api/internal/modules/restaurant"
	"restaurant-inventory-api/internal/modules/system"
	"restaurant-inventory-api/internal/platform/postgres"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Dependencies struct {
	Postgres Pinger
	Mongo    Pinger
	Database *postgres.Client
}

func New(cfg config.Config, deps Dependencies) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	metricsCollector := middleware.NewMetricsCollector()
	r.HandleMethodNotAllowed = true
	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.Metrics(metricsCollector),
		middleware.SecurityHeaders(),
		middleware.BodyLimit(cfg.Security.MaxBodyBytes),
		middleware.RateLimit(cfg.RateLimit.Enabled, cfg.RateLimit.RequestsPerMinute),
		middleware.CORS(cfg.AllowedOrigins),
	)
	if deps.Database != nil && deps.Database.Pool != nil {
		r.Use(middleware.Audit(deps.Database.Pool))
	}

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, apperrors.NotFound("route not found"))
	})

	r.NoMethod(func(c *gin.Context) {
		response.Error(c, apperrors.MethodNotAllowed("method not allowed"))
	})

	healthHandler := health.NewHandler(cfg, time.Now(), []health.Check{
		{
			Name: "postgres",
			Ping: pingFunc(deps.Postgres),
		},
		{
			Name: "mongo",
			Ping: pingFunc(deps.Mongo),
		},
	})

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
	system.RegisterRoutes(r, metricsCollector)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.API)
		auth.RegisterRoutes(v1, deps.Database, cfg.Auth)
		restaurant.RegisterRoutes(v1, deps.Database)
		authMiddleware := auth.Authenticate(cfg.Auth)
		product.RegisterRoutes(v1, deps.Database, authMiddleware)
		inventory.RegisterRoutes(v1, deps.Database, authMiddleware)
		invtxn.RegisterRoutes(v1, deps.Database, authMiddleware)
		order.RegisterRoutes(v1, deps.Database, authMiddleware)
		payment.RegisterRoutes(v1, deps.Database, authMiddleware)
		notification.RegisterRoutes(v1, deps.Database, authMiddleware)
		report.RegisterRoutes(v1, deps.Database, authMiddleware)
	}

	return r
}

func pingFunc(pinger Pinger) func(ctx context.Context) error {
	if pinger == nil {
		return nil
	}

	return pinger.Ping
}
