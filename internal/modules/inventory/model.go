package inventory

import "time"

type Item struct {
	ID                 string    `json:"id"`
	RestaurantID       string    `json:"restaurant_id"`
	BranchID           *string   `json:"branch_id,omitempty"`
	ProductID          string    `json:"product_id"`
	Unit               string    `json:"unit"`
	CurrentStock       float64   `json:"current_stock"`
	MinimumStock       float64   `json:"minimum_stock"`
	ReorderLevel       float64   `json:"reorder_level"`
	AllowNegativeStock bool      `json:"allow_negative_stock"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
