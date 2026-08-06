package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/devrapture/pod-events/internal/dto"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DashboardSummaryRepository interface {
	GetDashboardSummary(ctx context.Context, userID uuid.UUID) (*dto.DashboardSummaryDTO, error)
}

type dasboardSummaryRepository struct {
	db        *gorm.DB
	tokenRepo TokenRepository
}

func NewDashboardSummaryRepository(db *gorm.DB, tokenRepo TokenRepository) DashboardSummaryRepository {
	return &dasboardSummaryRepository{
		db:        db,
		tokenRepo: tokenRepo,
	}
}

func (r *dasboardSummaryRepository) GetDashboardSummary(ctx context.Context, userID uuid.UUID) (*dto.DashboardSummaryDTO, error) {
	weekStart := time.Now().AddDate(0, 0, -7)

	var podcastTracked int64
	var activeChannels int64
	var newEpisodesThisWeek int64
	var notificationSent int64
	recentEpisodesResponse := make([]dto.RecentEpisodesResponse, 0)

	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Model(&models.Subscription{}).Where("user_id = ?", userID).Count(&podcastTracked).Error; err != nil {
		return nil, err
	}

	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Model(&models.NotificationChannel{}).Where("user_id = ? AND is_active = true", userID).Count(&activeChannels).Error; err != nil {
		return nil, err
	}

	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Model(&models.Episode{}).Joins("JOIN subscriptions ON subscriptions.podcast_show_id = episodes.podcast_show_id").Where("subscriptions.user_id = ? AND episodes.created_at >= ?", userID, weekStart).Count(&newEpisodesThisWeek).Error; err != nil {
		return nil, err
	}

	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Model(&models.NotificationLog{}).Where("user_id = ? AND status = ?", userID, models.NotificationStatusSent).Count(&notificationSent).Error; err != nil {
		return nil, err
	}

	hasSpotifyToken := false
	_, err := r.tokenRepo.GetByUserID(ctx, userID)
	if err == nil {
		hasSpotifyToken = true
	} else if !errors.Is(err, apperrors.ErrorSpotifyTokenNotFound) {
		return nil, err
	}

	setupItems := []dto.DashboardItems{
		{
			Key:       "connect_spotify",
			Label:     "Connect Spotify",
			Completed: hasSpotifyToken,
		},
		{
			Key:       "import_spotify",
			Label:     "Import Spotify Podcasts",
			Completed: hasSpotifyToken,
		},
		{
			Key:       "subscribe_podcast",
			Label:     "Subscribe to at least one podcast",
			Completed: podcastTracked > 0,
		},
		{
			Key:       "add_channel",
			Label:     "Add a notification channel",
			Completed: activeChannels > 0,
		},
		{
			Key:       "receive_notification",
			Label:     "Receive first notification",
			Completed: notificationSent > 0,
		},
	}

	if err := dbFromCtx(ctx, r.db).WithContext(ctx).
		Model(&models.NotificationLog{}).Select(`
		episodes.name AS episode_title,
		podcast_shows.name AS podcast_name,
		episodes.release_date AS release_date,
		episodes.duration_ms AS duration,
		episodes.spotify_url AS url
	`).Joins("JOIN episodes ON episodes.id = notification_logs.episode_id").
		Joins("JOIN podcast_shows ON podcast_shows.id = episodes.podcast_show_id").
		Where(
			"notification_logs.user_id = ? AND notification_logs.status = ?",
			userID,
			models.NotificationStatusSent,
		).
		Group("episodes.id, podcast_shows.id").
		Order("MAX(notification_logs.sent_at) DESC").
		Limit(5).
		Scan(&recentEpisodesResponse).Error; err != nil {
		return nil, err
	}
	completed := 0
	total := len(setupItems)
	percent := 0

	for _, item := range setupItems {
		if item.Completed {
			completed++
		}
	}
	if total > 0 {
		percent = int(float64(completed) / float64(total) * 100)
	}

	return &dto.DashboardSummaryDTO{
		Setup: dto.DashboardSetupResponse{
			Completed: completed,
			Total:     total,
			Percent:   percent,
			Items:     setupItems,
		},
		Stats: dto.DashboardStatsResponse{
			PodcastsTracked:     podcastTracked,
			ActiveChannels:      activeChannels,
			NewEpisodesThisWeek: newEpisodesThisWeek,
			NotificationSent:    notificationSent,
		},
		RecentEpisodes: recentEpisodesResponse,
	}, nil
}
