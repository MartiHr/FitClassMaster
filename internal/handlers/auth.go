package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strings"
)

// AuthHandler handles authentication-related requests such as login, registration, and logout.
type AuthHandler struct {
	AuthService *services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// RegisterPage renders the registration page.
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":     "Register | FitClassMaster",
		"FirstName": "",
		"LastName":  "",
		"Email":     "",
	}

	templates.SmartRender(w, r, "register", "", data)
}

// RegisterPost handles the submission of the registration form.
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

	// Basic validation.
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

	// Attempt to register the new user.
	if err := h.AuthService.Register(user, password); err != nil {
		data["Error"] = err.Error()
		templates.SmartRender(w, r, "register", "register_form", data)
		return
	}

	// Display success message and clear form fields.
	data["Success"] = "Registration successful. You can now log in."
	data["FirstName"] = ""
	data["LastName"] = ""
	data["Email"] = email

	templates.SmartRender(w, r, "register", "register_form", data)
}

// LoginPage renders the login page.
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.SmartRender(w, r, "login", "", map[string]any{"Title": "Login | FitClassMaster"})
}

// LoginPost handles user login requests.
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

	// Validate user credentials.
	user, err := h.AuthService.Login(email, password)
	if err != nil {
		data["Error"] = "Invalid credentials"
		templates.SmartRender(w, r, "login", "", data)
		return
	}

	// Create a new session for the authenticated user.
	if err := auth.SaveUserSession(w, r, user); err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	// Support HTMX client-side redirection.
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout terminates the user's session.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = auth.ClearUserSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
