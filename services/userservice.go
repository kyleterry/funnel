package services

import (
	"github.com/jinzhu/gorm"
	"github.com/kyleterry/funnel/data"
)

type UserService struct {
	db *data.FunnelDB
}

func NewUserService(db *data.FunnelDB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Exists(user *data.User) (bool, error) {
	if err := s.db.Conn.First(&user, user.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
