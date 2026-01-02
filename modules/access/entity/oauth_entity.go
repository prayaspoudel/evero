package entity
package entity

import "time"

// OAuth2Client represents an OAuth2 application
type OAuth2Client struct {
	ID           string   `gorm:"column:id;primaryKey"`
	ClientID     string   `gorm:"column:client_id;uniqueIndex;not null"`
	ClientSecret string   `gorm:"column:client_secret;not null"` // Hashed
	Name         string   `gorm:"column:name;not null"`
	Description  *string  `gorm:"column:description"`
	RedirectURIs string   `gorm:"column:redirect_uris;type:text"` // JSON array as string
	GrantTypes   string   `gorm:"column:grant_types;type:text"`   // JSON array as string
	Scopes       string   `gorm:"column:scopes;type:text"`        // JSON array as string
	OwnerID      string   `gorm:"column:owner_id;not null"`
	LogoURL      *string  `gorm:"column:logo_url"`
	Active       bool     `gorm:"column:active;default:true"`
	CreatedAt    int64    `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    int64    `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (oc *OAuth2Client) TableName() string {
	return "sso_oauth_clients"
}

// OAuth2AuthorizationCode represents an authorization code
type OAuth2AuthorizationCode struct {
	ID          string     `gorm:"column:id;primaryKey"`
	Code        string     `gorm:"column:code;uniqueIndex;not null"`
	ClientID    string     `gorm:"column:client_id;not null"`
	UserID      string     `gorm:"column:user_id;not null"`
	RedirectURI string     `gorm:"column:redirect_uri;not null"`
	Scopes      string     `gorm:"column:scopes;type:text"` // JSON array as string
	ExpiresAt   time.Time  `gorm:"column:expires_at;not null"`
	UsedAt      *time.Time `gorm:"column:used_at"`
	CreatedAt   int64      `gorm:"column:created_at;autoCreateTime:milli"`
}

func (oac *OAuth2AuthorizationCode) TableName() string {
	return "sso_oauth_authorization_codes"
}

// OAuth2Token represents an OAuth2 access token
type OAuth2Token struct {














}	return "sso_oauth_tokens"func (ot *OAuth2Token) TableName() string {}	CreatedAt    int64      `gorm:"column:created_at;autoCreateTime:milli"`	RevokedAt    *time.Time `gorm:"column:revoked_at"`	ExpiresAt    time.Time  `gorm:"column:expires_at;not null"`	Scopes       string     `gorm:"column:scopes;type:text"` // JSON array as string	UserID       string     `gorm:"column:user_id;not null"`	ClientID     string     `gorm:"column:client_id;not null"`	RefreshToken *string    `gorm:"column:refresh_token"`	AccessToken  string     `gorm:"column:access_token;uniqueIndex;not null"`	ID           string     `gorm:"column:id;primaryKey"`