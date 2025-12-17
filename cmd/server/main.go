package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/md-sanim-mia/golang-first-project/internal/config"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.CoonectDB()
	defer config.CloseDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {

		return c.String(http.StatusOK, "Ecom project server is running!")
	})

	e.GET("/health", func(c echo.Context) error {

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "Ok",
			"message": "Server is healthy",
		})
	})

	log.Printf("🚀 Server starting on port %s", ":1323")
	e.Logger.Fatal(e.Start(":1323"))

}
