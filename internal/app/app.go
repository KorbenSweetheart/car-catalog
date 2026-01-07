package app

import (
	"context"
	"log/slog"
	"time"
	"viewer/internal/clients/webapi"
	"viewer/internal/config"
	"viewer/internal/lib/e"
	"viewer/internal/repository/memory"
	"viewer/pkg/logger"
)

// This struct holds your entire running application
type App struct {
	logger slog.Logger
	client webapi.Client
	repo   memory.Repository

	// Usecases
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

func Run(cfg *config.Config) error { // maybe make is a method for App struct ???

	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	// Logger: slog
	log := logger.New(cfg.Env)

	log.Info("starting car viewer", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Client
	client := webapi.New(
		cfg.Client.Host,
		cfg.Client.Timeout,
	)

	// Repository
	repo := memory.New(log, client)

	// Initial load (Startup phase)
	// TODO: think about backoff and failcount
	log.Info("performing initial data load...")

	startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := repo.Refresh(startupCtx); err != nil {
		log.Error("initial data load failed", slog.Any("error", err))
		return e.Wrap("initial data load failed", err)
	}

	log.Info("initial data load complete")

	// Start the Refresh Ticker
	go func() {
		ticker := time.NewTicker(cfg.Repo.RefreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-appCtx.Done(): // Stop signal: If Run() finishes, this case triggers
				return
			case <-ticker.C: // Ticker signal: Normal work
				log.Info("start refreshing data cache...")

				refreshCtx, cancel := context.WithTimeout(context.Background(), cfg.Client.Timeout*3)

				if err := repo.Refresh(refreshCtx); err != nil {
					log.Error("background refresh failed", slog.Any("error", err))
				} else {
					log.Info("data refreshed successfully")
				}
				cancel()
			}
		}
	}()

	// Usecase (Service)
	carUsecase := usecase.NewCarModel(client)

	// parse templates

	// construct new router

	// construct new server

	// run server

	// graceful shutdown

	return nil
}
