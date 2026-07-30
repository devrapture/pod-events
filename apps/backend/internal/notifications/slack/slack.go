package slack

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
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (n *Notifier) Send(ctx context.Context, message notifications.NotificationMessage) error {
	payload := map[string]interface{}{
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": ":headphones: *A new podcast episode is available!*",
				},
			},
			{
				"type":     "section",
				"block_id": "episode_details",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s*\n\n%s\n\n<%s|Listen on Spotify>", message.EpisodeTitle, utils.Truncate(message.Description, 300), message.SpotifyURL),
				},
				"accessory": map[string]interface{}{
					"type":      "image",
					"image_url": message.ImageURL,
					"alt_text":  "Podcast episode cover",
				},
			},
			{
				"type":     "section",
				"block_id": "episode_metadata",
				"fields": []map[string]interface{}{
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Podcast*\n%s", message.ShowName),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Published*\n%s", message.PublishedAt.Format("Jan 02, 2006")),
					},
				},
			},
			{
				"type": "context",
				"elements": []map[string]interface{}{
					{
						"type": "mrkdwn",
						"text": "<https://pod-event.vercel.app|PodEvents> • Built with :heart: by <https://x.com/devrappy|DevRapture>",
					},
				},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		n.logger.Error("failed to marshal slack payload", zap.Error(err))
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webHookURL, bytes.NewReader(body))
	if err != nil {
		n.logger.Error("failed to create slack request", zap.Error(err))
		return fmt.Errorf("failed to create slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := n.http.Do(req)
	if err != nil {
		n.logger.Error("failed to send slack notification", zap.Error(err))
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		n.logger.Error("slack returned unexpected status", zap.Int("status", res.StatusCode))
		return fmt.Errorf("slack returned non-200 status: %d", res.StatusCode)
	}
	return nil
}

func (n *Notifier) Type() string {
	return string(models.ChannelTypeSlack)
}
