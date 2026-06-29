package payment

import (
	"context"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct{ repo Repository }

func NewService(repo Repository) Service { return Service{repo: repo} }

func (s Service) Create(ctx context.Context, req CreatePaymentRequest) (Payment, error) {
	result, err := s.repo.Create(ctx, req)
	if IsNotFound(err) {
		return result, apperrors.NotFound("order not found")
	}
	return result, mapErr(err)
}

func (s Service) CompleteMockUPI(ctx context.Context, id string, req CompleteMockUPIRequest) (Payment, error) {
	result, err := s.repo.CompleteMockUPI(ctx, id, req)
	if IsNotFound(err) {
		return result, apperrors.NotFound("mock upi payment not found")
	}
	return result, mapErr(err)
}

func (s Service) Get(ctx context.Context, id string) (Payment, error) {
	result, err := s.repo.Get(ctx, id)
	if IsNotFound(err) {
		return result, apperrors.NotFound("payment not found")
	}
	return result, mapErr(err)
}

func (s Service) List(ctx context.Context, restaurantID string) ([]Payment, error) {
	result, err := s.repo.List(ctx, restaurantID)
	return result, mapErr(err)
}

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	if IsInvalidInput(err) {
		return apperrors.BadRequest("invalid id", nil)
	}
	return err
}
