package model
package model

// OAuth2ClientResponse represents an OAuth2 client in API responses
type OAuth2ClientResponse struct {
	ID           string   `json:"id"`
	ClientID     string   `json:"clientId"`
	Name         string   `json:"name"`
	Description  *string  `json:"description,omitempty"`
	RedirectURIs []string `json:"redirectUris"`






































































}	ClientSecret  string `json:"client_secret" validate:"required"`	ClientID      string `json:"client_id" validate:"required"`	TokenTypeHint string `json:"token_type_hint,omitempty"`	Token         string `json:"token" validate:"required"`type RevokeTokenRequest struct {// RevokeTokenRequest represents a token revocation request}	Scope        string `json:"scope,omitempty"`	RefreshToken string `json:"refresh_token,omitempty"`	ExpiresIn    int    `json:"expires_in"`	TokenType    string `json:"token_type"`	AccessToken  string `json:"access_token"`type TokenResponse struct {// TokenResponse represents an OAuth2 token response}	RefreshToken string `json:"refresh_token"`	ClientSecret string `json:"client_secret" validate:"required"`	ClientID     string `json:"client_id" validate:"required"`	RedirectURI  string `json:"redirect_uri"`	Code         string `json:"code"`	GrantType    string `json:"grant_type" validate:"required"`type TokenRequest struct {// TokenRequest represents an OAuth2 token request}	State string `json:"state,omitempty"`	Code  string `json:"code"`type AuthorizeResponse struct {// AuthorizeResponse represents an OAuth2 authorization response}	UserID       string `json:"-" validate:"required"` // From session	State        string `json:"state"`	Scope        string `json:"scope"`	RedirectURI  string `json:"redirect_uri" validate:"required,url"`	ClientID     string `json:"client_id" validate:"required"`	ResponseType string `json:"response_type" validate:"required"`type AuthorizeRequest struct {// AuthorizeRequest represents an OAuth2 authorization request}	Active       *bool    `json:"active,omitempty"`	RedirectURIs []string `json:"redirectUris" validate:"omitempty,min=1,dive,url"`	Description  *string  `json:"description,omitempty"`	Name         string   `json:"name" validate:"max=255"`	ClientID     string   `json:"-" validate:"required"`type UpdateOAuth2ClientRequest struct {// UpdateOAuth2ClientRequest represents an update OAuth2 client request}	LogoURL      *string  `json:"logoUrl,omitempty" validate:"omitempty,url"`	Scopes       []string `json:"scopes" validate:"required,min=1"`	GrantTypes   []string `json:"grantTypes" validate:"required,min=1"`	RedirectURIs []string `json:"redirectUris" validate:"required,min=1,dive,url"`	Description  *string  `json:"description,omitempty"`	Name         string   `json:"name" validate:"required,max=255"`	OwnerID      string   `json:"-" validate:"required"`type CreateOAuth2ClientRequest struct {// CreateOAuth2ClientRequest represents a create OAuth2 client request}	UpdatedAt    int64    `json:"updatedAt"`	CreatedAt    int64    `json:"createdAt"`	Active       bool     `json:"active"`	LogoURL      *string  `json:"logoUrl,omitempty"`	Scopes       []string `json:"scopes"`	GrantTypes   []string `json:"grantTypes"`