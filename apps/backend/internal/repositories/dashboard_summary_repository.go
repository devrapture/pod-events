package repositories

import (
	"context"
	"time"

	"github.com/devrapture/pod-events/internal/dto"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DashboardSummaryRepository interface {
	GetDashboardSummary(ctx context.Context, userID uuid.UUID) (*dto.DashboardSummaryDTO, error)
}

type dasboardSummaryRepository struct {
	db *gorm.DB
}

func NewDashboardSummaryRepository(db *gorm.DB) DashboardSummaryRepository {
	return &dasboardSummaryRepository{
		db: db,
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

	setupItems := []dto.DashboardItems{
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
		Model(&models.Episode{}).Select(`
		episodes.name AS episode_title,
		podcast_shows.name AS podcast_name,
		episodes.release_date AS release_date,
		episodes.duration_ms AS duration,
		episodes.spotify_url AS url
	`).Joins("JOIN subscriptions ON subscriptions.podcast_show_id = episodes.podcast_show_id AND subscriptions.deleted_at IS NULL").
		Joins("JOIN podcast_shows ON podcast_shows.id = episodes.podcast_show_id").
		Where("subscriptions.user_id = ?", userID).
		Order("episodes.release_date DESC").
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
