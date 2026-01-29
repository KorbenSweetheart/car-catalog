package httpserver

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"

	"gitea.kood.tech/ivanandreev/viewer/internal/controller/http/handlers"
	"gitea.kood.tech/ivanandreev/viewer/internal/controller/http/middleware"
	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
)

// Maybe implement as struct in app???
type CarStorage interface {
	Car(ctx context.Context, ID int) (domain.Car, error)
	Cars(ctx context.Context) ([]domain.Car, error)
	RandomCars(ctx context.Context) ([]domain.Car, error)
	Catalog(ctx context.Context, filters domain.FilterOptions) ([]domain.Car, error)
	Metadata(ctx context.Context) (domain.Metadata, error)
}

func NewRouter(log *slog.Logger, tmplts map[string]*template.Template, storage CarStorage) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		mux,
		log,
		tmplts,
		storage,
	)

	reqID := middleware.NewReqIDMiddleware(log)
	logMw := middleware.NewLoggingMiddleware(log)
	recoverMw := middleware.NewRecoveringMiddleware(log)
	handler := middleware.Chain(mux, reqID, logMw, recoverMw)
	return handler
}

// func newMiddleware(log *slog.Logger) func(h http.Handler) http.Handler

func addRoutes(mux *http.ServeMux, logger *slog.Logger, tmplts map[string]*template.Template, storage CarStorage) {

	homeHandler := handlers.NewHomeHandler(logger, tmplts, storage)
	carHandler := handlers.NewCarHandler(logger, tmplts, storage)
	catalogHandler := handlers.NewCatalogHandler(logger, tmplts, storage)
	notFoundHandler := handlers.NewNotFoundHandler(logger, tmplts)

	mux.HandleFunc("GET /{$}", homeHandler.Index)
	mux.HandleFunc("GET /catalog/{id}", carHandler.Index)
	mux.HandleFunc("GET /catalog", catalogHandler.Index)
	mux.HandleFunc("/", notFoundHandler.NotFound)

	// Load static
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Action handlers
	// mux.Handle("POST /encoder", handlers.HandleEncoder(logger, proc, tmplts))
}
