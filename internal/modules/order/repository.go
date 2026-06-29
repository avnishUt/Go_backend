package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return Repository{db: db} }

func (r Repository) Create(ctx context.Context, req CreateOrderRequest, userID *string) (Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback(ctx)

	orderNumber := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	var created Order
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (restaurant_id, branch_id, user_id, order_number, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, restaurant_id, branch_id, user_id, order_number, status, payment_status,
			subtotal, tax_total, total, notes, created_at, updated_at
	`, req.RestaurantID, req.BranchID, userID, orderNumber, req.Notes).Scan(
		&created.ID, &created.RestaurantID, &created.BranchID, &created.UserID, &created.OrderNumber,
		&created.Status, &created.PaymentStatus, &created.Subtotal, &created.TaxTotal, &created.Total,
		&created.Notes, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return Order{}, err
	}

	var subtotal, taxTotal, total float64
	items := make([]OrderItem, 0, len(req.Items))
	for _, input := range req.Items {
		var product struct {
			Name          string
			Price         float64
			TaxPercentage float64
		}
		err = tx.QueryRow(ctx, `
			SELECT name, price, tax_percentage
			FROM products
			WHERE id = $1 AND restaurant_id = $2 AND status = 'active'
		`, input.ProductID, req.RestaurantID).Scan(&product.Name, &product.Price, &product.TaxPercentage)
		if err != nil {
			return Order{}, err
		}

		lineSubtotal := product.Price * input.Quantity
		lineTax := lineSubtotal * product.TaxPercentage / 100
		lineTotal := lineSubtotal + lineTax
		subtotal += lineSubtotal
		taxTotal += lineTax
		total += lineTotal

		var item OrderItem
		err = tx.QueryRow(ctx, `
			INSERT INTO order_items (
				order_id, product_id, product_name, quantity, unit_price, tax_percentage,
				line_subtotal, line_tax, line_total
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, order_id, product_id, product_name, quantity, unit_price,
				tax_percentage, line_subtotal, line_tax, line_total, created_at
		`, created.ID, input.ProductID, product.Name, input.Quantity, product.Price, product.TaxPercentage,
			lineSubtotal, lineTax, lineTotal).Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.UnitPrice, &item.TaxPercentage, &item.LineSubtotal, &item.LineTax,
			&item.LineTotal, &item.CreatedAt)
		if err != nil {
			return Order{}, err
		}
		items = append(items, item)
	}

	err = tx.QueryRow(ctx, `
		UPDATE orders
		SET subtotal = $2, tax_total = $3, total = $4, updated_at = now()
		WHERE id = $1
		RETURNING subtotal, tax_total, total, updated_at
	`, created.ID, subtotal, taxTotal, total).Scan(&created.Subtotal, &created.TaxTotal, &created.Total, &created.UpdatedAt)
	if err != nil {
		return Order{}, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO order_status_history (order_id, status, changed_by, note)
		VALUES ($1, $2, $3, $4)
	`, created.ID, created.Status, userID, "order created")
	if err != nil {
		return Order{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Order{}, err
	}
	created.Items = items
	return created, nil
}

func (r Repository) List(ctx context.Context, restaurantID string) ([]Order, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, branch_id, user_id, order_number, status, payment_status,
			subtotal, tax_total, total, notes, created_at, updated_at
		FROM orders
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]Order, 0)
	for rows.Next() {
		var item Order
		if err := scanOrder(rows, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r Repository) Get(ctx context.Context, id string) (Order, error) {
	var result Order
	err := r.db.QueryRow(ctx, `
		SELECT id, restaurant_id, branch_id, user_id, order_number, status, payment_status,
			subtotal, tax_total, total, notes, created_at, updated_at
		FROM orders
		WHERE id = $1
	`, id).Scan(&result.ID, &result.RestaurantID, &result.BranchID, &result.UserID, &result.OrderNumber,
		&result.Status, &result.PaymentStatus, &result.Subtotal, &result.TaxTotal, &result.Total,
		&result.Notes, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return result, err
	}
	items, err := r.items(ctx, id)
	if err != nil {
		return result, err
	}
	result.Items = items
	return result, nil
}

func (r Repository) UpdateStatus(ctx context.Context, id string, req UpdateStatusRequest, userID *string) (Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback(ctx)

	var result Order
	err = tx.QueryRow(ctx, `
		UPDATE orders
		SET status = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, restaurant_id, branch_id, user_id, order_number, status, payment_status,
			subtotal, tax_total, total, notes, created_at, updated_at
	`, id, req.Status).Scan(&result.ID, &result.RestaurantID, &result.BranchID, &result.UserID, &result.OrderNumber,
		&result.Status, &result.PaymentStatus, &result.Subtotal, &result.TaxTotal, &result.Total,
		&result.Notes, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return Order{}, err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO order_status_history (order_id, status, changed_by, note)
		VALUES ($1, $2, $3, $4)
	`, id, req.Status, userID, req.Note)
	if err != nil {
		return Order{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Order{}, err
	}
	return result, nil
}

func (r Repository) items(ctx context.Context, orderID string) ([]OrderItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, order_id, product_id, product_name, quantity, unit_price,
			tax_percentage, line_subtotal, line_tax, line_total, created_at
		FROM order_items
		WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]OrderItem, 0)
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Quantity,
			&item.UnitPrice, &item.TaxPercentage, &item.LineSubtotal, &item.LineTax, &item.LineTotal,
			&item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanOrder(rows pgx.Rows, item *Order) error {
	return rows.Scan(&item.ID, &item.RestaurantID, &item.BranchID, &item.UserID, &item.OrderNumber,
		&item.Status, &item.PaymentStatus, &item.Subtotal, &item.TaxTotal, &item.Total,
		&item.Notes, &item.CreatedAt, &item.UpdatedAt)
}

func IsNotFound(err error) bool { return errors.Is(err, pgx.ErrNoRows) }
func pgCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}
func IsForeignKeyViolation(err error) bool { return pgCode(err, "23503") }
func IsInvalidInput(err error) bool        { return pgCode(err, "22P02") }
