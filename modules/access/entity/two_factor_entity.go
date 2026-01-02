package entity
package entity

import "time"

// TwoFactorMethod represents the type of 2FA method
type TwoFactorMethod string

const (
	TwoFactorMethodTOTP TwoFactorMethod = "totp"
	TwoFactorMethodSMS  TwoFactorMethod = "sms"
)

// TwoFactorStatus represents the status of 2FA for a user
type TwoFactorStatus string

const (
	TwoFactorStatusDisabled TwoFactorStatus = "disabled"
	TwoFactorStatusPending  TwoFactorStatus = "pending"
	TwoFactorStatusEnabled  TwoFactorStatus = "enabled"
)

// UserTwoFactor represents a user's 2FA settings
type UserTwoFactor struct {
	ID               string          `gorm:"column:id;primaryKey"`
	UserID           string          `gorm:"column:user_id;not null;uniqueIndex"`
	Method           TwoFactorMethod `gorm:"column:method;not null"`
	Secret           string          `gorm:"column:secret"`
	PhoneNumber      *string         `gorm:"column:phone_number"`
	Status           TwoFactorStatus `gorm:"column:status;default:'disabled'"`
	BackupCodesCount int             `gorm:"column:backup_codes_count;default:0"`
	VerifiedAt       *time.Time      `gorm:"column:verified_at"`
	CreatedAt        int64           `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt        int64           `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

















}	return "sso_backup_codes"func (bc *BackupCode) TableName() string {}	CreatedAt int64      `gorm:"column:created_at;autoCreateTime:milli"`	UsedAt    *time.Time `gorm:"column:used_at"`	Code      string     `gorm:"column:code;not null"` // Hashed	UserID    string     `gorm:"column:user_id;not null"`	ID        string     `gorm:"column:id;primaryKey"`type BackupCode struct {// BackupCode represents a 2FA backup code}	return "sso_user_two_factors"func (utf *UserTwoFactor) TableName() string {