package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *User) error {

	if s.DB == nil {
		return errors.New("database not initialized")
	}


	isExist := s.DB.Where("email=?", user.Email).First(&User{}).Error
	if isExist == nil {
		return errors.New("user already exists")
	}
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashed)

	return s.DB.Create(user).Error
}
