package auth

import (
	"context"
	"errors"

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

func (r Repository) CreateUser(ctx context.Context, req RegisterRequest, passwordHash string) (User, error) {
	var user User
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return user, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO users (restaurant_id, branch_id, full_name, email, phone, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, restaurant_id, branch_id, full_name, email, phone, status, created_at, updated_at
	`, req.RestaurantID, req.BranchID, req.FullName, req.Email, req.Phone, passwordHash).Scan(
		&user.ID, &user.RestaurantID, &user.BranchID, &user.FullName, &user.Email,
		&user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, restaurant_id)
		SELECT $1, id, $3 FROM roles WHERE name = $2
	`, user.ID, req.Role, req.RestaurantID)
	if err != nil {
		return user, err
	}

	if err := tx.Commit(ctx); err != nil {
		return user, err
	}

	user.Roles = []string{req.Role}
	return user, nil
}

func (r Repository) GetByEmailWithPassword(ctx context.Context, email string) (User, string, error) {
	var user User
	var passwordHash string
	err := r.db.QueryRow(ctx, `
		SELECT id, restaurant_id, branch_id, full_name, email, phone, password_hash, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.RestaurantID, &user.BranchID, &user.FullName, &user.Email,
		&user.Phone, &passwordHash, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return user, "", err
	}

	roles, err := r.RolesForUser(ctx, user.ID)
	if err != nil {
		return user, "", err
	}
	user.Roles = roles

	return user, passwordHash, nil
}

func (r Repository) GetByID(ctx context.Context, id string) (User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		SELECT id, restaurant_id, branch_id, full_name, email, phone, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID, &user.RestaurantID, &user.BranchID, &user.FullName, &user.Email,
		&user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	roles, err := r.RolesForUser(ctx, user.ID)
	if err != nil {
		return user, err
	}
	user.Roles = roles

	return user, nil
}

func (r Repository) RolesForUser(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT roles.name
		FROM user_roles
		JOIN roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = $1
		ORDER BY roles.name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]string, 0)
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
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
