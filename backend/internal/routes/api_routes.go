package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/handlers"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/middleware"
)

func MapAuthRoutes(router *echo.Group, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.GET("/google/config", authHandler.GoogleConfig)
		auth.POST("/google", authHandler.Login)
		auth.POST("/google/code", authHandler.GoogleCodeLogin)
		auth.POST("/dev", authHandler.DevLogin)
	}
}

func MapMovieRoutes(g *echo.Group, h *handlers.MovieHandler) {
	movies := g.Group("/movies")
	movies.GET("", h.GetMovies)
	movies.GET("/:id", h.GetMovieByID)
	movies.POST("", h.CreateMovie, middleware.AuthMiddleware, middleware.AdminMiddleware)
}

func MapShowtimeRoutes(g *echo.Group, h *handlers.ShowtimeHandler) {
	showtimes := g.Group("/showtimes")
	showtimes.GET("/movie/:movieID", h.GetShowtimesByMovie)
	showtimes.GET("/:id", h.GetShowtimeByID)
	showtimes.POST("", h.CreateShowtime, middleware.AuthMiddleware, middleware.AdminMiddleware)
}

func MapBookingRoutes(g *echo.Group, h *handlers.BookingHandler) {
	bookings := g.Group("/bookings")
	bookings.Use(middleware.AuthMiddleware)

	bookings.POST("/lock", h.LockSeat)
	bookings.POST("/confirm", h.ConfirmBooking)
	bookings.GET("/my", h.GetUserBookings)
}

func MapAdminRoutes(g *echo.Group, h *handlers.AdminHandler) {
	admin := g.Group("/admin", middleware.AuthMiddleware, middleware.AdminMiddleware)

	admin.GET("/bookings", h.GetBookings)
	admin.GET("/audit-logs", h.GetAuditLogs)
}

func MapWSRoutes(g *echo.Group, h *handlers.WSHandler) {
	g.GET("/ws", h.ServeWS)
}
