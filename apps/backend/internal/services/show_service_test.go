package services

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
