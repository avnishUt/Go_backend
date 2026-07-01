package notification

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return Repository{db: db} }

func (r Repository) Create(ctx context.Context, req CreateRequest) (Notification, error) {
	var item Notification
	err := r.db.QueryRow(ctx, `
		INSERT INTO notifications (restaurant_id, user_id, channel, recipient, subject, body)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, restaurant_id, user_id, channel, recipient, subject, body, status,
			attempts, last_error, sent_at, created_at, updated_at
	`, req.RestaurantID, req.UserID, req.Channel, req.Recipient, req.Subject, req.Body).Scan(
		&item.ID, &item.RestaurantID, &item.UserID, &item.Channel, &item.Recipient,
		&item.Subject, &item.Body, &item.Status, &item.Attempts, &item.LastError,
		&item.SentAt, &item.CreatedAt, &item.UpdatedAt,
	)
	return item, err
}

func (r Repository) List(ctx context.Context, restaurantID string) ([]Notification, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, user_id, channel, recipient, subject, body, status,
			attempts, last_error, sent_at, created_at, updated_at
		FROM notifications
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Notification, 0)
	for rows.Next() {
		var item Notification
		if err := scan(rows, &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r Repository) Mark(ctx context.Context, id string, req MarkSentRequest) (Notification, error) {
	var item Notification
	err := r.db.QueryRow(ctx, `
		UPDATE notifications
		SET status = $2,
			attempts = attempts + 1,
			last_error = $3,
			sent_at = CASE WHEN $2 = 'sent' THEN now() ELSE sent_at END,
			updated_at = now()
		WHERE id = $1
		RETURNING id, restaurant_id, user_id, channel, recipient, subject, body, status,
			attempts, last_error, sent_at, created_at, updated_at
	`, id, req.Status, req.LastError).Scan(
		&item.ID, &item.RestaurantID, &item.UserID, &item.Channel, &item.Recipient,
		&item.Subject, &item.Body, &item.Status, &item.Attempts, &item.LastError,
		&item.SentAt, &item.CreatedAt, &item.UpdatedAt,
	)
	return item, err
}

func (r Repository) Pending(ctx context.Context, limit int) ([]Notification, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, restaurant_id, user_id, channel, recipient, subject, body, status,
			attempts, last_error, sent_at, created_at, updated_at
		FROM notifications
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Notification, 0)
	for rows.Next() {
		var item Notification
		if err := scan(rows, &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scan(rows pgx.Rows, item *Notification) error {
	return rows.Scan(&item.ID, &item.RestaurantID, &item.UserID, &item.Channel, &item.Recipient,
		&item.Subject, &item.Body, &item.Status, &item.Attempts, &item.LastError,
		&item.SentAt, &item.CreatedAt, &item.UpdatedAt)
}

func IsNotFound(err error) bool { return errors.Is(err, pgx.ErrNoRows) }
func IsInvalidInput(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "22P02"
}
func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
