package httpserver

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"

	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

// NOTE: its better to do this by building a Template Cache (or a "Render Engine") and using Layout Inheritance.
// TODO: refactor to parse all templates at once, without manual changes

// Idea:
// tmpl := template.Must(template.New("").Funcs(template.FuncMap{
// 		"contains": contains,
// 	}).ParseGlob("static/templates/*.html"))

// ParseTemplates scans the directory and builds a map of ready-to-render templates.
// It combines: Layouts + Partials + [Specific Page]

func ParseTemplates(rootDir string, log *slog.Logger) (map[string]*template.Template, error) {
	const op = "internal.httpserver.ParseTemplates"

	log = log.With("op", op)

	cache := make(map[string]*template.Template)

	// 1. Helper Functions
	funcMap := template.FuncMap{
		"safe":         safe,
		"dict":         dict,
		"replaceParam": replaceParam,
		"toggleID":     toggleID,
		// Add more helpers here if needed later
	}

	// 2. Parse Shared Templates ONCE (Layouts + Partials)
	root := template.New("root").Funcs(funcMap)

	// A. Layouts
	layouts, err := filepath.Glob(filepath.Join(rootDir, "layouts", "*.html"))
	if err != nil {
		log.Error("failed to identify layouts", slog.Any("error", err))
		return nil, e.Wrap("failed to identify layouts", err)
	}
	if len(layouts) > 0 {
		if _, err := root.ParseFiles(layouts...); err != nil {
			return nil, e.Wrap("failed to parse layouts", err)
		}
	}

	// B. Partials
	partials, err := filepath.Glob(filepath.Join(rootDir, "partials", "*.html"))
	if err != nil {
		log.Error("failed to identify partials", slog.Any("error", err))
		return nil, e.Wrap("failed to identify partials", err)
	}
	if len(partials) > 0 {
		if _, err := root.ParseFiles(partials...); err != nil {
			return nil, e.Wrap("failed to parse partials", err)
		}
	}

	// 3. Parse Individual Pages
	pages, err := filepath.Glob(filepath.Join(rootDir, "pages", "*.html"))
	if err != nil {
		log.Error("failed to identify pages", slog.Any("error", err))
		return nil, e.Wrap("failed to identify pages", err)
	}

	for _, pagePath := range pages {
		name := filepath.Base(pagePath)

		// Optimisation: Clone the root (which already has layouts/partials).
		clone, err := root.Clone()
		if err != nil {
			log.Error("failed to clone root", slog.Any("error", err))
			return nil, e.Wrap("failed to clone root template", err)
		}

		// Only parse the specific page file into the clone
		if _, err := clone.ParseFiles(pagePath); err != nil {
			log.Error("failed to parse page template", "page", name, slog.Any("error", err))
			return nil, e.Wrap(fmt.Sprintf("failed to parse page %s", name), err)
		}

		tmpl := clone.Lookup(name)
		if tmpl == nil {
			log.Error("empty template", "page", name, slog.Any("error", err))
			return nil, fmt.Errorf("template %s not found after parsing", name)
		}

		cache[name] = tmpl
		log.Info("template cached", "name", name)
	}

	return cache, nil
}

func safe(s string) template.HTML {
	return template.HTML(s)
}

// dict creates a map from a list of key-value pairs.
// Usage: {{ dict "Key1" .Value1 "Key2" "Val2" }}
func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict: values must be in pairs")
	}
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key at index %d must be a string", i)
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// replaceParam updates a specific URL query parameter while preserving others.
func replaceParam(current url.Values, key, value string) string {
	// 1. Clone the current values
	newParams := make(url.Values)
	for k, v := range current {
		newParams[k] = v
	}

	// 2. Update or delete the specific key
	if value == "" {
		newParams.Del(key)
	} else {
		newParams.Set(key, value)
	}

	// 3. Encode and clean up commas
	encoded := newParams.Encode()
	return "?" + strings.ReplaceAll(encoded, "%2C", ",")
}

// toggleID adds an ID to a comma-separated list if missing, or removes it if present.
// Used for the "Add/Remove" logic in the comparison feature.
func toggleID(currentList string, id int) string {
	idStr := fmt.Sprintf("%d", id)

	if currentList == "" {
		return idStr
	}

	parts := strings.Split(currentList, ",")
	var kept []string
	exists := false

	for _, p := range parts {
		if p == idStr {
			exists = true // Found it, so we skip it (remove)
			continue
		}
		if p != "" {
			kept = append(kept, p)
		}
	}

	// If it didn't exist, append it (add)
	if !exists {
		kept = append(kept, idStr)
	}

	return strings.Join(kept, ",")
}
