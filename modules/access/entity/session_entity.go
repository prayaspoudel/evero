package entity
package entity

import "time"

// Session represents a user session
type Session struct {
	ID           string    `gorm:"column:id;primaryKey"`
	UserID       string    `gorm:"column:user_id;not null"`
	SessionToken string    `gorm:"column:session_token;uniqueIndex;not null"`
	IPAddress    string    `gorm:"column:ip_address"`
	UserAgent    string    `gorm:"column:user_agent"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null"`
	CreatedAt    int64     `gorm:"column:created_at;autoCreateTime:milli"`
}

func (s *Session) TableName() string {
	return "sso_sessions"
}
