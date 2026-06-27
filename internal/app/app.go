package app

import (
	"restaurant-inventory-api/internal/config"
	"restaurant-inventory-api/internal/http/router"
	"restaurant-inventory-api/internal/server"
)

func Run() error {
	cfg := config.Load()
	r := router.New(cfg)

	return server.Run(cfg, r)
}
