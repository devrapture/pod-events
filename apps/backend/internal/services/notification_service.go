package services

import (
	"context"
	"fmt"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/notifications"
	"github.com/devrapture/pod-events/internal/notifications/discord"
	"github.com/devrapture/pod-events/internal/notifications/slack"
	"github.com/devrapture/pod-events/internal/notifications/telegram"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationService interface {
	NotifyUser(ctx context.Context, userID uuid.UUID, episode *models.Episode, show *models.PodcastShow) error
}

type notificationService struct {
	logRepo     repositories.NotificationLogRepository
	logger      *zap.Logger
	cfg         *config.Config
	channelRepo repositories.ChannelRepository
}

func NewNotificationService(logRepo repositories.NotificationLogRepository, logger *zap.Logger, cfg *config.Config, channelRepo repositories.ChannelRepository) NotificationService {
	return &notificationService{
		logRepo:     logRepo,
		logger:      logger,
		cfg:         cfg,
		channelRepo: channelRepo,
	}
}

func (s *notificationService) NotifyUser(ctx context.Context, userID uuid.UUID, episode *models.Episode, show *models.PodcastShow) error {
	channels, err := s.channelRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get notification channels: %w", err)
	}
	if len(channels) == 0 {
		s.logger.Info("User has no notification channel configured", zap.String("user_id", userID.String()))
		return nil
	}

	msg := notifications.NotificationMessage{
		ShowName:     show.Name,
		EpisodeTitle: episode.Name,
		SpotifyURL:   episode.SpotifyURL,
		Description:  episode.Description,
		ImageURL:     show.ImageURL,
		PublishedAt:  episode.ReleaseDate,
	}

	for _, channel := range channels {
		alreadySent, err := s.logRepo.AlreadlySent(ctx, userID, episode.ID, channel.ChannelType)
		if err != nil {
			s.logger.Error(
				"failed to check duplicate notification",
				zap.String("user_id", userID.String()),
				zap.String("episode_id", episode.ID.String()),
				zap.String("channel_type", string(channel.ChannelType)),
				zap.Error(err),
			)
			continue
		}

		if alreadySent {
			s.logger.Info("Notification already sent", zap.String("user_id", userID.String()), zap.String("episode_id", episode.ID.String()), zap.String("channel_type", string(channel.ChannelType)))
			continue
		}
		notifier, err := s.buildNotifier(channel)
		if err != nil {
			s.logger.Error("failed to build notifier", zap.String("channel_type", string(channel.ChannelType)), zap.Error(err))
			s.saveLog(ctx, userID, episode.ID, channel.ChannelType, models.NotificationStatusFailed, err.Error())
			continue
		}
		sendCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		sendErr := notifier.Send(sendCtx, msg)
		cancel()

		if sendErr != nil {
			s.logger.Error(
				"notification send failed",
				zap.String("user_id", userID.String()),
				zap.String("episode_id", episode.ID.String()),
				zap.String("channel_type", string(channel.ChannelType)))
			s.saveLog(ctx, userID, episode.ID, channel.ChannelType, models.NotificationStatusFailed, sendErr.Error())
		} else {
			s.logger.Info(
				"notification sent",
				zap.String("user_id", userID.String()),
				zap.String("episode_id", episode.ID.String()),
				zap.String("channel_type", string(channel.ChannelType)),
			)
			s.saveLog(ctx, userID, episode.ID, channel.ChannelType, models.NotificationStatusSent, "")
		}
	}

	return nil
}

func (s *notificationService) buildNotifier(channel models.NotificationChannel) (notifications.Notifier, error) {
	switch channel.ChannelType {
	case models.ChannelTypeSlack:
		return slack.NewNotifier(channel.Destination, s.logger), nil
	case models.ChannelTypeDiscord:
		return discord.NewNotifier(channel.Destination, s.logger), nil
	case models.ChannelTypeTelegram:
		return telegram.NewNotifier(s.cfg), nil
	default:
		return nil, fmt.Errorf("unknown channel type: %s", channel.ChannelType)
	}
}

func (s *notificationService) saveLog(ctx context.Context, userID uuid.UUID, episodeID uuid.UUID, channelType models.ChannelType, status models.NotificationStatus, errMsg string) {
	now := time.Now().UTC()
	log := &models.NotificationLog{
		UserID:       userID,
		EpisodeID:    episodeID,
		ChannelType:  channelType,
		Status:       status,
		ErrorMessage: errMsg,
		Base:         models.Base{CreatedAt: now},
	}
	if status == models.NotificationStatusSent {
		log.SentAt = &now
	}
	if err := s.logRepo.Create(ctx, log); err != nil {
		s.logger.Error("failed to save notification log", zap.Error(err))
	}
}
