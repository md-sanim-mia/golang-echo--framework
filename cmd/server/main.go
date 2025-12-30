package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/md-sanim-mia/golang-first-project/internal/config"
	"github.com/md-sanim-mia/golang-first-project/internal/modules/auth"
	"github.com/md-sanim-mia/golang-first-project/internal/modules/product"
	"github.com/md-sanim-mia/golang-first-project/internal/modules/users"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.CoonectDB()
	defer config.CloseDB()
	if err := config.DB.AutoMigrate(&users.User{}); err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BodyLimit("20M"))
	e.Use(middleware.Secure())

	e.GET("/", func(c echo.Context) error {

		return c.String(http.StatusOK, "Ecom project server is running!")
	})

	e.GET("/health", func(c echo.Context) error {

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "Ok",
			"message": "Server is healthy",
		})
	})
	userService := &users.UserService{
		DB: config.DB,
	}

	userController := &users.UserController{
		Service: userService,
	}

	authService := &auth.AuthService{
		DB: config.DB,
	}

	authController := auth.NewAuthController(*authService)

	productService := &product.ProductService{
		DB: config.DB,
	}

	productContller := &product.ProductContller{
		Service: productService,
	}

	api := e.Group("/api/v1")

	users.UserRoutes(api, userController)

	auth.AuthRoute(api, authController)
	product.ProudctRoute(api, productContller)

	log.Printf("🚀 Server starting on port %s", ":1323")
	e.Logger.Fatal(e.Start(":1323"))

}
