package server

import (
	"context"
	goerrors "errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/juanMaAV92/go-echo-blueprint/config"
	echoerr "github.com/juanMaAV92/go-utils/errors/echo"
	"github.com/juanMaAV92/go-utils/logger"
	"github.com/juanMaAV92/go-utils/validator"
	"github.com/labstack/echo/v4"
)

// Server wraps an Echo instance with its dependencies.
type Server struct {
	Echo   *echo.Echo
	cfg    config.Config
	logger logger.Logger
}

// New creates an Echo instance with the standard middleware stack pre-configured.
func New(cfg config.Config, log logger.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = echoerr.HTTPErrorHandler
	e.Validator = validator.New()

	return &Server{
		Echo:   e,
		cfg:    cfg,
		logger: log,
	}
}

// Run starts the HTTP server and blocks until a shutdown signal is received.
// Returns any non-ErrServerClosed error.
func (s *Server) Run() error {
	errC := make(chan error, 1)
	s.listenForShutdown(errC)

	go func() {
		if err := s.Echo.Start(":" + s.cfg.Port); err != nil && !goerrors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return <-errC
}

func (s *Server) listenForShutdown(errC chan<- error) {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()
		stop()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.GracefulTime)
		defer cancel()

		s.logger.Info(shutdownCtx, "server.shutdown", "shutting down gracefully")
		s.Echo.Server.SetKeepAlivesEnabled(false)

		if err := s.Echo.Shutdown(shutdownCtx); err != nil {
			errC <- err
			return
		}
		close(errC)
	}()
}
