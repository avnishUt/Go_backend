package product

type CreateCategoryRequest struct {
	RestaurantID string  `json:"restaurant_id" binding:"required,uuid"`
	Name         string  `json:"name" binding:"required,min=2,max=120"`
	Slug         string  `json:"slug" binding:"required,min=2,max=80"`
	Description  *string `json:"description" binding:"omitempty,max=300"`
}

type CreateProductRequest struct {
	RestaurantID  string  `json:"restaurant_id" binding:"required,uuid"`
	CategoryID    *string `json:"category_id" binding:"omitempty,uuid"`
	Name          string  `json:"name" binding:"required,min=2,max=160"`
	SKU           string  `json:"sku" binding:"required,min=2,max=80"`
	Description   *string `json:"description" binding:"omitempty,max=500"`
	Price         float64 `json:"price" binding:"min=0"`
	TaxPercentage float64 `json:"tax_percentage" binding:"min=0,max=100"`
}

type UpdateProductRequest struct {
	CategoryID    *string  `json:"category_id" binding:"omitempty,uuid"`
	Name          *string  `json:"name" binding:"omitempty,min=2,max=160"`
	Description   *string  `json:"description" binding:"omitempty,max=500"`
	Price         *float64 `json:"price" binding:"omitempty,min=0"`
	TaxPercentage *float64 `json:"tax_percentage" binding:"omitempty,min=0,max=100"`
	Status        *string  `json:"status" binding:"omitempty,oneof=active inactive"`
}
