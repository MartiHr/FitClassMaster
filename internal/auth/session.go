package auth

import (
	"FitClassMaster/internal/config"
	"net/http"
	"strings"
)

const SessionName = "app_session"
const SessionUserIDKey = "user_id"
const SessionUserEmailKey = "user_email"
const SessionUsernameKey = "username"

// SaveUserSession stores the authenticated user's ID (uint) into the session.
func SaveUserSession(w http.ResponseWriter, r *http.Request, userID uint) error {
	session, _ := config.Store.Get(r, SessionName)
	session.Values[SessionUserIDKey] = userID
	return session.Save(r, w)
}

// SaveUserMeta stores email and derived username into the session.
func SaveUserMeta(w http.ResponseWriter, r *http.Request, email string) error {
	session, _ := config.Store.Get(r, SessionName)
	session.Values[SessionUserEmailKey] = email
	username := email
	if at := strings.Index(email, "@"); at > 0 {
		username = email[:at]
	}
	session.Values[SessionUsernameKey] = username
	return session.Save(r, w)
}

// ClearUserSession removes all user-related info from the session and persists the change.
func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := config.Store.Get(r, SessionName)
	delete(session.Values, SessionUserIDKey)
	delete(session.Values, SessionUserEmailKey)
	delete(session.Values, SessionUsernameKey)
	return session.Save(r, w)
}

// GetUserIDFromSession retrieves the user ID from session as uint.
func GetUserIDFromSession(r *http.Request) (uint, bool) {
	session, _ := config.Store.Get(r, SessionName)
	userID, ok := session.Values[SessionUserIDKey].(uint)
	return userID, ok
}

// GetUserEmailFromSession retrieves the user's email from the session.
func GetUserEmailFromSession(r *http.Request) (string, bool) {
	session, _ := config.Store.Get(r, SessionName)
	email, ok := session.Values[SessionUserEmailKey].(string)
	return email, ok
}

// GetUsernameFromSession retrieves the username (email prefix) from the session.
func GetUsernameFromSession(r *http.Request) (string, bool) {
	session, _ := config.Store.Get(r, SessionName)
	username, ok := session.Values[SessionUsernameKey].(string)
	return username, ok
}

// IsAuthenticated returns true if a user ID is present in the session.
func IsAuthenticated(r *http.Request) bool {
	_, ok := GetUserIDFromSession(r)
	return ok
}
