package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/notifications"
	"github.com/devrapture/pod-events/pkg/utils"
	"go.uber.org/zap"
)

const telegramAPIBaseURL = "https://api.telegram.org/bot"

type Notifier struct {
	botToken   string
	httpClient *http.Client
	baseURL    string
	chatID     int64
	logger     *zap.Logger
}

func NewNotifier(cfg *config.Config, chatID int64, logger *zap.Logger) *Notifier {
	return &Notifier{
		botToken: cfg.TelegramBotToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: fmt.Sprintf("%s%s", telegramAPIBaseURL, cfg.TelegramBotToken),
		chatID:  chatID,
		logger:  logger,
	}
}

func (n *Notifier) Send(ctx context.Context, message notifications.NotificationMessage) error {
	text := fmt.Sprintf(
		"🎙️ <b>A new podcast episode is available!</b>\n\n"+
			"<b>%s</b>\n\n"+
			"%s\n\n"+
			"🎧 <b>Podcast:</b> %s\n"+
			"📅 <b>Published:</b> %s\n\n"+
			"▶️ <a href=\"%s\">Listen on Spotify</a>",
		html.EscapeString(message.EpisodeTitle),
		html.EscapeString(utils.Truncate(message.Description, 300)),
		html.EscapeString(message.ShowName),
		message.PublishedAt.Format("Jan 02, 2006"),
		html.EscapeString(message.SpotifyURL),
	)

	payload := map[string]interface{}{
		"chat_id":    n.chatID,
		"text":       text,
		"parse_mode": "HTML",
	}
	url := fmt.Sprintf("%s/sendMessage", n.baseURL)
	return n.post(ctx, url, payload, nil)
}

func (n *Notifier) Type() string {
	return string(models.ChannelTypeTelegram)
}

func (n *Notifier) SendToChatID(ctx context.Context, text string, chatID int64) error {
	url := fmt.Sprintf("%s/sendMessage", n.baseURL)
	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}
	return n.post(ctx, url, payload, nil)
}

func (n *Notifier) post(ctx context.Context, url string, payload, result interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed or timed out: %w", err)
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read telegram response body: %w", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return apperrors.ErrInvalidTelegramBotToken
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status code %d", res.StatusCode)
	}

	if result != nil {
		if err := json.Unmarshal(resBody, result); err != nil {
			return fmt.Errorf("failed to parse telegram response: %w", err)
		}
	}
	return nil
}
