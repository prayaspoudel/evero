package repository
package repository

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PasswordResetTokenRepository struct {
	Repository[entity.PasswordResetToken]
	Log *logrus.Logger
}

func NewPasswordResetTokenRepository(log *logrus.Logger) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{
		Log: log,
	}
}

func (r *PasswordResetTokenRepository) FindByToken(db *gorm.DB, token *entity.PasswordResetToken, tokenString string) error {
	return db.Where("token = ? AND used = false AND expires_at > NOW()", tokenString).First(token).Error
}

func (r *PasswordResetTokenRepository) MarkAsUsed(db *gorm.DB, id string) error {
	return db.Model(&entity.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used", true).Error
}

type EmailVerificationTokenRepository struct {
	Repository[entity.EmailVerificationToken]
	Log *logrus.Logger
}

func NewEmailVerificationTokenRepository(log *logrus.Logger) *EmailVerificationTokenRepository {
	return &EmailVerificationTokenRepository{
		Log: log,
	}
}

func (r *EmailVerificationTokenRepository) FindByToken(db *gorm.DB, token *entity.EmailVerificationToken, tokenString string) error {
	return db.Where("token = ? AND verified = false AND expires_at > NOW()", tokenString).First(token).Error
}

func (r *EmailVerificationTokenRepository) MarkAsVerified(db *gorm.DB, id string) error {
	return db.Model(&entity.EmailVerificationToken{}).
		Where("id = ?", id).
		Update("verified", true).Error
}

type AuditLogRepository struct {
	Repository[entity.AuditLog]
	Log *logrus.Logger
}

func NewAuditLogRepository(log *logrus.Logger) *AuditLogRepository {
	return &AuditLogRepository{
		Log: log,
	}
}






















}	return logs, total, nil	}		return nil, 0, err	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&logs).Error; err != nil {	}		return nil, 0, err	if err := query.Count(&total).Error; err != nil {	}		query = query.Where("user_id = ?", userID)	if userID != "" {	query := db.Model(&entity.AuditLog{})	offset := (page - 1) * size	var total int64	var logs []entity.AuditLogfunc (r *AuditLogRepository) FindByUserID(db *gorm.DB, userID string, page int, size int) ([]entity.AuditLog, int64, error) {