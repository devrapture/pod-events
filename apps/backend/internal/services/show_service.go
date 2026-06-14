package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/dto"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
)

type ShowServices interface {
	GetUserSavedShows(ctx context.Context, userID uuid.UUID, query string) ([]dto.SavedShowResponse, error)
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

func (s *showServices) GetUserSavedShows(ctx context.Context, userID uuid.UUID, query string) ([]dto.SavedShowResponse, error) {
	cacheKey := fmt.Sprintf("user:%s:saved-shows", userID.String())
	if cachedShows, found := s.cache.Get(cacheKey); found {
		shows, ok := cachedShows.([]dto.SavedShowResponse)
		if ok {
			return s.filterShows(shows, query), nil
		}
	}
	accessToken, err := s.authService.GetValidAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	savedShows, err := s.spotifyClient.GetAllUserSavedShows(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	shows := savedShows.ToSavedShows()
	s.cache.Set(cacheKey, shows, gocache.DefaultExpiration)
	return s.filterShows(shows, query), nil
}

func (s *showServices) filterShows(shows []dto.SavedShowResponse, query string) []dto.SavedShowResponse {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return shows
	}
	filtered := make([]dto.SavedShowResponse, 0)
	for _, item := range shows {
		if strings.Contains(strings.ToLower(item.Name), query) ||
			strings.Contains(strings.ToLower(item.Description), query) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
