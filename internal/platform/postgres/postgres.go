package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"restaurant-inventory-api/internal/config"
)

type Client struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, cfg config.PostgresConfig) (*Client, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MinOpenConns)
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	client := &Client{Pool: pool}
	if err := client.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.Pool == nil {
		return fmt.Errorf("postgres client is not initialized")
	}

	return c.Pool.Ping(ctx)
}

func (c *Client) Close() {
	if c == nil || c.Pool == nil {
		return
	}

	c.Pool.Close()
}
