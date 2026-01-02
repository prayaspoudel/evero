package repository
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





































}	return db.Where("user_id = ?", userID).Delete(&entity.BackupCode{}).Errorfunc (r *BackupCodeRepository) DeleteByUserID(db *gorm.DB, userID string) error {}		Update("used_at", gorm.Expr("NOW()")).Error		Where("id = ?", id).	return db.Model(&entity.BackupCode{}).func (r *BackupCodeRepository) MarkAsUsed(db *gorm.DB, id string) error {}	return codes, err	err := db.Where("user_id = ? AND used_at IS NULL", userID).Find(&codes).Error	var codes []entity.BackupCodefunc (r *BackupCodeRepository) FindByUserID(db *gorm.DB, userID string) ([]entity.BackupCode, error) {}	}		Log: log,	return &BackupCodeRepository{func NewBackupCodeRepository(log *logrus.Logger) *BackupCodeRepository {}	Log *logrus.Logger	Repository[entity.BackupCode]type BackupCodeRepository struct {}	return db.Where("user_id = ?", userID).Delete(&entity.UserTwoFactor{}).Errorfunc (r *TwoFactorRepository) DeleteByUserID(db *gorm.DB, userID string) error {}		Update("status", status).Error		Where("user_id = ?", userID).	return db.Model(&entity.UserTwoFactor{}).func (r *TwoFactorRepository) UpdateStatus(db *gorm.DB, userID string, status entity.TwoFactorStatus) error {