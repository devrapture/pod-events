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

type EpisodeRepository interface {
	Create(ctx context.Context, episode *models.Episode) error
}

type episodeRepository struct {
	db *gorm.DB
}

func NewEpisodeRepository(db *gorm.DB) EpisodeRepository {
	return &episodeRepository{
		db: db,
	}
}

func (r *episodeRepository) Create(ctx context.Context, episode *models.Episode) error {
	result := r.db.WithContext(ctx).Create(episode)
	if result.Error != nil {
		return fmt.Errorf("failed to create episode: %w", result.Error)
	}

	return nil
}

// GetByID fetches an episode by UUID.
func (r *episodeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Episode, error) {
	var episode models.Episode
	result := r.db.WithContext(ctx).First(&episode, "id = ?")
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrEpisodeNotFound
		}
		return nil, result.Error
	}
	return &episode, nil
}

// GetBySpotifyID checks if we've already saved this episode.
// Returns (nil, nil) if not found — meaning it's a new episode.
func (r *episodeRepository) GetBySpotifyID(ctx context.Context, spotifyEpisodeID string) (*models.Episode, error) {
	var episode models.Episode
	result := r.db.WithContext(ctx).First(&episode, "spotify_episode_id = ?", spotifyEpisodeID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &episode, nil
}
