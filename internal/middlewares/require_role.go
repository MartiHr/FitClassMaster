package middlewares

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/models"
	"net/http"
)

// RequireRole is a middleware factory that returns a middleware ensuring
// the authenticated user has one of the allowed roles.
func RequireRole(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve user role from session.
			userRole, ok := auth.GetUserRoleFromSession(r)

			// Check if the user's role matches any of the allowed roles.
			isAllowed := false
			for _, role := range allowedRoles {
				if ok && userRole == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				// Handle HTMX requests for clean client-side redirection.
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login?error=unauthorized")
					w.WriteHeader(http.StatusForbidden)
					return
				}
				// Standard HTTP redirect for regular requests.
				http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
