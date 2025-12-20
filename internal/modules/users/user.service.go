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

func (s *UserService) GetAllUsers() ([]User, error) {

	var users []User

	result := s.DB.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (s *UserService) GetsingleUserById(id uint) (*User, error) {

	var user User

	result := s.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil

}

func (s *UserService) UpdateUser(id uint, payload User) (*User, error) {

	var user User

	result := s.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	err := s.DB.Model(&user).Updates(payload).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
