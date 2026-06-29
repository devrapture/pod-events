package repositories

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChannelRepository interface {
	Create(ctx context.Context, channel *models.NotificationChannel) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.NotificationChannel, error)
	ToggleActive(ctx context.Context, userID, channelID uuid.UUID, isActive bool) error
	Delete(ctx context.Context, userID, channelID uuid.UUID) error
}

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{
		db: db,
	}
}

func (r *channelRepository) Create(ctx context.Context, channel *models.NotificationChannel) error {
	result := r.db.WithContext(ctx).Create(channel)
	if result.Error != nil {
		return fmt.Errorf("failed to create notification channel: %w", result.Error)
	}
	return nil
}

func (r *channelRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.NotificationChannel, error) {
	var channel []models.NotificationChannel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&channel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get notification channel: %w", err)
	}
	return channel, nil
}

func (r *channelRepository) ToggleActive(ctx context.Context, userID, channelID uuid.UUID, isActive bool) error {
	result := r.db.WithContext(ctx).
		Model(&models.NotificationChannel{}).
		Where("id = ? AND user_id = ?", channelID, userID).
		Update("is_active", isActive)
	if result.Error != nil {
		return fmt.Errorf("failed to toggle notification channel: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrChannelIDNotFound
	}
	return nil
}

func (r *channelRepository) Delete(ctx context.Context, userID, channelID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", channelID, userID).Delete(&models.NotificationChannel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification channel: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrChannelIDNotFound
	}
	return nil
}
