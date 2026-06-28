package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) Handler {
	return Handler{service: service}
}

func NewUnavailableHandler() Handler {
	return Handler{}
}

func (h Handler) Register(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid register payload", err.Error()))
		return
	}

	result, err := h.service.Register(c.Request.Context(), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusCreated, "user registered", result)
}

func (h Handler) Login(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid login payload", err.Error()))
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "login successful", result)
}

func (h Handler) Me(c *gin.Context) {
	if !h.available(c) {
		return
	}
	claimsValue, exists := c.Get(ContextClaimsKey)
	if !exists {
		response.Error(c, apperrors.Unauthorized("authentication required"))
		return
	}
	claims := claimsValue.(Claims)

	user, err := h.service.Me(c.Request.Context(), claims.UserID)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "user fetched", user)
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
