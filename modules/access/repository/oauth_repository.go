package repository
package repository

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OAuth2ClientRepository struct {
	Repository[entity.OAuth2Client]
	Log *logrus.Logger
}

func NewOAuth2ClientRepository(log *logrus.Logger) *OAuth2ClientRepository {
	return &OAuth2ClientRepository{
		Log: log,
	}
}

func (r *OAuth2ClientRepository) FindByClientID(db *gorm.DB, client *entity.OAuth2Client, clientID string) error {
	return db.Where("client_id = ?", clientID).First(client).Error
}

func (r *OAuth2ClientRepository) FindByOwnerID(db *gorm.DB, ownerID string) ([]entity.OAuth2Client, error) {
	var clients []entity.OAuth2Client
	err := db.Where("owner_id = ?", ownerID).Find(&clients).Error
	return clients, err
}

type OAuth2AuthCodeRepository struct {
	Repository[entity.OAuth2AuthorizationCode]
	Log *logrus.Logger
}

func NewOAuth2AuthCodeRepository(log *logrus.Logger) *OAuth2AuthCodeRepository {
	return &OAuth2AuthCodeRepository{
		Log: log,
	}
}



































}	return db.Where("expires_at <= NOW()").Delete(&entity.OAuth2Token{}).Errorfunc (r *OAuth2TokenRepository) DeleteExpired(db *gorm.DB) error {}		Update("revoked_at", gorm.Expr("NOW()")).Error		Where("access_token = ?", accessToken).	return db.Model(&entity.OAuth2Token{}).func (r *OAuth2TokenRepository) RevokeByAccessToken(db *gorm.DB, accessToken string) error {}	return db.Where("access_token = ? AND revoked_at IS NULL AND expires_at > NOW()", accessToken).First(token).Errorfunc (r *OAuth2TokenRepository) FindByAccessToken(db *gorm.DB, token *entity.OAuth2Token, accessToken string) error {}	}		Log: log,	return &OAuth2TokenRepository{func NewOAuth2TokenRepository(log *logrus.Logger) *OAuth2TokenRepository {}	Log *logrus.Logger	Repository[entity.OAuth2Token]type OAuth2TokenRepository struct {}		Update("used_at", gorm.Expr("NOW()")).Error		Where("id = ?", id).	return db.Model(&entity.OAuth2AuthorizationCode{}).func (r *OAuth2AuthCodeRepository) MarkAsUsed(db *gorm.DB, id string) error {}	return db.Where("code = ? AND used_at IS NULL AND expires_at > NOW()", code).First(authCode).Errorfunc (r *OAuth2AuthCodeRepository) FindByCode(db *gorm.DB, authCode *entity.OAuth2AuthorizationCode, code string) error {