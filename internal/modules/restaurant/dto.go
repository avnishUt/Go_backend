package restaurant

type CreateRestaurantRequest struct {
	Name      string  `json:"name" binding:"required,min=2,max=120"`
	Slug      string  `json:"slug" binding:"required,min=2,max=80"`
	LegalName *string `json:"legal_name" binding:"omitempty,max=160"`
	Email     *string `json:"email" binding:"omitempty,email,max=160"`
	Phone     *string `json:"phone" binding:"omitempty,max=30"`
}

type UpdateRestaurantRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=2,max=120"`
	LegalName *string `json:"legal_name" binding:"omitempty,max=160"`
	Email     *string `json:"email" binding:"omitempty,email,max=160"`
	Phone     *string `json:"phone" binding:"omitempty,max=30"`
	Status    *string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

type CreateBranchRequest struct {
	Name         string  `json:"name" binding:"required,min=2,max=120"`
	Slug         string  `json:"slug" binding:"required,min=2,max=80"`
	Email        *string `json:"email" binding:"omitempty,email,max=160"`
	Phone        *string `json:"phone" binding:"omitempty,max=30"`
	AddressLine1 *string `json:"address_line1" binding:"omitempty,max=200"`
	AddressLine2 *string `json:"address_line2" binding:"omitempty,max=200"`
	City         *string `json:"city" binding:"omitempty,max=80"`
	State        *string `json:"state" binding:"omitempty,max=80"`
	Country      *string `json:"country" binding:"omitempty,len=2"`
	PostalCode   *string `json:"postal_code" binding:"omitempty,max=20"`
}

type UpdateSettingsRequest struct {
	Currency                *string                `json:"currency" binding:"omitempty,len=3"`
	Timezone                *string                `json:"timezone" binding:"omitempty,max=80"`
	TaxPercentage           *float64               `json:"tax_percentage" binding:"omitempty,min=0,max=100"`
	ServiceChargePercentage *float64               `json:"service_charge_percentage" binding:"omitempty,min=0,max=100"`
	Metadata                map[string]interface{} `json:"metadata"`
}
