package v1

import (
	"net/http"
	"strconv"

	"github.com/MisterCodo/plantjournal/store"
	"github.com/labstack/echo/v4"
)

// registerPlantRoutes registers all routes related to plants.
func (a *APIV1Service) registerPlantRoutes(g *echo.Group) {
	g.GET("/plants", a.GetPlants)
	g.GET("/plants/:id", a.GetPlantByID)
	g.POST("/plants", a.CreatePlant)
}

// GetPlants returns the list of all plants.
func (a *APIV1Service) GetPlants(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch all plants from database.
	plants, err := a.Store.GetPlants(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch all plants").SetInternal(err)
	}

	return c.Render(http.StatusOK, "plants.html", plants)
}

// GetPlantsByID returns the plant with given id from the request.
func (a *APIV1Service) GetPlantByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}

	// Fetch plant by id from store.
	p, err := a.Store.GetPlantByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch plant").SetInternal(err)
	}

	return c.Render(http.StatusOK, "plant.html", p)
}

// CreatePlant creates a new blank plant.
func (a *APIV1Service) CreatePlant(c echo.Context) error {
	ctx := c.Request().Context()

	_, err := a.Store.CreatePlant(ctx, &store.Plant{Name: "New Plant"})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create new plant").SetInternal(err)
	}

	// Refresh list of plants.
	plants, err := a.Store.GetPlants(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to refresh plants list").SetInternal(err)
	}

	// TODO: Fix selected plant from list. Currently focus is removed but a plant details section is still showed.

	return c.Render(http.StatusOK, "plants.html", plants)
}
