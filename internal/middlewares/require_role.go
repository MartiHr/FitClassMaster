package middlewares

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/models"
	"net/http"
)

func RequireRole(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get role from session
			userRole, ok := auth.GetUserRoleFromSession(r)

			// Check if the user's role is in the 'allowed' slice
			isAllowed := false
			for _, role := range allowedRoles {
				if ok && userRole == string(role) {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				// If HTMX, send a specific header for a clean redirect
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login?error=unauthorized")
					w.WriteHeader(http.StatusForbidden)
					return
				}
				http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
