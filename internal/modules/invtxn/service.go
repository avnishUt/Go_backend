package invtxn

import (
	"context"
	"errors"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct{ repo Repository }

func NewService(repo Repository) Service { return Service{repo: repo} }

func (s Service) Create(ctx context.Context, req CreateTransactionRequest, createdBy *string) (Transaction, error) {
	result, err := s.repo.Create(ctx, req, createdBy)
	if IsNotFound(err) {
		return result, apperrors.NotFound("inventory item not found")
	}
	if IsInvalidInput(err) {
		return result, apperrors.BadRequest("invalid id", nil)
	}
	if errors.Is(err, ErrNegativeStock) {
		return result, apperrors.BadRequest("inventory stock cannot be negative", nil)
	}
	return result, err
}

func (s Service) List(ctx context.Context, restaurantID string) ([]Transaction, error) {
	result, err := s.repo.List(ctx, restaurantID)
	if IsInvalidInput(err) {
		return result, apperrors.BadRequest("invalid restaurant id", nil)
	}
	return result, err
}
