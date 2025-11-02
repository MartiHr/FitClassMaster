package config

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitSessionStore() {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		// Provide a development fallback to avoid startup failure
		// NOTE: replace with a secure env var in production
		log.Println("[WARN] SESSION_KEY not set; using insecure development key")
		key = "dev-insecure-session-key-please-set-ENV"
	}

	Store = sessions.NewCookieStore([]byte(key))
	Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,     // set true for HTTPS
		MaxAge:   86400 * 7, // one week
	}
}
