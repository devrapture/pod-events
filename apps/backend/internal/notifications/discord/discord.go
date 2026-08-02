package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/notifications"
	"github.com/devrapture/pod-events/pkg/utils"
	"go.uber.org/zap"
)

type Notifier struct {
	webHookURL string
	http       *http.Client
	logger     *zap.Logger
}

func NewNotifier(webHookURL string, logger *zap.Logger) *Notifier {
	return &Notifier{
		webHookURL: webHookURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

type Discord struct {
	Content string         `json:"content"`
	Embeds  []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	Color       int        `json:"color"`
	Fields      []Fields   `json:"fields,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
	TimeStamp   string     `json:"timestamp"`
	ThumbNail   *ThumbNail `json:"thumbnail,omitempty"`
}

type Footer struct {
	Text string `json:"text"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type ThumbNail struct {
	URL string `json:"url"`
}

func (n *Notifier) Send(ctx context.Context, message notifications.NotificationMessage) error {
	spotifyGreen := 1_947_988
	payload := Discord{
		Content: "🎧 **A new podcast episode is available!**",
		Embeds: []DiscordEmbed{
			{
				Title:       "🎙️ " + message.EpisodeTitle,
				Description: utils.Truncate(message.Description, 300),
				URL:         message.SpotifyURL,
				Color:       spotifyGreen,
				Fields: []Fields{
					{
						Name:   "Podcast",
						Value:  message.ShowName,
						Inline: true,
					},
					{
						Name:   "Published",
						Value:  message.PublishedAt.Format("Jan 02, 2006"),
						Inline: true,
					},
					{
						Name:  "Listen on Spotify",
						Value: fmt.Sprintf("[Play episode](%s)", message.SpotifyURL),
					},
				},
				Footer: &Footer{
					Text: "PodEvents • Built with ❤️ by devrapture",
				},
				TimeStamp: message.PublishedAt.Format(time.RFC3339),
				ThumbNail: &ThumbNail{
					URL: message.ImageURL,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		n.logger.Error("failed to marshal discord payload", zap.Error(err))
		return fmt.Errorf("marshal discord payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webHookURL, bytes.NewReader(body))
	if err != nil {
		n.logger.Error("failed to create discord request", zap.Error(err))
		return fmt.Errorf("create discord request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := n.http.Do(req)
	if err != nil {
		n.logger.Error("failed to send discord notification", zap.Error(err))
		return fmt.Errorf("send discord notification: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		n.logger.Error("discord returned unexpected status", zap.Int("status", res.StatusCode))
		return fmt.Errorf("discord returned unexpected status: %d", res.StatusCode)
	}
	return nil
}

func (n *Notifier) Type() string {
	return string(models.ChannelTypeDiscord)
}
