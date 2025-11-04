package auth

import (
	"FitClassMaster/internal/config"
	"net/http"
)

const SessionName = "app_session"
const SessionUserIDKey = "user_id"

// SaveUserSession stores the authenticated user's ID (uint) into the session.
func SaveUserSession(w http.ResponseWriter, r *http.Request, userID uint) error {
	session, _ := config.Store.Get(r, SessionName)
	session.Values[SessionUserIDKey] = userID
	return session.Save(r, w)
}

// ClearUserSession removes the user ID from the session and persists the change.
func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := config.Store.Get(r, SessionName)
	delete(session.Values, SessionUserIDKey)
	return session.Save(r, w)
}

// GetUserIDFromSession retrieves the user ID from session as uint.
func GetUserIDFromSession(r *http.Request) (uint, bool) {
	session, _ := config.Store.Get(r, SessionName)
	userID, ok := session.Values[SessionUserIDKey].(uint)
	return userID, ok
}
