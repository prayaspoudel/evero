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

func (r *RefreshTokenRepository) FindByToken(db *gorm.DB, token *entity.RefreshToken, tokenStr string) error {
	return db.Where("token = ? AND revoked_at IS NULL", tokenStr).First(token).Error
}

func (r *RefreshTokenRepository) RevokeByToken(db *gorm.DB, tokenStr string) error {
	return db.Model(&entity.RefreshToken{}).
		Where("token = ?", tokenStr).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}

func (r *RefreshTokenRepository) RevokeByUserID(db *gorm.DB, userID string) error {
	return db.Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}

func (r *RefreshTokenRepository) DeleteExpired(db *gorm.DB) error {
	return db.Where("expires_at < NOW()").Delete(&entity.RefreshToken{}).Error
}
