package users

import "github.com/labstack/echo/v4"

func UserRoutes(e *echo.Group, controller *UserController) {
	e.POST("/users", controller.CreateUserHandler)
	e.GET("/users", controller.GetAllUsers)

	e.GET("/users/:id", controller.GetsingleUserById)
	e.PATCH("/users/:id", controller.UpdateUserHandler)

	e.DELETE("/users/:id", controller.DeleteUserHandiler)
}
