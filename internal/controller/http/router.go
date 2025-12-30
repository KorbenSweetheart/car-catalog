package http

import (
	"html/template"
	"log/slog"
	"net/http"
	"viewer/internal/controller/http/handlers"
	"viewer/internal/controller/http/middleware"
)

func NewRouter(log *slog.Logger, tmplts map[string]*template.Template) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		mux,
		log,
		tmplts,
	)

	reqID := middleware.NewReqIDMiddleware(log)
	logMw := middleware.NewLoggingMiddleware(log)
	recoverMw := middleware.NewRecoveringMiddleware(log)
	handler := middleware.Chain(mux, reqID, logMw, recoverMw)
	return handler
}

// func newMiddleware(log *slog.Logger) func(h http.Handler) http.Handler

func addRoutes(mux *http.ServeMux, logger *slog.Logger, tmplts map[string]*template.Template) {

	// Page handlers
	mux.Handle("GET /{$}", handlers.HandleHome(logger, tmplts)) // handle must return http.Handler
	mux.Handle("GET /static/", handlers.HandleStatic(logger))
	mux.Handle("GET /", handlers.HandleNotFound(logger, tmplts))

	// Action handlers
	// mux.Handle("POST /decoder", handlers.HandleDecoder(logger, proc, tmplts))
	// mux.Handle("POST /encoder", handlers.HandleEncoder(logger, proc, tmplts))
}
