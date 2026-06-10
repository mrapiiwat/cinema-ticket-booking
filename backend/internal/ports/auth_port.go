package ports

import (
	"context"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
)

type UserRepositoryPort interface {
	FindByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type AuthServicePort interface {
	VerifyGoogleToken(ctx context.Context, tokenStr string) (*model.User, error)
}