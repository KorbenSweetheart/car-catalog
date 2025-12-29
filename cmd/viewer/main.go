package main

import (
	"log/slog"
	"os"
	"viewer/internal/clients/carapi"
	"viewer/internal/config"
	"viewer/internal/usecase"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config: json
	cfg := config.MustLoad()

	// init logger: slog
	log := setupLogger(cfg.Env)

	log.Info("starting car viewer", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// construct new Client
	client := carapi.New(
		cfg.Client.Host,
		cfg.Client.Timeout,
	)

	// construct Usecase (Service)
	carUsecase := usecase.NewCarModel(client)

	// parse templates

	// construct new router

	// construct new server

	// maybe: some storage or implement zero-copy

	// run server

}

// maybe move to pkg so we can reuse it.
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
