package repositories

import (
	"context"
	"fmt"
	"time"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"gorm.io/gorm"
)

type TelegramConnectionRepository interface {
	Create(ctx context.Context, conn *models.TelegramConnection) error
	DeleteExpired(ctx context.Context) error
	Consume(ctx context.Context, tokenHash string) (*models.TelegramConnection, error)
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

func (r *telegramConnectionRepository) Consume(ctx context.Context, tokenHash string) (*models.TelegramConnection, error) {
	var t models.TelegramConnection
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&t, "token_hash = ?", tokenHash).Error; err != nil {
			return err
		}
		if t.Consumed || time.Now().After(t.ExpiresAt) {
			tx.Delete(&t)
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&t).Update("consumed", true).Error
	})

	if err != nil {
		return nil, fmt.Errorf("failed to consume telegram connection: %w", err)
	}
	return &t, nil
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
