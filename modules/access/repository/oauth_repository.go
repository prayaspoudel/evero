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
	return db.Where("client_id = ? AND active = ?", clientID, true).First(client).Error
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

func (r *OAuth2AuthCodeRepository) FindByCode(db *gorm.DB, code *entity.OAuth2AuthorizationCode, codeStr string) error {
	return db.Where("code = ? AND used_at IS NULL AND expires_at > NOW()", codeStr).First(code).Error
}

type OAuth2TokenRepository struct {
	Repository[entity.OAuth2Token]
	Log *logrus.Logger
}

func NewOAuth2TokenRepository(log *logrus.Logger) *OAuth2TokenRepository {
	return &OAuth2TokenRepository{
		Log: log,
	}
}

func (r *OAuth2TokenRepository) FindByAccessToken(db *gorm.DB, token *entity.OAuth2Token, accessToken string) error {
	return db.Where("access_token = ? AND revoked_at IS NULL AND expires_at > NOW()", accessToken).First(token).Error
}

func (r *OAuth2TokenRepository) RevokeByAccessToken(db *gorm.DB, accessToken string) error {
	return db.Model(&entity.OAuth2Token{}).
		Where("access_token = ?", accessToken).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}
