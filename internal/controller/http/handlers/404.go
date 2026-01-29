package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type NotFoundHandler struct {
	log    *slog.Logger
	tmplts map[string]*template.Template
}

func NewNotFoundHandler(log *slog.Logger, tmplts map[string]*template.Template) *NotFoundHandler {
	return &NotFoundHandler{
		log:    log,
		tmplts: tmplts,
	}
}

func (h *NotFoundHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.NotFound"

	RenderError(w, h.tmplts, h.log, http.StatusNotFound)
}
