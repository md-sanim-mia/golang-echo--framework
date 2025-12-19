package users

import "github.com/labstack/echo/v4"

func UserRoutes(e *echo.Group, controller *UserController) {
	e.POST("/users", controller.CreateUserHandler)
}
