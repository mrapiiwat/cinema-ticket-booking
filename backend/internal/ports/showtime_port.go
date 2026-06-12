package ports

import (
	"context"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
)

type ShowtimeRepositoryPort interface {
	FindByMovieID(ctx context.Context, movieID string) ([]model.Showtime, error)
	FindByID(ctx context.Context, id string) (*model.Showtime, error)
	Create(ctx context.Context, showtime *model.Showtime) error
	UpdateSeatStatus(ctx context.Context, showtimeID string, seatNo string, status string, lockedBy *string, lockedUntil *time.Time) error
	LockSeatIfAvailable(ctx context.Context, showtimeID string, seatNo string, userID string, lockedUntil time.Time) error
	ReleaseSeatLock(ctx context.Context, showtimeID string, seatNo string, userID string) error
	ReleaseExpiredLocks(ctx context.Context, showtimeID string, now time.Time) error
}

type ShowtimeServicePort interface {
	GetShowtimesByMovie(ctx context.Context, movieID string) ([]model.Showtime, error)
	GetShowtimeByID(ctx context.Context, id string) (*model.Showtime, error)
	CreateShowtime(ctx context.Context, showtime *model.Showtime) error
}
