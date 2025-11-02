package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var Tmpl *template.Template

func Init() {
	// Parse all templates together so they can reference each other
	Tmpl = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))
}

func Render(w http.ResponseWriter, name string, data any) {
	// Render the shared layout and let it include the requested page template by name.
	// We pass `{Data, Page}` so the layout can access shared fields like title via `.Data`.
	err := Tmpl.ExecuteTemplate(w, "layout", struct {
		Data any
		Page string
	}{
		Data: data,
		Page: name, // e.g. "home" or "register"
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
