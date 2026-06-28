package inventory

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return Repository{db: db} }

func (r Repository) CreateItem(ctx context.Context, req CreateItemRequest) (Item, error) {
	var item Item
	err := r.db.QueryRow(ctx, `
		INSERT INTO inventory_items (
			restaurant_id, branch_id, product_id, unit, current_stock,
			minimum_stock, reorder_level, allow_negative_stock
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, restaurant_id, branch_id, product_id, unit, current_stock,
			minimum_stock, reorder_level, allow_negative_stock, status, created_at, updated_at
	`, req.RestaurantID, req.BranchID, req.ProductID, req.Unit, req.CurrentStock,
		req.MinimumStock, req.ReorderLevel, req.AllowNegativeStock).Scan(
		&item.ID, &item.RestaurantID, &item.BranchID, &item.ProductID, &item.Unit,
		&item.CurrentStock, &item.MinimumStock, &item.ReorderLevel, &item.AllowNegativeStock,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	return item, err
}

func (r Repository) ListItems(ctx context.Context, restaurantID string) ([]Item, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, branch_id, product_id, unit, current_stock,
			minimum_stock, reorder_level, allow_negative_stock, status, created_at, updated_at
		FROM inventory_items
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := scanItem(rows, &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r Repository) GetItem(ctx context.Context, id string) (Item, error) {
	var item Item
	err := r.db.QueryRow(ctx, `
		SELECT id, restaurant_id, branch_id, product_id, unit, current_stock,
			minimum_stock, reorder_level, allow_negative_stock, status, created_at, updated_at
		FROM inventory_items
		WHERE id = $1
	`, id).Scan(&item.ID, &item.RestaurantID, &item.BranchID, &item.ProductID, &item.Unit,
		&item.CurrentStock, &item.MinimumStock, &item.ReorderLevel, &item.AllowNegativeStock,
		&item.Status, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r Repository) UpdateItem(ctx context.Context, id string, req UpdateItemRequest) (Item, error) {
	var item Item
	err := r.db.QueryRow(ctx, `
		UPDATE inventory_items
		SET unit = COALESCE($2::text, unit),
			minimum_stock = COALESCE($3::numeric, minimum_stock),
			reorder_level = COALESCE($4::numeric, reorder_level),
			allow_negative_stock = COALESCE($5::boolean, allow_negative_stock),
			status = COALESCE($6::text, status),
			updated_at = now()
		WHERE id = $1
		RETURNING id, restaurant_id, branch_id, product_id, unit, current_stock,
			minimum_stock, reorder_level, allow_negative_stock, status, created_at, updated_at
	`, id, req.Unit, req.MinimumStock, req.ReorderLevel, req.AllowNegativeStock, req.Status).Scan(
		&item.ID, &item.RestaurantID, &item.BranchID, &item.ProductID, &item.Unit,
		&item.CurrentStock, &item.MinimumStock, &item.ReorderLevel, &item.AllowNegativeStock,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	return item, err
}

func (r Repository) AdjustStock(ctx context.Context, id string, delta float64) (Item, error) {
	var item Item
	err := r.db.QueryRow(ctx, `
		UPDATE inventory_items
		SET current_stock = current_stock + $2,
			updated_at = now()
		WHERE id = $1
		RETURNING id, restaurant_id, branch_id, product_id, unit, current_stock,
			minimum_stock, reorder_level, allow_negative_stock, status, created_at, updated_at
	`, id, delta).Scan(&item.ID, &item.RestaurantID, &item.BranchID, &item.ProductID, &item.Unit,
		&item.CurrentStock, &item.MinimumStock, &item.ReorderLevel, &item.AllowNegativeStock,
		&item.Status, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func scanItem(rows pgx.Rows, item *Item) error {
	return rows.Scan(&item.ID, &item.RestaurantID, &item.BranchID, &item.ProductID, &item.Unit,
		&item.CurrentStock, &item.MinimumStock, &item.ReorderLevel, &item.AllowNegativeStock,
		&item.Status, &item.CreatedAt, &item.UpdatedAt)
}

func IsNotFound(err error) bool { return errors.Is(err, pgx.ErrNoRows) }

func pgCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}

func IsUniqueViolation(err error) bool     { return pgCode(err, "23505") }
func IsForeignKeyViolation(err error) bool { return pgCode(err, "23503") }
func IsInvalidInput(err error) bool        { return pgCode(err, "22P02") }
func IsCheckViolation(err error) bool      { return pgCode(err, "23514") }
