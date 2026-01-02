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

func (r *PasswordResetTokenRepository) FindByToken(db *gorm.DB, token *entity.PasswordResetToken, tokenStr string) error {
	return db.Where("token = ? AND used_at IS NULL AND expires_at > NOW()", tokenStr).First(token).Error
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

func (r *EmailVerificationTokenRepository) FindByToken(db *gorm.DB, token *entity.EmailVerificationToken, tokenStr string) error {
	return db.Where("token = ? AND used_at IS NULL AND expires_at > NOW()", tokenStr).First(token).Error
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

func (r *AuditLogRepository) FindByUserID(db *gorm.DB, logs *[]entity.AuditLog, userID string, limit int) error {
	return db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(logs).Error
}
