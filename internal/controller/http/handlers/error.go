package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

func RenderError(w http.ResponseWriter, tmplts map[string]*template.Template, log *slog.Logger, code int) {
	const op = "handlers.common.RenderError"

	log = log.With(
		slog.String("op", op),
	)

	// Default to the generic error template
	tmplName := "404.html"
	title := "Page Not Found | RedCars"
	heading := "Page Not Found"
	message := "It looks like you've taken a wrong turn. We can get you back on the road."

	// Implement if we don't have a 500.html or make dynamique error template
	if code == http.StatusInternalServerError {
		tmplName = "maintenance.html"
		title = "Internal Server Error | RedCars"
		heading = "Internal Server Error"
		message = "We are currently experiencing technical difficulties. Our mechanics are working on it."
	}

	data := map[string]any{
		"Title":   title,
		"Heading": heading,
		"Message": message,
	}

	tmpl, ok := tmplts[tmplName]
	if !ok {
		log.Error("template not found", "name", tmplName)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Error("failed to render error template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	buf.WriteTo(w)
}
