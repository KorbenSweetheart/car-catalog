package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

const page = "404.html"

func HandleNotFound(log *slog.Logger, tmplts map[string]*template.Template) http.Handler {
	const op = "handlers.404.handle.NotFound"

	log = log.With(
		slog.String("op", op),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notFoundTmpl, ok := tmplts["404"]
		if !ok {
			log.Error("template not found in map", "name", "404")
			http.Error(w, "Internal Configuration Error", http.StatusInternalServerError)
			return
		}

		data := map[string]string{
			"HTTPResponse": "404: NOT FOUND",
		}

		var buf bytes.Buffer

		if err := notFoundTmpl.Execute(&buf, data); err != nil {
			log.Error("failed to render template", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		buf.WriteTo(w)
	})
}
