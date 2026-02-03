package main

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/handlers"
	"FitClassMaster/internal/middlewares"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"FitClassMaster/internal/websockets"

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
		&models.Enrollment{},
		&models.Exercise{},
		&models.WorkoutPlan{},
		&models.WorkoutExercise{},
		&models.WorkoutSession{},
		&models.SessionLog{}); err != nil {
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
	exerciseRepo := repositories.NewExerciseRepo()
	workoutRepo := repositories.NewWorkoutRepo()
	sessionRepo := repositories.NewSessionRepo()

	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	classService := services.NewClassService(classRepo)
	enrollService := services.NewEnrollmentService(enrollRepo, classRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	workoutService := services.NewWorkoutService(workoutRepo)
	sessionService := services.NewSessionService(sessionRepo)

	// Initialize and Run the Hub
	hub := websockets.NewHub(sessionService)
	go hub.Run() // Run the hub in a separate goroutine

	// Handlers
	authH := handlers.NewAuthHandler(authService)
	userH := handlers.NewUserHandler(userService)
	classH := handlers.NewClassHandler(classService, enrollService)
	enrollH := handlers.NewEnrollmentHandler(enrollService)
	dashboardH := handlers.NewDashboardHandler(enrollService, sessionService)
	exerciseH := handlers.NewExerciseHandler(exerciseService)
	workoutH := handlers.NewWorkoutHandler(workoutService, exerciseService)
	sessionH := handlers.NewSessionHandler(sessionService)

	homeH := handlers.NewHomeHandler()
	aboutH := handlers.NewAboutHandler()

	wsHandler := handlers.NewWSHandler(hub)

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

		r.Get("/exercises", exerciseH.List)

		r.Get("/workout-plans", workoutH.List)
		r.Get("/workout-plans/{id}", workoutH.DetailsPage)

		r.Post("/sessions/start", sessionH.Start)
		r.Get("/sessions/{id}/perform", sessionH.PerformPage)

		r.Get("/ws/session/{id}", wsHandler.HandleSessionConnection)
		r.Post("/sessions/{id}/finish", sessionH.Finish)
	})

	// Staff tier (Trainer or Admin)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Use(middlewares.RequireRole(models.RoleTrainer, models.RoleAdmin))

		r.Post("/exercises", exerciseH.Create)
		r.Delete("/exercises/{id}", exerciseH.Delete)

		// Page to Create Plan
		r.Get("/workout-plans/create", workoutH.CreatePage)
		// HTMX Endpoint to get a new row
		r.Get("/workout-plans/add-row", workoutH.AddExerciseRow)
		// Form Submission
		r.Post("/workout-plans/create", workoutH.CreatePost)
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
