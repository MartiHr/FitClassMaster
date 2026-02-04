// Package auth provides utilities for session management and user authentication.
package auth

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

const (
	// SessionName is the name of the session cookie.
	SessionName = "app_session"
	// SessionUserIDKey is the key used to store the user's ID in the session.
	SessionUserIDKey = "user_id"
	// SessionUserEmailKey is the key used to store the user's email in the session.
	SessionUserEmailKey = "user_email"
	// SessionUsernameKey is the key used to store the derived username in the session.
	SessionUsernameKey = "username"
	// SessionUserRoleKey is the key used to store the user's role in the session.
	SessionUserRoleKey = "user_role"
)

// SaveUserSession populates the session with user data and saves it to the response writer.
func SaveUserSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	// Retrieve the session from the store.
	session, _ := config.Store.Get(r, SessionName)

	// Populate session values.
	SetUserID(session, user.ID)
	SetUserMeta(session, user.Email)
	SetUserRole(session, user.Role)

	// Persist the session.
	return session.Save(r, w)
}

// SetUserID sets the user ID in the session values.
func SetUserID(session *sessions.Session, userID uint) {
	session.Values[SessionUserIDKey] = userID
}

// SetUserMeta sets the user email and derives a username for the session.
func SetUserMeta(session *sessions.Session, email string) {
	session.Values[SessionUserEmailKey] = email
	username := email
	if at := strings.Index(email, "@"); at > 0 {
		username = email[:at]
	}
	session.Values[SessionUsernameKey] = username
}

// SetUserRole sets the user role in the session values.
func SetUserRole(session *sessions.Session, role models.Role) {
	session.Values[SessionUserRoleKey] = string(role)
}

// ClearUserSession removes user data from the session and expires the cookie.
func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := config.Store.Get(r, SessionName)

	delete(session.Values, SessionUserIDKey)
	delete(session.Values, SessionUserEmailKey)
	delete(session.Values, SessionUsernameKey)
	delete(session.Values, SessionUserRoleKey)

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// GetUserIDFromSession retrieves the user ID from the session.
func GetUserIDFromSession(r *http.Request) (uint, bool) {
	session, _ := config.Store.Get(r, SessionName)
	val := session.Values[SessionUserIDKey]

	switch v := val.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	default:
		return 0, false
	}
}

// GetUserEmailFromSession retrieves the user email from the session.
func GetUserEmailFromSession(r *http.Request) (string, bool) {
	session, _ := config.Store.Get(r, SessionName)
	email, ok := session.Values[SessionUserEmailKey].(string)
	return email, ok
}

// GetUsernameFromSession retrieves the username from the session.
func GetUsernameFromSession(r *http.Request) (string, bool) {
	session, _ := config.Store.Get(r, SessionName)
	username, ok := session.Values[SessionUsernameKey].(string)
	return username, ok
}

// GetUserRoleFromSession retrieves the user role from the session.
func GetUserRoleFromSession(r *http.Request) (models.Role, bool) {
	session, _ := config.Store.Get(r, SessionName)
	val, ok := session.Values[SessionUserRoleKey].(string)
	return models.Role(val), ok
}

// IsAuthenticated checks if the current request is associated with a logged-in user.
func IsAuthenticated(r *http.Request) bool {
	_, ok := GetUserIDFromSession(r)
	return ok
}
