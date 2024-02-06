package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Plant contains the details of a plant.
type Plant struct {
	ID          int32
	Name        string
	Lighting    string
	Watering    string
	Fertilizing string
	Toxicity    string
	Notes       string
	// TODO maintenance item.
	// TODO image.
}

// registerPlantRoutes registers all routes related to plants.
func (a *APIV1Service) registerPlantRoutes(g *echo.Group) {
	g.GET("/plant", a.GetPlants)
	g.GET("/plant/:id", a.GetPlantByID)
}

// GetPlants returns the list of all plants.
func (a *APIV1Service) GetPlants(c echo.Context) error {
	// ctx := c.Request().Context()

	// TODO: fetch plants from database.

	// TODO: return html.

	return c.JSON(http.StatusOK, nil)
}

// GetPlantsByID returns the plant with given id from the request.
func (a *APIV1Service) GetPlantByID(c echo.Context) error {
	// ctx := c.Request().Context()

	// id, err := util.ConvertStringToInt32(c.Param("id"))
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "invalid plant id").SetInternal(err)
	// }

	// TODO: fetch plant by id from database. Delete hardcoded plant below.
	p := Plant{
		ID:          0,
		Name:        "Neon Pothos",
		Lighting:    "Bright indirect light, tolerates shade",
		Watering:    "Water when top 2 or 3 inches are dry",
		Fertilizing: "Weak fertilizer during spring and summer",
		Toxicity:    "Toxic to humans, cats and dogs",
		Notes:       "Cutting from mother plant",
	}

	data := map[string]interface{}{
		"Name":        p.Name,
		"Lighting":    p.Lighting,
		"Watering":    p.Watering,
		"Fertilizing": p.Fertilizing,
		"Toxicity":    p.Toxicity,
		"Notes":       p.Notes,
	}
	return c.Render(http.StatusOK, "plant.html", data)
}
