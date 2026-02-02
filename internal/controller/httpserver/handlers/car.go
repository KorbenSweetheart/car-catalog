package handlers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"gitea.kood.tech/ivanandreev/viewer/internal/controller/httpserver/cookies"
	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
)

type CarUsecase interface {
	Car(ctx context.Context, ID int) (domain.Car, error)
	RecommendedCars(ctx context.Context, IDs []int) ([]domain.Car, error)
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

	log := h.log.With(
		slog.String("op", op),
	)

	// We use the request context to support cancellation/timeouts
	ctx := r.Context()

	// 1. Checking that ID in URL is valid
	idStr := r.PathValue("id")
	ID, err := strconv.Atoi(idStr)
	if err != nil || ID < 1 {
		log.Error("invalid car id", "input", idStr, slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusNotFound)
		return
	}

	// 2. Get car by ID
	car, err := h.uc.Car(ctx, ID)
	if err != nil {
		log.Warn("car not found", "id", ID, slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusNotFound)
		return
	}

	// --- COOKIES LOGIC ---
	// Update the 'viewed_cars' cookie with the current ID (Stack/Recency)
	cookies.TrackViewedCar(w, r, car.ID, log)

	viewedCarIDs := cookies.ViewedCarIDs(r, log)
	// Get Personalized cars (excl current ID if its top viewed car)
	recommendedCars, err := h.uc.RecommendedCars(ctx, viewedCarIDs)
	if err != nil {
		log.Error("failed to load recommended cars", slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusInternalServerError)
		return
	}

	// 3. Get 4 random cars
	popularCars, err := h.uc.RandomCars(ctx)
	if err != nil {
		log.Error("failed to load popular cars", slog.Any("error", err))
		RenderError(w, h.tmplts, log, http.StatusInternalServerError)
		return
	}

	// Placeholder: Hardcoded list of experts
	experts := []Expert{
		{Name: "Mika Häkkinen", Title: "Car salesman", Location: "Vantaa", Email: "mika.hakkinen@redcars.fi", Phone: "+358 50 222 3333", ImageURL: "https://i.pravatar.cc/150?img=57"},
		{Name: "Kimi Räikkönen", Title: "Car salesman", Location: "Espoo", Email: "kimi.raikkonen@redcars.fi", Phone: "+358 50 444 5555", ImageURL: "https://i.pravatar.cc/150?img=52"},
		{Name: "Valtteri Bottas", Title: "Car salesman", Location: "Nastola", Email: "valtteri.bottas@redcars.fi", Phone: "+358 50 666 7777", ImageURL: "https://i.pravatar.cc/150?img=14"},
		{Name: "Keke Rosberg", Title: "Car salesman", Location: "Solna", Email: "keke.rosberg@redcars.fi", Phone: "+358 50 123 4567", ImageURL: "https://i.pravatar.cc/150?img=69"},
	}

	// Prepare Data for Template
	data := map[string]any{
		"Title":       fmt.Sprintf("%s %d - HP: %d, %s | RedCar Oy", car.Name, car.Year, car.Specs.HP, car.Specs.Transmission),
		"Car":         car,
		"PopularCars": popularCars, // TODO: display recentrly viewed cars or the same category/vendor
		"Experts":     experts,     // Placeholder
	}

	// Render
	tmpl, ok := h.tmplts["car.html"]
	if !ok {
		log.Error("template not found", "name", "car.html")
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
