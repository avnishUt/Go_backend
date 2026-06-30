package report

import (
	"context"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct{ repo Repository }

func NewService(repo Repository) Service { return Service{repo: repo} }

func (s Service) SalesSummary(ctx context.Context, restaurantID string) (SalesSummary, error) {
	summary, err := s.repo.SalesSummary(ctx, restaurantID)
	if IsInvalidInput(err) {
		return summary, apperrors.BadRequest("invalid restaurant id", nil)
	}
	return summary, err
}

func (s Service) LowStock(ctx context.Context, restaurantID string) ([]LowStockItem, error) {
	items, err := s.repo.LowStock(ctx, restaurantID)
	if IsInvalidInput(err) {
		return items, apperrors.BadRequest("invalid restaurant id", nil)
	}
	return items, err
}
