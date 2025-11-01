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
	err := Tmpl.ExecuteTemplate(w, "layout", struct {
		Data any
		Page string
	}{
		Data: data,
		Page: name, // tells layout which content to show
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
