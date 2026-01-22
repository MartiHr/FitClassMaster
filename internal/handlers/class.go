package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strconv"
)

type ClassHandler struct {
	ClassService      *services.ClassService
	EnrollmentService *services.EnrollmentService
}

func NewClassHandler(s *services.ClassService, es *services.EnrollmentService) *ClassHandler {
	return &ClassHandler{
		ClassService:      s,
		EnrollmentService: es,
	}
}

func (h *ClassHandler) ClassesPage(w http.ResponseWriter, r *http.Request) {
	// Get the current user from session
	userID, _ := auth.GetUserIDFromSession(r)

	// Fetch all available classes
	classesWithStatus, err := h.ClassService.GetClassesForUser(userID)
	if err != nil {
		http.Error(w, "Classes not found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title":   "Fitness Classes | FitClassMaster",
		"Classes": classesWithStatus,
	}

	templates.SmartRender(w, r, "classes", "", data)
}

func (h *ClassHandler) ClassDetailsPage(w http.ResponseWriter, r *http.Request) {
	// Extract ID
	idStr := r.PathValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	classID := uint(id)

	// Get Class Roster Data (Using ClassService)
	class, err := h.ClassService.GetFullDetails(classID)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	// Get Current User Info
	userID, _ := auth.GetUserIDFromSession(r)
	username, _ := auth.GetUsernameFromSession(r)

	// Check if current user is enrolled
	isEnrolled, _ := h.EnrollmentService.IsUserEnrolled(userID, classID)

	data := map[string]any{
		"Title":      class.Name + " | Details",
		"Class":      class,
		"Username":   username,
		"IsEnrolled": isEnrolled,        // Boolean for the button state
		"Roster":     class.Enrollments, // List for the "Assigned Members" view
	}

	templates.SmartRender(w, r, "class_details", "", data)
}
