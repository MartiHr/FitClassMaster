package main

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	fileServer := http.FileServer(http.Dir("internal/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Routes
	r.Get("/", homeHandler)
	r.Get("/htmx/hello", helloHtmx)

	log.Println("✅ Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "home.gohtml", map[string]any{
		"Title": "FitClassMaster — Home",
	})
}

func helloHtmx(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from HTMX! ✅"))
}

func render(w http.ResponseWriter, tmpl string, data any) {
	err := templates.ExecuteTemplate(w, "layout.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
