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
	BatchGetOrCreate(ctx context.Context, shows []*models.PodcastShow) error
	GetOrCreate(ctx context.Context, show *models.PodcastShow) error
	GetAllTracked(ctx context.Context) ([]models.PodcastShow, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.PodcastShow, error)
	GetBySpotifyID(ctx context.Context, spotifyShowID string) (*models.PodcastShow, error)
	UpdateLatestEpisode(ctx context.Context, showID uuid.UUID, episodeID string, publishedAt interface{}) error
	// UpdateLatestEpisodeIfNull sets the latest-episode watermark only when it is currently null.
	// RowsAffected may be 0 if another writer already seeded the show; that is not an error.
	UpdateLatestEpisodeIfNull(ctx context.Context, showID uuid.UUID, episodeID string, publishedAt interface{}) error
}

type showRepository struct {
	db *gorm.DB
}

func NewShowRepository(db *gorm.DB) ShowRepository {
	return &showRepository{
		db: db,
	}
}

func (r *showRepository) BatchGetOrCreate(ctx context.Context, shows []*models.PodcastShow) error {
	if len(shows) == 0 {
		return nil
	}
	err := dbFromCtx(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: "spotify_show_id",
		}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"description",
			"image_url",
			"spotify_url",
		}),
	}).Create(shows).Error

	if err != nil {
		return fmt.Errorf("failed to upsert podcast show: %w", err)
	}

	spotifyIDs := make([]string, len(shows))
	for i, s := range shows {
		spotifyIDs[i] = s.SpotifyShowID
	}

	var existing []models.PodcastShow
	if err := dbFromCtx(ctx, r.db).WithContext(ctx).
		Where("spotify_show_id IN ?", spotifyIDs).
		Find(&existing).Error; err != nil {
		return fmt.Errorf("failed to refetch podcast shows: %w", err)
	}

	bySpotifyID := make(map[string]models.PodcastShow, len(existing))
	for _, e := range existing {
		bySpotifyID[e.SpotifyShowID] = e
	}

	for _, s := range shows {
		full, ok := bySpotifyID[s.SpotifyShowID]
		if !ok {
			return fmt.Errorf("podcast show %s missing after upsert", s.SpotifyShowID)
		}
		*s = full
	}

	return nil
}

// GetOrCreate finds a show by Spotify ID or creates it if it doesn't exist.
// This is used when a user subscribes to a show — we need to ensure the show
// record exists in our DB before creating the subscription.
func (r *showRepository) GetOrCreate(ctx context.Context, show *models.PodcastShow) error {
	err := dbFromCtx(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: "spotify_show_id",
		}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"description",
			"image_url",
			"spotify_url",
		}),
	}).Create(show).Error
	if err != nil {
		return fmt.Errorf("failed to upsert podcast show: %w", err)
	}

	existing, err := r.GetBySpotifyID(ctx, show.SpotifyShowID)
	if err != nil {
		return err
	}
	*show = *existing
	return nil
}

// GetAllTracked returns all shows that have at least one active subscription.
func (r *showRepository) GetAllTracked(ctx context.Context) ([]models.PodcastShow, error) {
	var shows []models.PodcastShow
	result := dbFromCtx(ctx, r.db).WithContext(ctx).
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
	result := dbFromCtx(ctx, r.db).WithContext(ctx).First(&show, "id = ?", id)
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
	result := dbFromCtx(ctx, r.db).WithContext(ctx).First(&show, "spotify_show_id = ?", spotifyShowID)
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
	result := dbFromCtx(ctx, r.db).WithContext(ctx).
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

// UpdateLatestEpisodeIfNull sets latest_episode_id / latest_episode_published_at only when
// latest_episode_id is currently NULL. Used on subscribe to seed a baseline without
// overwriting a watermark already established by another subscriber or the cron job.
func (r *showRepository) UpdateLatestEpisodeIfNull(ctx context.Context, showID uuid.UUID, episodeID string, publishedAt interface{}) error {
	result := dbFromCtx(ctx, r.db).WithContext(ctx).
		Model(&models.PodcastShow{}).
		Where("id = ? AND latest_episode_id IS NULL", showID).
		Updates(map[string]interface{}{
			"latest_episode_id":           episodeID,
			"latest_episode_published_at": publishedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to seed latest episode: %w", result.Error)
	}
	return nil
}
