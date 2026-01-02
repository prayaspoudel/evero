package repository

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	Repository[entity.Company]
	Log *logrus.Logger
}

func NewCompanyRepository(log *logrus.Logger) *CompanyRepository {
	return &CompanyRepository{
		Log: log,
	}
}

func (r *CompanyRepository) FindByUserID(db *gorm.DB, companies *[]entity.Company, userID string) error {
	return db.Table("sso_companies c").
		Joins("INNER JOIN sso_user_companies uc ON c.id = uc.company_id").
		Where("uc.user_id = ?", userID).
		Find(companies).Error
}

func (r *CompanyRepository) FindByDomain(db *gorm.DB, company *entity.Company, domain string) error {
	return db.Where("domain = ? AND is_active = ?", domain, true).First(company).Error
}
