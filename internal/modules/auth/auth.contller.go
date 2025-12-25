package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthContller struct {
	service AuthService
}

func NewAuthController(service AuthService) *AuthContller {
	return &AuthContller{
		service: service,
	}
}
func (c *AuthContller) loginUser(ctx echo.Context) error {

	payload := new(Payload)

	if err := ctx.Bind(payload); err != nil {

		return ctx.JSON(http.StatusBadRequest, echo.Map{

			"message": "Email and password ar required",
		})
	}

	if payload.Email == "" || payload.Password == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "Email and password are required",
		})
	}

	res, err := c.service.LoginUser(payload)

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "User login successfully !",

		"data": res.Token,
	})

}
