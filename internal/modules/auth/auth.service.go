package auth

import "gorm.io/gorm"

type AuthService struct {
	*gorm.DB
}

func (s AuthService) LoginUser()
