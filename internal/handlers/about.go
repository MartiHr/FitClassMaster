package handlers

import (
	"FitClassMaster/internal/templates"
	"net/http"
)

type AboutHandler struct{}

func NewAboutHandler() *AboutHandler {
	return &AboutHandler{}
}

func (h *AboutHandler) About(w http.ResponseWriter, r *http.Request) {
	templates.SmartRender(w, r, "about", "", map[string]any{"Title": "About Us"})
}
