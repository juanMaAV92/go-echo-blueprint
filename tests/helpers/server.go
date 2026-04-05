package helpers

import (
	"time"

	"github.com/juanMaAV92/go-echo-blueprint/config"
	"github.com/juanMaAV92/go-echo-blueprint/server"
	"github.com/juanMaAV92/go-utils/logger"
)

// NewTestServer builds a Server configured for testing.
// No .env loading, no real telemetry — safe to call from any test.
func NewTestServer() *server.Server {
	cfg := config.Config{
		ServiceName:  config.ServiceName,
		Environment:  "test",
		Port:         "8080",
		GracefulTime: 5 * time.Second,
	}
	log := logger.New(cfg.ServiceName, cfg.Environment)
	return server.New(cfg, log)
}
