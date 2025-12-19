package users

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewuserSerive(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *User) error {
	hashed, _ := bcrypt.GenerateFromPassword(
		[]byte(user.Password), 12,
	)
	user.Password = string(hashed)

	result := s.DB.Create(user).Error

	return result
}
