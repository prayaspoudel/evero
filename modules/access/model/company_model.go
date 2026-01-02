package model

// CompanyResponse represents company data in API responses
type CompanyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	Industry  string `json:"industry,omitempty"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

// CreateCompanyRequest represents a create company request
type CreateCompanyRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"email,max=255"`
	Industry string `json:"industry" validate:"max=100"`
}

// UpdateCompanyRequest represents an update company request
type UpdateCompanyRequest struct {
	CompanyID string `json:"-" validate:"required"`
	Name      string `json:"name" validate:"max=255"`
	Email     string `json:"email" validate:"email,max=255"`
	Industry  string `json:"industry" validate:"max=100"`
	Status    string `json:"status" validate:"oneof=active inactive"`
}

// GetCompanyRequest represents a get company request
type GetCompanyRequest struct {
	CompanyID string `json:"-" validate:"required"`
}

// ListCompaniesRequest represents a list companies request
type ListCompaniesRequest struct {
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
	Status string `json:"status" validate:"omitempty,oneof=active inactive"`
}

// AddUserToCompanyRequest represents a request to add a user to a company
type AddUserToCompanyRequest struct {
	UserID    string `json:"userId" validate:"required"`
	CompanyID string `json:"companyId" validate:"required"`
	Role      string `json:"role" validate:"required"`
	IsPrimary bool   `json:"isPrimary"`
}

// RemoveUserFromCompanyRequest represents a request to remove a user from a company
type RemoveUserFromCompanyRequest struct {
	UserID    string `json:"userId" validate:"required"`
	CompanyID string `json:"companyId" validate:"required"`
}
