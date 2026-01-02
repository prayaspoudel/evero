package entity

import "time"

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id;not null"`
	Token     string    `gorm:"column:token;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	Used      bool      `gorm:"column:used;default:false"`
	CreatedAt int64     `gorm:"column:created_at;autoCreateTime:milli"`
}

func (prt *PasswordResetToken) TableName() string {
	return "sso_password_reset_tokens"
}

// EmailVerificationToken represents an email verification token
type EmailVerificationToken struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id;not null"`
	Token     string    `gorm:"column:token;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	Verified  bool      `gorm:"column:verified;default:false"`
	CreatedAt int64     `gorm:"column:created_at;autoCreateTime:milli"`
}

func (evt *EmailVerificationToken) TableName() string {
	return "sso_email_verification_tokens"
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string `gorm:"column:id;primaryKey"`
	UserID    string `gorm:"column:user_id"`
	Action    string `gorm:"column:action;not null"`
	Resource  string `gorm:"column:resource"`
	Details   string `gorm:"column:details;type:jsonb"` // JSON details
	IPAddress string `gorm:"column:ip_address"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
}

func (al *AuditLog) TableName() string {
	return "sso_audit_logs"
}
