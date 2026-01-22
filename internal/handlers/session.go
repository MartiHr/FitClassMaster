package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SessionHandler struct {
	Service *services.SessionService
}

func NewSessionHandler(s *services.SessionService) *SessionHandler {
	return &SessionHandler{Service: s}
}

// Start handles POST /sessions/start
// It takes a plan_id, creates a session, and redirects to the perform page
func (h *SessionHandler) Start(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	planIDStr := r.FormValue("plan_id")
	planID, _ := strconv.ParseUint(planIDStr, 10, 32)
	userID, _ := auth.GetUserIDFromSession(r)

	// Create the session
	session, err := h.Service.StartSession(userID, uint(planID))
	if err != nil {
		http.Error(w, "Failed to start session", http.StatusInternalServerError)
		return
	}

	// Redirect to the "Player" view
	redirectURL := "/sessions/" + strconv.FormatUint(uint64(session.ID), 10) + "/perform"
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// PerformPage handles GET /sessions/{id}/perform
// This is the main UI where the user actually does the workout
func (h *SessionHandler) PerformPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	session, err := h.Service.GetDetails(uint(id))
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title":   "Perform Workout",
		"Session": session,
	}

	templates.SmartRender(w, r, "session_perform", "", data)
}
