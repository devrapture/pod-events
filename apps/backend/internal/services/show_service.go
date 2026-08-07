package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/devrapture/pod-events/internal/dto"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/errgroup"
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
	cache                  *gocache.Cache
	subscriptionRepository repositories.SubscriptionRepository
	showRepository         repositories.ShowRepository
}

const (
	maxSpotifyShowIDs                = 50
	maxConcurrentSpotifyShowRequests = 8
	spotifyShowCachePrefix           = "spotify:show:"
)

func NewShowServices(sc *spotify.SpotifyClient, authService AuthService, cache *gocache.Cache, subscriptionRepository repositories.SubscriptionRepository, showRepository repositories.ShowRepository) ShowServices {
	return &showServices{
		spotifyClient:          sc,
		authService:            authService,
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
	s.cacheSavedShows(savedShows)

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
	s.cacheSearchResults(result)

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
	showsData, err := s.fetchShowsFromSpotify(ctx, userID, uniqueIDs)
	if err != nil {
		if errors.Is(err, apperrors.ErrSpotifyResourceNotFound) {
			return nil, apperrors.ErrPodcastShowNotFound
		}
		return nil, fmt.Errorf("%w: %w", apperrors.ErrSpotifyUnavailable, err)
	}

	showModels := make([]*models.PodcastShow, 0, len(showsData))

	for _, sd := range showsData {
		showModels = append(showModels, &models.PodcastShow{
			SpotifyShowID: sd.ID,
			Name:          sd.Name,
			Description:   sd.Description,
			ImageURL:      sd.ImageURL(),
			SpotifyURL:    sd.ExternalURLs.Spotify,
		})
	}

	if err := s.showRepository.BatchGetOrCreate(ctx, showModels); err != nil {
		return nil, err
	}

	// Seed latest-episode watermark for first-time shows so the next cron run
	// does not treat the show's current (already-published) episode as new.
	// Only writes when latest_episode_id is still null.
	if err := s.seedLatestEpisodesIfNull(ctx, userID, showModels); err != nil {
		return nil, err
	}

	subscriptions := make([]models.Subscription, 0, len(showModels))
	for _, show := range showModels {
		subscriptions = append(subscriptions, models.Subscription{
			UserID:        userID,
			PodcastShowID: show.ID,
			PodcastShow:   *show,
		})
	}

	if err := s.subscriptionRepository.CreateBatch(ctx, subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// seedLatestEpisodesIfNull fetches each show's current latest episode from Spotify and
// stores it as the watermark when latest_episode_id is null. Shows that already have a
// watermark are skipped. Shows with no episodes are left null (cron handles that case).
func (s *showServices) seedLatestEpisodesIfNull(ctx context.Context, userID uuid.UUID, shows []*models.PodcastShow) error {
	needsSeed := make([]*models.PodcastShow, 0, len(shows))
	for _, show := range shows {
		if show.LatestEpisodeID == nil {
			needsSeed = append(needsSeed, show)
		}
	}
	if len(needsSeed) == 0 {
		return nil
	}

	accessToken, err := s.authService.GetValidAccessToken(ctx, userID)
	if err != nil {
		return err
	}

	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(maxConcurrentSpotifyShowRequests)

	for _, show := range needsSeed {
		show := show
		group.Go(func() error {
			latest, err := s.spotifyClient.GetShowLatestEpisode(groupCtx, accessToken, show.SpotifyShowID)
			if err != nil {
				if errors.Is(err, apperrors.ErrSpotifyResourceNotFound) {
					// Show has no episodes yet — leave watermark null.
					return nil
				}
				var rateLimitErr *spotify.RateLimitError
				if errors.As(err, &rateLimitErr) {
					return err
				}
				if errors.Is(err, apperrors.ErrSpotifyAuthorizationRequired) {
					return err
				}
				return fmt.Errorf("%w: %w", apperrors.ErrSpotifyUnavailable, err)
			}
			if latest == nil {
				return nil
			}

			var publishedAt interface{}
			var publishedAtPtr *time.Time
			if parsed, parseErr := latest.ParsedReleaseDate(); parseErr == nil {
				publishedAt = parsed
				t := parsed
				publishedAtPtr = &t
			}

			if err := s.showRepository.UpdateLatestEpisodeIfNull(groupCtx, show.ID, latest.ID, publishedAt); err != nil {
				return err
			}

			episodeID := latest.ID
			show.LatestEpisodeID = &episodeID
			show.LatestEpisodePublishedAt = publishedAtPtr
			return nil
		})
	}

	return group.Wait()
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

func (s *showServices) fetchShowsFromSpotify(ctx context.Context, userID uuid.UUID, showIDs []string) ([]*spotify.SpotifyShow, error) {
	result := make([]*spotify.SpotifyShow, len(showIDs))
	missingIndexes := make([]int, 0, len(showIDs))

	for i, spotifyShowID := range showIDs {
		if cachedShow, found := s.cachedSpotifyShow(spotifyShowID); found {
			result[i] = cachedShow
			continue
		}
		missingIndexes = append(missingIndexes, i)
	}

	if len(missingIndexes) == 0 {
		return result, nil
	}

	accessToken, err := s.authService.GetValidAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(maxConcurrentSpotifyShowRequests)

	for _, index := range missingIndexes {
		index := index
		spotifyShowID := showIDs[index]

		group.Go(func() error {
			spotifyShow, err := s.spotifyClient.GetShow(groupCtx, accessToken, spotifyShowID)
			if err != nil {
				if errors.Is(err, apperrors.ErrSpotifyResourceNotFound) {
					return apperrors.ErrPodcastShowNotFound
				}
				return err
			}

			result[index] = spotifyShow
			s.cacheSpotifyShow(*spotifyShow)
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *showServices) cacheSavedShows(response *spotify.SpotifySavedShowsResponse) {
	for _, item := range response.Items {
		s.cacheSpotifyShow(item.Show)
	}
}

func (s *showServices) cacheSearchResults(result *spotify.ShowSearchResult) {
	for _, show := range result.Shows.Items {
		s.cacheSpotifyShow(show)
	}
}

func (s *showServices) cacheSpotifyShow(show spotify.SpotifyShow) {
	s.cache.Set(spotifyShowCachePrefix+show.ID, show, gocache.DefaultExpiration)
}

func (s *showServices) cachedSpotifyShow(spotifyShowID string) (*spotify.SpotifyShow, bool) {
	value, found := s.cache.Get(spotifyShowCachePrefix + spotifyShowID)
	if !found {
		return nil, false
	}

	show, ok := value.(spotify.SpotifyShow)
	if !ok {
		s.cache.Delete(spotifyShowCachePrefix + spotifyShowID)
		return nil, false
	}

	return &show, true
}
