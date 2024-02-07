package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// registerPlantRoutes registers all routes related to plants.
func (a *APIV1Service) registerPlantRoutes(g *echo.Group) {
	g.GET("/plants", a.GetPlants)
	g.GET("/plants/:id", a.GetPlantByID)
}

// GetPlants returns the list of all plants.
func (a *APIV1Service) GetPlants(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch all plants from database.
	plants, err := a.Store.GetPlants(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch all plants").SetInternal(err)
	}

	// TODO: return html.

	return c.Render(http.StatusOK, "plants.html", plants)
}

// GetPlantsByID returns the plant with given id from the request.
func (a *APIV1Service) GetPlantByID(c echo.Context) error {
	ctx := c.Request().Context()

	parsed, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}
	id := int32(parsed)

	// Fetch plant by id from store.
	p, err := a.Store.GetPlantByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch plant").SetInternal(err)
	}

	// Set data for html template.
	data := map[string]interface{}{
		"ID":          p.ID,
		"Name":        p.Name,
		"Lighting":    p.Lighting,
		"Watering":    p.Watering,
		"Fertilizing": p.Fertilizing,
		"Toxicity":    p.Toxicity,
		"Notes":       p.Notes,
	}

	return c.Render(http.StatusOK, "plant.html", data)
}
