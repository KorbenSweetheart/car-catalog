package app

import (
	"context"
	"log/slog"
	"time"
	"viewer/internal/clients/webapi"
	"viewer/internal/config"
	httpserver "viewer/internal/controller/http"
	"viewer/internal/lib/e"
	"viewer/internal/repository/memory"
	"viewer/internal/usecase/carstore"
	"viewer/pkg/logger"
)

// This struct holds your entire running application
type App struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		log: logger.New(cfg.Env),
	}
}

func (app *App) Run() error {

	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	// Logger: slog
	app.log.Info("starting car viewer", slog.String("env", app.cfg.Env))
	app.log.Debug("debug messages are enabled")

	// Client
	client := webapi.New(
		app.cfg.Client.Host,
		app.cfg.Client.Timeout,
	)

	// Repository - Storage layer
	repo := memory.New(app.log, client)

	// Initial load (Startup phase)
	// TODO: think about backoff and failcount
	app.log.Info("performing initial data load...")

	startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := repo.Refresh(startupCtx); err != nil {
		app.log.Error("initial data load failed", slog.Any("error", err))
		return e.Wrap("initial data load failed", err)
	}

	app.log.Info("initial data load complete")

	// Start the Refresh Ticker and background memory refresh each 10 min
	go app.startBackgroundRefresh(appCtx, repo)

	// Usecase (CarStore) - business logic layer
	carStore := carstore.New(app.log, repo)

	// parse templates
	templates, err := httpserver.ParseTemplates(app.cfg.HTTPServer.TemplatesPath, app.log)
	if err != nil {
		app.log.Error("Failed to parse templates", slog.Any("error", err))
		return e.Wrap("failed to parse templates", err)
	}

	// construct new router -> Transport layer
	router := httpserver.NewRouter(app.log, templates, carStore)

	// construct new server
	httpServer := httpserver.NewHTTPServer(router, app.cfg)

	// run server
	if err := httpserver.RunServer(appCtx, app.log, app.cfg, httpServer, 5*time.Second); err != nil {
		app.log.Error("failed to start server", slog.Any("error", err))
		return err
	}

	return nil
}

func (app *App) startBackgroundRefresh(ctx context.Context, repo *memory.Repository) {
	ticker := time.NewTicker(app.cfg.Repo.RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // Stop signal: If Run() finishes, this case triggers
			return
		case <-ticker.C: // Ticker signal: Normal work
			app.log.Info("start refreshing data cache...")

			refreshCtx, cancel := context.WithTimeout(context.Background(), app.cfg.Client.Timeout*3)

			if err := repo.Refresh(refreshCtx); err != nil {
				app.log.Error("background refresh failed", slog.Any("error", err))
			} else {
				app.log.Info("data refreshed successfully")
			}
			cancel()
		}
	}
}
