package auth

import (
	"FitClassMaster/internal/config"
	"net/http"
)

const SessionName = "app_session"
const SessionUserIDKey = "user_id"

func saveUserSession(w http.ResponseWriter, r *http.Request, userID string) error {
	session, _ := config.Store.Get(r, SessionName)
	session.Values[SessionUserIDKey] = userID
	return session.Save(r, w)
}

func clearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := config.Store.Get(r, SessionName)
	delete(session.Values, SessionUserIDKey)
	return session.Save(r, w)
}

func GetUserIDFromSession(r *http.Request) (string, bool) {
	session, _ := config.Store.Get(r, SessionName)
	userID, exists := session.Values[SessionUserIDKey].(string)
	return userID, exists
}
