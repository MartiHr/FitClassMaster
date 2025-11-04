package main

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/handlers"
	"FitClassMaster/internal/middlewares"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/templates"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Init DB & session
	config.InitDB()
	config.InitSessionStore()

	// Init template parsing
	templates.Init()

	// Auto migrate
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Router setup
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)    // chi logger
	r.Use(middleware.Recoverer) // chi recoverer
	//r.Use(middlewares.LoadSession)

	// Static files
	fileServer := http.FileServer(http.Dir("internal/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Handlers
	homeH := handlers.NewHomeHandler()
	authH := handlers.NewAuthHandler()

	// Public routes
	r.Get("/", homeH.Home)
	r.Get("/htmx/hello", homeH.HelloHtmx)
	r.Get("/register", authH.RegisterPage)
	r.Post("/register", authH.RegisterPost)
	r.Get("/login", authH.LoginPage)
	r.Post("/login", authH.LoginPost)
	r.Post("/logout", authH.Logout)

	// Protected example route
	r.With(middlewares.RequireAuth).Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to your dashboard!"))
	})

	// Run server
	log.Println("✅ Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
