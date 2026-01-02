package repository

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	Repository[entity.RefreshToken]
	Log *logrus.Logger
}

func NewRefreshTokenRepository(log *logrus.Logger) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		Log: log,
	}
}

func (r *RefreshTokenRepository) FindByToken(db *gorm.DB, token *entity.RefreshToken, tokenString string) error {
	return db.Where("token = ? AND revoked = false AND expires_at > NOW()", tokenString).First(token).Error
}

func (r *RefreshTokenRepository) RevokeByUserID(db *gorm.DB, userID string) error {
	return db.Model(&entity.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeByToken(db *gorm.DB, tokenString string) error {
	return db.Model(&entity.RefreshToken{}).
		Where("token = ?", tokenString).
		Update("revoked", true).Error
}

func (r *RefreshTokenRepository) DeleteExpired(db *gorm.DB) error {
	return db.Where("expires_at <= NOW()").Delete(&entity.RefreshToken{}).Error
}
