package payment

type CreatePaymentRequest struct {
	OrderID string `json:"order_id" binding:"required,uuid"`
	Method  string `json:"method" binding:"required,oneof=mock_upi cash card"`
}

type CompleteMockUPIRequest struct {
	ProviderReference *string `json:"provider_reference" binding:"omitempty,max=120"`
}
