package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"restaurant-inventory-api/internal/config"
)

type Client struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(ctx context.Context, cfg config.MongoConfig) (*Client, error) {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.URL))
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	client := &Client{
		Client:   mongoClient,
		Database: mongoClient.Database(cfg.Database),
	}

	if err := client.Ping(ctx); err != nil {
		_ = mongoClient.Disconnect(ctx)
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	return client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("mongo client is not initialized")
	}

	return c.Client.Ping(ctx, readpref.Primary())
}

func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.Client == nil {
		return nil
	}

	return c.Client.Disconnect(ctx)
}
