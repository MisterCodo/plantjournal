package v1

import (
	"fmt"
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
	g.PUT("/plants/:id", a.UpdatePlant)
	g.PUT("/plants/:id/water", a.WaterPlant)
	g.PUT("/plants/:id/fertilize", a.FertilizePlant)
	g.PUT("/plants/:id/note", a.AddNoteToPlant)
	g.PUT("/plants/:id/actions/:day", a.UpdateAction)
	g.DELETE("/plants/:id", a.DeletePlant)
	g.DELETE("/plants/:id/actions/:day", a.DeleteAction)
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

	// TODO: Doing a full page reload for now, but not ideal behavior.
	c.Response().Header().Set("HX-Redirect", "/")

	return c.NoContent(http.StatusOK)
}

type UpdatePlantRequest struct {
	Name        string `form:"name"`
	Lighting    string `form:"lighting"`
	Watering    string `form:"watering"`
	Fertilizing string `form:"fertilizing"`
	Toxicity    string `form:"toxicity"`
	Notes       string `form:"notes"`
}

// UpdatePlant updates the plant in the database.
func (a *APIV1Service) UpdatePlant(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch plant details from request.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}
	r := &UpdatePlantRequest{}
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not get details for plant update request").SetInternal(err)
	}

	// Set plant for store update.
	p := &store.Plant{
		ID:          id,
		Name:        r.Name,
		Lighting:    r.Lighting,
		Watering:    r.Watering,
		Fertilizing: r.Fertilizing,
		Toxicity:    r.Toxicity,
		Notes:       r.Notes,
	}

	// Update plant in database.
	err = a.Store.UpdatePlant(ctx, p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plant").SetInternal(err)
	}

	// Update plant label in plant list.
	return c.String(http.StatusOK, fmt.Sprintf("(%d) %s", p.ID, p.Name))
}

// DeletePlant deletes the plant with passed id value from the database.
func (a *APIV1Service) DeletePlant(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch plant details from request.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}

	// Delete plant from database.
	err = a.Store.DeletePlant(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete plant").SetInternal(err)
	}

	// Trigger a full page reload.
	c.Response().Header().Set("HX-Redirect", "/")

	return c.NoContent(http.StatusOK)
}

// WaterPlant saves an action of watering on today's date for the plant with passed id. If plant was already
// watered today, nothing is changed in the database.
func (a *APIV1Service) WaterPlant(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch plant details from request.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}

	// Upsert watering action in database.
	err = a.Store.WaterPlant(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to water plant").SetInternal(err)
	}

	// Refresh entire plant to see latest plant actions.
	p, err := a.Store.GetPlantByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch plant").SetInternal(err)
	}

	return c.Render(http.StatusOK, "plant.html", p)
}

// FertilizePlant saves an action of fertilizing on today's date for the plant with passed id. If plant was already
// fertilized today, nothing is changed in the database.
func (a *APIV1Service) FertilizePlant(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch plant details from request.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}

	// Upsert fertilizing action in database.
	err = a.Store.FertilizePlant(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fertilize plant").SetInternal(err)
	}

	// Refresh entire plant to see latest plant actions.
	p, err := a.Store.GetPlantByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch plant").SetInternal(err)
	}

	return c.Render(http.StatusOK, "plant.html", p)
}

// AddNoteToPlant saves an action on today's date for the plant with passed id.
func (a *APIV1Service) AddNoteToPlant(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch plant details from request.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}

	// Upsert add note to plant action in database.
	err = a.Store.AddNoteToPlant(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to add note to plant").SetInternal(err)
	}

	// Refresh entire plant to see latest plant actions.
	p, err := a.Store.GetPlantByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch plant").SetInternal(err)
	}

	return c.Render(http.StatusOK, "plant.html", p)
}

type UpdateActionRequest struct {
	Watered    string `form:"watered"`
	Fertilized string `form:"ferlilized"`
	Notes      string `form:"notes"`
}

// UpdateAction updates a plant maintenance action for the plant with passed id and the passed day.
func (a *APIV1Service) UpdateAction(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch plant details from request.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}

	// Fetch action date from request.
	actionDay := c.Param("day")

	r := &UpdateActionRequest{}
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not get details for action update request").SetInternal(err)
	}

	// Set action for store update.
	action := &store.Action{
		PlantID:    id,
		Day:        actionDay,
		Watered:    true,
		Fertilized: true,
		Notes:      r.Notes,
	}

	// Update action in database.
	err = a.Store.UpdateAction(ctx, action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update action").SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

// DeleteAction deletes a maintenance action for the plant with passed id and the passed day.
func (a *APIV1Service) DeleteAction(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch plant details from request.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	}

	// Fetch action date from request.
	actionDay := c.Param("day")

	// Delete plant from database.
	err = a.Store.DeleteAction(ctx, id, actionDay)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete plant action").SetInternal(err)
	}

	// Refresh entire plant to see latest plant actions.
	p, err := a.Store.GetPlantByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch plant").SetInternal(err)
	}

	return c.Render(http.StatusOK, "plant.html", p)
}
