package handlers

import (
	"bytes"
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
)

type CatalogUsecase interface {
	Catalog(ctx context.Context, filters domain.FilterOptions) ([]domain.Car, error)
	Metadata(ctx context.Context) (domain.Metadata, error)
	// We expect a new method that accepts filters
}

type CatalogHandler struct {
	log    *slog.Logger
	uc     CatalogUsecase
	tmplts map[string]*template.Template
}

func NewCatalogHandler(log *slog.Logger, tmplts map[string]*template.Template, uc CatalogUsecase) *CatalogHandler {
	return &CatalogHandler{log: log, uc: uc, tmplts: tmplts}
}

func (h *CatalogHandler) Index(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.catalog.Index"

	ctx := r.Context()

	// Parse FilterOptions
	q := r.URL.Query()

	filters := domain.FilterOptions{
		Transmission: q.Get("transmission"),
		Drivetrain:   q.Get("drivetrain"),
	}

	// Helper to safely parse integers (defaults to 0 if empty/invalid)
	filters.ManufacturerID, _ = strconv.Atoi(q.Get("manufacturer_id"))
	filters.CategoryID, _ = strconv.Atoi(q.Get("category_id"))
	filters.MinYear, _ = strconv.Atoi(q.Get("min_year"))
	filters.MinHP, _ = strconv.Atoi(q.Get("min_hp"))

	// Fetch Data (Cars & Metadata for Dropdowns)
	cars, err := h.uc.Catalog(ctx, filters)
	if err != nil {
		h.log.Error("failed to load catalog", "op", op, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Load to display filters in a sidebar
	metadata, err := h.uc.Metadata(ctx)
	if err != nil {
		h.log.Error("failed to load metadata", "op", op, "error", err)
		// We continue, just with empty dropdowns
	}

	// 3. Render
	data := map[string]any{
		"Title":    "Catalog | RedCar Oy",
		"Cars":     cars,
		"Metadata": metadata,
		"Filters":  filters, // Pass back so we can "pre-fill" the form inputs
	}

	tmpl, ok := h.tmplts["catalog.html"]
	if !ok {
		h.log.Error("template not found", "name", "catalog.html")
		http.Error(w, "Configuration Error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		h.log.Error("failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
