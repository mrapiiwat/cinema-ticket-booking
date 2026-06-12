package registry

import (
	"context"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/handlers"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/repositories"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	AuthHandler     *handlers.AuthHandler
	MovieHandler    *handlers.MovieHandler
	ShowtimeHandler *handlers.ShowtimeHandler
	BookingHandler  *handlers.BookingHandler
	WSHandler       *handlers.WSHandler
	AdminHandler    *handlers.AdminHandler
}

func NewContainer(db *mongo.Database) *Container {

	userRepo := repositories.NewUserRepository()
	movieRepo := repositories.NewMovieRepository()
	showtimeRepo := repositories.NewShowtimeRepository()
	bookingRepo := repositories.NewBookingRepository()
	auditRepo := repositories.NewAuditRepository()
	lockRepo := repositories.NewLockRepository()
	mqRepo := repositories.NewMQRepository()

	authService := services.NewAuthService(userRepo)
	movieService := services.NewMovieService(movieRepo)
	showtimeService := services.NewShowtimeService(showtimeRepo)
	auditService := services.NewAuditService(auditRepo)
	wsHub := services.NewWSHub()
	bookingService := services.NewBookingService(bookingRepo, showtimeRepo, lockRepo, mqRepo, wsHub)

	go wsHub.Run()
	go services.StartAuditLogConsumer(mqRepo, auditService, wsHub)
	go bookingService.StartExpirationWorker(context.Background())

	authHandler := handlers.NewAuthHandler(authService)
	movieHandler := handlers.NewMovieHandler(movieService)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	wsHandler := handlers.NewWSHandler(wsHub)
	adminHandler := handlers.NewAdminHandler(bookingService, auditService)

	return &Container{
		AuthHandler:     authHandler,
		MovieHandler:    movieHandler,
		ShowtimeHandler: showtimeHandler,
		BookingHandler:  bookingHandler,
		WSHandler:       wsHandler,
		AdminHandler:    adminHandler,
	}
}
