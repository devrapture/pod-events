package repositories

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	appcrypto "github.com/devrapture/pod-events/pkg/crypto"
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
	db            *gorm.DB
	encryptionKey string
}

func NewChannelRepository(db *gorm.DB, encryptionKey string) ChannelRepository {
	return &channelRepository{
		db:            db,
		encryptionKey: encryptionKey,
	}
}

func (r *channelRepository) Create(
	ctx context.Context,
	channel *models.NotificationChannel,
) error {
	toSave := *channel

	if r.isWebhookChannel(toSave.ChannelType) {
		fingerprint, err := appcrypto.FingerprintText(
			toSave.Destination,
			r.encryptionKey,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to fingerprint webhook: %w",
				err,
			)
		}

		encryptedWebhook, err := r.encryptWebhook(
			toSave.Destination,
		)
		if err != nil {
			return err
		}

		toSave.Destination = encryptedWebhook
		toSave.DestinationFingerprint = fingerprint
	}

	result := dbFromCtx(ctx, r.db).
		WithContext(ctx).
		Create(&toSave)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return apperrors.ErrNotificationChannelAlreadyExists
		}
		return fmt.Errorf(
			"failed to create notification channel: %w",
			result.Error,
		)
	}

	channel.Base = toSave.Base
	return nil
}

func (r *channelRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.NotificationChannel, error) {
	var channel []models.NotificationChannel
	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Where("user_id = ?", userID).Find(&channel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get notification channel: %w", err)
	}
	return channel, nil
}

func (r *channelRepository) ToggleActive(ctx context.Context, userID, channelID uuid.UUID, isActive bool) error {
	result := dbFromCtx(ctx, r.db).WithContext(ctx).
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
	result := dbFromCtx(ctx, r.db).WithContext(ctx).Where("id = ? AND user_id = ?", channelID, userID).Delete(&models.NotificationChannel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification channel: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrChannelIDNotFound
	}
	return nil
}

func (r *channelRepository) isWebhookChannel(channelType models.ChannelType) bool {
	return channelType == models.ChannelTypeDiscord || channelType == models.ChannelTypeSlack
}

func (r *channelRepository) encryptWebhook(webhook string) (string, error) {
	encrypted, err := appcrypto.EncryptText(webhook, r.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt webhook: %w", err)
	}
	return encrypted, nil
}
