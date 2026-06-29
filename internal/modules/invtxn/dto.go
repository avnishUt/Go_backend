package invtxn

type CreateTransactionRequest struct {
	InventoryItemID string  `json:"inventory_item_id" binding:"required,uuid"`
	Type            string  `json:"type" binding:"required,oneof=stock_in stock_out adjustment wastage return"`
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`
	Reason          *string `json:"reason" binding:"omitempty,max=300"`
	ReferenceType   *string `json:"reference_type" binding:"omitempty,max=50"`
	ReferenceID     *string `json:"reference_id" binding:"omitempty,uuid"`
}
