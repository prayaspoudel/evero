package entity
package entity

























}	return "sso_users"func (u *User) TableName() string {}	Company       *Company   `gorm:"foreignKey:company_id;references:id"`	UpdatedAt     int64      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`	CreatedAt     int64      `gorm:"column:created_at;autoCreateTime:milli"`	LastLoginIP   string     `gorm:"column:last_login_ip"`	LastLoginAt   *time.Time `gorm:"column:last_login_at"`	EmailVerified bool       `gorm:"column:email_verified;default:false"`	IsVerified    bool       `gorm:"column:is_verified;default:false"`	IsActive      bool       `gorm:"column:is_active;default:true"`	Role          string     `gorm:"column:role"`	CompanyID     string     `gorm:"column:company_id"`	LastName      string     `gorm:"column:last_name"`	FirstName     string     `gorm:"column:first_name"`	PasswordHash  string     `gorm:"column:password_hash;not null"`	Email         string     `gorm:"column:email;uniqueIndex;not null"`	ID            string     `gorm:"column:id;primaryKey"`type User struct {// User represents the SSO user entityimport "time"