package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return Repository{db: db} }

func (r Repository) Create(ctx context.Context, req CreatePaymentRequest) (Payment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Payment{}, err
	}
	defer tx.Rollback(ctx)

	var order struct {
		RestaurantID string
		Total        float64
	}
	err = tx.QueryRow(ctx, `SELECT restaurant_id, total FROM orders WHERE id = $1`, req.OrderID).Scan(&order.RestaurantID, &order.Total)
	if err != nil {
		return Payment{}, err
	}

	var result Payment
	err = tx.QueryRow(ctx, `
		INSERT INTO payments (restaurant_id, order_id, amount, method, status)
		VALUES ($1, $2, $3, $4, 'pending')
		RETURNING id, restaurant_id, order_id, amount, method, status, provider_reference, created_at, updated_at
	`, order.RestaurantID, req.OrderID, order.Total, req.Method).Scan(&result.ID, &result.RestaurantID, &result.OrderID,
		&result.Amount, &result.Method, &result.Status, &result.ProviderReference, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return Payment{}, err
	}

	_, err = tx.Exec(ctx, `UPDATE orders SET payment_status = 'pending', updated_at = now() WHERE id = $1`, req.OrderID)
	if err != nil {
		return Payment{}, err
	}

	if err := r.recordTx(ctx, tx, result.ID, result.Status, result.Amount, result.ProviderReference, map[string]interface{}{"event": "payment_created"}); err != nil {
		return Payment{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Payment{}, err
	}
	return result, nil
}

func (r Repository) CompleteMockUPI(ctx context.Context, id string, req CompleteMockUPIRequest) (Payment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Payment{}, err
	}
	defer tx.Rollback(ctx)

	ref := req.ProviderReference
	if ref == nil || *ref == "" {
		value := fmt.Sprintf("MOCKUPI-%d", time.Now().UnixNano())
		ref = &value
	}

	var result Payment
	err = tx.QueryRow(ctx, `
		UPDATE payments
		SET status = 'completed', provider_reference = $2, updated_at = now()
		WHERE id = $1 AND method = 'mock_upi'
		RETURNING id, restaurant_id, order_id, amount, method, status, provider_reference, created_at, updated_at
	`, id, ref).Scan(&result.ID, &result.RestaurantID, &result.OrderID, &result.Amount, &result.Method,
		&result.Status, &result.ProviderReference, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return Payment{}, err
	}

	_, err = tx.Exec(ctx, `UPDATE orders SET payment_status = 'paid', updated_at = now() WHERE id = $1`, result.OrderID)
	if err != nil {
		return Payment{}, err
	}

	if err := r.recordTx(ctx, tx, result.ID, result.Status, result.Amount, result.ProviderReference, map[string]interface{}{"event": "mock_upi_completed"}); err != nil {
		return Payment{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Payment{}, err
	}
	return result, nil
}

func (r Repository) Get(ctx context.Context, id string) (Payment, error) {
	var result Payment
	err := r.db.QueryRow(ctx, `
		SELECT id, restaurant_id, order_id, amount, method, status, provider_reference, created_at, updated_at
		FROM payments WHERE id = $1
	`, id).Scan(&result.ID, &result.RestaurantID, &result.OrderID, &result.Amount, &result.Method,
		&result.Status, &result.ProviderReference, &result.CreatedAt, &result.UpdatedAt)
	return result, err
}

func (r Repository) List(ctx context.Context, restaurantID string) ([]Payment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, order_id, amount, method, status, provider_reference, created_at, updated_at
		FROM payments WHERE restaurant_id = $1 ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]Payment, 0)
	for rows.Next() {
		var item Payment
		if err := rows.Scan(&item.ID, &item.RestaurantID, &item.OrderID, &item.Amount, &item.Method,
			&item.Status, &item.ProviderReference, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

type txExec interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func (r Repository) recordTx(ctx context.Context, tx txExec, paymentID string, status string, amount float64, providerReference *string, raw map[string]interface{}) error {
	rawJSON, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO payment_transactions (payment_id, status, amount, provider_reference, raw_response)
		VALUES ($1, $2, $3, $4, $5)
	`, paymentID, status, amount, providerReference, string(rawJSON))
	return err
}

func IsNotFound(err error) bool { return errors.Is(err, pgx.ErrNoRows) }
func pgCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}
func IsInvalidInput(err error) bool { return pgCode(err, "22P02") }
