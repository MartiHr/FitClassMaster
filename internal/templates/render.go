package templates

import (
	"FitClassMaster/internal/auth"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// baseTmpl holds only the shared layout (and any shared partials if added later).
var baseTmpl *template.Template

// cache holds a fully built template set per page: layout + page file.
var cache map[string]*template.Template

// devMode enables per-request parsing (hot reload) when DEV_TEMPLATES=1
var devMode bool

// Init builds the template cache at startup: one set per page (layout + page).
func Init() {
	devMode = os.Getenv("DEV_TEMPLATES") == "1"

	// Parse the layout once
	baseTmpl = template.Must(template.ParseFiles(
		filepath.Join("internal", "templates", "layout.gohtml"),
	))

	// In dev mode, skip building the cache to allow hot reload
	if devMode {
		return
	}

	// Build the per-page cache
	cache = make(map[string]*template.Template)

	// Find all page files (exclude layout.gohtml)
	pattern := filepath.Join("internal", "templates", "*.gohtml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if filepath.Base(f) == "layout.gohtml" {
			continue
		}

		// Page name is file base without extension (e.g., home.gohtml -> "home")
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		// Clone base and parse the page into it so it gets that page’s `content` override
		cl, err := baseTmpl.Clone()
		if err != nil {
			panic(err)
		}
		if _, err := cl.ParseFiles(f); err != nil {
			panic(err)
		}

		cache[name] = cl
	}
}

// SmartRender automatically detects HTMX + can fall back to fragment
func SmartRender(w http.ResponseWriter, r *http.Request, page string, fragment string, data any) {
	// Ensure we have a map[string]any to enrich
	var m map[string]any
	switch v := data.(type) {
	case nil:
		m = map[string]any{}
	case map[string]any:
		m = v
	default:
		m = map[string]any{"Data": v}
	}

	// Inject auth context for templates
	m["IsAuthenticated"] = auth.IsAuthenticated(r)
	if username, ok := auth.GetUsernameFromSession(r); ok {
		m["Username"] = username
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX && fragment != "" {
		renderFragment(w, page, fragment, m)
		return
	}

	render(w, page, m)
}

// Render renders a full page using the cache in prod, or per-request parse in dev.
func render(w http.ResponseWriter, name string, data any) {
	if devMode {
		// Dev: clone + parse page each request for hot reload
		cl, err := baseTmpl.Clone()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pagePath := filepath.Join("internal", "templates", name+".gohtml")
		cl, err = cl.ParseFiles(pagePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cl.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prod: use cached set
	t, ok := cache[name]
	if !ok {
		http.Error(w, fmt.Sprintf("template not found: %s", name), http.StatusNotFound)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderFragment executes a named template (e.g., partial) from the cached set of a page
// in prod, or from a per-request parsed set in dev.
func renderFragment(w http.ResponseWriter, pageName, tmplName string, data any) {
	if devMode {
		cl, err := baseTmpl.Clone()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pagePath := filepath.Join("internal", "templates", pageName+".gohtml")
		cl, err = cl.ParseFiles(pagePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cl.ExecuteTemplate(w, tmplName, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	t, ok := cache[pageName]
	if !ok {
		http.Error(w, fmt.Sprintf("template not found: %s", pageName), http.StatusNotFound)
		return
	}
	if err := t.ExecuteTemplate(w, tmplName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
