package dto

type SavedShowResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	TotalEpisodes int    `json:"total_episodes"`
	ImageURL      string `json:"image_url"`
	SpotifyURL    string `json:"spotify_url"`
	AddedAt       string `json:"added_at"`
}
