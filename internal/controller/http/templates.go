package httpserver

import (
	"fmt"
	"html/template"
	"log/slog"
	"path/filepath"
	"viewer/internal/lib/e"
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
	const op = "render.ParseTemplates"
	log = log.With("op", op)

	cache := make(map[string]*template.Template)

	// 1. Helper Functions
	// Useful for templates (e.g., {{ .Specs.Drivetrain | firstLetter }})
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) },
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

		// OPTIMIZATION: Clone the root (which already has layouts/partials).
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
