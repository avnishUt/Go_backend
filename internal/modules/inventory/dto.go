package inventory

type CreateItemRequest struct {
	RestaurantID       string  `json:"restaurant_id" binding:"required,uuid"`
	BranchID           *string `json:"branch_id" binding:"omitempty,uuid"`
	ProductID          string  `json:"product_id" binding:"required,uuid"`
	Unit               string  `json:"unit" binding:"required,max=30"`
	CurrentStock       float64 `json:"current_stock" binding:"min=0"`
	MinimumStock       float64 `json:"minimum_stock" binding:"min=0"`
	ReorderLevel       float64 `json:"reorder_level" binding:"min=0"`
	AllowNegativeStock bool    `json:"allow_negative_stock"`
}

type UpdateItemRequest struct {
	Unit               *string  `json:"unit" binding:"omitempty,max=30"`
	MinimumStock       *float64 `json:"minimum_stock" binding:"omitempty,min=0"`
	ReorderLevel       *float64 `json:"reorder_level" binding:"omitempty,min=0"`
	AllowNegativeStock *bool    `json:"allow_negative_stock"`
	Status             *string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

type AdjustStockRequest struct {
	Delta float64 `json:"delta" binding:"required"`
}
