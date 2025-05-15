package app

import (
	"context"
	"os"
	"os/signal"
	"parallel_download_from_many_urls/internal/adapter"
	"parallel_download_from_many_urls/internal/handler"
	"parallel_download_from_many_urls/internal/service"
	"parallel_download_from_many_urls/pkg/logger"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"parallel_download_from_many_urls/internal/config"
	"parallel_download_from_many_urls/internal/server"
)

type App struct {
	cfg    *config.Config
	log    *zerolog.Logger
	server *server.Server
}

func New() *App {
	// Логгер
	log := logger.InitLogger()

	cfg := config.MustReadConfig(log)

	adapterObj := adapter.New(log, cfg)
	serviceObj := service.New(cfg, log, adapterObj)
	handlerObj := handler.New(cfg, log, serviceObj)

	srv := server.New(cfg.Server, handlerObj)

	return &App{
		cfg:    cfg,
		log:    log,
		server: srv,
	}
}

func (app *App) Run() {
	stopChanel := make(chan os.Signal, 1)
	signal.Notify(stopChanel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		app.log.Info().Msg("Input System start working")
		if err := app.server.Run(); err != nil {
			app.log.Error().Msgf("Could not listen on port: %s, Error: %v", app.cfg.Server.Port, err)
		}
	}()

	<-stopChanel
	app.log.Info().Msg("Shutting down server gracefully ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.ShutDown(ctx); err != nil {
		app.log.Fatal().Msgf("server shutdown Error: %v", err)
	}
	
	app.log.Info().Msg("Server exited properly")
}
