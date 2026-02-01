package handlers

import (
	"bytes"
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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

	log := h.log.With(
		slog.String("op", op),
	)

	ctx := r.Context()

	// Parse FilterOptions
	q := r.URL.Query()

	filters := domain.FilterOptions{
		Transmission: q.Get("transmission"),
		Drivetrain:   q.Get("drivetrain"),
		SearchQuery:  q.Get("q"),
	}

	// Helper to safely parse integers (defaults to 0 if empty/invalid)
	filters.ManufacturerID, _ = strconv.Atoi(q.Get("manufacturer_id"))
	filters.CategoryID, _ = strconv.Atoi(q.Get("category_id"))
	filters.MinYear, _ = strconv.Atoi(q.Get("min_year"))
	filters.MinHP, _ = strconv.Atoi(q.Get("min_hp"))

	// Comparisson logic
	compareIDsStr := q.Get("compare_ids")
	selectedMap := make(map[int]bool)
	count := 0

	if compareIDsStr != "" {
		parts := strings.Split(compareIDsStr, ",")
		for _, p := range parts {
			// Clean whitespace just in case
			p = strings.TrimSpace(p)
			if id, err := strconv.Atoi(p); err == nil && id > 0 {
				selectedMap[id] = true
				count++
			}
		}
	}

	// Flag to disable "Add" buttons if full
	limitReached := count >= 3

	// Fetch Data (Cars & Metadata for Dropdowns)
	cars, err := h.uc.Catalog(ctx, filters)
	if err != nil {
		log.Error("failed to load catalog", slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusInternalServerError)
		return
	}

	// Load to display filters in a sidebar
	metadata, err := h.uc.Metadata(ctx)
	if err != nil {
		log.Error("failed to load metadata", slog.Any("error", err))
		// We continue, just with empty dropdowns
	}

	// 3. Render
	data := map[string]any{
		"Title":        "Catalog | RedCar Oy",
		"Cars":         cars,
		"Metadata":     metadata,
		"Filters":      filters,       // Pass back so we can "pre-fill" the form inputs
		"CompareIDs":   compareIDsStr, // The raw string for generating links
		"SelectedMap":  selectedMap,   // To visually mark selected cars
		"LimitReached": limitReached,  // To block add for comparisson button
		"Params":       r.URL.Query(),
	}

	tmpl, ok := h.tmplts["catalog.html"]
	if !ok {
		log.Error("template not found", "name", "catalog.html")
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
	buf.WriteTo(w)
}
