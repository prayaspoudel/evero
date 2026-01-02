package model

// CompanyResponse represents company data in API responses
type CompanyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain,omitempty"`
	IsActive  bool   `json:"isActive"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

// CreateCompanyRequest represents a request to create a company
type CreateCompanyRequest struct {
	Name   string  `json:"name" validate:"required,max=255"`
	Domain *string `json:"domain,omitempty" validate:"omitempty,max=255"`
}

// UpdateCompanyRequest represents a request to update a company
type UpdateCompanyRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Domain   *string `json:"domain,omitempty" validate:"omitempty,max=255"`
	IsActive *bool   `json:"isActive,omitempty"`
}
