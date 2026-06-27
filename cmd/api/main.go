package main

import (
	"log"

	"restaurant-inventory-api/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
