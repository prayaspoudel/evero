package entity

// Company represents a company/organization entity
type Company struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name;not null"`
	Domain    string `gorm:"column:domain"`
	Email     string `gorm:"column:email"`
	Industry  string `gorm:"column:industry"`
	IsActive  bool   `gorm:"column:is_active;default:true"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (c *Company) TableName() string {
	return "sso_companies"
}
