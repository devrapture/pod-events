package dto

// SavedShowResponse represents a saved podcast show from Spotify.
type SavedShowResponse struct {
	ID            string `json:"id" example:"4GvhKq1eEYj7JzAKUwWmNc"`
	Name          string `json:"name" example:"The Joe Rogan Experience"`
	Description   string `json:"description" example:"Joe Rogan's podcast"`
	TotalEpisodes int    `json:"total_episodes" example:"200"`
	ImageURL      string `json:"image_url" example:"https://i.scdn.co/image/abc123"`
	SpotifyURL    string `json:"spotify_url" example:"https://open.spotify.com/show/abc123"`
	AddedAt       string `json:"added_at" example:"2024-01-15T10:30:00Z"`
}
