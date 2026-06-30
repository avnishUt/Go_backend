package report

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "restaurant-inventory-api/internal/errors"
	"restaurant-inventory-api/internal/http/response"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) Handler { return Handler{service: service} }
func NewUnavailableHandler() Handler      { return Handler{} }

func (h Handler) SalesSummary(c *gin.Context) {
	if !h.available(c) {
		return
	}
	summary, err := h.service.SalesSummary(c.Request.Context(), c.Query("restaurant_id"))
	if h.handleError(c, err) {
		return
	}
	response.OK(c, http.StatusOK, "sales summary fetched", summary)
}

func (h Handler) LowStockCSV(c *gin.Context) {
	if !h.available(c) {
		return
	}
	items, err := h.service.LowStock(c.Request.Context(), c.Query("restaurant_id"))
	if h.handleError(c, err) {
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=low-stock.csv")
	writer := csv.NewWriter(c.Writer)
	_ = writer.Write([]string{"product_name", "sku", "current_stock", "minimum_stock", "reorder_level", "unit"})
	for _, item := range items {
		_ = writer.Write([]string{
			item.ProductName,
			item.SKU,
			strconv.FormatFloat(item.CurrentStock, 'f', 3, 64),
			strconv.FormatFloat(item.MinimumStock, 'f', 3, 64),
			strconv.FormatFloat(item.ReorderLevel, 'f', 3, 64),
			item.Unit,
		})
	}
	writer.Flush()
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
