package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := flag.String("database", os.Getenv("POSTGRES_URL"), "postgres database URL")
	path := flag.String("path", "migrations/postgres", "migration files path")
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

	if err := ensureSchemaMigrations(ctx, pool); err != nil {
		log.Fatalf("ensure schema migrations: %v", err)
	}

	files, err := upMigrationFiles(*path)
	if err != nil {
		log.Fatalf("read migrations: %v", err)
	}

	for _, file := range files {
		version := migrationVersion(file)
		applied, err := isApplied(ctx, pool, version)
		if err != nil {
			log.Fatalf("check migration %s: %v", version, err)
		}
		if applied {
			log.Printf("skip %s", filepath.Base(file))
			continue
		}

		sql, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("read migration %s: %v", file, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			log.Fatalf("begin migration %s: %v", version, err)
		}

		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			_ = tx.Rollback(ctx)
			log.Fatalf("apply migration %s: %v", version, err)
		}

		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			_ = tx.Rollback(ctx)
			log.Fatalf("record migration %s: %v", version, err)
		}

		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("commit migration %s: %v", version, err)
		}

		log.Printf("applied %s", filepath.Base(file))
	}
}

func ensureSchemaMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	return err
}

func upMigrationFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		files = append(files, filepath.Join(path, entry.Name()))
	}

	sort.Strings(files)
	return files, nil
}

func migrationVersion(file string) string {
	base := filepath.Base(file)
	return strings.TrimSuffix(base, ".up.sql")
}

func isApplied(ctx context.Context, pool *pgxpool.Pool, version string) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM schema_migrations WHERE version = $1
		)
	`, version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query schema_migrations: %w", err)
	}

	return exists, nil
}
