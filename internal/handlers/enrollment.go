package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"net/http"
	"strconv"
)

type EnrollmentHandler struct {
	Service *services.EnrollmentService
}

func NewEnrollmentHandler(s *services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{Service: s}
}

func (h *EnrollmentHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	classIDStr := r.FormValue("class_id")
	classID, _ := strconv.ParseUint(classIDStr, 10, 32)

	// Get UserID from session
	userID, _ := auth.GetUserIDFromSession(r)

	err := h.Service.EnrollUser(userID, uint(classID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
        <button class="btn-enroll" style="background-color: #28a745;" disabled>
            <i class="fas fa-check"></i> Enrolled
        </button>
    `))
}
