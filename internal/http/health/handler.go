package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/config"
	"restaurant-inventory-api/internal/http/response"
)

type Check struct {
	Name string
	Ping func(ctx context.Context) error
}

type Handler struct {
	cfg       config.Config
	startedAt time.Time
	checks    []Check
}

func NewHandler(cfg config.Config, startedAt time.Time, checks []Check) Handler {
	return Handler{
		cfg:       cfg,
		startedAt: startedAt,
		checks:    checks,
	}
}

func (h Handler) Service(c *gin.Context) {
	response.OK(c, http.StatusOK, "service is healthy", h.payload("service"))
}

func (h Handler) API(c *gin.Context) {
	response.OK(c, http.StatusOK, "api is healthy", h.payload("api"))
}

func (h Handler) payload(scope string) gin.H {
	return gin.H{
		"app":        h.cfg.AppName,
		"env":        h.cfg.Env,
		"scope":      scope,
		"started_at": h.startedAt.UTC().Format(time.RFC3339),
		"uptime":     time.Since(h.startedAt).String(),
		"checks":     h.checkStatuses(),
	}
}

func (h Handler) checkStatuses() gin.H {
	statuses := gin.H{}

	for _, check := range h.checks {
		if check.Ping == nil {
			statuses[check.Name] = gin.H{
				"status": "disabled",
			}
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := check.Ping(ctx)
		cancel()

		if err != nil {
			statuses[check.Name] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			continue
		}

		statuses[check.Name] = gin.H{
			"status": "healthy",
		}
	}

	return statuses
}
