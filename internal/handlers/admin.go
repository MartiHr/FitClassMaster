package handlers

import (
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	UserService *services.UserService
}

func NewAdminHandler(us *services.UserService) *AdminHandler {
	return &AdminHandler{UserService: us}
}

// ManageUsers renders the user table
func (h *AdminHandler) ManageUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.ListAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title": "User Management",
		"Users": users,
	}

	templates.SmartRender(w, r, "admin_users", "", data)
}

// ToggleRole handles promoting/demoting
func (h *AdminHandler) ToggleRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	action := r.FormValue("action") // "promote" or "demote"

	var err error
	if action == "promote" {
		err = h.UserService.PromoteUser(uint(id))
	} else {
		err = h.UserService.DemoteUser(uint(id))
	}

	if err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// DeleteUser handles permanent removal
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.UserService.DeleteUser(uint(id)); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
