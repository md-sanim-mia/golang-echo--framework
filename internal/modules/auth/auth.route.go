package auth

import "github.com/labstack/echo/v4"

func AuthRoute(e *echo.Group, controller *AuthContller) {

	e.POST("/auth/login", controller.loginUser)
}
