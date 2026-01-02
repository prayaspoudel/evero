package repository

import (
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SessionRepository struct {
	Repository[entity.Session]
	Log *logrus.Logger
}

func NewSessionRepository(log *logrus.Logger) *SessionRepository {
	return &SessionRepository{
		Log: log,
	}
}

func (r *SessionRepository) FindByToken(db *gorm.DB, session *entity.Session, token string) error {
	return db.Where("session_token = ?", token).First(session).Error
}

func (r *SessionRepository) DeleteByToken(db *gorm.DB, token string) error {
	return db.Where("session_token = ?", token).Delete(&entity.Session{}).Error
}

func (r *SessionRepository) DeleteExpired(db *gorm.DB) error {
	return db.Where("expires_at < NOW()").Delete(&entity.Session{}).Error
}
