package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strconv"
)

type ExerciseHandler struct {
	Service *services.ExerciseService
}

func NewExerciseHandler(s *services.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{Service: s}
}

// List handles GET /exercises
func (h *ExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}

	// Check role to determine if we show "Edit/Delete" buttons in the view
	userRole, _ := auth.GetUserRoleFromSession(r)
	canManage := (userRole == models.RoleTrainer || userRole == models.RoleAdmin)

	data := map[string]any{
		"Title":     "Exercise Library",
		"Exercises": exercises,
		"CanManage": canManage,
	}

	templates.SmartRender(w, r, "exercises", "", data)
}

// Create handles POST /exercises (Trainer/Admin only)
func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := h.Service.Create(
		r.FormValue("name"),
		r.FormValue("description"),
		r.FormValue("muscle_group"),
		r.FormValue("equipment"),
		r.FormValue("video_url"),
	)
	if err != nil {
		http.Error(w, "Failed to create exercise", http.StatusInternalServerError)
		return
	}

	// Redirect back to list
	http.Redirect(w, r, "/exercises", http.StatusSeeOther)
}

// Delete handles DELETE /exercises/{id} (Trainer/Admin only)
func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	err := h.Service.Delete(uint(id))
	if err != nil {
		http.Error(w, "Failed to delete exercise", http.StatusInternalServerError)
		return
	}

	// HTMX: Return empty 200 to remove element, or trigger a redirect
	w.WriteHeader(http.StatusOK)
}
