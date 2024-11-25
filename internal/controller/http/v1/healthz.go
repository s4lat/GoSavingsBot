package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthzRoutes struct{}

func NewHealthzRoutes() HealthzRoutes {
	return HealthzRoutes{}
}

type EmptyResponse struct{}

// getHealthz godoc
//
//	@Summary    health check
//	@Tags		metrics
//	@Produce	json
//
//	@Success	200		{object}	EmptyResponse
//	@Router		/healthz [GET]
func (cr *HealthzRoutes) getHealthz(c echo.Context) error {
	return c.JSON(http.StatusOK, struct{}{})
}
