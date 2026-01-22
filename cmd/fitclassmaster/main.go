package main

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/handlers"
	"FitClassMaster/internal/middlewares"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"FitClassMaster/internal/services"
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
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.Enrollment{}); err != nil {
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

	// Repositories
	userRepo := repositories.NewUserRepo()
	classRepo := repositories.NewClassRepo()
	enrollRepo := repositories.NewEnrollmentRepo()

	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	classService := services.NewClassService(classRepo)
	enrollService := services.NewEnrollmentService(enrollRepo, classRepo)

	// Handlers
	authH := handlers.NewAuthHandler(authService)
	userH := handlers.NewUserHandler(userService)
	classH := handlers.NewClassHandler(classService, enrollService)
	enrollH := handlers.NewEnrollmentHandler(enrollService)
	dashboardH := handlers.NewDashboardHandler(enrollService)

	homeH := handlers.NewHomeHandler()
	aboutH := handlers.NewAboutHandler()

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", homeH.Home)
		//r.Get("/htmx/hello", homeH.HelloHtmx)
		r.Get("/register", authH.RegisterPage)
		r.Post("/register", authH.RegisterPost)
		r.Get("/login", authH.LoginPage)
		r.Post("/login", authH.LoginPost)
		r.Get("/about", aboutH.About)

	})

	// Member tier (Any Logged-in User to start with)
	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Use(middlewares.RequireRole(models.RoleMember, models.RoleTrainer, models.RoleAdmin))

		r.Post("/logout", authH.Logout)

		r.Get("/dashboard", dashboardH.Dashboard)
		r.Get("/profile", userH.ProfilePage)
		r.Post("/profile/update", userH.UpdateProfile)
		r.Post("/profile/update-password", userH.UpdatePassword)
		r.Get("/classes", classH.ClassesPage)

		r.Get("/dashboard", dashboardH.Dashboard)
		r.Post("/enrollments", enrollH.Enroll)
		r.Delete("/enrollments/{id}", enrollH.Cancel)

		r.Get("/classes/{id}", classH.ClassDetailsPage)
	})

	// Staff tier (Trainer or Admin)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Use(middlewares.RequireRole(models.RoleTrainer, models.RoleAdmin))

		// r.Get("/manage-programs", programH.List)
	})

	// Admin tier (Admin Only)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Use(middlewares.RequireRole(models.RoleAdmin))

		// r.Get("/admin/users", adminH.ManageUsers)
		// r.Get("/admin", adminHandler)

	})

	// Run server
	log.Println("✅ Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
