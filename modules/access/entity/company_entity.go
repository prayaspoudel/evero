package entity

// Company represents a company/organization entity
type Company struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name;not null"`
	Email     string `gorm:"column:email"`
	Industry  string `gorm:"column:industry"`
	Status    string `gorm:"column:status;default:'active'"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (c *Company) TableName() string {
	return "sso_companies"
}
