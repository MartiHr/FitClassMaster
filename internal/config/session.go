package config

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
)

// Store is the global session store (cookie-based).
var Store *sessions.CookieStore

// InitSessionStore initializes the Gorilla session store.
// It uses the SESSION_KEY environment variable for encryption.
func InitSessionStore() {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		// Provide a development fallback to avoid startup failure.
		// NOTE: In production, always set a secure SESSION_KEY environment variable.
		log.Println("[WARN] SESSION_KEY not set; using insecure development key")
		key = "dev-insecure-session-key-please-set-ENV"
	}

	Store = sessions.NewCookieStore([]byte(key))
	Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,     // Set to true if using HTTPS.
		MaxAge:   86400 * 7, // Session expires after one week.
	}
}
