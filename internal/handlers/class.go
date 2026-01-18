package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
)

type ClassHandler struct {
	ClassService *services.ClassService
}

func NewClassHandler(s *services.ClassService) *ClassHandler {
	return &ClassHandler{ClassService: s}
}

//
//func (h *ClassHandler) ClassesPage(w http.ResponseWriter, r *http.Request) {
//	classes, _ := h.ClassService.GetAvailableClasses()
//
//	data := map[string]any{
//		"Title":   "Fitness Classes | FitClassMaster",
//		"Classes": classes,
//	}
//	templates.SmartRender(w, r, "classes", "", data)
//}

func (h *ClassHandler) ClassesPage(w http.ResponseWriter, r *http.Request) {
	// Get the current user from session
	userID, _ := auth.GetUserIDFromSession(r)

	// Fetch all available classes
	classesWithStatus, err := h.ClassService.GetClassesForUser(userID)
	if err != nil {
		http.Error(w, "Failed to load classes", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":   "Fitness Classes | FitClassMaster",
		"Classes": classesWithStatus,
	}

	templates.SmartRender(w, r, "classes", "", data)
}
