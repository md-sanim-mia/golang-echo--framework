package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId uint `json:"user_id"`

	Email string `json:"email"`

	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(userId uint, email string, role string) (string, error) {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "kfdskalj40r9sadfmsdaofmas"
	}

	claims := JWTClaims{
		UserId: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
