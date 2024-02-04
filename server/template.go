package server

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// Template implements the echo.Renderer interface, allowing the use of html templates
type Template struct {
	templates *template.Template
}

// Render executes the template with given name and passed data
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
