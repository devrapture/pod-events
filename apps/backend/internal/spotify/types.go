package spotify

import (
	"fmt"
	"time"

	"github.com/devrapture/pod-events/internal/dto"
)

// TokenResponse is Spotify's response to a token request.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`    // seconds until expiry (usually 3600 = 1 hour)
	RefreshToken string `json:"refresh_token"` // may be empty on refresh responses
	Scope        string `json:"scope"`
}

// ExpiresAt converts ExpiresIn seconds to an absolute time.Time
func (t *TokenResponse) ExpiresAt() time.Time {
	return time.Now().UTC().Add(time.Duration(t.ExpiresIn) * time.Second)
}

type SpotifyUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
}

func (u *SpotifyUser) AvatarURL() string {
	if len(u.Images) > 0 {
		return u.Images[0].URL
	}
	return ""
}

type RateLimitError struct {
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("spotify rate limit exceeded, retry after %d seconds", e.RetryAfter)
}

type SpotifySavedShowsResponse struct {
	Href     string                 `json:"href"`
	Limit    int                    `json:"limit"`
	Next     string                 `json:"next"`
	Offset   int                    `json:"offset"`
	Previous string                 `json:"previous"`
	Total    int                    `json:"total"`
	Items    []SpotifySavedShowItem `json:"items"`
}

type SpotifySavedShowItem struct {
	AddedAt string      `json:"added_at"`
	Show    SpotifyShow `json:"show"`
}

type SpotifyShow struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Publisher     string `json:"publisher"`
	TotalEpisodes int    `json:"total_episodes"`
	Images        []struct {
		URL string `json:"url"`
	} `json:"images"`
}

func (s *SpotifyShow) ImageURL() string {
	if len(s.Images) > 0 {
		return s.Images[0].URL
	}
	return ""
}

func (s *SpotifySavedShowsResponse) ToSavedShows() []dto.SavedShowResponse {
	shows := make([]dto.SavedShowResponse, 0, len(s.Items))

	for _, item := range s.Items {
		shows = append(shows, dto.SavedShowResponse{
			ID:            item.Show.ID,
			Name:          item.Show.Name,
			Publisher:     item.Show.Publisher,
			AddedAt:       item.AddedAt,
			ImageURL:      item.Show.ImageURL(),
			TotalEpisodes: item.Show.TotalEpisodes,
		})
	}
	return shows
}
