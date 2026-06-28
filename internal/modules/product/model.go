package product

import "time"

type Category struct {
	ID           string    `json:"id"`
	RestaurantID string    `json:"restaurant_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Product struct {
	ID            string    `json:"id"`
	RestaurantID  string    `json:"restaurant_id"`
	CategoryID    *string   `json:"category_id,omitempty"`
	Name          string    `json:"name"`
	SKU           string    `json:"sku"`
	Description   *string   `json:"description,omitempty"`
	Price         float64   `json:"price"`
	TaxPercentage float64   `json:"tax_percentage"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
