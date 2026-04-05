package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler defines the HTTP interface for health operations.
type Handler interface {
	Check(c echo.Context) error
}

type handler struct {
	svc Service
}

// NewHandler returns a Handler backed by the given Service.
func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

// RegisterRoutes attaches health routes to the provided Echo group.
func RegisterRoutes(g *echo.Group, h Handler) {
	g.GET("/health", h.Check)
}

func (h *handler) Check(c echo.Context) error {
	return c.JSON(http.StatusOK, h.svc.Check(c.Request().Context()))
}
