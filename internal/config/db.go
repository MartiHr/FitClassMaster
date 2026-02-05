// Package config handles application configuration including database and session setup.
package config

import (
	"FitClassMaster/internal/models"
	"encoding/gob"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
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

	var dialector gorm.Dialector

	// SMART SWITCH:
	// If the URL starts with "postgres", we use the Cloud driver.
	// Otherwise, we default to SQL Server (your local Docker/Windows setup).
	if strings.HasPrefix(dsn, "postgres") {
		log.Println("🐘 Connecting to PostgreSQL (Cloud Mode)...")
		dialector = postgres.Open(dsn)
	} else {
		log.Println("🖥️ Connecting to SQL Server (Local Mode)...")
		dialector = sqlserver.Open(dsn)
	}

	// Open connection using the selected dialector
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
}

func RunMigrations() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.Enrollment{},
		&models.Exercise{},
		&models.WorkoutPlan{},
		&models.WorkoutExercise{},
		&models.WorkoutSession{},
		&models.SessionLog{},
		&models.Conversation{},
		&models.Message{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
