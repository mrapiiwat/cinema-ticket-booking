package ports

import (
	"context"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
)

type BookingRepositoryPort interface {
	Create(ctx context.Context, booking *model.Booking) error
	FindByID(ctx context.Context, id string) (*model.Booking, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	FindByUserID(ctx context.Context, userID string) ([]model.Booking, error)
	FindAll(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error)
	FindExpiredPending(ctx context.Context, now time.Time) ([]model.Booking, error)
}

type BookingServicePort interface {
	LockSeats(ctx context.Context, showtimeID string, seats []string, userID string) (*model.Booking, error)
	ConfirmBooking(ctx context.Context, bookingID, userID string) (*model.Booking, error)
	GetUserBookings(ctx context.Context, userID string) ([]model.Booking, error)
	GetAllBookings(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error)
	StartExpirationWorker(ctx context.Context)
}
