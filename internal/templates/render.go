package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var Tmpl *template.Template

func Init() {
	Tmpl = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))
}

func Render(w http.ResponseWriter, name string, data any) {
	err := Tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
