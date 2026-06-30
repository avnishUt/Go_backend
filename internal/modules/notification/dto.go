package notification

type CreateRequest struct {
	RestaurantID *string `json:"restaurant_id" binding:"omitempty,uuid"`
	UserID       *string `json:"user_id" binding:"omitempty,uuid"`
	Channel      string  `json:"channel" binding:"required,oneof=email sms push in_app"`
	Recipient    string  `json:"recipient" binding:"required,max=180"`
	Subject      *string `json:"subject" binding:"omitempty,max=180"`
	Body         string  `json:"body" binding:"required,max=1000"`
}

type MarkSentRequest struct {
	Status    string  `json:"status" binding:"required,oneof=sent failed"`
	LastError *string `json:"last_error" binding:"omitempty,max=500"`
}
