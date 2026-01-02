package converter

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/prayaspoudel/modules/access/model"
)

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

func CompanyToResponse(company *entity.Company) *model.CompanyResponse {
	return &model.CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Domain:    company.Domain,
		IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}

func OAuth2ClientToResponse(client *entity.OAuth2Client) *model.OAuth2ClientResponse {
	return &model.OAuth2ClientResponse{
		ID:          client.ID,
		ClientID:    client.ClientID,
		Name:        client.Name,
		Description: client.Description,
		OwnerID:     client.OwnerID,
		LogoURL:     client.LogoURL,
		Active:      client.Active,
		CreatedAt:   client.CreatedAt,
		UpdatedAt:   client.UpdatedAt,
	}
}
