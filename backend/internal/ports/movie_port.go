package ports

import (
	"context"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
)

type MovieRepositoryPort interface {
	FindAll(ctx context.Context) ([]model.Movie, error)
	FindByID(ctx context.Context, id string) (*model.Movie, error)
	Create(ctx context.Context, movie *model.Movie) error
}

type MovieServicePort interface {
	GetMovies(ctx context.Context) ([]model.Movie, error)
	GetMovieByID(ctx context.Context, id string) (*model.Movie, error)
	CreateMovie(ctx context.Context, movie *model.Movie) error
}
