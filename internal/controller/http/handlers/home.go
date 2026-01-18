package handlers

import (
	"bytes"
	"context"
	"html/template"
	"log/slog"
	"net/http"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
)

type HomeUsecase interface {
	RandomCars(ctx context.Context) ([]domain.Car, error)
}

type HomeHandler struct {
	log    *slog.Logger
	uc     HomeUsecase
	tmplts map[string]*template.Template
}

func NewHomeHandler(log *slog.Logger, uc HomeUsecase, tmplts map[string]*template.Template) *HomeHandler {
	return &HomeHandler{
		log:    log,
		uc:     uc,
		tmplts: tmplts,
	}
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.home.Index"

	// We use the request context to support cancellation/timeouts
	ctx := r.Context()

	// 1. Fetch Data via Usecase
	cars, err := h.uc.RandomCars(ctx)
	if err != nil {
		h.log.Error("failed to load home data", "op", op, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 2. Prepare Data for Template
	data := map[string]any{
		"Title": "Home - CarViewer",
		"Cars":  cars, // Passed to {{range .Cars}}
	}

	// 3. Render
	tmpl, ok := h.tmplts["home.html"]
	if !ok {
		h.log.Error("template not found", "op", op, "name", "home.html")
		http.Error(w, "Configuration Error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		h.log.Error("failed to render template", "op", op, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
