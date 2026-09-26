package services

import (
	"context"
	"testing"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type authTestTokenRepository struct {
	token         *models.SpotifyToken
	requestedUser uuid.UUID
}

func (r *authTestTokenRepository) Upsert(context.Context, *models.SpotifyToken) error {
	return nil
}

func (r *authTestTokenRepository) UpdateAccessToken(context.Context, uuid.UUID, string, time.Time) error {
	return nil
}

func (r *authTestTokenRepository) GetByUserID(_ context.Context, userID uuid.UUID) (*models.SpotifyToken, error) {
	r.requestedUser = userID
	return r.token, nil
}

type authTestUserRepository struct {
	user *models.User
}

func (r *authTestUserRepository) Create(context.Context, *models.User) error { return nil }
func (r *authTestUserRepository) Update(context.Context, *models.User) error { return nil }
func (r *authTestUserRepository) GetByID(context.Context, uuid.UUID) (*models.User, error) {
	return r.user, nil
}
func (r *authTestUserRepository) GetByGoogleUserID(context.Context, string) (*models.User, error) {
	return nil, nil
}
func (r *authTestUserRepository) GetBySpotifyUserID(context.Context, string) (*models.User, error) {
	return nil, nil
}
func (r *authTestUserRepository) GetByEmail(_ context.Context, email string) (*models.User, error) {
	if r.user != nil && r.user.Email == email {
		return r.user, nil
	}
	return nil, nil
}

func TestGetSpotifyAccessTokenUsesConfiguredOwner(t *testing.T) {
	ownerID := uuid.New()
	tokenRepo := &authTestTokenRepository{
		token: &models.SpotifyToken{
			UserID:      ownerID,
			AccessToken: "owner-access-token",
			ExpiresAt:   time.Now().Add(time.Hour),
		},
	}
	userRepo := &authTestUserRepository{user: &models.User{
		Base:  models.Base{ID: ownerID},
		Email: "owner@example.com",
	}}
	service := NewAuthService(
		&config.Config{OwnerEmail: "owner@example.com"},
		tokenRepo,
		userRepo,
		nil,
		nil,
		gocache.New(time.Minute, time.Minute),
		zap.NewNop(),
	)

	accessToken, err := service.GetSpotifyAccessToken(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "owner-access-token", accessToken)
	assert.Equal(t, ownerID, tokenRepo.requestedUser)
}

func TestIsOwnerMatchesEmailCaseInsensitively(t *testing.T) {
	ownerID := uuid.New()
	service := NewAuthService(
		&config.Config{OwnerEmail: "owner@example.com"},
		&authTestTokenRepository{},
		&authTestUserRepository{user: &models.User{
			Base:  models.Base{ID: ownerID},
			Email: "Owner@Example.com",
		}},
		nil,
		nil,
		gocache.New(time.Minute, time.Minute),
		zap.NewNop(),
	)

	isOwner, err := service.IsOwner(context.Background(), ownerID)

	require.NoError(t, err)
	assert.True(t, isOwner)
}
