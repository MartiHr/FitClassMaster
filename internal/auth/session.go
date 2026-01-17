package auth

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

const (
	SessionName         = "app_session"
	SessionUserIDKey    = "user_id"
	SessionUserEmailKey = "user_email"
	SessionUsernameKey  = "username"
	SessionUserRoleKey  = "user_role"
)

func SaveUserSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	// Open session once
	session, _ := config.Store.Get(r, SessionName)

	// Use helpers to populate the session object (No encryption/saving happens yet)
	SetUserID(session, user.ID)
	SetUserMeta(session, user.Email)
	SetUserRole(session, user.Role)

	// Encrypt and Write Header ONCE
	return session.Save(r, w)
}

func SetUserID(session *sessions.Session, userID uint) {
	session.Values[SessionUserIDKey] = userID
}

func SetUserMeta(session *sessions.Session, email string) {
	session.Values[SessionUserEmailKey] = email
	username := email
	if at := strings.Index(email, "@"); at > 0 {
		username = email[:at]
	}
	session.Values[SessionUsernameKey] = username
}

func SetUserRole(session *sessions.Session, role models.Role) {
	session.Values[SessionUserRoleKey] = string(role)
}

func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := config.Store.Get(r, SessionName)

	delete(session.Values, SessionUserIDKey)
	delete(session.Values, SessionUserEmailKey)
	delete(session.Values, SessionUsernameKey)
	delete(session.Values, SessionUserRoleKey)

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

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

func GetUserEmailFromSession(r *http.Request) (string, bool) {
	session, _ := config.Store.Get(r, SessionName)
	email, ok := session.Values[SessionUserEmailKey].(string)
	return email, ok
}

func GetUsernameFromSession(r *http.Request) (string, bool) {
	session, _ := config.Store.Get(r, SessionName)
	username, ok := session.Values[SessionUsernameKey].(string)
	return username, ok
}

func GetUserRoleFromSession(r *http.Request) (models.Role, bool) {
	session, _ := config.Store.Get(r, SessionName)
	val, ok := session.Values[SessionUserRoleKey].(string)
	return models.Role(val), ok
}

func IsAuthenticated(r *http.Request) bool {
	_, ok := GetUserIDFromSession(r)
	return ok
}
