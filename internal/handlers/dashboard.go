package handlers

import (
	"FitClassMaster/internal/templates"
	"net/http"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Dashboard | FitClassMaster",
	}

	templates.SmartRender(w, r, "dashboard", "", data)
}
