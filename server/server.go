package server

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Temporary counter variable, just for testing
var counter = 0

// Server implements a plant journal server
type Server struct {
	e      *echo.Echo
	logger *slog.Logger
}

// NewServer loads templates, sets the server and registers routes
func NewServer(ctx context.Context, logger *slog.Logger) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Set logger middleware to slog
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	// Load html templates
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Set new server
	s := &Server{
		e:      e,
		logger: logger,
	}

	// Serve static content
	e.Static("/static", "static")

	// Serve healthz endpoint for systems which perform health checks
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Serve other routes
	e.GET("/", HomeHandler)
	e.POST("/increase", IncreaseHandler)

	return s, nil
}

func HomeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"CounterValue": counter,
	})
}

func IncreaseHandler(c echo.Context) error {
	counter++
	data := map[string]interface{}{
		"CounterValue": counter,
	}
	return c.Render(http.StatusOK, "counter.html", data)
}

// Start starts serving the server
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("started server")
	return s.e.Start(":8080")
}

// Shutdown closes the server
func (s *Server) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	defer s.logger.Info("closed server")

	err := s.e.Shutdown(ctx)
	if err != nil {
		s.logger.Error("failed to gracefully close server", "err", err)
	}
}