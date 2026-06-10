package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type authService struct {
	userRepo ports.UserRepositoryPort
}

func NewAuthService(userRepo ports.UserRepositoryPort) ports.AuthServicePort {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) VerifyGoogleToken(ctx context.Context, tokenStr string) (*model.User, error) {
	if tokenStr == "" {
		return nil, errors.New("invalid google token")
	}

	httpClient := &http.Client{}
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, errors.New("failed to initialize google auth service")
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(tokenStr)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, errors.New("google token validation failed or expired")
	}

	googleID := tokenInfo.UserId 
	email := tokenInfo.Email
	
	name := email

	if googleID == "" {
		return nil, errors.New("invalid token metadata from google")
	}

	existingUser, err := s.userRepo.FindByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		newUser := &model.User{
			GoogleID: googleID,
			Name:     name,
			Email:    email,
			Role:     "USER", 
		}
		
		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return nil, err
		}
		return newUser, nil
	}

	return existingUser, nil
}