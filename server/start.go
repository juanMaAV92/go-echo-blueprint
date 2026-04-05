package server

import (
	"context"
	"os"

	"github.com/juanMaAV92/go-echo-blueprint/config"
	"github.com/juanMaAV92/go-echo-blueprint/internal/health"
	"github.com/juanMaAV92/go-utils/logger"
	"github.com/juanMaAV92/go-utils/telemetry"
)

const (
	exitFailStartup = iota + 1
	exitFailRuntime
)

// Start is the application entry point.
// It loads configuration, initialises observability, wires dependencies, and runs the server.
func Start() {
	ctx := context.Background()
	cfg := config.Load()

	log := logger.New(cfg.ServiceName, cfg.Environment)

	shutdown, err := telemetry.InitTelemetry(ctx, cfg.Telemetry)
	if err != nil {
		log.Fatal(ctx, "start.telemetry", "failed to initialise telemetry", "error", err)
		os.Exit(exitFailStartup)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Error(ctx, "start.telemetry", "error shutting down telemetry", "error", err)
		}
	}()

	srv := New(cfg, log)

	// Wire dependencies — add new services/handlers here
	handlers := Handlers{
		Health: health.NewHandler(health.NewService()),
	}
	srv.RegisterRoutes(handlers)

	log.Info(ctx, "start.server", "server starting", "port", cfg.Port, "env", cfg.Environment)

	if err := srv.Run(); err != nil {
		log.Fatal(ctx, "start.server", "server stopped with error", "error", err)
		os.Exit(exitFailRuntime)
	}
}
