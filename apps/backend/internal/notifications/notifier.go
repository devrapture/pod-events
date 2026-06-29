package notifications

import (
	"context"
	"time"
)

type Notifier interface {
	Send(ctx context.Context, notificationMessage NotificationMessage) error
	Type() string
}

type NotificationMessage struct {
	ShowName     string
	EpisodeTitle string
	Description  string
	SpotifyURL   string
	ImageURL     string
	PublishedAt  time.Time
}

// Format the message consistently for text-based channels (Telegram, WhatsApp)
func FormatTextMessage(msg NotificationMessage) string {
	return "🎙️ *New Episode: " + msg.ShowName + "*\n\n" +
		"*" + msg.EpisodeTitle + "*\n\n" +
		truncate(msg.Description, 200) + "\n\n" +
		"▶️ Listen: " + msg.SpotifyURL
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)

	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen]) + "..."
}
