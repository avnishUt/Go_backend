package restaurant

import (
	"context"
	"strings"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) CreateRestaurant(ctx context.Context, req CreateRestaurantRequest) (Restaurant, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Slug = normalizeSlug(req.Slug)

	restaurant, err := s.repo.CreateRestaurant(ctx, req)
	if IsUniqueViolation(err) {
		return restaurant, apperrors.Conflict("restaurant slug already exists", nil)
	}
	return restaurant, normalizeDBError("create restaurant", err)
}

func (s Service) ListRestaurants(ctx context.Context) ([]Restaurant, error) {
	restaurants, err := s.repo.ListRestaurants(ctx)
	return restaurants, normalizeDBError("list restaurants", err)
}

func (s Service) GetRestaurant(ctx context.Context, id string) (Restaurant, error) {
	restaurant, err := s.repo.GetRestaurant(ctx, id)
	if IsInvalidInput(err) {
		return restaurant, apperrors.BadRequest("invalid restaurant id", nil)
	}
	if IsNotFound(err) {
		return restaurant, apperrors.NotFound("restaurant not found")
	}
	return restaurant, normalizeDBError("get restaurant", err)
}

func (s Service) UpdateRestaurant(ctx context.Context, id string, req UpdateRestaurantRequest) (Restaurant, error) {
	restaurant, err := s.repo.UpdateRestaurant(ctx, id, req)
	if IsInvalidInput(err) {
		return restaurant, apperrors.BadRequest("invalid restaurant id", nil)
	}
	if IsNotFound(err) {
		return restaurant, apperrors.NotFound("restaurant not found")
	}
	return restaurant, normalizeDBError("update restaurant", err)
}

func (s Service) CreateBranch(ctx context.Context, restaurantID string, req CreateBranchRequest) (Branch, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Slug = normalizeSlug(req.Slug)

	branch, err := s.repo.CreateBranch(ctx, restaurantID, req)
	if IsInvalidInput(err) {
		return branch, apperrors.BadRequest("invalid restaurant id", nil)
	}
	if IsForeignKeyViolation(err) {
		return branch, apperrors.NotFound("restaurant not found")
	}
	if IsUniqueViolation(err) {
		return branch, apperrors.Conflict("branch slug already exists for restaurant", nil)
	}
	return branch, normalizeDBError("create branch", err)
}

func (s Service) ListBranches(ctx context.Context, restaurantID string) ([]Branch, error) {
	branches, err := s.repo.ListBranches(ctx, restaurantID)
	if IsInvalidInput(err) {
		return branches, apperrors.BadRequest("invalid restaurant id", nil)
	}
	return branches, normalizeDBError("list branches", err)
}

func (s Service) GetSettings(ctx context.Context, restaurantID string) (Settings, error) {
	settings, err := s.repo.GetSettings(ctx, restaurantID)
	if IsInvalidInput(err) {
		return settings, apperrors.BadRequest("invalid restaurant id", nil)
	}
	if IsNotFound(err) {
		return settings, apperrors.NotFound("restaurant settings not found")
	}
	return settings, normalizeDBError("get restaurant settings", err)
}

func (s Service) UpdateSettings(ctx context.Context, restaurantID string, req UpdateSettingsRequest) (Settings, error) {
	settings, err := s.repo.UpdateSettings(ctx, restaurantID, req)
	if IsInvalidInput(err) {
		return settings, apperrors.BadRequest("invalid restaurant id", nil)
	}
	if IsNotFound(err) {
		return settings, apperrors.NotFound("restaurant settings not found")
	}
	return settings, normalizeDBError("update restaurant settings", err)
}

func normalizeSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "-")
	return value
}
