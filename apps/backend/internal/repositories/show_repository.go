package repositories

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShowRepository interface {
	GetOrCreate(ctx context.Context, show *models.PodcastShow) error 
	GetAllTracked(ctx context.Context) ([]models.PodcastShow, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.PodcastShow, error)
	GetBySpotifyID(ctx context.Context, spotifyShowID string) (*models.PodcastShow, error)
	UpdateLatestEpisode(ctx context.Context, showID uuid.UUID, episodeID string, publishedAt interface{}) error
}

type showRepository struct {
	db *gorm.DB
}

func NewShowRepository(db *gorm.DB) ShowRepository {
	return &showRepository{
		db: db,
	}
}

// GetOrCreate finds a show by Spotify ID or creates it if it doesn't exist.
// This is used when a user subscribes to a show — we need to ensure the show
// record exists in our DB before creating the subscription.
func (r *showRepository) GetOrCreate(ctx context.Context, show *models.PodcastShow) error {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: "spotify_show_id",
		}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"description",
			"image_url",
			"spotify_url",
			"latest_episode_id",
			"latest_episode_published_at",
		}),
	}).Create(show).Error
	if err != nil {
		return fmt.Errorf("failed to upsert podcast show: %w", err)
	}
	return nil
}

// GetAllTracked returns all shows that have at least one active subscription.
func (r *showRepository) GetAllTracked(ctx context.Context) ([]models.PodcastShow, error) {
	var shows []models.PodcastShow
	result := r.db.WithContext(ctx).
		Where("id IN (SELECT DISTINCT podcast_show_id FROM subscriptions)").
		Find(&shows)
	if result.Error != nil {
		return nil, result.Error
	}
	return shows, nil
}

// GetByID fetches a podcast show by UUID.
func (r *showRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PodcastShow, error) {
	var show models.PodcastShow
	result := r.db.WithContext(ctx).First(&show, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPodcastShowNotFound
		}
		return nil, result.Error
	}
	return &show, nil
}

// GetBySpotifyID fetches a show by Spotify show ID.
func (r *showRepository) GetBySpotifyID(ctx context.Context, spotifyShowID string) (*models.PodcastShow, error) {
	var show models.PodcastShow
	result := r.db.WithContext(ctx).First(&show, "spotify_show_id = ?", spotifyShowID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPodcastShowNotFound
		}
		return nil, result.Error
	}
	return &show, nil
}

// UpdateLatestEpisode updates the show's record of what the newest episode is.
func (r *showRepository) UpdateLatestEpisode(ctx context.Context, showID uuid.UUID, episodeID string, publishedAt interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&models.PodcastShow{}).
		Where("id = ?", showID).
		Updates(map[string]interface{}{
			"latest_episode_id":           episodeID,
			"latest_episode_published_at": publishedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update latest episode: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrPodcastShowNotFound
	}
	return nil
}
