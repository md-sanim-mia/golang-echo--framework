package users

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/md-sanim-mia/golang-first-project/internal/middleware"
)

func UserRoutes(e *echo.Group, controller *UserController) {
	jwtSecret := os.Getenv("JWT_SECRET")
	e.POST("/users", controller.CreateUserHandler)

	protected := e.Group("/users", middleware.AuthMiddleware(jwtSecret))
	protected.GET("", controller.GetAllUsers, middleware.RoleMiddleware("USER"))

	protected.GET("/:id", controller.GetsingleUserById)
	protected.PATCH("/:id", controller.UpdateUserHandler)

	protected.DELETE("/:id", controller.DeleteUserHandiler)
}
