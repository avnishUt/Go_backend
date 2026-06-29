package invtxn

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return Repository{db: db} }

func (r Repository) Create(ctx context.Context, req CreateTransactionRequest, createdBy *string) (Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Transaction{}, err
	}
	defer tx.Rollback(ctx)

	var item struct {
		RestaurantID       string
		BranchID           *string
		ProductID          string
		CurrentStock       float64
		AllowNegativeStock bool
	}
	err = tx.QueryRow(ctx, `
		SELECT restaurant_id, branch_id, product_id, current_stock, allow_negative_stock
		FROM inventory_items
		WHERE id = $1
		FOR UPDATE
	`, req.InventoryItemID).Scan(&item.RestaurantID, &item.BranchID, &item.ProductID, &item.CurrentStock, &item.AllowNegativeStock)
	if err != nil {
		return Transaction{}, err
	}

	delta := req.Quantity
	if req.Type == "stock_out" || req.Type == "wastage" {
		delta = -req.Quantity
	}
	stockAfter := item.CurrentStock + delta
	if !item.AllowNegativeStock && stockAfter < 0 {
		return Transaction{}, ErrNegativeStock
	}

	_, err = tx.Exec(ctx, `
		UPDATE inventory_items
		SET current_stock = $2, updated_at = now()
		WHERE id = $1
	`, req.InventoryItemID, stockAfter)
	if err != nil {
		return Transaction{}, err
	}

	var result Transaction
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_transactions (
			restaurant_id, branch_id, inventory_item_id, product_id, type, quantity,
			stock_before, stock_after, reason, reference_type, reference_id, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, restaurant_id, branch_id, inventory_item_id, product_id, type, quantity,
			stock_before, stock_after, reason, reference_type, reference_id, created_by, created_at
	`, item.RestaurantID, item.BranchID, req.InventoryItemID, item.ProductID, req.Type, req.Quantity,
		item.CurrentStock, stockAfter, req.Reason, req.ReferenceType, req.ReferenceID, createdBy).Scan(
		&result.ID, &result.RestaurantID, &result.BranchID, &result.InventoryItemID, &result.ProductID,
		&result.Type, &result.Quantity, &result.StockBefore, &result.StockAfter, &result.Reason,
		&result.ReferenceType, &result.ReferenceID, &result.CreatedBy, &result.CreatedAt,
	)
	if err != nil {
		return Transaction{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Transaction{}, err
	}
	return result, nil
}

func (r Repository) List(ctx context.Context, restaurantID string) ([]Transaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, branch_id, inventory_item_id, product_id, type, quantity,
			stock_before, stock_after, reason, reference_type, reference_id, created_by, created_at
		FROM inventory_transactions
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Transaction, 0)
	for rows.Next() {
		var item Transaction
		if err := scan(rows, &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scan(rows pgx.Rows, item *Transaction) error {
	return rows.Scan(&item.ID, &item.RestaurantID, &item.BranchID, &item.InventoryItemID, &item.ProductID,
		&item.Type, &item.Quantity, &item.StockBefore, &item.StockAfter, &item.Reason,
		&item.ReferenceType, &item.ReferenceID, &item.CreatedBy, &item.CreatedAt)
}

var ErrNegativeStock = errors.New("negative stock is not allowed")

func IsNotFound(err error) bool { return errors.Is(err, pgx.ErrNoRows) }
func IsInvalidInput(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "22P02"
}
