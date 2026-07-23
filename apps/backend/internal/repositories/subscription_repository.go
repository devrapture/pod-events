package repositories

import (
	"context"
	"errors"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	CreateBatch(ctx context.Context, subscriptions []models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error)
	GetSubscribedUsersByShowID(ctx context.Context, showID uuid.UUID) ([]models.User, error)
	GetByUserAndShow(ctx context.Context, userID, showID uuid.UUID) (*models.Subscription, error)
	Delete(ctx context.Context, userID, subscriptionID uuid.UUID) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

// Create creates a new subscription.
// Returns an error if the user is already subscribed to this show.
func (r *subscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	existing, err := r.GetByUserAndShow(ctx, subscription.UserID, subscription.PodcastShowID)

	if err != nil {
		return err
	}
	if existing != nil {
		return apperrors.ErrSubscriptionAlreadyExists
	}

	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Create(subscription).Error; err != nil {
		if isUniqueViolation(err) {
			return apperrors.ErrSubscriptionAlreadyExists
		}
		return err
	}
	return nil
}

// CreateBatch creates all subscriptions in one atomic insert.
func (r *subscriptionRepository) CreateBatch(ctx context.Context, subscriptions []models.Subscription) error {
	if len(subscriptions) == 0 {
		return nil
	}

	if err := dbFromCtx(ctx, r.db).WithContext(ctx).Create(&subscriptions).Error; err != nil {
		if isUniqueViolation(err) {
			return apperrors.ErrSubscriptionAlreadyExists
		}
		return err
	}

	return nil
}

// GetByID fetches a subscription by UUID.
func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	result := dbFromCtx(ctx, r.db).WithContext(ctx).Where("id = ?", id).Preload("PodcastShow").First(&subscription)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &subscription, nil
}

// GetByUserID returns all subscriptions for a user, with show data loaded.
func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	var subscription []models.Subscription
	result := dbFromCtx(ctx, r.db).WithContext(ctx).Where("user_id = ?", userID).Preload("PodcastShow").Find(&subscription)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return subscription, nil
}

// GetSubscribedUsersByShowID returns all users subscribed to a given show.
// Used by the cron job to know who to notify about a new episode.
func (r *subscriptionRepository) GetSubscribedUsersByShowID(ctx context.Context, showID uuid.UUID) ([]models.User, error) {
	var users []models.User
	result := dbFromCtx(ctx, r.db).WithContext(ctx).Joins("JOIN subscriptions ON subscriptions.user_id = users.id").Where("subscriptions.podcast_show_id = ?", showID).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetByUserAndShow checks if a specific user is subscribed to a specific show.
// Returns (nil, nil) if not subscribed.
func (r *subscriptionRepository) GetByUserAndShow(ctx context.Context, userID, showID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	result := dbFromCtx(ctx, r.db).WithContext(ctx).Where("podcast_show_id = ? AND user_id = ?", showID, userID).First(&subscription)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &subscription, nil
}

// Delete removes a subscription by ID, verifying the user owns it.
func (r *subscriptionRepository) Delete(ctx context.Context, userID, subscriptionID uuid.UUID) error {
	result := dbFromCtx(ctx, r.db).WithContext(ctx).Where("user_id = ? AND id = ?", userID, subscriptionID).Unscoped().Delete(&models.Subscription{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrSubscriptionNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
