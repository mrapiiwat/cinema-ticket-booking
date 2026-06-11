package registry

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/routes"
)

func SetupRoutes(app *echo.Echo, container *Container) {
	api := app.Group("/api/v1")

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	routes.MapAuthRoutes(api, container.AuthHandler)

}
