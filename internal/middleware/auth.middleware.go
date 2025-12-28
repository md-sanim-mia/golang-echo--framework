package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecrect string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {

				return c.JSON(http.StatusUnauthorized, echo.Map{

					"message": "Missing authorization header",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{

					"error": "Invalid or expired token",
				})
			}
			fmt.Println("token", tokenString)

			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {

				return []byte(jwtSecrect), nil
			})

			if err != nil || !token.Valid {

				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			claims, ok := token.Claims.(*JWTClaims)

			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			fmt.Printf("JWT Claims: %+v\n", claims)

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("full_name", claims.FullName)
			c.Set("role", claims.Role)

			return next(c)

		}
	}
}
