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

type CompareUsecase interface {
	Car(ctx context.Context, id int) (domain.Car, error)
}

type CompareHandler struct {
	log    *slog.Logger
	uc     CompareUsecase
	tmplts map[string]*template.Template
}

func NewCompareHandler(log *slog.Logger, tmplts map[string]*template.Template, uc CompareUsecase) *CompareHandler {
	return &CompareHandler{log: log, uc: uc, tmplts: tmplts}
}

func (h *CompareHandler) Index(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.compare.Index"

	log := h.log.With("op", op)

	ctx := r.Context()

	// 1. Parse & Validate IDs
	idsStr := r.URL.Query().Get("ids")
	var validIDs []int
	var cleanIDStrings []string // We keep this to pass back to the template

	// duplicates check
	seen := make(map[int]bool)

	if idsStr != "" {
		for _, p := range strings.Split(idsStr, ",") {
			// Clean up whitespace and validate
			if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil && id > 0 {
				// Only append if we haven't seen this ID yet
				if !seen[id] {
					seen[id] = true
					validIDs = append(validIDs, id)
					cleanIDStrings = append(cleanIDStrings, strconv.Itoa(id))
				}
			}
		}
	}

	// If the user manually requests more than 3 cars, we simply ignore the extras.
	maxCars := 3
	if len(validIDs) > maxCars {
		validIDs = validIDs[:maxCars]
		cleanIDStrings = cleanIDStrings[:maxCars]
	}

	// 2. Fetch Cars
	// Note: If you have many users, fetching one-by-one (N+1) is inefficient.
	// Ideally, your Usecase would support h.uc.Cars(ctx, validIDs) to do 1 DB query.
	// For max 4 cars, this is acceptable.
	var cars []domain.Car
	maxHP := 0
	newestCar := 0

	for _, id := range validIDs {
		car, err := h.uc.Car(ctx, id)
		if err != nil {
			log.Warn("failed to fetch car for comparison", "id", id, slog.Any("error", err))
			continue
		}

		// Check for max HP values while fetching
		if car.Specs.HP > maxHP {
			maxHP = car.Specs.HP
		}

		// Check for latest year
		if car.Year > newestCar {
			newestCar = car.Year
		}

		cars = append(cars, car)
	}

	// 3. Render
	// We pass "CompareIDs" (clean string) so the template helper can calculate removals.
	data := map[string]any{
		"Title":      "Compare Vehicles | RedCars",
		"Cars":       cars,
		"CompareIDs": strings.Join(cleanIDStrings, ","),
		"MaxHP":      maxHP,
		"NewestCar":  newestCar,
	}

	tmpl, ok := h.tmplts["compare.html"]
	if !ok {
		log.Error("template not found", "name", "compare.html")
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
