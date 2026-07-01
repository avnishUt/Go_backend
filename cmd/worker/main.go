package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"restaurant-inventory-api/internal/modules/notification"
)

func main() {
	databaseURL := flag.String("database", os.Getenv("POSTGRES_URL"), "postgres database URL")
	limit := flag.Int("limit", 10, "notifications to process per tick")
	interval := flag.Duration("interval", 0, "repeat interval; 0 runs once")
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

	repo := notification.NewRepository(pool)
	service := notification.NewService(repo)

	for {
		processed, err := service.ProcessPending(context.Background(), *limit)
		if err != nil {
			log.Fatalf("process notifications: %v", err)
		}
		log.Printf("processed_notifications=%d", processed)

		if *interval <= 0 {
			return
		}
		time.Sleep(*interval)
	}
}
