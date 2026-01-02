package model

// OAuth2ClientResponse represents an OAuth2 client in API responses
type OAuth2ClientResponse struct {
	ID           string   `json:"id"`
	ClientID     string   `json:"clientId"`
	Name         string   `json:"name"`
	Description  *string  `json:"description,omitempty"`
	RedirectURIs []string `json:"redirectUris"`
	GrantTypes   []string `json:"grantTypes"`
	Scopes       []string `json:"scopes"`
	OwnerID      string   `json:"ownerId"`
	LogoURL      *string  `json:"logoUrl,omitempty"`
	Active       bool     `json:"active"`
	CreatedAt    int64    `json:"createdAt"`
	UpdatedAt    int64    `json:"updatedAt"`
}

// CreateOAuth2ClientRequest represents a request to create an OAuth2 client
type CreateOAuth2ClientRequest struct {
	Name         string   `json:"name" validate:"required,max=255"`
	Description  *string  `json:"description,omitempty"`
	RedirectURIs []string `json:"redirectUris" validate:"required,min=1"`
	GrantTypes   []string `json:"grantTypes" validate:"required,min=1"`
	Scopes       []string `json:"scopes" validate:"required,min=1"`
	LogoURL      *string  `json:"logoUrl,omitempty" validate:"omitempty,url"`
}

// UpdateOAuth2ClientRequest represents a request to update an OAuth2 client
type UpdateOAuth2ClientRequest struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,max=255"`
	Description  *string  `json:"description,omitempty"`
	RedirectURIs []string `json:"redirectUris,omitempty"`
	GrantTypes   []string `json:"grantTypes,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	LogoURL      *string  `json:"logoUrl,omitempty" validate:"omitempty,url"`
	Active       *bool    `json:"active,omitempty"`
}

// AuthorizeRequest represents an OAuth2 authorization request
type AuthorizeRequest struct {
	ResponseType string `json:"responseType" validate:"required"`
	ClientID     string `json:"clientId" validate:"required"`
	RedirectURI  string `json:"redirectUri" validate:"required,url"`
	Scope        string `json:"scope,omitempty"`
	State        string `json:"state,omitempty"`
}

// AuthorizeResponse represents an OAuth2 authorization response
type AuthorizeResponse struct {
	Code  string `json:"code"`
	State string `json:"state,omitempty"`
}

// TokenRequest represents an OAuth2 token request
type TokenRequest struct {
	GrantType    string `json:"grantType" validate:"required"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ClientID     string `json:"clientId" validate:"required"`
	ClientSecret string `json:"clientSecret" validate:"required"`
	RedirectURI  string `json:"redirectUri,omitempty"`
}

// TokenResponse represents an OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Scope        string `json:"scope,omitempty"`
}
