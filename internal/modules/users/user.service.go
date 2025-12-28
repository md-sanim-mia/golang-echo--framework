package users

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/md-sanim-mia/golang-first-project/internal/builder"
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

	fmt.Println("ORIGINAL PASSWORD BEFORE HASHING:", user.Password)

	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),

		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	fmt.Println("HASHED PASSWORD:", string(hashed))

	user.Password = string(hashed)

	return s.DB.Create(user).Error
}

func (s *UserService) GetAllUsers(c echo.Context) ([]User, *builder.PaginationMeta, error) {

	var users []User

	qd := builder.NewQueryBuilder(s.DB, c).Model(&User{}).Search([]string{"fullName", "email"}).Filter().Sort().Paginate()

	// result := s.DB.Find(&users)

	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	// return users, nil

	if err := qd.Execute(&users); err != nil {
		return nil, nil, err
	}

	// Get pagination meta
	meta, err := qd.CountTotal()
	if err != nil {
		return nil, nil, err
	}

	return users, meta, nil
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
func (s *UserService) DeleteUser(id uint) (*User, error) {
	var user User

	// check user exists
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	fmt.Println("check user id", id)

	// correct delete
	if err := s.DB.Delete(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
