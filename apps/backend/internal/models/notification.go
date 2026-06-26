package models

import "github.com/google/uuid"

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
	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid; not null;index"`
	ChannelType ChannelType `json:"channel_type" gorm:"not null"`
	Destination string      `json:"destination" gorm:"not null"`
	IsActive    bool        `json:"is_active" gorm:"type:bool; default:true"`
	Label       string      `json:"label"` // optional human-readable label, e.g. "#dev-alerts"

	User User `json:"-" gorm:"foreignkey:UserID"`
}

func (c ChannelType) IsValid() bool {
	switch c {
	case ChannelTypeWhatsApp, ChannelTypeTelegram, ChannelTypeSlack, ChannelTypeDiscord:
		return true
	}
	return false
}
