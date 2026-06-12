package services

import (
	"context"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type movieService struct {
	movieRepo ports.MovieRepositoryPort
}

func NewMovieService(movieRepo ports.MovieRepositoryPort) ports.MovieServicePort {
	return &movieService{
		movieRepo: movieRepo,
	}
}

func (s *movieService) GetMovies(ctx context.Context) ([]model.Movie, error) {
	return s.movieRepo.FindAll(ctx)
}

func (s *movieService) GetMovieByID(ctx context.Context, id string) (*model.Movie, error) {
	return s.movieRepo.FindByID(ctx, id)
}

func (s *movieService) CreateMovie(ctx context.Context, movie *model.Movie) error {
	return s.movieRepo.Create(ctx, movie)
}
