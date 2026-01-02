package entity

import "time"

// User represents the SSO user entity
type User struct {
	ID            string     `gorm:"column:id;primaryKey"`
	Email         string     `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash  string     `gorm:"column:password_hash;not null"`
	FirstName     string     `gorm:"column:first_name"`
	LastName      string     `gorm:"column:last_name"`
	CompanyID     string     `gorm:"column:company_id"`
	Role          string     `gorm:"column:role"`
	IsActive      bool       `gorm:"column:is_active;default:true"`
	IsVerified    bool       `gorm:"column:is_verified;default:false"`
	EmailVerified bool       `gorm:"column:email_verified;default:false"`
	LastLoginAt   *time.Time `gorm:"column:last_login_at"`
	LastLoginIP   string     `gorm:"column:last_login_ip"`
	CreatedAt     int64      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt     int64      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Company       *Company   `gorm:"foreignKey:company_id;references:id"`
}

func (u *User) TableName() string {
	return "sso_users"
}
