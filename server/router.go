package server

import (
	"github.com/juanMaAV92/go-echo-blueprint/internal/health"
	"github.com/juanMaAV92/go-utils/middleware/identity"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// Handlers aggregates all HTTP handlers.
// Add new handlers here as the service grows.
type Handlers struct {
	Health health.Handler
}

// RegisterRoutes attaches global middleware and all route groups to the Echo instance.
func (s *Server) RegisterRoutes(h Handlers) {
	s.registerMiddleware()

	api := s.Echo.Group("/" + s.cfg.ServiceName)

	// Public routes — no identity required
	health.RegisterRoutes(api, h.Health)

	// Protected routes — identity middleware applied at group level
	// Extend HeaderConfig.Extra for any project-specific headers
	protected := api.Group("", identity.Middleware(identity.HeaderConfig{}))
	_ = protected // register feature routes here: feature.RegisterRoutes(protected, h.Feature)
}

func (s *Server) registerMiddleware() {
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.RequestID())
	s.Echo.Use(otelecho.Middleware(s.cfg.ServiceName))
}
