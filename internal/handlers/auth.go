package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strings"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// RegisterPage (GET)
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":     "Register | FitClassMaster",
		"FirstName": "",
		"LastName":  "",
		"Email":     "",
	}

	templates.SmartRender(w, r, "register", "", data)
}

// RegisterPost (POST)
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
		templates.SmartRender(w, r, "register", "register_form", data)
		return
	}

	user := &models.User{
		FirstName: first,
		LastName:  last,
		Email:     email,
		Role:      models.RoleMember,
	}

	if err := h.AuthService.Register(user, password); err != nil {
		data["Error"] = err.Error()
		templates.SmartRender(w, r, "register", "register_form", data)
		return
	}

	// success
	data["Success"] = "Registration successful. You can now log in."
	// Clear the fields except email (optional)
	data["FirstName"] = ""
	data["LastName"] = ""
	data["Email"] = email

	templates.SmartRender(w, r, "register", "register_form", data)
}

// LoginPage (GET)
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.SmartRender(w, r, "login", "", map[string]any{"Title": "Login | FitClassMaster"})
}

// LoginPost (POST)
func (h *AuthHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid form"))
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	data := map[string]any{
		"Email": email,
	}

	if email == "" || password == "" {
		data["Error"] = "all fields are required"
		templates.SmartRender(w, r, "login", "", data)
		return
	}

	user, err := h.AuthService.Login(email, password)
	if err != nil {
		data["Error"] = "Invalid credentials"
		templates.SmartRender(w, r, "login", "", data)
		return
	}

	// Save user ID in session
	if err := auth.SaveUserSession(w, r, user); err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	// If this was an HTMX request, instruct htmx to perform a full redirect
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout (POST)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the user session using auth helpers
	_ = auth.ClearUserSession(w, r)
	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
