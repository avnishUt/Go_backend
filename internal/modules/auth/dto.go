package auth

type RegisterRequest struct {
	RestaurantID *string `json:"restaurant_id" binding:"omitempty,uuid"`
	BranchID     *string `json:"branch_id" binding:"omitempty,uuid"`
	FullName     string  `json:"full_name" binding:"required,min=2,max=120"`
	Email        string  `json:"email" binding:"required,email,max=160"`
	Phone        *string `json:"phone" binding:"omitempty,max=30"`
	Password     string  `json:"password" binding:"required,min=8,max=80"`
	Role         string  `json:"role" binding:"required,oneof=super_admin admin store_manager user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	User        User   `json:"user"`
}
