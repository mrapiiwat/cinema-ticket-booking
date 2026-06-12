package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/security"
	xoauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

var googleScopes = []string{"openid", "email", "profile"}

type authService struct {
	userRepo ports.UserRepositoryPort
}

type googleProfile struct {
	ID            string
	Email         string
	Name          string
	EmailVerified bool
}

func NewAuthService(userRepo ports.UserRepositoryPort) ports.AuthServicePort {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) VerifyGoogleToken(ctx context.Context, tokenStr string) (*model.AuthResponse, error) {
	if tokenStr == "" {
		return nil, errors.New("invalid google token")
	}
	if googleClientID() == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID is not configured")
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
	if !matchesGoogleAudience(tokenInfo.Audience, tokenInfo.IssuedTo) {
		return nil, errors.New("google token audience does not match GOOGLE_CLIENT_ID")
	}
	if !tokenInfo.VerifiedEmail {
		return nil, errors.New("google email is not verified")
	}

	return s.authenticateGoogleProfile(ctx, googleProfile{
		ID:            tokenInfo.UserId,
		Email:         tokenInfo.Email,
		Name:          tokenInfo.Email,
		EmailVerified: tokenInfo.VerifiedEmail,
	})
}

func (s *authService) ExchangeGoogleCode(ctx context.Context, code string, redirectURI string) (*model.AuthResponse, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("invalid google authorization code")
	}

	config, err := googleOAuthConfig(redirectURI)
	if err != nil {
		return nil, err
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("google code exchange failed: %w", err)
	}

	idToken, _ := token.Extra("id_token").(string)
	if idToken != "" {
		return s.VerifyGoogleToken(ctx, idToken)
	}

	profile, err := fetchGoogleProfile(ctx, config, token)
	if err != nil {
		return nil, err
	}
	return s.authenticateGoogleProfile(ctx, profile)
}

func (s *authService) GoogleConfig() map[string]interface{} {
	clientID := googleClientID()
	return map[string]interface{}{
		"client_id":         clientID,
		"configured":        clientID != "" && googleClientSecret() != "",
		"scope":             strings.Join(googleScopes, " "),
		"dev_login_enabled": devLoginEnabled(),
	}
}

func (s *authService) DevLogin(ctx context.Context, role string) (*model.AuthResponse, error) {
	if !devLoginEnabled() {
		return nil, errors.New("dev login is disabled")
	}

	if strings.ToUpper(role) != model.RoleAdmin {
		role = model.RoleUser
	} else {
		role = model.RoleAdmin
	}

	googleID := "dev-user"
	email := "demo.user@cinepass.local"
	name := "Demo User"
	if role == model.RoleAdmin {
		googleID = "dev-admin"
		email = "admin@cinepass.local"
		name = "Admin Staff"
	}

	user, err := s.userRepo.FindByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user = &model.User{
			GoogleID: googleID,
			Name:     name,
			Email:    email,
			Role:     role,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	return buildAuthResponse(user)
}

func (s *authService) authenticateGoogleProfile(ctx context.Context, profile googleProfile) (*model.AuthResponse, error) {
	if profile.ID == "" || profile.Email == "" {
		return nil, errors.New("invalid profile metadata from google")
	}
	if !profile.EmailVerified {
		return nil, errors.New("google email is not verified")
	}

	name := strings.TrimSpace(profile.Name)
	if name == "" {
		name = profile.Email
	}

	user, err := s.userRepo.FindByGoogleID(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	role := roleForEmail(profile.Email)
	if user == nil {
		user = &model.User{
			GoogleID: profile.ID,
			Name:     name,
			Email:    profile.Email,
			Role:     role,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else {
		user.Name = name
		user.Email = profile.Email
		user.Role = role
		if err := s.userRepo.UpdateGoogleProfile(ctx, user); err != nil {
			return nil, err
		}
	}

	return buildAuthResponse(user)
}

func buildAuthResponse(user *model.User) (*model.AuthResponse, error) {
	token, err := security.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, err
	}
	return &model.AuthResponse{Token: token, User: user}, nil
}

func fetchGoogleProfile(ctx context.Context, config *xoauth2.Config, token *xoauth2.Token) (googleProfile, error) {
	client := config.Client(ctx, token)
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return googleProfile{}, errors.New("failed to initialize google userinfo service")
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return googleProfile{}, errors.New("google userinfo request failed")
	}

	emailVerified := true
	if userInfo.VerifiedEmail != nil {
		emailVerified = *userInfo.VerifiedEmail
	}

	return googleProfile{
		ID:            userInfo.Id,
		Email:         userInfo.Email,
		Name:          userInfo.Name,
		EmailVerified: emailVerified,
	}, nil
}

func googleOAuthConfig(redirectURI string) (*xoauth2.Config, error) {
	clientID := googleClientID()
	clientSecret := googleClientSecret()
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be configured")
	}

	redirectURI = strings.TrimSpace(redirectURI)
	if redirectURI == "" {
		redirectURI = strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"))
	}
	if redirectURI == "" {
		redirectURI = "postmessage"
	}

	return &xoauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       googleScopes,
		Endpoint:     google.Endpoint,
	}, nil
}

func matchesGoogleAudience(audience string, issuedTo string) bool {
	clientID := googleClientID()
	return audience == clientID || issuedTo == clientID
}

func googleClientID() string {
	return strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
}

func googleClientSecret() string {
	return strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
}

func devLoginEnabled() bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("ENABLE_DEV_LOGIN")))
	if raw == "" {
		return true
	}
	return raw == "1" || raw == "true" || raw == "yes" || raw == "on"
}

func roleForEmail(email string) string {
	adminEmails := strings.Split(os.Getenv("ADMIN_EMAILS"), ",")
	for _, item := range adminEmails {
		if strings.EqualFold(strings.TrimSpace(item), email) {
			return model.RoleAdmin
		}
	}
	return model.RoleUser
}
