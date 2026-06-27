package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restaurant-inventory-api/internal/config"
	"restaurant-inventory-api/internal/http/response"
)

type Handler struct {
	cfg       config.Config
	startedAt time.Time
}

func NewHandler(cfg config.Config, startedAt time.Time) Handler {
	return Handler{
		cfg:       cfg,
		startedAt: startedAt,
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
	}
}
