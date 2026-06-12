package ports

import (
	"context"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
)

type UserRepositoryPort interface {
	FindByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	UpdateGoogleProfile(ctx context.Context, user *model.User) error
}

type AuthServicePort interface {
	VerifyGoogleToken(ctx context.Context, tokenStr string) (*model.AuthResponse, error)
	ExchangeGoogleCode(ctx context.Context, code string, redirectURI string) (*model.AuthResponse, error)
	DevLogin(ctx context.Context, role string) (*model.AuthResponse, error)
	GoogleConfig() map[string]interface{}
}
