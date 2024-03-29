package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	v1 "github.com/MisterCodo/plantjournal/api/v1"
	"github.com/MisterCodo/plantjournal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server implements a plant journal server.
type Server struct {
	e      *echo.Echo
	logger *slog.Logger
	config *Config
	store  *store.Store
}

// NewServer loads templates, establishes persistent data storage, sets the server and registers routes.
func NewServer(ctx context.Context, logger *slog.Logger, config *Config) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Set logger middleware to slog.
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	// Load html templates.
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Set persistent data storage.
	store, err := store.NewStore(config.DB)
	if err != nil {
		return nil, err
	}

	// Set new server.
	s := &Server{
		e:      e,
		logger: logger,
		config: config,
		store:  store,
	}

	// Serve static content.
	e.Static("/static", "static")

	// Serve healthz endpoint for systems which perform health checks.
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Serve home.
	e.GET("/", HomeHandler)

	// Serve API V1 routes.
	rootGroup := e.Group("")
	apiV1Service := v1.NewAPIV1Service(s.store)
	apiV1Service.Register(rootGroup)

	return s, nil
}

// HomeHandler serves the home.html page.
func HomeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{})
}

// Start starts serving the server.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("started server")
	return s.e.Start(fmt.Sprintf("%s:%d", s.config.Addr, s.config.Port))
}

// Shutdown closes the server.
func (s *Server) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	defer s.store.Close()
	defer s.logger.Info("closed server")

	err := s.e.Shutdown(ctx)
	if err != nil {
		s.logger.Error("failed to gracefully close server", "err", err)
	}
}
