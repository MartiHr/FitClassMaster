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
		log.Fatal("SESSION_KEY env var not set (32 or 64 bytes recommended)")
	}

	Store = sessions.NewCookieStore([]byte(key))
	Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,     // set true for HTTPS
		MaxAge:   86400 * 7, // one week
	}
}
