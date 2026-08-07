package models

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type ChannelType string

const (
	ChannelTypeSlack    ChannelType = "slack_webhook"
	ChannelTypeDiscord  ChannelType = "discord_webhook"
	ChannelTypeTelegram ChannelType = "telegram"
	ChannelTypeWhatsApp ChannelType = "whatsapp"
)

// ChannelType determines how we interpret Destination:
//   - "slack_webhook"    → Destination is the Slack webhook URL
//   - "discord_webhook"  → Destination is the Discord webhook URL
//   - "telegram"         → Destination is the Telegram chat ID (string of digits)
//   - "whatsapp"         → Destination is the WhatsApp phone number (E.164 format)

type NotificationChannel struct {
	Base
	UserID                 uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_user_channel_destination,where:deleted_at IS NULL"`
	ChannelType            ChannelType `json:"channel_type" gorm:"not null;uniqueIndex:idx_user_channel_destination,where:deleted_at IS NULL"`
	Destination            string      `json:"-" gorm:"not null"`
	DestinationFingerprint string      `json:"-" gorm:"type:char(64);not null;uniqueIndex:idx_user_channel_destination,where:deleted_at IS NULL"`
	IsActive               bool        `json:"is_active" gorm:"type:bool;default:true"`
	Label                  string      `json:"label"`
	User                   User        `json:"-" gorm:"foreignkey:UserID"`
}

func (c ChannelType) IsValid() bool {
	switch c {
	case ChannelTypeWhatsApp, ChannelTypeTelegram, ChannelTypeSlack, ChannelTypeDiscord:
		return true
	}
	return false
}

func (c ChannelType) IsWebhook() bool {
	return c == ChannelTypeSlack || c == ChannelTypeDiscord
}

func (c ChannelType) ValidateDestination(destination string) error {
	switch c {
	case ChannelTypeSlack:
		if !c.isValidSlackWebhook(destination) {
			return fmt.Errorf("invalid slack webhook")
		}
	case ChannelTypeDiscord:
		if !c.isValidDiscordWebhook(destination) {
			return fmt.Errorf("invalid discord webhook")
		}
	}
	return nil
}

func (c ChannelType) isValidSlackWebhook(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if u.Scheme != "https" {
		return false
	}
	if !strings.EqualFold(u.Hostname(), "hooks.slack.com") {
		return false
	}

	parts := splitPath(u.Path)

	if len(parts) == 0 || parts[0] != "services" {
		return false
	}

	return true
}

func (c ChannelType) isValidDiscordWebhook(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if u.Scheme != "https" {
		return false
	}

	if !strings.EqualFold(u.Hostname(), "discord.com") {
		return false
	}

	return true
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}

	return strings.Split(path, "/")
}
