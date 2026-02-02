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

func NewHomeHandler(log *slog.Logger, tmplts map[string]*template.Template, uc HomeUsecase) *HomeHandler {
	return &HomeHandler{
		log:    log,
		uc:     uc,
		tmplts: tmplts,
	}
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.home.Index"

	log := h.log.With(
		slog.String("op", op),
	)

	// We use the request context to support cancellation/timeouts
	ctx := r.Context()

	// 1. Fetch Data via Usecase
	popularCars, err := h.uc.RandomCars(ctx)
	if err != nil {
		log.Error("failed to load home data", slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusInternalServerError)
		return
	}

	familyCars, err := h.uc.RandomCars(ctx)
	if err != nil {
		log.Error("failed to load home data", slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusInternalServerError)
		return
	}

	// 2. Prepare Data for Template
	data := map[string]any{
		"Title":       "Home - RedCar Oy",
		"PopularCars": popularCars, // Passed to {{range .Cars}}
		"FamilyCars":  familyCars,
	}

	// 3. Render
	tmpl, ok := h.tmplts["home.html"]
	if !ok {
		log.Error("template not found", "name", "home.html")
		RenderError(w, h.tmplts, log, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Error("failed to render template", slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
