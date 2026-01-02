package entity

// UserCompany represents the relationship between users and companies
type UserCompany struct {
	ID        string `gorm:"column:id;primaryKey"`
	UserID    string `gorm:"column:user_id;not null"`
	CompanyID string `gorm:"column:company_id;not null"`
	Role      string `gorm:"column:role"`
	IsPrimary bool   `gorm:"column:is_primary;default:false"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
}

func (uc *UserCompany) TableName() string {
	return "sso_user_companies"
}
