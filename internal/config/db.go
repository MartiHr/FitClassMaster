// Package config handles application configuration including database and session setup.
package config

import (
	"encoding/gob"
	"log"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// DB is the global database connection pool.
var DB *gorm.DB

// InitDB initializes the database connection using the DB_DSN environment variable.
// It also registers types for session serialization.
func InitDB() {
	// Register uint so it can be stored in the session (used for UserID).
	gob.Register(uint(0))

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN env var not set")
	}

	// Open connection to SQL Server.
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
}
