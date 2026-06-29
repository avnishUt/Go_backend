package order

type CreateOrderRequest struct {
	RestaurantID string                   `json:"restaurant_id" binding:"required,uuid"`
	BranchID     *string                  `json:"branch_id" binding:"omitempty,uuid"`
	Notes        *string                  `json:"notes" binding:"omitempty,max=500"`
	Items        []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required,uuid"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateStatusRequest struct {
	Status string  `json:"status" binding:"required,oneof=pending accepted preparing ready completed cancelled"`
	Note   *string `json:"note" binding:"omitempty,max=300"`
}
