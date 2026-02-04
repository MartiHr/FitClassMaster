// Package handlers contains the HTTP controllers for the application.
package handlers

import (
	"FitClassMaster/internal/templates"
	"net/http"
)

// HomeHandler handles requests to the public home page.
type HomeHandler struct{}

// NewHomeHandler creates a new instance of HomeHandler.
func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// Home renders the home page.
func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Home | FitClassMaster",
	}

	templates.SmartRender(w, r, "home", "", data)
}
