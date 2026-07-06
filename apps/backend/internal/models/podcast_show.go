package models

import "time"

// PodcastShow represents a podcast on Spotify.
type PodcastShow struct {
	Base
	SpotifyShowID            string     `json:"spotify_show_id" gorm:"uniqueIndex;not null"`
	Name                     string     `json:"name" gorm:"type:text;not null"`
	Description              string     `json:"description" gorm:"type:text;not null"`
	ImageURL                 string     `json:"image_url" gorm:"type:text;not null"`
	SpotifyURL               string     `json:"spotify_url" gorm:"type:text;not null"`
	LatestEpisodeID          *string    `json:"latest_episode_id,omitempty" gorm:"type:text"`
	LatestEpisodePublishedAt *time.Time `json:"latest_episode_published_at,omitempty"`

	Episode      []Episode      `json:"-" gorm:"foreignKey:PodcastShowID"`
	Subscription []Subscription `json:"-" gorm:"foreignKey:PodcastShowID"`
}
