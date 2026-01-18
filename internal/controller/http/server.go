package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"gitea.kood.tech/ivanandreev/viewer/internal/config"
)

func NewHTTPServer(handler http.Handler, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
}

func RunServer(ctx context.Context, log *slog.Logger, cfg *config.Config, server *http.Server, shutdownTimeout time.Duration) error {
	const op = "internal.httpserver.server.Run"

	log = log.With(
		slog.String("op", op),
	)

	shutdownCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Creating channel to capture only startup errors
	serverErr := make(chan error, 1)

	go func() {
		log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", slog.Any("error", err))
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server failed to start: %w", err)
	case <-shutdownCtx.Done():
		log.Info("shutdown signal received, starting graceful shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Error("graceful shutdown timed out, forcing close", slog.Any("error", err))

		return server.Close()
	}

	log.Info("shutdown complete gracefully")

	return nil
}
