package handlers

import (
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

func (h *ClassHandler) ClassesPage(w http.ResponseWriter, r *http.Request) {
	classes, _ := h.ClassService.GetAvailableClasses()

	data := map[string]any{
		"Title":   "Fitness Classes | FitClassMaster",
		"Classes": classes,
	}
	templates.SmartRender(w, r, "classes", "", data)
}
