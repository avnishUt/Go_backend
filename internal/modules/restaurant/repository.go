package restaurant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{db: db}
}

func (r Repository) CreateRestaurant(ctx context.Context, req CreateRestaurantRequest) (Restaurant, error) {
	var restaurant Restaurant

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return restaurant, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO restaurants (name, slug, legal_name, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, slug, legal_name, email, phone, status, created_at, updated_at
	`, req.Name, req.Slug, req.LegalName, req.Email, req.Phone).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Slug,
		&restaurant.LegalName,
		&restaurant.Email,
		&restaurant.Phone,
		&restaurant.Status,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)
	if err != nil {
		return restaurant, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO restaurant_settings (restaurant_id)
		VALUES ($1)
	`, restaurant.ID)
	if err != nil {
		return restaurant, err
	}

	if err := tx.Commit(ctx); err != nil {
		return restaurant, err
	}

	return restaurant, nil
}

func (r Repository) ListRestaurants(ctx context.Context) ([]Restaurant, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, slug, legal_name, email, phone, status, created_at, updated_at
		FROM restaurants
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := make([]Restaurant, 0)
	for rows.Next() {
		var restaurant Restaurant
		if err := scanRestaurant(rows, &restaurant); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, rows.Err()
}

func (r Repository) GetRestaurant(ctx context.Context, id string) (Restaurant, error) {
	var restaurant Restaurant
	err := r.db.QueryRow(ctx, `
		SELECT id, name, slug, legal_name, email, phone, status, created_at, updated_at
		FROM restaurants
		WHERE id = $1
	`, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Slug,
		&restaurant.LegalName,
		&restaurant.Email,
		&restaurant.Phone,
		&restaurant.Status,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	return restaurant, err
}

func (r Repository) UpdateRestaurant(ctx context.Context, id string, req UpdateRestaurantRequest) (Restaurant, error) {
	var restaurant Restaurant
	err := r.db.QueryRow(ctx, `
		UPDATE restaurants
		SET
			name = COALESCE($2::text, name),
			legal_name = COALESCE($3::text, legal_name),
			email = COALESCE($4::text, email),
			phone = COALESCE($5::text, phone),
			status = COALESCE($6::text, status),
			updated_at = now()
		WHERE id = $1
		RETURNING id, name, slug, legal_name, email, phone, status, created_at, updated_at
	`, id, req.Name, req.LegalName, req.Email, req.Phone, req.Status).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Slug,
		&restaurant.LegalName,
		&restaurant.Email,
		&restaurant.Phone,
		&restaurant.Status,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	return restaurant, err
}

func (r Repository) CreateBranch(ctx context.Context, restaurantID string, req CreateBranchRequest) (Branch, error) {
	country := "IN"
	if req.Country != nil {
		country = *req.Country
	}

	var branch Branch
	err := r.db.QueryRow(ctx, `
		INSERT INTO restaurant_branches (
			restaurant_id, name, slug, email, phone, address_line1, address_line2,
			city, state, country, postal_code
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, restaurant_id, name, slug, email, phone, address_line1, address_line2,
			city, state, country, postal_code, status, created_at, updated_at
	`, restaurantID, req.Name, req.Slug, req.Email, req.Phone, req.AddressLine1, req.AddressLine2,
		req.City, req.State, country, req.PostalCode).Scan(
		&branch.ID,
		&branch.RestaurantID,
		&branch.Name,
		&branch.Slug,
		&branch.Email,
		&branch.Phone,
		&branch.AddressLine1,
		&branch.AddressLine2,
		&branch.City,
		&branch.State,
		&branch.Country,
		&branch.PostalCode,
		&branch.Status,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	return branch, err
}

func (r Repository) ListBranches(ctx context.Context, restaurantID string) ([]Branch, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, name, slug, email, phone, address_line1, address_line2,
			city, state, country, postal_code, status, created_at, updated_at
		FROM restaurant_branches
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branches := make([]Branch, 0)
	for rows.Next() {
		var branch Branch
		if err := scanBranch(rows, &branch); err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}

	return branches, rows.Err()
}

func (r Repository) GetSettings(ctx context.Context, restaurantID string) (Settings, error) {
	var settings Settings
	var metadata []byte

	err := r.db.QueryRow(ctx, `
		SELECT restaurant_id, currency, timezone, tax_percentage, service_charge_percentage,
			metadata, created_at, updated_at
		FROM restaurant_settings
		WHERE restaurant_id = $1
	`, restaurantID).Scan(
		&settings.RestaurantID,
		&settings.Currency,
		&settings.Timezone,
		&settings.TaxPercentage,
		&settings.ServiceChargePercentage,
		&metadata,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return settings, err
	}

	if len(metadata) > 0 {
		if err := json.Unmarshal(metadata, &settings.Metadata); err != nil {
			return settings, err
		}
	}
	if settings.Metadata == nil {
		settings.Metadata = map[string]interface{}{}
	}

	return settings, nil
}

func (r Repository) UpdateSettings(ctx context.Context, restaurantID string, req UpdateSettingsRequest) (Settings, error) {
	metadata, err := json.Marshal(req.Metadata)
	if err != nil {
		return Settings{}, err
	}

	var settings Settings
	var metadataRaw []byte

	err = r.db.QueryRow(ctx, `
		UPDATE restaurant_settings
		SET
			currency = COALESCE($2::text, currency),
			timezone = COALESCE($3::text, timezone),
			tax_percentage = COALESCE($4::numeric, tax_percentage),
			service_charge_percentage = COALESCE($5::numeric, service_charge_percentage),
			metadata = COALESCE($6::jsonb, metadata),
			updated_at = now()
		WHERE restaurant_id = $1
		RETURNING restaurant_id, currency, timezone, tax_percentage, service_charge_percentage,
			metadata, created_at, updated_at
	`, restaurantID, req.Currency, req.Timezone, req.TaxPercentage, req.ServiceChargePercentage, nullableJSON(req.Metadata, metadata)).Scan(
		&settings.RestaurantID,
		&settings.Currency,
		&settings.Timezone,
		&settings.TaxPercentage,
		&settings.ServiceChargePercentage,
		&metadataRaw,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return settings, err
	}

	if err := json.Unmarshal(metadataRaw, &settings.Metadata); err != nil {
		return settings, err
	}

	return settings, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

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

func nullableJSON(value map[string]interface{}, encoded []byte) interface{} {
	if value == nil {
		return nil
	}

	return string(encoded)
}

func scanRestaurant(rows pgx.Rows, restaurant *Restaurant) error {
	return rows.Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Slug,
		&restaurant.LegalName,
		&restaurant.Email,
		&restaurant.Phone,
		&restaurant.Status,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)
}

func scanBranch(rows pgx.Rows, branch *Branch) error {
	return rows.Scan(
		&branch.ID,
		&branch.RestaurantID,
		&branch.Name,
		&branch.Slug,
		&branch.Email,
		&branch.Phone,
		&branch.AddressLine1,
		&branch.AddressLine2,
		&branch.City,
		&branch.State,
		&branch.Country,
		&branch.PostalCode,
		&branch.Status,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)
}

func normalizeDBError(operation string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", operation, err)
}
