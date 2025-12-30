package http

import (
	"html/template"
	"log/slog"
)

// NOTE: its better to do this by building a Template Cache (or a "Render Engine") and using Layout Inheritance.
// TODO: refactor to parse all templates at once, without manual changes

// Idea:
// tmpl := template.Must(template.New("").Funcs(template.FuncMap{
// 		"contains": contains,
// 	}).ParseGlob("static/templates/*.html"))

func ParseTemplates(log *slog.Logger) (map[string]*template.Template, error) {

	tmplts := make(map[string]*template.Template)

	// Helper to reduce copy-paste error handling
	parse := func(name, path string) error {
		t, err := template.ParseFiles(path)
		if err != nil {
			log.Error("failed to parse template", "path", path, slog.Any("error", err))
			return err
		}
		tmplts[name] = t
		return nil
	}

	if err := parse("layout", "static/templates/layout.html"); err != nil {
		return nil, err
	}

	if err := parse("home", "static/templates/home.html"); err != nil {
		return nil, err
	}

	if err := parse("404", "static/templates/404.html"); err != nil {
		return nil, err
	}

	return tmplts, nil
}
