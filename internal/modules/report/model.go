package report

type SalesSummary struct {
	RestaurantID    string  `json:"restaurant_id"`
	TotalOrders     int64   `json:"total_orders"`
	PaidOrders      int64   `json:"paid_orders"`
	GrossSales      float64 `json:"gross_sales"`
	TaxTotal        float64 `json:"tax_total"`
	PaymentReceived float64 `json:"payment_received"`
}

type LowStockItem struct {
	ProductName  string  `json:"product_name"`
	SKU          string  `json:"sku"`
	CurrentStock float64 `json:"current_stock"`
	MinimumStock float64 `json:"minimum_stock"`
	ReorderLevel float64 `json:"reorder_level"`
	Unit         string  `json:"unit"`
}
