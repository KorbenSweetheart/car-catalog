package handlers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
)

type CarUsecase interface {
	Car(ctx context.Context, ID int) (domain.Car, error)
	RandomCars(ctx context.Context) ([]domain.Car, error) // TODO: display cars from the same category/brand or most viewed based on cookies
}

type CarHandler struct {
	log    *slog.Logger
	uc     CarUsecase
	tmplts map[string]*template.Template
}

// Temporary view model for the template
type Expert struct {
	Name     string
	Title    string
	Location string
	Email    string
	Phone    string
	ImageURL string
}

func NewCarHandler(log *slog.Logger, tmplts map[string]*template.Template, uc CarUsecase) *CarHandler {
	return &CarHandler{
		log:    log,
		uc:     uc,
		tmplts: tmplts,
	}
}

func (h *CarHandler) Index(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.car.Index"

	// We use the request context to support cancellation/timeouts
	ctx := r.Context()

	// 1. Fetch Data via Usecase

	idStr := r.PathValue("id")
	ID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error("failed to convert id to integer", "op", op, "input", idStr, "error", err)
		http.Error(w, "Invalid Car ID", http.StatusBadRequest)
		return
	}

	car, err := h.uc.Car(ctx, ID)
	if err != nil {
		h.log.Error("failed to load car by id", "op", op, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	popularCars, err := h.uc.RandomCars(ctx)
	if err != nil {
		h.log.Error("failed to load home data", "op", op, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Placeholder: Hardcoded list of experts
	experts := []Expert{
		{Name: "Mika Häkkinen", Title: "Car salesman", Location: "Vantaa", Email: "mika.hakkinen@redcars.fi", Phone: "+358 50 222 3333", ImageURL: "https://i.pravatar.cc/150?img=57"},
		{Name: "Kimi Räikkönen", Title: "Car salesman", Location: "Espoo", Email: "kimi.raikkonen@redcars.fi", Phone: "+358 50 444 5555", ImageURL: "https://i.pravatar.cc/150?img=52"},
		{Name: "Valtteri Bottas", Title: "Car salesman", Location: "Nastola", Email: "valtteri.bottas@redcars.fi", Phone: "+358 50 666 7777", ImageURL: "https://i.pravatar.cc/150?img=14"},
		{Name: "Keke Rosberg", Title: "Car salesman", Location: "Solna", Email: "keke.rosberg@redcars.fi", Phone: "+358 50 123 4567", ImageURL: "https://i.pravatar.cc/150?img=69"},
	}

	// 2. Prepare Data for Template
	data := map[string]any{
		"Title":       fmt.Sprintf("%s %d - HP: %d, %s | RedCar Oy", car.Name, car.Year, car.Specs.HP, car.Specs.Transmission),
		"Car":         car,
		"PopularCars": popularCars, // TODO: display recentrly viewed cars or the same category/vendor
		"Experts":     experts,     // Placeholder
	}

	// 3. Render
	tmpl, ok := h.tmplts["car.html"]
	if !ok {
		h.log.Error("template not found", "op", op, "name", "car.html")
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
