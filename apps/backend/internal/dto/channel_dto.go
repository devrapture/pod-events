package dto

import "github.com/devrapture/pod-events/internal/models"

// CreateChannelRequest represents a request to create a notification channel.
type CreateChannelRequest struct {
	// Channel type: slack_webhook, discord_webhook, whatsapp, or telegram
	ChannelType models.ChannelType `json:"channel_type" binding:"required,oneof=slack_webhook discord_webhook whatsapp telegram" example:"slack_webhook"`
	// Destination URL or identifier for the channel
	Destination string `json:"destination" binding:"required" example:"https://hooks.slack.com/services/..."`
	// Optional label to identify the channel
	Label string `json:"label" example:"Work Slack"`
}

// ToggleActiveChannelRequest represents a request to enable or disable a channel.
type ToggleActiveChannelRequest struct {
	IsActive *bool `json:"is_active" binding:"required" example:"true"`
}
