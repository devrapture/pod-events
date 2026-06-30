package dto

import "github.com/devrapture/pod-events/internal/models"

type CreateChannelRequest struct {
	ChannelType models.ChannelType `json:"channel_type" binding:"required,oneof=slack_webhook discord_webhook whatsapp telegram"`
	Destination string             `json:"destination" binding:"required"`
	Label       string             `json:"label"`
}

type ToggleActiveChannelRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}
