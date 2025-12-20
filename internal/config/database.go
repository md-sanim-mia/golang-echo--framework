package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func CoonectDB() {

	dbPath := os.Getenv("DATABASE_URL")

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	}

	db, err := gorm.Open(postgres.Open(dbPath), config)

	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	DB = db
	pSqlDb, err := db.DB()
	if err != nil {

		log.Fatal("❌ Database connection failed:", err)
	}

	pSqlDb.SetMaxIdleConns(10)
	pSqlDb.SetMaxOpenConns(100)
	pSqlDb.SetConnMaxLifetime(time.Hour)

	if err := pSqlDb.Ping(); err != nil {

		log.Fatal("❌ Database connection failed:", err)
	}
	log.Println("✅ PostgreSQL connected successfully!")

}

func GetDB() *gorm.DB {
	return DB
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
