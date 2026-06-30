package models

// Episode represents a single podcast episode.
type Episode struct {
	Base
	PodcastShowID    string `json:"podcast_show_id" gorm:"not null;index"`
	SpotifyEpisodeID string `json:"spotify_episode_id" gorm:"uniqueIndex;not null"`
	AudioPreviewURL  string `json:"audio_preview_url"`
	SpotifyURL       string `json:"spotify_url"`
	DurationMs       int    `json:"duration_ms"`
	ReleaseDate      string `json:"release_date"`
	ImageURL         string `json:"image_url" `
	Name             string `json:"name"`

	PodcastShow PodcastShow `json:"-" gorm:"foreignKey:PodcastShowID"`
}
