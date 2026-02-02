package app

import (
	"context"
	"log/slog"
	"time"

	"gitea.kood.tech/ivanandreev/viewer/internal/config"
	"gitea.kood.tech/ivanandreev/viewer/internal/controller/httpserver"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/adapter"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
	"gitea.kood.tech/ivanandreev/viewer/internal/repository/webapi"
	"gitea.kood.tech/ivanandreev/viewer/internal/usecase/carstore"
	"gitea.kood.tech/ivanandreev/viewer/pkg/cache"
	"gitea.kood.tech/ivanandreev/viewer/pkg/httpclient"
	"gitea.kood.tech/ivanandreev/viewer/pkg/logger"
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
	client := httpclient.New(
		app.cfg.Client.Host,
		app.cfg.Client.Timeout,
	)

	// Repository - Storage layer (WebAPI storage)
	repo := webapi.New(
		app.log,
		client,
	)

	// Cache
	cache := cache.New(app.cfg.Cache.DefaultExpiration, app.cfg.Cache.CleanupInterval)
	app.log.Info("launched cache janitor in goroutine")

	// Cache adapter to wire cache keys-value, to domain structs
	cacheAdapter := adapter.NewAdapter(cache, app.log)

	// Usecase (CarStore) - business logic layer
	carStore := carstore.New(app.log, repo, cacheAdapter)

	// parse templates
	templates, err := httpserver.ParseTemplates(app.cfg.HTTPServer.TemplatesPath, app.log)
	if err != nil {
		app.log.Error("Failed to parse templates", slog.Any("error", err))
		return e.Wrap("failed to parse templates", err)
	}

	// Router -> Transport layer
	router := httpserver.NewRouter(app.log, templates, carStore)

	// Server
	// TODO: maybe move to pkg as well.
	httpServer := httpserver.NewHTTPServer(router, app.cfg)

	// run server
	if err := httpserver.RunServer(appCtx, app.log, app.cfg, httpServer, 5*time.Second); err != nil {
		app.log.Error("failed to start server", slog.Any("error", err))
		return err
	}

	return nil
}
