package spotify

import (
	"fmt"
	"time"
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
