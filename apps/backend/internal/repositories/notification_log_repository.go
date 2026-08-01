package repositories

import (
	"context"
	"fmt"

	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationLogRepository interface {
	Create(ctx context.Context, log *models.NotificationLog) error
	AlreadySent(ctx context.Context, userID, episodeID uuid.UUID, channelType models.ChannelType) (bool, error)
}

type notificationLogRepository struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) NotificationLogRepository {
	return &notificationLogRepository{
		db: db,
	}
}

func (r *notificationLogRepository) Create(ctx context.Context, log *models.NotificationLog) error {
	result := dbFromCtx(ctx, r.db).WithContext(ctx).Create(log)
	if result.Error != nil {
		return fmt.Errorf("failed to create notification log: %w", result.Error)
	}
	return nil
}

func (r *notificationLogRepository) AlreadySent(ctx context.Context, userID, episodeID uuid.UUID, channelType models.ChannelType) (bool, error) {
	var count int64
	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Model(&models.NotificationLog{}).
		Where("user_id = ? AND episode_id = ? AND channel_type = ? AND status = ?", userID, episodeID, channelType, models.NotificationStatusSent).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if notification has already been sent: %w", err)
	}
	return count > 0, nil
}
