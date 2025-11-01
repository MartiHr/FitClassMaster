package config

import (
	"log"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN env var not set")
	}

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nill {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
}
