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
	r.Use(middleware.Logger)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))

	// Routes
	r.Get("/", homeHandler)
	r.Get("/htmx/hello", helloHtmx)

	log.Println("✅ Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
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
