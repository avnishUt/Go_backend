package product

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

func (h Handler) CreateCategory(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid category payload", err.Error()))
		return
	}
	result, err := h.service.CreateCategory(c.Request.Context(), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusCreated, "category created", result)
}

func (h Handler) ListCategories(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.service.ListCategories(c.Request.Context(), c.Query("restaurant_id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "categories fetched", result)
}

func (h Handler) CreateProduct(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid product payload", err.Error()))
		return
	}
	result, err := h.service.CreateProduct(c.Request.Context(), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusCreated, "product created", result)
}

func (h Handler) ListProducts(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.service.ListProducts(c.Request.Context(), c.Query("restaurant_id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "products fetched", result)
}

func (h Handler) GetProduct(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.service.GetProduct(c.Request.Context(), c.Param("id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "product fetched", result)
}

func (h Handler) UpdateProduct(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid product payload", err.Error()))
		return
	}
	result, err := h.service.UpdateProduct(c.Request.Context(), c.Param("id"), req)
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "product updated", result)
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
