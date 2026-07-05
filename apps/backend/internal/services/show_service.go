package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/dto"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
)

type ShowServices interface {
	GetUserSavedShows(ctx context.Context, userID uuid.UUID, query string) ([]dto.SavedShowResponse, error)
	SearchShows(ctx context.Context, userID uuid.UUID, query string, limit, offset int) ([]dto.SavedShowResponse, error)
	GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error)
	Subscribe(ctx context.Context, userID uuid.UUID, spotifyShowID string) (*models.Subscription, error)
	Unsubscribe(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) error
}

type showServices struct {
	spotifyClient          *spotify.SpotifyClient
	authService            AuthService
	cfg                    *config.Config
	cache                  *gocache.Cache
	subscriptionRepository repositories.SubscriptionRepository
	showRepository         repositories.ShowRepository
}

func NewShowServices(sc *spotify.SpotifyClient, authService AuthService, cfg *config.Config, cache *gocache.Cache, subscriptionRepository repositories.SubscriptionRepository, showRepository repositories.ShowRepository) ShowServices {
	return &showServices{
		spotifyClient:          sc,
		authService:            authService,
		cfg:                    cfg,
		cache:                  cache,
		subscriptionRepository: subscriptionRepository,
		showRepository:         showRepository,
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
	filtered := make([]dto.SavedShowResponse, 0, len(shows))
	for _, item := range shows {
		if strings.Contains(strings.ToLower(item.Name), query) ||
			strings.Contains(strings.ToLower(item.Description), query) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (s *showServices) SearchShows(ctx context.Context, userID uuid.UUID, query string, limit, offset int) ([]dto.SavedShowResponse, error) {
	accessToken, err := s.authService.GetValidAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	result, err := s.spotifyClient.SearchShows(ctx, accessToken, query, limit, offset)
	if err != nil {
		return nil, err
	}
	shows := result.ToSavedShows()
	return shows, nil
}

// GetSubscriptions returns all of a user's subscriptions.
func (s *showServices) GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	return s.subscriptionRepository.GetByUserID(ctx, userID)
}

func (s *showServices) Subscribe(ctx context.Context, userID uuid.UUID, spotifyShowID string) (*models.Subscription, error) {
	accessToken, err := s.authService.GetValidAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	spotifyShow, err := s.spotifyClient.GetShow(ctx, accessToken, spotifyShowID)
	if err != nil {
		return nil, err
	}

	show := &models.PodcastShow{
		SpotifyShowID: spotifyShow.ID,
		Name:          spotifyShow.Name,
		Description:   spotifyShow.Description,
		ImageURL:      spotifyShow.ImageURL(),
		SpotifyURL:    spotifyShow.ExternalURLs.Spotify,
	}
	if err := s.showRepository.GetOrCreate(ctx, show); err != nil {
		return nil, err
	}

	subscription := &models.Subscription{
		UserID:        userID,
		PodcastShowID: show.ID,
	}
	if err := s.subscriptionRepository.Create(ctx, subscription); err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *showServices) Unsubscribe(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) error {
	return s.subscriptionRepository.Delete(ctx, userID, subscriptionID)
}
