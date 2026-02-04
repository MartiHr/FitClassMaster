// Package main is the entry point for the FitClassMaster application.
// It initializes the database, sessions, templates, and starts the HTTP server.
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
	"os"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.Enrollment{},
		&models.Exercise{},
		&models.WorkoutPlan{},
		&models.WorkoutExercise{},
		&models.WorkoutSession{},
		&models.SessionLog{},
		&models.Conversation{},
		&models.Message{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Create a new chi router.
	r := chi.NewRouter()

	// Register global middlewares.
	r.Use(middleware.Logger)    // Log each HTTP request.
	r.Use(middleware.Recoverer) // Recover from panics without crashing the server.

	// Serve static files (CSS, JS, images) from the internal/static directory.
	fileServer := http.FileServer(http.Dir("internal/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Initialize Repositories (Data Access Layer).
	userRepo := repositories.NewUserRepo()
	classRepo := repositories.NewClassRepo()
	enrollRepo := repositories.NewEnrollmentRepo()
	exerciseRepo := repositories.NewExerciseRepo()
	workoutRepo := repositories.NewWorkoutRepo()
	sessionRepo := repositories.NewSessionRepo()
	messageRepo := repositories.NewMessageRepo()

	// Initialize Services (Business Logic Layer).
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	classService := services.NewClassService(classRepo)
	enrollService := services.NewEnrollmentService(enrollRepo, classRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	workoutService := services.NewWorkoutService(workoutRepo)
	sessionService := services.NewSessionService(sessionRepo)
	messageService := services.NewMessageService(messageRepo)

	// Initialize and Run the WebSocket Hubs for real-time features.
	hub := websockets.NewHub(sessionService)
	go hub.Run() // Hub for workout sessions.

	chatHub := websockets.NewChatHub()
	go chatHub.Run() // Hub for real-time messaging.

	// Initialize Handlers (Controller Layer).
	homeH := handlers.NewHomeHandler()
	aboutH := handlers.NewAboutHandler()
	authH := handlers.NewAuthHandler(authService)
	userH := handlers.NewUserHandler(userService)
	enrollH := handlers.NewEnrollmentHandler(enrollService)
	sessionH := handlers.NewSessionHandler(sessionService)
	exerciseH := handlers.NewExerciseHandler(exerciseService)
	adminH := handlers.NewAdminHandler(userService)
	msgH := handlers.NewMessageHandler(messageService, userService, chatHub)
	wsH := handlers.NewWSHandler(hub)
	workoutH := handlers.NewWorkoutHandler(workoutService, exerciseService)
	classH := handlers.NewClassHandler(classService, enrollService)
	dashboardH := handlers.NewDashboardHandler(enrollService, sessionService)

	// Define Public Routes.
	r.Group(func(r chi.Router) {
		r.Get("/", homeH.Home)
		r.Get("/register", authH.RegisterPage)
		r.Post("/register", authH.RegisterPost)
		r.Get("/login", authH.LoginPage)
		r.Post("/login", authH.LoginPost)
		r.Get("/about", aboutH.About)
	})

	// Define Protected Routes (Require Authentication).
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Use(middlewares.RequireRole(models.RoleMember, models.RoleTrainer, models.RoleAdmin))

		r.Post("/logout", authH.Logout)
		r.Get("/dashboard", dashboardH.Dashboard)
		r.Get("/profile", userH.ProfilePage)
		r.Post("/profile/update", userH.UpdateProfile)
		r.Post("/profile/update-password", userH.UpdatePassword)

		r.Get("/classes", classH.ClassesPage)
		r.Get("/classes/{id}", classH.ClassDetailsPage)

		r.Post("/enrollments", enrollH.Enroll)
		r.Delete("/enrollments/{id}", enrollH.Cancel)

		r.Get("/exercises", exerciseH.List)

		r.Get("/workout-plans", workoutH.List)
		r.Get("/workout-plans/{id}", workoutH.DetailsPage)

		r.Post("/sessions/start", sessionH.Start)
		r.Get("/sessions/{id}/perform", sessionH.PerformPage)

		r.Get("/ws/session/{id}", wsH.HandleSessionConnection)
		r.Post("/sessions/{id}/finish", sessionH.Finish)

		r.Get("/messages", msgH.InboxPage)
		r.Get("/messages/{id}", msgH.ThreadPage)
		r.Post("/messages/send", msgH.SendPost)
		r.Get("/messages/start/{userID}", msgH.StartChat)
		r.Get("/ws/chat", msgH.ServeWS)
	})

	// Define Staff Routes (Trainer or Admin Role).
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Use(middlewares.RequireRole(models.RoleTrainer, models.RoleAdmin))

		r.Post("/exercises", exerciseH.Create)
		r.Delete("/exercises/{id}", exerciseH.Delete)

		r.Get("/workout-plans/create", workoutH.CreatePage)
		r.Get("/workout-plans/add-row", workoutH.AddExerciseRow)
		r.Post("/workout-plans/create", workoutH.CreatePost)
		r.Get("/workout-plans/{id}/edit", workoutH.EditPage)
		r.Post("/workout-plans/{id}/edit", workoutH.UpdatePost)
		r.Post("/workout-plans/{id}/delete", workoutH.DeletePost)

		r.Get("/classes/create", classH.CreatePage)
		r.Post("/classes/create", classH.CreatePost)
		r.Get("/classes/{id}/edit", classH.EditPage)
		r.Post("/classes/{id}/edit", classH.UpdatePost)
		r.Post("/classes/{id}/cancel", classH.Cancel)
	})

	// Define Admin Only Routes.
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Use(middlewares.RequireRole(models.RoleAdmin))

		r.Get("/admin/users", adminH.ManageUsers)
		r.Post("/admin/users/{id}/role", adminH.ToggleRole)
		r.Post("/admin/users/{id}/delete", adminH.DeleteUser)
	})

	// Get the PORT from the environment (Render sets this automatically).
	// If it is empty (running locally), default to "8080".
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✅ Server running on port %s", port)

	// Listen on the dynamic port
	sErr := http.ListenAndServe(":"+port, r)
	if sErr != nil {
		log.Fatalf("Server failed: %v", sErr)
	}
}
