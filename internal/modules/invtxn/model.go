package invtxn

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	RestaurantID    string    `json:"restaurant_id"`
	BranchID        *string   `json:"branch_id,omitempty"`
	InventoryItemID string    `json:"inventory_item_id"`
	ProductID       string    `json:"product_id"`
	Type            string    `json:"type"`
	Quantity        float64   `json:"quantity"`
	StockBefore     float64   `json:"stock_before"`
	StockAfter      float64   `json:"stock_after"`
	Reason          *string   `json:"reason,omitempty"`
	ReferenceType   *string   `json:"reference_type,omitempty"`
	ReferenceID     *string   `json:"reference_id,omitempty"`
	CreatedBy       *string   `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
