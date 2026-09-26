package googleauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const userInfoURL = "https://openidconnect.googleapis.com/v1/userinfo"

type Profile struct {
	ID            string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

type Client struct {
	oauthConfig *oauth2.Config
	httpClient  *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "email", "profile"},
		},
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) AuthorizationURL(state string) string {
	return c.oauthConfig.AuthCodeURL(state)
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (*Profile, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, c.httpClient)
	token, err := c.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange google code: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create google userinfo request: %w", err)
	}
	token.SetAuthHeader(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get google userinfo: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read google userinfo: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get google userinfo: %s: %s", response.Status, string(body))
	}

	var profile Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("decode google userinfo: %w", err)
	}
	if profile.ID == "" || profile.Email == "" || !profile.EmailVerified {
		return nil, fmt.Errorf("google account must have a verified email")
	}
	return &profile, nil
}
