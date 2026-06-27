package app

import (
	"context"
	"log"
	"time"

	"restaurant-inventory-api/internal/config"
	"restaurant-inventory-api/internal/http/router"
	"restaurant-inventory-api/internal/platform/mongo"
	"restaurant-inventory-api/internal/platform/postgres"
	"restaurant-inventory-api/internal/server"
)

func Run() error {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pgClient *postgres.Client
	if cfg.Postgres.URL != "" {
		client, err := postgres.Connect(ctx, cfg.Postgres)
		if err != nil {
			return err
		}
		defer client.Close()
		pgClient = client
		log.Print("postgres connected")
	}

	var mongoClient *mongo.Client
	if cfg.Mongo.URL != "" {
		client, err := mongo.Connect(ctx, cfg.Mongo)
		if err != nil {
			return err
		}
		defer client.Close(context.Background())
		mongoClient = client
		log.Print("mongo connected")
	}

	deps := router.Dependencies{}
	if pgClient != nil {
		deps.Postgres = pgClient
		deps.Database = pgClient
	}
	if mongoClient != nil {
		deps.Mongo = mongoClient
	}

	r := router.New(cfg, deps)

	return server.Run(cfg, r)
}
