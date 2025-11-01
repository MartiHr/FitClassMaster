package templates

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var Tmpl *template.Template

func Init() {
	Tmpl = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))
}

func Render(w http.ResponseWriter, name string, data any) {
	log.Println("Render:", name)

	err := Tmpl.ExecuteTemplate(w, "layout.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
