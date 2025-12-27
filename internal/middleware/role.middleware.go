package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			userRole := c.Get("role").(string)

			for _, role := range allowedRoles {

				if userRole == role {

					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "You don't have permission to access this resource",
			})
		}
	}
}
