package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

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

	if err := req.ChannelType.ValidateDestination(req.Destination); err != nil {
		return nil, err
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
		if !s.isValidSlackWebhook(destination) {
			return apperrors.ErrInvalidSlackWebhook
		}
	case models.ChannelTypeDiscord:
		if !s.isValidDiscordWebhook(destination) {
			return apperrors.ErrInvalidDiscordWebhook
		}
	case models.ChannelTypeTelegram:
		if _, err := strconv.ParseInt(destination, 10, 64); err != nil {
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

func (s *channelServices) isValidSlackWebhook(rawURL string) bool {
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

func (s *channelServices) isValidDiscordWebhook(rawURL string) bool {
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
