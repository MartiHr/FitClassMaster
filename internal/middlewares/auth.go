// Package middlewares contains HTTP middleware functions for the application.
package middlewares

import (
	"FitClassMaster/internal/auth"
	"net/http"
)

// RequireAuth is a middleware that ensures the user is authenticated.
// If not authenticated, it redirects the user to the login page.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for user ID in the session.
		_, ok := auth.GetUserIDFromSession(r)

		if !ok {
			// Redirect to login if session is missing or invalid.
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
