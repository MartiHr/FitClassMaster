package handlers

import (
	"FitClassMaster/internal/templates"
	"net/http"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Home | FitClassMaster",
	}

	templates.SmartRender(w, r, "home", "", data)
}
