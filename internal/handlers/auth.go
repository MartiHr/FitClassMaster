package handlers

import (
	"FitClassMaster/internal/repositories"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"encoding/gob"
	"html/template"
	"net/http"
)

var authTemplates *template.Template

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
	templates.Render(w, "register.gohtml", map[string]any{
		"Title": "Register | FitClassMaster",
	})
}

//
