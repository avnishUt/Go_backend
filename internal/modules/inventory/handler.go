package inventory

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

func (h Handler) CreateItem(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid inventory payload", err.Error()))
		return
	}
	result, err := h.service.CreateItem(c.Request.Context(), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusCreated, "inventory item created", result)
}

func (h Handler) ListItems(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.service.ListItems(c.Request.Context(), c.Query("restaurant_id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "inventory items fetched", result)
}

func (h Handler) GetItem(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.service.GetItem(c.Request.Context(), c.Param("id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "inventory item fetched", result)
}

func (h Handler) UpdateItem(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid inventory payload", err.Error()))
		return
	}
	result, err := h.service.UpdateItem(c.Request.Context(), c.Param("id"), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "inventory item updated", result)
}

func (h Handler) AdjustStock(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid stock payload", err.Error()))
		return
	}
	result, err := h.service.AdjustStock(c.Request.Context(), c.Param("id"), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "inventory stock adjusted", result)
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
