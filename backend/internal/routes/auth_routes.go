package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/handlers"
)

func MapAuthRoutes(router *echo.Group, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/google", authHandler.Login)
	}
}