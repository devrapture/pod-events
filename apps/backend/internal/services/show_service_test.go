package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMapSubscribeSpotifyError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		wantIs error
	}{
		{
			name:   "resource not found maps to podcast show not found",
			err:    apperrors.ErrSpotifyResourceNotFound,
			wantIs: apperrors.ErrPodcastShowNotFound,
		},
		{
			name:   "podcast show not found preserved",
			err:    apperrors.ErrPodcastShowNotFound,
			wantIs: apperrors.ErrPodcastShowNotFound,
		},
		{
			name:   "token not found preserved",
			err:    apperrors.ErrorSpotifyTokenNotFound,
			wantIs: apperrors.ErrorSpotifyTokenNotFound,
		},
		{
			name:   "authorization required preserved",
			err:    apperrors.ErrSpotifyAuthorizationRequired,
			wantIs: apperrors.ErrSpotifyAuthorizationRequired,
		},
		{
			name:   "rate limit preserved",
			err:    &spotify.RateLimitError{RetryAfter: 3},
			wantIs: nil, // checked via errors.As below
		},
		{
			name:   "already unavailable preserved",
			err:    apperrors.ErrSpotifyUnavailable,
			wantIs: apperrors.ErrSpotifyUnavailable,
		},
		{
			name:   "unknown error wrapped as unavailable",
			err:    errors.New("connection reset"),
			wantIs: apperrors.ErrSpotifyUnavailable,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mapSubscribeSpotifyError(tt.err)

			if tt.name == "rate limit preserved" {
				var rateLimitErr *spotify.RateLimitError
				require.True(t, errors.As(got, &rateLimitErr))
				assert.Equal(t, 3, rateLimitErr.RetryAfter)
				return
			}

			require.ErrorIs(t, got, tt.wantIs)
			if tt.wantIs == apperrors.ErrSpotifyUnavailable && !errors.Is(tt.err, apperrors.ErrSpotifyUnavailable) {
				assert.Contains(t, got.Error(), tt.err.Error())
			}
		})
	}
}

func TestMapSubscribeSpotifyErrorWrappedNotFound(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("%w: gone", apperrors.ErrSpotifyResourceNotFound)
	got := mapSubscribeSpotifyError(err)
	require.ErrorIs(t, got, apperrors.ErrPodcastShowNotFound)
}

type stubTokenRepository struct{}

func (stubTokenRepository) Upsert(context.Context, *models.SpotifyToken) error { return nil }
func (stubTokenRepository) UpdateAccessToken(context.Context, uuid.UUID, string, time.Time) error {
	return nil
}
func (stubTokenRepository) GetByUserID(context.Context, uuid.UUID) (*models.SpotifyToken, error) {
	return nil, nil
}
func (stubTokenRepository) Delete(context.Context, uuid.UUID) error { return nil }

func TestSubscribeMissingToken(t *testing.T) {
	authService := NewAuthService(
		&config.Config{},
		stubTokenRepository{},
		nil,
		nil,
		gocache.New(gocache.DefaultExpiration, gocache.DefaultExpiration),
		zap.NewNop(),
	)
	svc := NewShowServices(
		nil,
		authService,
		gocache.New(gocache.DefaultExpiration, gocache.DefaultExpiration),
		nil,
		nil,
		zap.NewNop(),
	)

	_, err := svc.Subscribe(context.Background(), uuid.New(), []string{"uncached-spotify-show"})

	require.ErrorIs(t, err, apperrors.ErrorSpotifyTokenNotFound)
}
