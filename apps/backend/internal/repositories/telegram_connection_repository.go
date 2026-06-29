package repositories

import (
	"context"
	"fmt"
	"time"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&t).
		Clauses(clause.Returning{}).
		Where("token_hash = ? AND consumed = ? AND expires_at > ?", tokenHash, false, now).
		Update("consumed", true)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to consume telegram connection: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		if err := r.db.WithContext(ctx).
			Unscoped().
			Where("token_hash = ? AND (consumed = ? OR expires_at <= ?)", tokenHash, true, now).
			Delete(&models.TelegramConnection{}).Error; err != nil {
			return nil, fmt.Errorf("failed to delete stale telegram connection: %w", err)
		}
		return nil, fmt.Errorf("failed to consume telegram connection: %w", gorm.ErrRecordNotFound)
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
