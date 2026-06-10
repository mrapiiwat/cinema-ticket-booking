package registry

import (
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/handlers"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/repositories"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/services"
)

type Container struct {
	AuthHandler *handlers.AuthHandler
}

func NewContainer(db *mongo.Database) *Container {
	
	userRepo := repositories.NewUserRepository() 
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	return &Container{
		AuthHandler: authHandler,
	}
}