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

func (r *CompanyRepository) FindByID(db *gorm.DB, company *entity.Company, id string) error {
	return db.Where("id = ?", id).First(company).Error
}

func (r *CompanyRepository) FindByUserID(db *gorm.DB, userID string) ([]entity.Company, error) {
	var companies []entity.Company
	err := db.Table("sso_companies c").
		Joins("INNER JOIN sso_user_companies uc ON c.id = uc.company_id").
		Where("uc.user_id = ?", userID).
		Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) List(db *gorm.DB, page int, size int, status string) ([]entity.Company, int64, error) {
	var companies []entity.Company
	var total int64

	offset := (page - 1) * size
	query := db.Model(&entity.Company{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(size).Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}
