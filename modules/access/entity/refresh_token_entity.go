package entity

import "time"

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id;not null"`
	Token     string    `gorm:"column:token;uniqueIndex;not null"`
	ClientID  string    `gorm:"column:client_id"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt int64     `gorm:"column:created_at;autoCreateTime:milli"`
	Revoked   bool      `gorm:"column:revoked;default:false"`
}

func (rt *RefreshToken) TableName() string {
	return "sso_refresh_tokens"
}
