package models

import (
	"time"

	"github.com/google/uuid"
)

// Episode represents a single podcast episode.
type Episode struct {
	Base
	PodcastShowID    uuid.UUID `json:"podcast_show_id" gorm:"type:uuid;not null;index"`
	SpotifyEpisodeID string    `json:"spotify_episode_id" gorm:"uniqueIndex;not null"`
	SpotifyURL       string    `json:"spotify_url"`
	DurationMs       int       `json:"duration_ms"`
	ReleaseDate      time.Time `json:"release_date"`
	ImageURL         string    `json:"image_url" `
	Name             string    `json:"name"`
	Description      string    `json:"description"`

	PodcastShow PodcastShow `json:"-" gorm:"foreignKey:PodcastShowID"`
}
