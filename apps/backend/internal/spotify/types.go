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

type SpotifyEpisode struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	DurationMs           int    `json:"duration_ms"`
	ReleaseDate          string `json:"release_date"`
	ReleaseDatePrecision string `json:"release_date_precision"`
}

func (e *SpotifyEpisode) ImageURL() string {
	if len(e.Images) > 0 {
		return e.Images[0].URL
	}
	return ""
}

func (e *SpotifyEpisode) ParsedReleaseDate() (time.Time, error) {
	var parseErr error
	for _, layout := range []string{"2006", "2006-01", "2006-01-02"} {
		var parsed time.Time
		parsed, parseErr = time.Parse(layout, e.ReleaseDate)
		if parseErr == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("parse release date %q: %w", e.ReleaseDate, parseErr)
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
	Description   string `json:"description"`
	TotalEpisodes int    `json:"total_episodes"`
	ExternalURLs  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images []struct {
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
			Description:   item.Show.Description,
			AddedAt:       item.AddedAt,
			ImageURL:      item.Show.ImageURL(),
			TotalEpisodes: item.Show.TotalEpisodes,
			SpotifyURL:    item.Show.ExternalURLs.Spotify,
		})
	}
	return shows
}

type ShowSearchResult struct {
	Shows struct {
		Items []SpotifyShow `json:"items"`
		Total int           `json:"total"`
	} `json:"shows"`
}

func (s *ShowSearchResult) ToSavedShows() []dto.SavedShowResponse {
	shows := make([]dto.SavedShowResponse, 0, len(s.Shows.Items))

	for _, item := range s.Shows.Items {
		shows = append(shows, dto.SavedShowResponse{
			ID:            item.ID,
			Name:          item.Name,
			Description:   item.Description,
			ImageURL:      item.ImageURL(),
			TotalEpisodes: item.TotalEpisodes,
			SpotifyURL:    item.ExternalURLs.Spotify,
		})
	}
	return shows
}
