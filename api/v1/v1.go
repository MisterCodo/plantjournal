package v1

import (
	"github.com/labstack/echo/v4"
)

type APIV1Service struct {
}

// NewAPIV1Service returns a new APIV1Service.
func NewAPIV1Service() *APIV1Service {
	return &APIV1Service{}
}

// Register registers echo APIV1Service routes.
func (a *APIV1Service) Register(rootGroup *echo.Group) {
	apiV1ServiceGroup := rootGroup.Group("/api/v1")
	a.registerPlantRoutes(apiV1ServiceGroup)
}
