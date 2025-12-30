package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

// The GET / endpoint returns the main page of the web interface.
// The GET / endpoint returns HTTP200.

func HandleHome(log *slog.Logger, tmplts map[string]*template.Template) http.Handler {
	const op = "handlers.home.handle.Home"

	log = log.With(
		slog.String("op", op),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		homeTmpl, ok := tmplts["home"]
		if !ok {
			log.Error("template not found in map", "name", "home")
			http.Error(w, "Internal Configuration Error", http.StatusInternalServerError)
			return
		}

		// Dynamic data to pass to HTML for test
		// data := map[string]string{
		// 	"input":        "",
		// 	"HTTPResponse": "200: OK",
		// }

		var buf bytes.Buffer

		if err := homeTmpl.Execute(&buf, nil); err != nil {
			log.Error("failed to render template", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		buf.WriteTo(w)
	})
}
