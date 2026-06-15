package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	"go.uber.org/zap"
)

type SpotifyClient struct {
	cfg        *config.Config
	httpClient *http.Client
	logger     *zap.Logger
}

const (
	spotifyAPIBase  = "https://api.spotify.com/v1"
	spotifyTokenURL = "https://accounts.spotify.com/api/token"
	spotifyAuthURL  = "https://accounts.spotify.com/authorize"
)

func NewSpotifyClient(cfg *config.Config, logger *zap.Logger) *SpotifyClient {
	return &SpotifyClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (c *SpotifyClient) AuthorizationURL(state string) string {
	params := url.Values{}
	params.Set("client_id", c.cfg.SpotifyClientID)
	params.Set("state", state)
	params.Set("redirect_uri", c.cfg.SpotifyRedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join([]string{
		"user-read-email",   // Read user's email and profile
		"user-read-private", // Read subscription status
		"user-library-read", // Read saved shows
		"user-follow-read",  // Read followed podcasts
	}, " "))

	return fmt.Sprintf("%s/?%s", spotifyAuthURL, params.Encode())
}

func (c *SpotifyClient) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.cfg.SpotifyRedirectURL)

	return c.requestToken(ctx, data)
}

func (c *SpotifyClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.cfg.SpotifyClientID)

	return c.requestToken(ctx, data)
}

func (c *SpotifyClient) GetCurrentUser(ctx context.Context, accessToken string) (*SpotifyUser, error) {
	var user SpotifyUser
	if err := c.get(ctx, accessToken, "me", &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *SpotifyClient) SearchShows(ctx context.Context, accessToken, query string, limit, offset int) (*ShowSearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", "show")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	endpoint := "search?" + params.Encode()
	var result ShowSearchResult
	if err := c.get(ctx, accessToken, endpoint, &result); err != nil {
		return nil, fmt.Errorf("error searching for shows: %w", err)
	}
	return &result, nil
}

// get podcasts saved by a user on spotify
func (c *SpotifyClient) GetUserSavedShows(ctx context.Context, accessToken string, offset, limit int) (*SpotifySavedShowsResponse, error) {
	endpoint := fmt.Sprintf("me/shows?offset=%d&limit=%d", offset, limit)
	var show SpotifySavedShowsResponse
	if err := c.get(ctx, accessToken, endpoint, &show); err != nil {
		return nil, err
	}
	return &show, nil
}

func (c *SpotifyClient) GetAllUserSavedShows(ctx context.Context, accessToken string) (*SpotifySavedShowsResponse, error) {
	const limit = 50
	all := &SpotifySavedShowsResponse{
		Limit: limit,
		Items: []SpotifySavedShowItem{},
	}

	firstPage := true
	for offset := 0; ; offset += limit {
		page, err := c.GetUserSavedShows(ctx, accessToken, offset, limit)
		if err != nil {
			return nil, err
		}
		if firstPage {
			all.Href = page.Href
			all.Total = page.Total
			firstPage = false
		}

		all.Items = append(all.Items, page.Items...)
		if page.Next == "" || len(all.Items) >= page.Total {
			break
		}
	}
	return all, nil
}

func (c *SpotifyClient) requestToken(ctx context.Context, data url.Values) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, spotifyTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.cfg.SpotifyClientID, c.cfg.SpotifyClientSecret)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()
	var result TokenResponse
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get token: %s", res.Status)
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result, nil
}

func (c *SpotifyClient) get(ctx context.Context, accessToken, endpoint string, target interface{}) error {
	endpoint = strings.TrimPrefix(endpoint, "/")
	url := fmt.Sprintf("%s/%s", spotifyAPIBase, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request to %s: %w", endpoint, err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request to %s: %w", endpoint, err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", endpoint, err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		return nil

	case http.StatusTooManyRequests:
		retryAfter := 1
		if v := res.Header.Get("Retry-After"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				retryAfter = n
			}
		}
		c.logger.Warn("spotify rate limit exceeded", zap.Int("retry_after", retryAfter))
		return &RateLimitError{RetryAfter: retryAfter}

	case http.StatusNotFound:
		return fmt.Errorf("spotify resource not found: %s", string(body))

	case http.StatusUnauthorized:
		return fmt.Errorf("spotify access token expired or invalid: %s", string(body))

	default:
		return fmt.Errorf("failed to get %s: %s: %s", endpoint, res.Status, string(body))
	}
}
