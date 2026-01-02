package repository

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}

func (r *UserRepository) FindByID(db *gorm.DB, user *entity.User, id string) error {
	return db.Preload("Company").Where("id = ?", id).First(user).Error
}

func (r *UserRepository) UpdateLastLogin(db *gorm.DB, userID string, ip string) error {
	return db.Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error
}

func (r *UserRepository) UpdateEmailVerified(db *gorm.DB, userID string) error {
	return db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("email_verified", true).Error
}
