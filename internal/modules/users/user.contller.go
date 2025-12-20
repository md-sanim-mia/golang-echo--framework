package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	Service *UserService
}

func NewUserController(service *UserService) *UserController {
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

func (c *UserController) GetAllUsers(ctx echo.Context) error {

	users, err := c.Service.GetAllUsers()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{

			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, users)

}

func (c *UserController) GetsingleUserById(ctx echo.Context) error {

	idParam := ctx.Param("id")
	fmt.Println("check id......................", idParam)
	id64, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid user Id",
		})
	}

	user, err := c.Service.GetsingleUserById(uint(id64))

	if err != nil {

		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (c *UserController) UpdateUserHandler(ctx echo.Context) error {

	idParam := ctx.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid user id",
		})
	}


	user,err:=c.Service.UpdateUser(uint(id64),)

}
