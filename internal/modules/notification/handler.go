package notification

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/response"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) Handler { return Handler{service: service} }
func NewUnavailableHandler() Handler      { return Handler{} }

func (h Handler) Create(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid notification payload", err.Error()))
		return
	}
	item, err := h.service.Create(c.Request.Context(), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusCreated, "notification queued", item)
}

func (h Handler) List(c *gin.Context) {
	if !h.available(c) {
		return
	}
	items, err := h.service.List(c.Request.Context(), c.Query("restaurant_id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "notifications fetched", items)
}

func (h Handler) Mark(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req MarkSentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid notification status payload", err.Error()))
		return
	}
	item, err := h.service.Mark(c.Request.Context(), c.Param("id"), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "notification updated", item)
}

func (h Handler) available(c *gin.Context) bool {
	if h.service != nil {
		return true
	}
	response.Error(c, apperrors.ServiceUnavailable("postgres database is not configured"))
	return false
}

func (h Handler) handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var appErr apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr)
		return true
	}
	response.Error(c, apperrors.Internal("internal server error"))
	return true
}
