package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService interface {
	GenerateState() (string, error)
	GetAuthorizationURL(state string) string
	HandleCallback(ctx context.Context, code, state, cookieState string) (*models.User, string, error)
	GetValidAccessToken(ctx context.Context, userID uuid.UUID) (string, error)
}

type authService struct {
	cfg             *config.Config
	tokenRepository repositories.TokenRepository
	userRepository  repositories.UserRepository
	spotifyClient   *spotify.SpotifyClient
	logger          *zap.Logger
}

func NewAuthService(cfg *config.Config, tr repositories.TokenRepository, ur repositories.UserRepository, sc *spotify.SpotifyClient, logger *zap.Logger) AuthService {
	return &authService{
		cfg:             cfg,
		tokenRepository: tr,
		userRepository:  ur,
		spotifyClient:   sc,
		logger:          logger,
	}
}

// GenerateState creates a cryptographically random state string for CSRF protection.
// Store this in a cookie before redirecting to Spotify.
func (s *authService) GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (s *authService) GetAuthorizationURL(state string) string {
	return s.spotifyClient.AuthorizationURL(state)
}

func (s *authService) HandleCallback(ctx context.Context, code, state, cookieState string) (*models.User, string, error) {
	if state != cookieState || state == "" {
		return nil, "", fmt.Errorf("invalid state parameter - possible CSRF attack")
	}

	tokenResp, err := s.spotifyClient.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("exchange spotify code: %w", err)
	}
	spotifyUser, err := s.spotifyClient.GetCurrentUser(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("get spotify user: %w", err)
	}

	user, err := s.userRepository.GetBySpotifyUserID(ctx, spotifyUser.ID)
	if err != nil {
		return nil, "", fmt.Errorf("find user by spotify id: %w", err)
	}
	if user == nil {
		user, err = s.userRepository.GetByEmail(ctx, spotifyUser.Email)
		if err != nil {
			return nil, "", fmt.Errorf("find user by email: %w", err)
		}
	}

	if user == nil {
		user = &models.User{
			Name:          spotifyUser.DisplayName,
			Email:         spotifyUser.Email,
			SpotifyUserID: spotifyUser.ID,
			AvatarURL:     spotifyUser.AvatarURL(),
		}

		if err := s.userRepository.Create(ctx, user); err != nil {
			return nil, "", fmt.Errorf("create user: %w", err)
		}
		s.logger.Info("new user created via spotify oauth", zap.String("user_id", user.ID.String()), zap.String("email", user.Email))
	} else {
		// Existing user — update their Spotify info in case it changed
		user.Name = spotifyUser.DisplayName
		user.AvatarURL = spotifyUser.AvatarURL()
		user.SpotifyUserID = spotifyUser.ID

		if err := s.userRepository.Update(ctx, user); err != nil {
			return nil, "", fmt.Errorf("update user: %w", err)
		}
	}

	spotifyToken := &models.SpotifyToken{
		UserID:       user.ID,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    tokenResp.ExpiresAt(),
	}

	if tokenResp.RefreshToken == "" {
		exiting, _ := s.tokenRepository.GetByUserID(ctx, user.ID)
		if exiting != nil {
			spotifyToken.RefreshToken = exiting.RefreshToken
		}
	}
	if err := s.tokenRepository.Upsert(ctx, spotifyToken); err != nil {
		return nil, "", fmt.Errorf("save spotify tokens: %w", err)
	}

	// TODO implement jwt
	s.logger.Info("user authenticated via spotify", zap.String("user_id", user.ID.String()))
	// return user, sessionToken, nil
	return user, "", nil
}

// GetValidAccessToken returns a valid Spotify access token for a user.
// If the stored token is expired, it automatically refreshes it.
// This is the function you call before every Spotify API request.
func (s *authService) GetValidAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := s.tokenRepository.GetByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get spotify token: %w", err)
	}

	if token == nil {
		return "", fmt.Errorf("no spotify token found for user")
	}
	// Token is still valid — return it directly
	if !token.IsExpired() {
		return token.AccessToken, nil
	}
	s.logger.Info("refreshing expired spotify token", zap.String("userID", userID.String()))
	refreshedToken, err := s.spotifyClient.RefreshAccessToken(ctx, token.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh spotify token: %w", err)
	}

	newExpiredAt := refreshedToken.ExpiresAt()
	if err := s.tokenRepository.UpdateAccessToken(ctx, userID, refreshedToken.AccessToken, newExpiredAt); err != nil {
		return "", fmt.Errorf("saved refresh token: %w", err)
	}
	return refreshedToken.AccessToken, nil
}
