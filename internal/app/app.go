package app

import (
	"log/slog"
	"viewer/internal/config"
	carapi "viewer/internal/repository/api"
	"viewer/pkg/logger"
)

func Run(cfg *config.Config) { // maybe make is a method for App struct ???

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

// This struct holds your entire running application
type App struct {
	// HomeUsecase    *home.Usecase
	// CatalogUsecase *catalog.Usecase
	// Add others here...
}

// New initializes everything in one place
// func New(apiHost string) *App {
// 1. Init the shared adapter once
// client := carapi.New(apiHost)

// 2. Init all use cases, injecting the client
// return &App{
// 	HomeUsecase:    home.New(client),
// 	CatalogUsecase: catalog.New(client),
// }
// }
