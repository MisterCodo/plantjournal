package server

import (
	"context"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Temporary counter variable, just for testing
var counter = 0

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Server implements a plant journal server
type Server struct {
	e *echo.Echo
}

// NewServer loads templates, sets the server and registers routes
func NewServer(ctx context.Context) (*Server, error) {
	e := echo.New()

	// Load html templates
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Set new server
	s := &Server{
		e: e,
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
	return s.e.Start(":8080")
}

// Shutdown closes the server
func (s *Server) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.e.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to gracefully close server", "err", err)
	}
}
