package inventory

import (
	"context"
	"strings"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct{ repo Repository }

func NewService(repo Repository) Service { return Service{repo: repo} }

func (s Service) CreateItem(ctx context.Context, req CreateItemRequest) (Item, error) {
	req.Unit = strings.TrimSpace(req.Unit)
	item, err := s.repo.CreateItem(ctx, req)
	return item, mapDBError(err)
}

func (s Service) ListItems(ctx context.Context, restaurantID string) ([]Item, error) {
	items, err := s.repo.ListItems(ctx, restaurantID)
	return items, mapDBError(err)
}

func (s Service) GetItem(ctx context.Context, id string) (Item, error) {
	item, err := s.repo.GetItem(ctx, id)
	if IsNotFound(err) {
		return item, apperrors.NotFound("inventory item not found")
	}
	return item, mapDBError(err)
}

func (s Service) UpdateItem(ctx context.Context, id string, req UpdateItemRequest) (Item, error) {
	item, err := s.repo.UpdateItem(ctx, id, req)
	if IsNotFound(err) {
		return item, apperrors.NotFound("inventory item not found")
	}
	return item, mapDBError(err)
}

func (s Service) AdjustStock(ctx context.Context, id string, req AdjustStockRequest) (Item, error) {
	item, err := s.repo.AdjustStock(ctx, id, req.Delta)
	if IsNotFound(err) {
		return item, apperrors.NotFound("inventory item not found")
	}
	return item, mapDBError(err)
}

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if IsInvalidInput(err) {
		return apperrors.BadRequest("invalid id", nil)
	}
	if IsUniqueViolation(err) {
		return apperrors.Conflict("inventory item already exists for product scope", nil)
	}
	if IsForeignKeyViolation(err) {
		return apperrors.BadRequest("restaurant, branch, or product does not exist", nil)
	}
	if IsCheckViolation(err) {
		return apperrors.BadRequest("inventory stock cannot be negative", nil)
	}
	return err
}
