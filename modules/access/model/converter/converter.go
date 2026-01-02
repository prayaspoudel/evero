package converter

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/prayaspoudel/modules/access/model"
)

// UserToResponse converts a User entity to UserResponse model
func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		CompanyID:     user.CompanyID,
		Role:          user.Role,
		IsActive:      user.IsActive,
		IsVerified:    user.IsVerified,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// CompanyToResponse converts a Company entity to CompanyResponse model
func CompanyToResponse(company *entity.Company) *model.CompanyResponse {
	return &model.CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Email:     company.Email,
		Industry:  company.Industry,
		Status:    company.Status,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}

// OAuth2ClientToResponse converts an OAuth2Client entity to OAuth2ClientResponse model
func OAuth2ClientToResponse(client *entity.OAuth2Client) *model.OAuth2ClientResponse {
	// Note: You would need to parse the JSON strings back to arrays
	// This is a simplified version
	return &model.OAuth2ClientResponse{
		ID:          client.ID,
		ClientID:    client.ClientID,
		Name:        client.Name,
		Description: client.Description,
		LogoURL:     client.LogoURL,
		Active:      client.Active,
		CreatedAt:   client.CreatedAt,
		UpdatedAt:   client.UpdatedAt,
	}
}
