package order

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/response"
	"restaurant-inventory-api/internal/modules/auth"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) Handler { return Handler{service: service} }
func NewUnavailableHandler() Handler      { return Handler{} }

func (h Handler) Create(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid order payload", err.Error()))
		return
	}
	claimsValue, _ := c.Get(auth.ContextClaimsKey)
	claims := claimsValue.(auth.Claims)
	result, err := h.service.Create(c.Request.Context(), req, &claims.UserID)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusCreated, "order created", result)
}

func (h Handler) List(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.service.List(c.Request.Context(), c.Query("restaurant_id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "orders fetched", result)
}

func (h Handler) Get(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "order fetched", result)
}

func (h Handler) UpdateStatus(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid status payload", err.Error()))
		return
	}
	claimsValue, _ := c.Get(auth.ContextClaimsKey)
	claims := claimsValue.(auth.Claims)
	result, err := h.service.UpdateStatus(c.Request.Context(), c.Param("id"), req, &claims.UserID)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "order status updated", result)
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
