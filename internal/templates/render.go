// Package templates handles HTML template parsing and rendering.
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

// baseTmpl holds the shared layout and any shared partials.
var baseTmpl *template.Template

// cache holds fully built template sets per page: layout + specific page file.
var cache map[string]*template.Template

// devMode enables per-request parsing (hot reload) when DEV_TEMPLATES=1 is set.
var devMode bool

// Init initializes the template engine, building the cache if not in development mode.
func Init() {
	devMode = os.Getenv("DEV_TEMPLATES") == "1"

	// Parse the shared layout.
	baseTmpl = template.Must(template.ParseFiles(
		filepath.Join("internal", "templates", "layout.gohtml"),
	))

	// In development mode, we skip building the cache to allow for hot reloading of templates.
	if devMode {
		return
	}

	// Initialize the template cache.
	cache = make(map[string]*template.Template)

	// Identify all template files in the directory, excluding the layout.
	pattern := filepath.Join("internal", "templates", "*.gohtml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if filepath.Base(f) == "layout.gohtml" {
			continue
		}

		// The page name is derived from the filename without extension.
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		// Clone the base layout and parse the specific page template into it.
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

// SmartRender handles both full-page and HTMX fragment rendering, injecting session context automatically.
func SmartRender(w http.ResponseWriter, r *http.Request, page string, fragment string, data any) {
	// Ensure data is available as a map for injection.
	var m map[string]any
	switch v := data.(type) {
	case nil:
		m = map[string]any{}
	case map[string]any:
		m = v
	default:
		m = map[string]any{"Data": v}
	}

	// Inject authentication and user context for use in templates.
	m["IsAuthenticated"] = auth.IsAuthenticated(r)
	if username, ok := auth.GetUsernameFromSession(r); ok {
		m["Username"] = username
	}

	if role, ok := auth.GetUserRoleFromSession(r); ok {
		m["UserRole"] = role
	}

	// If HTMX fragment is requested and provided, render only that part.
	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX && fragment != "" {
		renderFragment(w, page, fragment, m)
		return
	}

	// Otherwise, render the full page.
	render(w, page, m)
}

// render executes a full-page template using either the cache or per-request parsing.
func render(w http.ResponseWriter, name string, data any) {
	if devMode {
		// Hot reload: re-parse templates for every request.
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

	// Use pre-cached template set for production.
	t, ok := cache[name]
	if !ok {
		http.Error(w, fmt.Sprintf("template not found: %s", name), http.StatusNotFound)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderFragment executes a specific named template within a page's template set.
func renderFragment(w http.ResponseWriter, pageName, tmplName string, data any) {
	if devMode {
		// Hot reload for fragments.
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

	// Execute specific fragment from the cache.
	t, ok := cache[pageName]
	if !ok {
		http.Error(w, fmt.Sprintf("template not found: %s", pageName), http.StatusNotFound)
		return
	}
	if err := t.ExecuteTemplate(w, tmplName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
