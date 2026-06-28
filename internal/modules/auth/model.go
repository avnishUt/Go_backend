package auth

import "time"

type User struct {
	ID           string    `json:"id"`
	RestaurantID *string   `json:"restaurant_id,omitempty"`
	BranchID     *string   `json:"branch_id,omitempty"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        *string   `json:"phone,omitempty"`
	Status       string    `json:"status"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Claims struct {
	UserID       string   `json:"user_id"`
	RestaurantID *string  `json:"restaurant_id,omitempty"`
	BranchID     *string  `json:"branch_id,omitempty"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
}
