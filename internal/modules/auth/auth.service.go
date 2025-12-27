package auth

import (
	"errors"
	"fmt"

	"github.com/md-sanim-mia/golang-first-project/internal/modules/users"
	"github.com/md-sanim-mia/golang-first-project/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	*gorm.DB
}

type LoginResponse struct {
	Token string
}

type Payload struct {
	Email    string
	Password string
}

func (s AuthService) LoginUser(payload *Payload) (*LoginResponse, error) {

	fmt.Println("auth password", payload)
	var user users.User

	if err := s.DB.Where("email=?", payload.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	fmt.Print("user check", user)

	fmt.Println("login password:", payload.Password)
	fmt.Println("hashed password:", user.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {

		fmt.Println("bcrypt error:", err)
		return nil, errors.New("invalid email or password !")

	}
	fmt.Printf("check user data after login: %+v\n", user)
	token, err := utils.GenerateJwtToken(user.ID, user.Email, user.Role.String(), user.FullName)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{

		Token: token,
	}, nil

}
