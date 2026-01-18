package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
)

type DashboardHandler struct {
	EnrollmentService *services.EnrollmentService
}

func NewDashboardHandler(s *services.EnrollmentService) *DashboardHandler {
	return &DashboardHandler{EnrollmentService: s}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromSession(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username, _ := auth.GetUsernameFromSession(r)

	// Fetch only classes this user is enrolled in
	myClasses, err := h.EnrollmentService.GetMySchedule(userID)
	if err != nil {
		// Handle database or service errors
		http.Error(w, "Could not load your schedule", http.StatusInternalServerError)
		return
	}

	// Prepare data for the dashboard.gohtml template
	data := map[string]any{
		"Title":    "My Schedule | FitClassMaster",
		"Username": username,  // Needed for "Welcome, {{.Username}}" in template
		"Classes":  myClasses, // The slice of models.Class
	}

	templates.SmartRender(w, r, "dashboard", "", data)
}
