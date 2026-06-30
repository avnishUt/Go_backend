package report

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return Repository{db: db} }

func (r Repository) SalesSummary(ctx context.Context, restaurantID string) (SalesSummary, error) {
	var summary SalesSummary
	err := r.db.QueryRow(ctx, `
		SELECT
			$1::uuid AS restaurant_id,
			COUNT(o.id) AS total_orders,
			COUNT(o.id) FILTER (WHERE o.payment_status = 'paid') AS paid_orders,
			COALESCE(SUM(o.total), 0) AS gross_sales,
			COALESCE(SUM(o.tax_total), 0) AS tax_total,
			COALESCE((
				SELECT SUM(p.amount)
				FROM payments p
				WHERE p.restaurant_id = $1 AND p.status = 'completed'
			), 0) AS payment_received
		FROM orders o
		WHERE o.restaurant_id = $1
	`, restaurantID).Scan(&summary.RestaurantID, &summary.TotalOrders, &summary.PaidOrders,
		&summary.GrossSales, &summary.TaxTotal, &summary.PaymentReceived)
	return summary, err
}

func (r Repository) LowStock(ctx context.Context, restaurantID string) ([]LowStockItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.name, p.sku, i.current_stock, i.minimum_stock, i.reorder_level, i.unit
		FROM inventory_items i
		JOIN products p ON p.id = i.product_id
		WHERE i.restaurant_id = $1
			AND i.status = 'active'
			AND i.current_stock <= i.reorder_level
		ORDER BY i.current_stock ASC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]LowStockItem, 0)
	for rows.Next() {
		var item LowStockItem
		if err := rows.Scan(&item.ProductName, &item.SKU, &item.CurrentStock,
			&item.MinimumStock, &item.ReorderLevel, &item.Unit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func IsInvalidInput(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "22P02"
}
