package repository

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TwoFactorRepository struct {
	Repository[entity.UserTwoFactor]
	Log *logrus.Logger
}

func NewTwoFactorRepository(log *logrus.Logger) *TwoFactorRepository {
	return &TwoFactorRepository{
		Log: log,
	}
}

func (r *TwoFactorRepository) FindByUserID(db *gorm.DB, twoFactor *entity.UserTwoFactor, userID string) error {
	return db.Where("user_id = ?", userID).First(twoFactor).Error
}

type BackupCodeRepository struct {
	Repository[entity.BackupCode]
	Log *logrus.Logger
}

func NewBackupCodeRepository(log *logrus.Logger) *BackupCodeRepository {
	return &BackupCodeRepository{
		Log: log,
	}
}

func (r *BackupCodeRepository) FindByUserIDAndCode(db *gorm.DB, code *entity.BackupCode, userID, codeStr string) error {
	return db.Where("user_id = ? AND code = ? AND used_at IS NULL", userID, codeStr).First(code).Error
}

func (r *BackupCodeRepository) DeleteByUserID(db *gorm.DB, userID string) error {
	return db.Where("user_id = ?", userID).Delete(&entity.BackupCode{}).Error
}
