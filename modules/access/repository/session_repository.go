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
	return db.Where("session_token = ? AND expires_at > NOW()", token).First(session).Error
}

func (r *SessionRepository) FindByUserID(db *gorm.DB, userID string) ([]entity.Session, error) {
	var sessions []entity.Session
	err := db.Where("user_id = ? AND expires_at > NOW()", userID).Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) DeleteByUserID(db *gorm.DB, userID string) error {
	return db.Where("user_id = ?", userID).Delete(&entity.Session{}).Error
}

func (r *SessionRepository) DeleteExpired(db *gorm.DB) error {
	return db.Where("expires_at <= NOW()").Delete(&entity.Session{}).Error
}
