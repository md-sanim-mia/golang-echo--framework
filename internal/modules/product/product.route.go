package product

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/md-sanim-mia/golang-first-project/internal/middleware"
)

func ProudctRoute(e *echo.Group, contller *ProductContller) {

	jwtSecret := os.Getenv("JWT_SECRET")

	//protucted route ------------------------
	protucted := e.Group("/products", middleware.AuthMiddleware(jwtSecret))
	protucted.POST("", contller.CreateProduct, middleware.RoleMiddleware("USER"))

}
