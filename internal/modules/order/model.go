package order

import "time"

type Order struct {
	ID            string      `json:"id"`
	RestaurantID  string      `json:"restaurant_id"`
	BranchID      *string     `json:"branch_id,omitempty"`
	UserID        *string     `json:"user_id,omitempty"`
	OrderNumber   string      `json:"order_number"`
	Status        string      `json:"status"`
	PaymentStatus string      `json:"payment_status"`
	Subtotal      float64     `json:"subtotal"`
	TaxTotal      float64     `json:"tax_total"`
	Total         float64     `json:"total"`
	Notes         *string     `json:"notes,omitempty"`
	Items         []OrderItem `json:"items,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Quantity      float64   `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	TaxPercentage float64   `json:"tax_percentage"`
	LineSubtotal  float64   `json:"line_subtotal"`
	LineTax       float64   `json:"line_tax"`
	LineTotal     float64   `json:"line_total"`
	CreatedAt     time.Time `json:"created_at"`
}
