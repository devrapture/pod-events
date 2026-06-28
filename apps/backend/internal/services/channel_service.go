package services

import (
	"context"
	"fmt"

	"github.com/devrapture/pod-events/internal/dto"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/google/uuid"
)

type ChannelServices interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateChannelRequest) (*models.NotificationChannel, error) 
	GetByUser(ctx context.Context, userID uuid.UUID) ([]models.NotificationChannel, error)
	DeleteChannel(ctx context.Context, userID, channelID uuid.UUID) error
	ToggleActive(ctx context.Context, userID, channelID uuid.UUID, isActive bool) error
}

type channelServices struct {
	channelRepo repositories.ChannelRepository
}

func NewChannelServices(channelRepo repositories.ChannelRepository) ChannelServices {
	return &channelServices{
		channelRepo: channelRepo,
	}
}

func (s *channelServices) Create(ctx context.Context, userID uuid.UUID, req dto.CreateChannelRequest) (*models.NotificationChannel, error) {
	if !req.ChannelType.IsValid() {
		return nil, fmt.Errorf("invalid channel_type %q — must be one of: slack_webhook, discord_webhook, telegram, whatsapp", req.ChannelType)
	}

	if err := s.validateDestination(req.ChannelType, req.Destination); err != nil {
		return nil, err
	}
	channel := &models.NotificationChannel{
		UserID:      userID,
		ChannelType: req.ChannelType,
		Destination: req.Destination,
		Label:       req.Label,
	}
	if err := s.channelRepo.Create(ctx, channel); err != nil {
		return nil, fmt.Errorf("failed to create notification channel: %w", err)
	}
	return channel, nil
}

// GetByUser returns all active channels for a user.
func (s *channelServices) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.NotificationChannel, error) {
	return s.channelRepo.GetByUserID(ctx, userID)
}

// Delete removes a notification channel, verifying ownership.
func (s *channelServices) DeleteChannel(ctx context.Context, userID, channelID uuid.UUID) error {
	return s.channelRepo.Delete(ctx, userID, channelID)
}

func (s *channelServices) ToggleActive(ctx context.Context, userID, channelID uuid.UUID, isActive bool) error {
	return s.channelRepo.ToggleActive(ctx, userID, channelID, isActive)
}

// validateDestination does basic format validation per channel type.
func (s *channelServices) validateDestination(channelType models.ChannelType, destination string) error {
	switch channelType {
	case models.ChannelTypeSlack:
		if len(destination) < 10 || destination[:8] != "https://" {
			return fmt.Errorf("slack destination must be a valid HTTPS webhook URL")
		}
	case models.ChannelTypeDiscord:
		if len(destination) < 10 || destination[:8] != "https://" {
			return fmt.Errorf("discord destination must be a valid HTTPS webhook URL")
		}
	case models.ChannelTypeTelegram:
		if len(destination) == 0 {
			return fmt.Errorf("telegram destination must be a chat_id (numeric string)")
		}
	case models.ChannelTypeWhatsApp:
		// Basic E.164 format check: starts with + and has at least 10 digits
		if len(destination) < 10 || destination[0] != '+' {
			return fmt.Errorf("whatsapp destination must be a phone number in E.164 format (e.g., +2348012345678)")
		}
	}
	return nil
}
