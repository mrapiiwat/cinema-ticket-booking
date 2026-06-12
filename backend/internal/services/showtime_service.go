package services

import (
	"context"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type showtimeService struct {
	showtimeRepo ports.ShowtimeRepositoryPort
}

func NewShowtimeService(showtimeRepo ports.ShowtimeRepositoryPort) ports.ShowtimeServicePort {
	return &showtimeService{
		showtimeRepo: showtimeRepo,
	}
}

func (s *showtimeService) GetShowtimesByMovie(ctx context.Context, movieID string) ([]model.Showtime, error) {
	showtimes, err := s.showtimeRepo.FindByMovieID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for _, showtime := range showtimes {
		_ = s.showtimeRepo.ReleaseExpiredLocks(ctx, showtime.ID.Hex(), now)
	}
	return s.showtimeRepo.FindByMovieID(ctx, movieID)
}

func (s *showtimeService) GetShowtimeByID(ctx context.Context, id string) (*model.Showtime, error) {
	if err := s.showtimeRepo.ReleaseExpiredLocks(ctx, id, time.Now()); err != nil {
		return nil, err
	}
	return s.showtimeRepo.FindByID(ctx, id)
}

func (s *showtimeService) CreateShowtime(ctx context.Context, showtime *model.Showtime) error {
	return s.showtimeRepo.Create(ctx, showtime)
}
