package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/database"
	"github.com/devrapture/pod-events/internal/dto"
	apperrors "github.com/devrapture/pod-events/internal/errors"
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
	Subscribe(ctx context.Context, userID uuid.UUID, spotifyShowIDs []string) ([]models.Subscription, error)
	Unsubscribe(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) error
}

type showServices struct {
	spotifyClient          *spotify.SpotifyClient
	authService            AuthService
	cfg                    *config.Config
	cache                  *gocache.Cache
	subscriptionRepository repositories.SubscriptionRepository
	showRepository         repositories.ShowRepository
	txManager              database.TransactionManager
}

const maxSpotifyShowIDs = 50

func NewShowServices(sc *spotify.SpotifyClient, authService AuthService, cfg *config.Config, cache *gocache.Cache, subscriptionRepository repositories.SubscriptionRepository, showRepository repositories.ShowRepository, txManager database.TransactionManager) ShowServices {
	return &showServices{
		spotifyClient:          sc,
		authService:            authService,
		cfg:                    cfg,
		cache:                  cache,
		subscriptionRepository: subscriptionRepository,
		showRepository:         showRepository,
		txManager:              txManager,
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

func (s *showServices) Subscribe(ctx context.Context, userID uuid.UUID, spotifyShowIDs []string) ([]models.Subscription, error) {
	uniqueIDs := cleanShowIDs(spotifyShowIDs)

	if len(uniqueIDs) == 0 {
		return nil, apperrors.ErrInvalidSpotifyShowIDs
	}

	if len(uniqueIDs) > maxSpotifyShowIDs {
		return nil, fmt.Errorf("%w: maximum is %d, got %d", apperrors.ErrTooManySpotifyShowIDs, maxSpotifyShowIDs, len(uniqueIDs))
	}
	accessToken, err := s.authService.GetValidAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	showsData, err := s.fetchShowsFromSpotify(ctx, accessToken, uniqueIDs)
	if err != nil {
		if errors.Is(err, apperrors.ErrSpotifyResourceNotFound) {
			return nil, apperrors.ErrPodcastShowNotFound
		}
		return nil, fmt.Errorf("%w: %w", apperrors.ErrSpotifyUnavailable, err)
	}

	modelsData := make([]showModelData, 0, len(showsData))

	for _, sd := range showsData {
		modelsData = append(modelsData, showModelData{
			model: &models.PodcastShow{
				SpotifyShowID: sd.spotifyShow.ID,
				Name:          sd.spotifyShow.Name,
				Description:   sd.spotifyShow.Description,
				ImageURL:      sd.spotifyShow.ImageURL(),
				SpotifyURL:    sd.spotifyShow.ExternalURLs.Spotify,
			},
			spotifyShowID: sd.originalID,
		})
	}

	showModels := make([]*models.PodcastShow, len(modelsData))
	for i, md := range modelsData {
		showModels[i] = md.model
	}

	if err := s.showRepository.BatchGetOrCreate(ctx, showModels); err != nil {
		return nil, err
	}

	var subscriptions []models.Subscription

	txErr := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, md := range modelsData {
			show := md.model

			subscription := &models.Subscription{
				UserID:        userID,
				PodcastShowID: show.ID,
			}

			if err := s.subscriptionRepository.Create(txCtx, subscription); err != nil {
				return err
			}
			subscriptions = append(subscriptions, *subscription)
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return subscriptions, nil
}

func (s *showServices) Unsubscribe(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) error {
	return s.subscriptionRepository.Delete(ctx, userID, subscriptionID)
}

func cleanShowIDs(showIDs []string) []string {
	seen := make(map[string]struct{}, len(showIDs))
	unique := make([]string, 0, len(showIDs))

	for _, id := range showIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique
}

type showModelData struct {
	model         *models.PodcastShow
	spotifyShowID string
}

type spotifyShowData struct {
	spotifyShow *spotify.SpotifyShow
	originalID  string
}

func (s *showServices) fetchShowsFromSpotify(ctx context.Context, accessToken string, showIDs []string) ([]spotifyShowData, error) {
	result := make([]spotifyShowData, 0, len(showIDs))
	for _, spotifyShowID := range showIDs {
		spotifyShow, err := s.spotifyClient.GetShow(ctx, accessToken, spotifyShowID)
		if err != nil {
			if errors.Is(err, apperrors.ErrSpotifyResourceNotFound) {
				return nil, apperrors.ErrPodcastShowNotFound
			}
			return nil, err
		}
		result = append(result, spotifyShowData{spotifyShow: spotifyShow, originalID: spotifyShowID})
	}
	return result, nil
}
