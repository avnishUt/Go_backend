package product

import (
	"context"
	"strings"

	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct{ repo Repository }

func NewService(repo Repository) Service { return Service{repo: repo} }

func (s Service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Slug = normalizeSlug(req.Slug)
	category, err := s.repo.CreateCategory(ctx, req)
	return category, mapDBError(err, "category already exists", "restaurant does not exist")
}

func (s Service) ListCategories(ctx context.Context, restaurantID string) ([]Category, error) {
	categories, err := s.repo.ListCategories(ctx, restaurantID)
	return categories, mapDBError(err, "", "")
}

func (s Service) CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.SKU = strings.ToUpper(strings.TrimSpace(req.SKU))
	product, err := s.repo.CreateProduct(ctx, req)
	return product, mapDBError(err, "product sku already exists", "restaurant or category does not exist")
}

func (s Service) ListProducts(ctx context.Context, restaurantID string) ([]Product, error) {
	products, err := s.repo.ListProducts(ctx, restaurantID)
	return products, mapDBError(err, "", "")
}

func (s Service) GetProduct(ctx context.Context, id string) (Product, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if IsNotFound(err) {
		return product, apperrors.NotFound("product not found")
	}
	return product, mapDBError(err, "", "")
}

func (s Service) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (Product, error) {
	product, err := s.repo.UpdateProduct(ctx, id, req)
	if IsNotFound(err) {
		return product, apperrors.NotFound("product not found")
	}
	return product, mapDBError(err, "product conflict", "category does not exist")
}

func mapDBError(err error, uniqueMsg string, fkMsg string) error {
	if err == nil {
		return nil
	}
	if IsInvalidInput(err) {
		return apperrors.BadRequest("invalid id", nil)
	}
	if IsUniqueViolation(err) && uniqueMsg != "" {
		return apperrors.Conflict(uniqueMsg, nil)
	}
	if IsForeignKeyViolation(err) && fkMsg != "" {
		return apperrors.BadRequest(fkMsg, nil)
	}
	return err
}

func normalizeSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	return strings.ReplaceAll(value, " ", "-")
}
