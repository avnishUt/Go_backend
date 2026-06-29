package payment

import "time"

type Payment struct {
	ID                string    `json:"id"`
	RestaurantID      string    `json:"restaurant_id"`
	OrderID           string    `json:"order_id"`
	Amount            float64   `json:"amount"`
	Method            string    `json:"method"`
	Status            string    `json:"status"`
	ProviderReference *string   `json:"provider_reference,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
