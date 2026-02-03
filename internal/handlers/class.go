package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

// CreatePage renders the form
func (h *ClassHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Schedule a Class",
	}

	templates.SmartRender(w, r, "class_create", "", data)
}

// CreatePost processes the form
func (h *ClassHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	trainerID, _ := auth.GetUserIDFromSession(r)

	// Basic Fields
	name := r.FormValue("name")
	description := r.FormValue("description")
	difficulty := r.FormValue("difficulty") // NEW
	capacity, _ := strconv.Atoi(r.FormValue("capacity"))

	// Time Calculation
	layout := "2006-01-02T15:04"
	start, err1 := time.Parse(layout, r.FormValue("start_time"))
	end, err2 := time.Parse(layout, r.FormValue("end_time"))

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Calculate Duration in Minutes
	duration := int(end.Sub(start).Minutes())
	if duration <= 0 {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// Call Service
	err = h.ClassService.CreateClass(name, description, difficulty, trainerID, start, duration, capacity)
	if err != nil {
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/classes", http.StatusSeeOther)
}

// Cancel handles deletion
func (h *ClassHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.ClassService.CancelClass(uint(id)); err != nil {
		http.Error(w, "Failed to cancel class", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/classes", http.StatusSeeOther)
}

// EditPage renders the form with existing class data
func (h *ClassHandler) EditPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	class, err := h.ClassService.GetFullDetails(uint(id))
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title": "Edit Class",
		"Class": class,
		// We need to format the time specifically for the HTML <input type="datetime-local">
		"FormattedStart": class.StartTime.Format("2006-01-02T15:04"),
		"FormattedEnd":   class.StartTime.Add(time.Duration(class.Duration) * time.Minute).Format("2006-01-02T15:04"),
	}
	templates.SmartRender(w, r, "class_edit", "", data)
}

// UpdatePost processes the changes
func (h *ClassHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Parse Fields
	name := r.FormValue("name")
	desc := r.FormValue("description")
	diff := r.FormValue("difficulty")
	cap, _ := strconv.Atoi(r.FormValue("capacity"))

	// Time Calculation
	layout := "2006-01-02T15:04"
	start, _ := time.Parse(layout, r.FormValue("start_time"))
	end, _ := time.Parse(layout, r.FormValue("end_time"))
	duration := int(end.Sub(start).Minutes())

	// Call Service
	err = h.ClassService.UpdateClass(uint(id), name, desc, diff, start, duration, cap)
	if err != nil {
		http.Error(w, "Failed to update class", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/classes", http.StatusSeeOther)
}
