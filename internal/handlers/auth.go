package handlers

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"encoding/gob"
	"net/http"
	"strings"
)

func init() {
	gob.Register(uint(0)) // so session can store user ID
}

type AuthHandler struct {
	AuthService *services.AuthService
	Repo        *repositories.UserRepo
}

func NewAuthHandler() *AuthHandler {
	repo := repositories.NewUserRepo()
	authService := services.NewAuthService(repo)

	return &AuthHandler{AuthService: authService, Repo: repo}
}

// RegisterPage (GET)
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":     "Register | FitClassMaster",
		"FirstName": "",
		"LastName":  "",
		"Email":     "",
	}
	templates.Render(w, "register", data)
}

// Register (POST)
func (h *AuthHandler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid form"))
		return
	}

	first := strings.TrimSpace(r.FormValue("first_name"))
	last := strings.TrimSpace(r.FormValue("last_name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	data := map[string]any{
		"Title":     "Register | FitClassMaster",
		"FirstName": first,
		"LastName":  last,
		"Email":     email,
	}

	if first == "" || last == "" || email == "" || password == "" {
		data["Error"] = "all fields are required"
		h.renderRegisterResponse(w, r, data)
		return
	}

	user := &models.User{FirstName: first, LastName: last, Email: email}
	if err := h.AuthService.Register(user, password); err != nil {
		data["Error"] = err.Error()
		h.renderRegisterResponse(w, r, data)
		return
	}

	// success
	data["Success"] = "Registration successful. You can now log in."
	// Clear the fields except email (optional)
	data["FirstName"] = ""
	data["LastName"] = ""
	data["Email"] = email

	h.renderRegisterResponse(w, r, data)
}

func (h *AuthHandler) renderRegisterResponse(w http.ResponseWriter, r *http.Request, data map[string]any) {
	// If it's an HTMX request, return only the form fragment to be swapped
	if r.Header.Get("HX-Request") == "true" {
		// Render just the register_form template from the register page file
		templates.RenderFragment(w, "register", "register_form", data)
		return
	}
	// Otherwise render the full page
	templates.Render(w, "register", data)
}
