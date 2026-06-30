package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	databaseURL := flag.String("database", os.Getenv("POSTGRES_URL"), "postgres database URL")
	restaurantName := flag.String("restaurant-name", "Seed Restaurant", "restaurant name")
	restaurantSlug := flag.String("restaurant-slug", "seed-restaurant", "restaurant slug")
	adminName := flag.String("admin-name", "Seed Admin", "admin full name")
	adminEmail := flag.String("admin-email", "seed-admin@example.com", "admin email")
	adminPassword := flag.String("admin-password", "password123", "admin password")
	flag.Parse()

	if *databaseURL == "" {
		log.Fatal("database URL is required; pass -database or set POSTGRES_URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, *databaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	restaurantID, err := upsertRestaurant(ctx, pool, *restaurantName, *restaurantSlug)
	if err != nil {
		log.Fatalf("seed restaurant: %v", err)
	}

	if err := seedAdmin(ctx, pool, restaurantID, *adminName, *adminEmail, *adminPassword); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	log.Printf("seed complete restaurant_id=%s admin_email=%s", restaurantID, strings.ToLower(*adminEmail))
}

func upsertRestaurant(ctx context.Context, pool *pgxpool.Pool, name string, slug string) (string, error) {
	var id string
	err := pool.QueryRow(ctx, `
		INSERT INTO restaurants (name, slug)
		VALUES ($1, $2)
		ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, updated_at = now()
		RETURNING id
	`, name, slug).Scan(&id)
	if err != nil {
		return "", err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO restaurant_settings (restaurant_id)
		VALUES ($1)
		ON CONFLICT (restaurant_id) DO NOTHING
	`, id)
	return id, err
}

func seedAdmin(ctx context.Context, pool *pgxpool.Pool, restaurantID string, name string, email string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	email = strings.ToLower(strings.TrimSpace(email))
	var userID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (restaurant_id, full_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE
		SET restaurant_id = EXCLUDED.restaurant_id,
			full_name = EXCLUDED.full_name,
			password_hash = EXCLUDED.password_hash,
			updated_at = now()
		RETURNING id
	`, restaurantID, name, email, string(hash)).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, restaurant_id)
		SELECT $1, id, $2 FROM roles WHERE name = 'admin'
		ON CONFLICT (user_id, role_id, restaurant_id) DO NOTHING
	`, userID, restaurantID)
	return err
}
