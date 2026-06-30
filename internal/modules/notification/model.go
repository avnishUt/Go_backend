package notification

import "time"

type Notification struct {
	ID           string     `json:"id"`
	RestaurantID *string    `json:"restaurant_id,omitempty"`
	UserID       *string    `json:"user_id,omitempty"`
	Channel      string     `json:"channel"`
	Recipient    string     `json:"recipient"`
	Subject      *string    `json:"subject,omitempty"`
	Body         string     `json:"body"`
	Status       string     `json:"status"`
	Attempts     int        `json:"attempts"`
	LastError    *string    `json:"last_error,omitempty"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
