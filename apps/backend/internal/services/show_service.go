package services

import (
	"context"
	"fmt"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/dto"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
)

type ShowServices interface {
	GetUserSavedShows(ctx context.Context, userID uuid.UUID) ([]dto.SavedShowResponse, error)
}

type showServices struct {
	spotifyClient *spotify.SpotifyClient
	authService   AuthService
	cfg           *config.Config
	cache         *gocache.Cache
}

func NewShowServices(sc *spotify.SpotifyClient, authService AuthService, cfg *config.Config, cache *gocache.Cache) ShowServices {
	return &showServices{
		spotifyClient: sc,
		authService:   authService,
		cfg:           cfg,
		cache:         cache,
	}
}

func (s *showServices) GetUserSavedShows(ctx context.Context, userID uuid.UUID) ([]dto.SavedShowResponse, error) {
	cacheKey := fmt.Sprintf("user:%s:saved-shows", userID.String())
	if cachedShows, found := s.cache.Get(cacheKey); found {
		shows, ok := cachedShows.([]dto.SavedShowResponse)
		if ok {
			return shows, nil
		}
	}
	accessToken, err := s.authService.GetValidAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	shows, err := s.spotifyClient.GetUserSavedShows(ctx, accessToken, userID)
	if err != nil {
		return nil, err
	}
	response := shows.ToSavedShows()
	s.cache.Set(cacheKey, response, gocache.DefaultExpiration)
	return response, nil
}
