package registry

import (
	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/routes"
)

func SetupRoutes(app *echo.Echo, container *Container) {
	api := app.Group("/api/v1")

	routes.MapAuthRoutes(api, container.AuthHandler)

}
