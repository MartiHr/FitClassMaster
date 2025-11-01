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
	templates.Render(w, "home.gohtml", map[string]any{
		"Title": "FitClassMaster — Home",
	})
}

func (h *HomeHandler) HelloHtmx(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from HTMX! ✅"))
}
