package notification

import (
	"context"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct{ repo Repository }

func NewService(repo Repository) Service { return Service{repo: repo} }

func (s Service) Create(ctx context.Context, req CreateRequest) (Notification, error) {
	item, err := s.repo.Create(ctx, req)
	return item, mapErr(err, "restaurant or user does not exist")
}

func (s Service) List(ctx context.Context, restaurantID string) ([]Notification, error) {
	items, err := s.repo.List(ctx, restaurantID)
	return items, mapErr(err, "")
}

func (s Service) Mark(ctx context.Context, id string, req MarkSentRequest) (Notification, error) {
	item, err := s.repo.Mark(ctx, id, req)
	if IsNotFound(err) {
		return item, apperrors.NotFound("notification not found")
	}
	return item, mapErr(err, "")
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
