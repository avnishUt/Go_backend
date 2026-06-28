package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return Repository{db: db} }

func (r Repository) CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, error) {
	var category Category
	err := r.db.QueryRow(ctx, `
		INSERT INTO product_categories (restaurant_id, name, slug, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, restaurant_id, name, slug, description, status, created_at, updated_at
	`, req.RestaurantID, req.Name, req.Slug, req.Description).Scan(
		&category.ID, &category.RestaurantID, &category.Name, &category.Slug,
		&category.Description, &category.Status, &category.CreatedAt, &category.UpdatedAt,
	)
	return category, err
}

func (r Repository) ListCategories(ctx context.Context, restaurantID string) ([]Category, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, name, slug, description, status, created_at, updated_at
		FROM product_categories
		WHERE restaurant_id = $1
		ORDER BY name
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.RestaurantID, &category.Name, &category.Slug,
			&category.Description, &category.Status, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r Repository) CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error) {
	var product Product
	err := r.db.QueryRow(ctx, `
		INSERT INTO products (restaurant_id, category_id, name, sku, description, price, tax_percentage)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, restaurant_id, category_id, name, sku, description, price, tax_percentage, status, created_at, updated_at
	`, req.RestaurantID, req.CategoryID, req.Name, req.SKU, req.Description, req.Price, req.TaxPercentage).Scan(
		&product.ID, &product.RestaurantID, &product.CategoryID, &product.Name, &product.SKU,
		&product.Description, &product.Price, &product.TaxPercentage, &product.Status,
		&product.CreatedAt, &product.UpdatedAt,
	)
	return product, err
}

func (r Repository) ListProducts(ctx context.Context, restaurantID string) ([]Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, category_id, name, sku, description, price, tax_percentage, status, created_at, updated_at
		FROM products
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := scanProduct(rows, &product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r Repository) GetProduct(ctx context.Context, id string) (Product, error) {
	var product Product
	err := r.db.QueryRow(ctx, `
		SELECT id, restaurant_id, category_id, name, sku, description, price, tax_percentage, status, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(&product.ID, &product.RestaurantID, &product.CategoryID, &product.Name, &product.SKU,
		&product.Description, &product.Price, &product.TaxPercentage, &product.Status,
		&product.CreatedAt, &product.UpdatedAt)
	return product, err
}

func (r Repository) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (Product, error) {
	var product Product
	err := r.db.QueryRow(ctx, `
		UPDATE products
		SET category_id = COALESCE($2::uuid, category_id),
			name = COALESCE($3::text, name),
			description = COALESCE($4::text, description),
			price = COALESCE($5::numeric, price),
			tax_percentage = COALESCE($6::numeric, tax_percentage),
			status = COALESCE($7::text, status),
			updated_at = now()
		WHERE id = $1
		RETURNING id, restaurant_id, category_id, name, sku, description, price, tax_percentage, status, created_at, updated_at
	`, id, req.CategoryID, req.Name, req.Description, req.Price, req.TaxPercentage, req.Status).Scan(
		&product.ID, &product.RestaurantID, &product.CategoryID, &product.Name, &product.SKU,
		&product.Description, &product.Price, &product.TaxPercentage, &product.Status,
		&product.CreatedAt, &product.UpdatedAt,
	)
	return product, err
}

func scanProduct(rows pgx.Rows, product *Product) error {
	return rows.Scan(&product.ID, &product.RestaurantID, &product.CategoryID, &product.Name, &product.SKU,
		&product.Description, &product.Price, &product.TaxPercentage, &product.Status,
		&product.CreatedAt, &product.UpdatedAt)
}

func IsNotFound(err error) bool { return errors.Is(err, pgx.ErrNoRows) }

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

func IsInvalidInput(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "22P02"
}
