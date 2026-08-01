package cron

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EpisodeChecker struct {
	subRepo       repositories.SubscriptionRepository
	episodeRepo   repositories.EpisodeRepository
	showRepo      repositories.ShowRepository
	authService   services.AuthService
	notifService  services.NotificationService
	logger        *zap.Logger
	spotifyClient *spotify.SpotifyClient
}

func NewEpisodeChecker(
	subRepo repositories.SubscriptionRepository,
	episodeRepo repositories.EpisodeRepository,
	showRepo repositories.ShowRepository,
	authService services.AuthService,
	notifService services.NotificationService,
	logger *zap.Logger,
	spotifyClient *spotify.SpotifyClient,
) *EpisodeChecker {
	return &EpisodeChecker{
		subRepo:       subRepo,
		episodeRepo:   episodeRepo,
		showRepo:      showRepo,
		authService:   authService,
		notifService:  notifService,
		logger:        logger,
		spotifyClient: spotifyClient,
	}
}

// CheckResult summarizes the outcome of a cron run.
type CheckResult struct {
	ShowsChecked      int
	NewEpisodes       int
	NotificationsSent int
	Errors            []string
}

func (c *EpisodeChecker) Run(ctx context.Context) (*CheckResult, error) {
	result := &CheckResult{}
	startTime := time.Now()
	c.logger.Info("cron job started: checking for new episodes")
	shows, err := c.showRepo.GetAllTracked(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tracked shows: %w", err)
	}
	c.logger.Info("found tracked shows", zap.Int("count", len(shows)))
	for _, show := range shows {
		if err := ctx.Err(); err != nil {
			return result, err
		}

		result.ShowsChecked++
		if err := c.checkShow(ctx, &show, result); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return result, ctxErr
			}
			c.logger.Error(
				"error checking show",
				zap.String("show_id", show.ID.String()),
				zap.String("show_name", show.Name), zap.Error(err),
			)
			result.Errors = append(result.Errors, fmt.Sprintf("show %s: %v", show.Name, err))
		}

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
	c.logger.Info(
		"cron job completed",
		zap.Int64("duration", time.Since(startTime).Milliseconds()),
		zap.Int("shows_checked", result.ShowsChecked),
		zap.Int("notifications_sent", result.NotificationsSent),
		zap.Int("new_episodes", result.NewEpisodes),
		zap.Int("errors", len(result.Errors)),
	)
	return result, nil
}

func (c *EpisodeChecker) checkShow(ctx context.Context, show *models.PodcastShow, result *CheckResult) error {
	accessToken, err := c.getAccessTokenForShow(ctx, show.ID)
	if err != nil {
		return fmt.Errorf("get access token: %w", err)
	}
	latestEpisode, err := c.spotifyClient.GetShowLatestEpisode(ctx, accessToken, show.SpotifyShowID)
	if err != nil {
		var rateLimiterErr *spotify.RateLimitError
		if errors.As(err, &rateLimiterErr) {
			c.logger.Warn(
				"rate limited by spotify",
				zap.String("show", show.Name),
				zap.Int("retry_after", rateLimiterErr.RetryAfter),
			)
			return nil
		}
		return fmt.Errorf("fetching latest episode from spotify: %w", err)
	}
	if latestEpisode == nil {
		c.logger.Info("show has no episodes on spotify", zap.String("show", show.Name))
		return nil
	}

	if show.LatestEpisodeID != nil && *show.LatestEpisodeID == latestEpisode.ID {
		c.logger.Info("no new episode", zap.String("show", show.Name))
		return nil
	}
	c.logger.Info(
		"new episode detected",
		zap.String("show", show.Name),
		zap.String("episode", latestEpisode.Name),
		zap.String("spotify_episode_id", latestEpisode.ID),
	)

	publishedAt, err := latestEpisode.ParsedReleaseDate()
	if err != nil {
		publishedAt = time.Now().UTC()
	}
	existingEpisode, err := c.episodeRepo.GetBySpotifyID(ctx, latestEpisode.ID)
	if err != nil {
		return fmt.Errorf("check existing episode: %w", err)
	}

	var episode *models.Episode
	if existingEpisode != nil {
		// Episode already in DB (from a previous failed run perhaps)
		episode = existingEpisode
	} else {
		episode = &models.Episode{
			SpotifyEpisodeID: latestEpisode.ID,
			PodcastShowID:    show.ID,
			Name:             latestEpisode.Name,
			Description:      latestEpisode.Description,
			SpotifyURL:       latestEpisode.ExternalURLs.Spotify,
			ImageURL:         latestEpisode.Images[0].URL,
			DurationMs:       latestEpisode.DurationMs,
			ReleaseDate:      publishedAt,
		}
		if err := c.episodeRepo.Create(ctx, episode); err != nil {
			return fmt.Errorf("save new episode: %w", err)
		}
		result.NewEpisodes++
	}

	if err := c.showRepo.UpdateLatestEpisode(ctx, show.ID, latestEpisode.ID, publishedAt); err != nil {
		c.logger.Error("failed to update show latest episode", zap.Error(err))
	}

	subscribers, err := c.subRepo.GetSubscribedUsersByShowID(ctx, show.ID)
	if err != nil {
		return fmt.Errorf("get subscribers: %w", err)
	}

	for _, user := range subscribers {
		if err := c.notifService.NotifyUser(ctx, user.ID, episode, show); err != nil {
			c.logger.Error(
				"failed to notify user",
				zap.String("user_id", user.ID.String()),
				zap.String("episode_id", latestEpisode.ID),
				zap.Error(err),
			)
		} else {
			result.NotificationsSent++
		}
	}
	return nil
}

func (c *EpisodeChecker) getAccessTokenForShow(ctx context.Context, showID uuid.UUID) (string, error) {
	subscribers, err := c.subRepo.GetSubscribedUsersByShowID(ctx, showID)
	if err != nil {
		return "", fmt.Errorf("get subscribers: %w", err)
	}

	if len(subscribers) == 0 {
		return "", fmt.Errorf("show has no subscribers")
	}

	for _, user := range subscribers {
		token, err := c.authService.GetValidAccessToken(ctx, user.ID)
		if err != nil {
			c.logger.Warn(
				"could not get token for subscriber",
				zap.String("user_id", user.ID.String()),
				zap.Error(err),
			)
			continue
		}
		return token, nil
	}
	return "", fmt.Errorf("no subscriber has a valid spotify token")
}
