package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

type SystemHandler struct {
	log    *slog.Logger
	tmplts map[string]*template.Template
}

func NewSystemHandler(log *slog.Logger, tmplts map[string]*template.Template) *SystemHandler {
	return &SystemHandler{
		log:    log,
		tmplts: tmplts,
	}
}

func (h *SystemHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.system.NotFound"

	data := map[string]any{
		"Title": "Page Not Found | RedCars",
	}

	tmpl, ok := h.tmplts["404.html"]
	if !ok {
		h.log.Error("template not found", "op", op, "name", "404.html")
		http.Error(w, "404 Page Not Found", http.StatusNotFound) // Fallback text
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		h.log.Error("failed to render template", "op", op, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	buf.WriteTo(w)
}
