// Package main is the entry point for the FitClassMaster application.
// It initializes the database, sessions, templates, and starts the HTTP server.
package main

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/server"
	"FitClassMaster/internal/templates"
	"os"

	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Database connection and Session store.
	config.InitDB()
	config.InitSessionStore()

	// Initialize template parsing for HTML rendering.
	templates.Init()

	// Auto-migrate database schemas to ensure they match the models.
	config.RunMigrations()

	// Initialize Server (Services, Repos, Routes)
	app := server.New()

	// Get the PORT from the environment (Render sets this automatically).
	// If it is empty (running locally), default to "8080".
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Listen on the dynamic port
	log.Printf("✅ Server running on port %s", port)
	if err := app.Run(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
