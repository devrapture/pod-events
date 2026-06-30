package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TelegramConnectionRepository interface {
	Create(ctx context.Context, conn *models.TelegramConnection) error
	DeleteExpired(ctx context.Context) error
	CompleteConnection(ctx context.Context, tokenHash string, chatID int64) (*models.NotificationChannel, error)
}

type telegramConnectionRepository struct {
	db *gorm.DB
}

func NewTelegramConnectionRepository(db *gorm.DB) TelegramConnectionRepository {
	return &telegramConnectionRepository{
		db: db,
	}
}

func (r *telegramConnectionRepository) Create(ctx context.Context, conn *models.TelegramConnection) error {
	result := r.db.WithContext(ctx).Create(conn)
	if result.Error != nil {
		return fmt.Errorf("failed to create telegram connection: %w", result.Error)
	}
	return nil
}

func (r *telegramConnectionRepository) CompleteConnection(ctx context.Context, tokenHash string, chatID int64) (*models.NotificationChannel, error) {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("token_hash = ? AND (consumed = ? OR expires_at <= ?)", tokenHash, true, now).
		Delete(&models.TelegramConnection{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete stale telegram connection: %w", err)
	}

	var channel models.NotificationChannel
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var conn models.TelegramConnection
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ? AND consumed = ? AND expires_at > ?", tokenHash, false, now).
			First(&conn).Error; err != nil {
			return err
		}

		var existing int64
		if err := tx.Model(&models.NotificationChannel{}).
			Where("user_id = ? AND channel_type = ?", conn.UserID, models.ChannelTypeTelegram).
			Count(&existing).Error; err != nil {
			return fmt.Errorf("failed to check existing telegram channel: %w", err)
		}
		if existing > 0 {
			return apperrors.ErrTelegramChannelAlreadyExists
		}

		channel = models.NotificationChannel{
			UserID:      conn.UserID,
			ChannelType: models.ChannelTypeTelegram,
			Destination: strconv.FormatInt(chatID, 10),
			IsActive:    true,
		}
		if err := tx.Create(&channel).Error; err != nil {
			return fmt.Errorf("failed to create notification channel: %w", err)
		}

		return tx.Model(&conn).Update("consumed", true).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to complete telegram connection: %w", err)
	}
	return &channel, nil
}

func (r *telegramConnectionRepository) DeleteExpired(ctx context.Context) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("expires_at < ? OR consumed = ?", time.Now(), true).
		Delete(&models.TelegramConnection{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired telegram connections: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNoExpiredTelegramConnectionsFound
	}
	return nil
}
