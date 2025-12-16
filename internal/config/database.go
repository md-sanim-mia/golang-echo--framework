package config

import (
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
func CoonectDB() {

	dbPath := os.Getenv("DATABASE_URL")

	config:=*&gorm.Config{
		Logger: logger
	}
}
