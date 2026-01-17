package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromSession(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.UserService.GetProfile(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title": "My Profile",
		"User":  user,
	}

	templates.SmartRender(w, r, "profile", "profile", data)
}

// UpdateProfile (POST)
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromSession(r)
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")

	err := h.UserService.UpdateProfile(userID, firstName, lastName)

	// Refresh user data through the service
	user, _ := h.UserService.GetProfile(userID)

	data := map[string]any{
		"User": user,
	}

	if err != nil {
		data["Error"] = "Could not update profile"
	} else {
		data["Success"] = "Profile updated successfully!"
	}

	// Uses SmartRender to only return the 'response_msg' fragment for HTMX
	templates.SmartRender(w, r, "profile", "response_msg", data)
}

// UpdatePassword (POST)
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromSession(r)
	current := r.FormValue("current_password")
	newPwd := r.FormValue("new_password")
	confirm := r.FormValue("confirm_password")

	data := map[string]any{}

	if newPwd != confirm {
		data["PassError"] = "New passwords do not match"
		templates.SmartRender(w, r, "profile", "password_msg", data)
		return
	}

	err := h.UserService.ChangePassword(userID, current, newPwd)
	if err != nil {
		data["PassError"] = err.Error()
	} else {
		data["PassSuccess"] = "Password changed successfully!"
	}

	// Uses SmartRender to only return the 'password_msg' fragment for HTMX
	templates.SmartRender(w, r, "profile", "password_msg", data)
}
