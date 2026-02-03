package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
)

type DashboardHandler struct {
	EnrollmentService *services.EnrollmentService
	SessionService    *services.SessionService
}

func NewDashboardHandler(s *services.EnrollmentService, ss *services.SessionService) *DashboardHandler {
	return &DashboardHandler{
		EnrollmentService: s,
		SessionService:    ss,
	}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromSession(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username, _ := auth.GetUsernameFromSession(r)
	role, _ := auth.GetUserRoleFromSession(r)

	data := map[string]any{
		"Title":    "Dashboard | FitClassMaster",
		"Username": username,
		"Role":     role,
	}

	// Check for interrupted session to show the Resume Banner
	activeSession, err := h.SessionService.GetActiveSession(userID)
	if err == nil && activeSession.ID != 0 {
		data["CurrentSession"] = activeSession
	}

	// Fetch Classes
	myClasses, err := h.EnrollmentService.GetMySchedule(userID)
	if err == nil {
		data["Classes"] = myClasses
	}

	// Fetch Live Sessions
	if role == "trainer" || role == "admin" {
		// Live
		activeSessions, _ := h.SessionService.ListActiveSessions(role, userID)
		data["ActiveSessions"] = activeSessions // This is why "Live" works

		// Client History
		clientHistory, _ := h.SessionService.GetTrainerHistory(userID)
		data["ClientHistory"] = clientHistory
	}

	// Fetch Workout History (For Members)
	history, err := h.SessionService.GetHistory(userID)
	if err == nil {
		data["History"] = history
	}

	templates.SmartRender(w, r, "dashboard", "", data)
}
