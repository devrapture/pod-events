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

	if err := r.db.WithContext(ctx).Model(&models.Subscription{}).Where("user_id = ?", userID).Count(&podcastTracked).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&models.NotificationChannel{}).Where("user_id = ? AND is_active = true", userID).Count(&activeChannels).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&models.Episode{}).Joins("JOIN subscriptions ON subscriptions.podcast_show_id = episodes.podcast_show_id").Where("subscriptions.user_id = ? AND episodes.created_at >= ?", userID, weekStart).Count(&newEpisodesThisWeek).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&models.NotificationLog{}).Where("user_id = ? AND status = ?", userID, models.NotificationStatusSent).Count(&notificationSent).Error; err != nil {
		return nil, err
	}

	setupItems := []dto.DashboardItems{
		{
			Label:     "Connect Spotify",
			Completed: true,
		},
		{
			Label:     "Subscribe to at least one podcast",
			Completed: podcastTracked > 0,
		},
		{
			Label:     "Add a notification channel",
			Completed: activeChannels > 0,
		},
		{
			Label:     "New Episodes This Week",
			Completed: newEpisodesThisWeek > 0,
		},
		{
			Label:     "Receive first notification",
			Completed: notificationSent > 0,
		},
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
		// percent = (completed / total) * 100
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
	}, nil
}
