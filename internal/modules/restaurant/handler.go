package restaurant

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

func (h Handler) CreateRestaurant(c *gin.Context) {
	if !h.available(c) {
		return
	}

	var req CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid restaurant payload", err.Error()))
		return
	}

	restaurant, err := h.service.CreateRestaurant(c.Request.Context(), req)
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusCreated, "restaurant created", restaurant)
}

func (h Handler) ListRestaurants(c *gin.Context) {
	if !h.available(c) {
		return
	}

	restaurants, err := h.service.ListRestaurants(c.Request.Context())
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusOK, "restaurants fetched", restaurants)
}

func (h Handler) GetRestaurant(c *gin.Context) {
	if !h.available(c) {
		return
	}

	restaurant, err := h.service.GetRestaurant(c.Request.Context(), c.Param("id"))
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusOK, "restaurant fetched", restaurant)
}

func (h Handler) UpdateRestaurant(c *gin.Context) {
	if !h.available(c) {
		return
	}

	var req UpdateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid restaurant payload", err.Error()))
		return
	}

	restaurant, err := h.service.UpdateRestaurant(c.Request.Context(), c.Param("id"), req)
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusOK, "restaurant updated", restaurant)
}

func (h Handler) CreateBranch(c *gin.Context) {
	if !h.available(c) {
		return
	}

	var req CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid branch payload", err.Error()))
		return
	}

	branch, err := h.service.CreateBranch(c.Request.Context(), c.Param("id"), req)
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusCreated, "branch created", branch)
}

func (h Handler) ListBranches(c *gin.Context) {
	if !h.available(c) {
		return
	}

	branches, err := h.service.ListBranches(c.Request.Context(), c.Param("id"))
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusOK, "branches fetched", branches)
}

func (h Handler) GetSettings(c *gin.Context) {
	if !h.available(c) {
		return
	}

	settings, err := h.service.GetSettings(c.Request.Context(), c.Param("id"))
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusOK, "restaurant settings fetched", settings)
}

func (h Handler) UpdateSettings(c *gin.Context) {
	if !h.available(c) {
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid settings payload", err.Error()))
		return
	}

	settings, err := h.service.UpdateSettings(c.Request.Context(), c.Param("id"), req)
	if h.handleError(c, err) {
		return
	}

	response.OK(c, http.StatusOK, "restaurant settings updated", settings)
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
