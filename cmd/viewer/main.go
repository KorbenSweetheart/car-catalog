package main

import (
	"log/slog"
	"viewer/internal/clients/carapi"
	"viewer/internal/config"
	"viewer/internal/usecase"
	"viewer/pkg/logger"
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
	// maybe move to appRun()
	log := logger.New(cfg.Env)

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
