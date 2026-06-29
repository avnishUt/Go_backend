package order

import (
	"context"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct{ repo Repository }

func NewService(repo Repository) Service { return Service{repo: repo} }

func (s Service) Create(ctx context.Context, req CreateOrderRequest, userID *string) (Order, error) {
	result, err := s.repo.Create(ctx, req, userID)
	return result, mapErr(err, "restaurant, branch, user, or product does not exist")
}

func (s Service) List(ctx context.Context, restaurantID string) ([]Order, error) {
	result, err := s.repo.List(ctx, restaurantID)
	return result, mapErr(err, "")
}

func (s Service) Get(ctx context.Context, id string) (Order, error) {
	result, err := s.repo.Get(ctx, id)
	if IsNotFound(err) {
		return result, apperrors.NotFound("order not found")
	}
	return result, mapErr(err, "")
}

func (s Service) UpdateStatus(ctx context.Context, id string, req UpdateStatusRequest, userID *string) (Order, error) {
	result, err := s.repo.UpdateStatus(ctx, id, req, userID)
	if IsNotFound(err) {
		return result, apperrors.NotFound("order not found")
	}
	return result, mapErr(err, "")
}

func mapErr(err error, fkMsg string) error {
	if err == nil {
		return nil
	}
	if IsInvalidInput(err) {
		return apperrors.BadRequest("invalid id", nil)
	}
	if IsForeignKeyViolation(err) && fkMsg != "" {
		return apperrors.BadRequest(fkMsg, nil)
	}
	return err
}
