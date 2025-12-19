package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	Service *UserService
}

func NewUsercontroller(service *UserService) *UserController {
	return &UserController{Service: service}
}

func (c *UserController) CreateUserHandler(ctx echo.Context) error {

	user := new(User)

	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request",
		})
	}

	if err := c.Service.CreateUser(user); err != nil {

		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, user)

}
