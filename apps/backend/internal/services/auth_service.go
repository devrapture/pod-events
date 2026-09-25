package services

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/googleauth"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/devrapture/pod-events/pkg/jwt"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

const (
	googleOAuthStateCachePrefix       = "google_oauth_state:"
	spotifyOwnerOAuthStateCachePrefix = "spotify_owner_oauth_state:"
	oauthStateTTL                     = 5 * time.Minute
	authExchangeCodeCachePrefix       = "auth_exchange_code:"
	authExchangeCodeTTL               = 1 * time.Minute
)

type AuthExchange struct {
	Token string
	User  *models.User
}

type AuthService interface {
	GenerateState() (string, error)
	RememberGoogleOAuthState(state string)
	GetGoogleAuthorizationURL(state string) string
	HandleGoogleCallback(ctx context.Context, code, state, browserState string) (*models.User, string, error)
	RememberSpotifyOwnerOAuthState(state string, userID uuid.UUID)
	GetSpotifyOwnerAuthorizationURL(state string) string
	HandleSpotifyOwnerCallback(ctx context.Context, code, state string) error
	IsOwner(ctx context.Context, userID uuid.UUID) (bool, error)
	GetSpotifyAccessToken(ctx context.Context) (string, error)
	CreateAuthExchangeCode(token string, user *models.User) (string, error)
	ConsumeAuthExchangeCode(code string) (*AuthExchange, bool)
}

type authService struct {
	cfg             *config.Config
	tokenRepository repositories.TokenRepository
	userRepository  repositories.UserRepository
	googleClient    *googleauth.Client
	spotifyClient   *spotify.SpotifyClient
	cache           *gocache.Cache
	logger          *zap.Logger
	oauthStateMu    sync.Mutex
	spotifyTokenMu  sync.Mutex
}

func NewAuthService(
	cfg *config.Config,
	tr repositories.TokenRepository,
	ur repositories.UserRepository,
	gc *googleauth.Client,
	sc *spotify.SpotifyClient,
	cache *gocache.Cache,
	logger *zap.Logger,
) AuthService {
	return &authService{
		cfg:             cfg,
		tokenRepository: tr,
		userRepository:  ur,
		googleClient:    gc,
		spotifyClient:   sc,
		cache:           cache,
		logger:          logger,
	}
}

func (s *authService) GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (s *authService) RememberGoogleOAuthState(state string) {
	s.cache.Set(googleOAuthStateCachePrefix+state, true, oauthStateTTL)
}

func (s *authService) GetGoogleAuthorizationURL(state string) string {
	return s.googleClient.AuthorizationURL(state)
}

func (s *authService) RememberSpotifyOwnerOAuthState(state string, userID uuid.UUID) {
	s.cache.Set(spotifyOwnerOAuthStateCachePrefix+state, userID, oauthStateTTL)
}

func (s *authService) GetSpotifyOwnerAuthorizationURL(state string) string {
	return s.spotifyClient.AuthorizationURL(state)
}

func (s *authService) consumeOAuthState(prefix, state string) (any, bool) {
	if state == "" {
		return nil, false
	}

	cacheKey := prefix + state
	s.oauthStateMu.Lock()
	defer s.oauthStateMu.Unlock()

	value, found := s.cache.Get(cacheKey)
	if !found {
		return nil, false
	}
	s.cache.Delete(cacheKey)
	return value, true
}

func (s *authService) HandleGoogleCallback(ctx context.Context, code, state, browserState string) (*models.User, string, error) {
	if state == "" || browserState == "" || subtle.ConstantTimeCompare([]byte(state), []byte(browserState)) != 1 {
		return nil, "", fmt.Errorf("invalid state parameter - possible CSRF attack")
	}
	if _, ok := s.consumeOAuthState(googleOAuthStateCachePrefix, state); !ok {
		return nil, "", fmt.Errorf("invalid state parameter - possible CSRF attack")
	}

	profile, err := s.googleClient.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", err
	}
	email := strings.ToLower(strings.TrimSpace(profile.Email))

	user, err := s.userRepository.GetByGoogleUserID(ctx, profile.ID)
	if err != nil {
		return nil, "", fmt.Errorf("find user by google id: %w", err)
	}
	if user == nil {
		user, err = s.userRepository.GetByEmail(ctx, email)
		if err != nil {
			return nil, "", fmt.Errorf("find user by email: %w", err)
		}
	}

	googleUserID := profile.ID
	if user == nil {
		user = &models.User{
			Name:         profile.Name,
			Email:        email,
			AvatarURL:    profile.Picture,
			GoogleUserID: &googleUserID,
		}
		if err := s.userRepository.Create(ctx, user); err != nil {
			return nil, "", fmt.Errorf("create user: %w", err)
		}
		s.logger.Info("new user created via google oauth", zap.String("user_id", user.ID.String()))
	} else {
		if user.GoogleUserID != nil && *user.GoogleUserID != profile.ID {
			return nil, "", fmt.Errorf("email is already linked to another google account")
		}
		user.Name = profile.Name
		user.Email = email
		user.AvatarURL = profile.Picture
		user.GoogleUserID = &googleUserID
		if err := s.userRepository.Update(ctx, user); err != nil {
			return nil, "", fmt.Errorf("update user: %w", err)
		}
	}

	sessionToken, err := jwt.GenerateJwt(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", fmt.Errorf("generate jwt: %w", err)
	}
	s.logger.Info("user authenticated via google", zap.String("user_id", user.ID.String()))
	return user, sessionToken, nil
}

func (s *authService) IsOwner(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return strings.EqualFold(strings.TrimSpace(user.Email), s.cfg.OwnerEmail), nil
}

func (s *authService) HandleSpotifyOwnerCallback(ctx context.Context, code, state string) error {
	value, ok := s.consumeOAuthState(spotifyOwnerOAuthStateCachePrefix, state)
	if !ok {
		return fmt.Errorf("invalid state parameter - possible CSRF attack")
	}
	ownerID, ok := value.(uuid.UUID)
	if !ok {
		return fmt.Errorf("invalid spotify owner authorization state")
	}
	isOwner, err := s.IsOwner(ctx, ownerID)
	if err != nil {
		return fmt.Errorf("verify spotify owner: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("spotify authorization is restricted to the application owner")
	}

	tokenResponse, err := s.spotifyClient.ExchangeCode(ctx, code)
	if err != nil {
		return fmt.Errorf("exchange spotify code: %w", err)
	}
	spotifyToken := &models.SpotifyToken{
		UserID:       ownerID,
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    tokenResponse.ExpiresAt(),
		Scope:        tokenResponse.Scope,
	}
	if tokenResponse.RefreshToken == "" {
		existing, loadErr := s.tokenRepository.GetByUserID(ctx, ownerID)
		if loadErr != nil && !errors.Is(loadErr, apperrors.ErrorSpotifyTokenNotFound) {
			return fmt.Errorf("load existing spotify token: %w", loadErr)
		}
		if existing != nil {
			spotifyToken.RefreshToken = existing.RefreshToken
		}
	}
	if spotifyToken.RefreshToken == "" {
		return fmt.Errorf("spotify did not return a refresh token")
	}
	if err := s.tokenRepository.Upsert(ctx, spotifyToken); err != nil {
		return fmt.Errorf("save spotify owner token: %w", err)
	}

	s.logger.Info("spotify owner credential updated", zap.String("user_id", ownerID.String()))
	return nil
}

// GetSpotifyAccessToken returns the single owner-authorized Spotify token used
// for catalog searches, subscription hydration, and episode polling for all users.
func (s *authService) GetSpotifyAccessToken(ctx context.Context) (string, error) {
	if s.cfg.OwnerEmail == "" || s.userRepository == nil {
		return "", apperrors.ErrorSpotifyTokenNotFound
	}
	owner, err := s.userRepository.GetByEmail(ctx, s.cfg.OwnerEmail)
	if err != nil {
		return "", fmt.Errorf("find spotify owner: %w", err)
	}
	if owner == nil {
		return "", apperrors.ErrorSpotifyTokenNotFound
	}

	s.spotifyTokenMu.Lock()
	defer s.spotifyTokenMu.Unlock()

	token, err := s.tokenRepository.GetByUserID(ctx, owner.ID)
	if err != nil {
		return "", fmt.Errorf("get spotify owner token: %w", err)
	}
	if token == nil {
		return "", apperrors.ErrorSpotifyTokenNotFound
	}
	if !token.IsExpired() {
		return token.AccessToken, nil
	}

	s.logger.Info("refreshing expired spotify owner token", zap.String("user_id", owner.ID.String()))
	refreshedToken, err := s.spotifyClient.RefreshAccessToken(ctx, token.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh spotify owner token: %w", err)
	}
	newRefreshToken := token.RefreshToken
	if refreshedToken.RefreshToken != "" {
		newRefreshToken = refreshedToken.RefreshToken
	}
	if err := s.tokenRepository.Upsert(ctx, &models.SpotifyToken{
		UserID:       owner.ID,
		AccessToken:  refreshedToken.AccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    refreshedToken.ExpiresAt(),
		Scope:        token.Scope,
	}); err != nil {
		return "", fmt.Errorf("save refreshed spotify owner token: %w", err)
	}
	return refreshedToken.AccessToken, nil
}

func (s *authService) CreateAuthExchangeCode(token string, user *models.User) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate auth exchange code: %w", err)
	}

	code := hex.EncodeToString(b)
	s.cache.Set(authExchangeCodeCachePrefix+code, &AuthExchange{Token: token, User: user}, authExchangeCodeTTL)
	return code, nil
}

func (s *authService) ConsumeAuthExchangeCode(code string) (*AuthExchange, bool) {
	if code == "" {
		return nil, false
	}

	cacheKey := authExchangeCodeCachePrefix + code
	value, found := s.cache.Get(cacheKey)
	if !found {
		return nil, false
	}
	s.cache.Delete(cacheKey)

	exchange, ok := value.(*AuthExchange)
	return exchange, ok
}
