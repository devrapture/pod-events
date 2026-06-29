package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/google/uuid"
)

type TelegramConnectionServices interface {
	CreateConnectLink(ctx context.Context, userID uuid.UUID) (string, error)
	CompleteConnection(ctx context.Context, token string, chatID int64) (*models.NotificationChannel, error)
}

type telegramConnectionService struct {
	telegramConnectionRepo repositories.TelegramConnectionRepository
	channelRepo            repositories.ChannelRepository
	botName                string
}

func NewTelegramConnectionService(telegramConnectionRepo repositories.TelegramConnectionRepository, channelRepo repositories.ChannelRepository, cfg *config.Config) TelegramConnectionServices {
	s := telegramConnectionService{
		telegramConnectionRepo: telegramConnectionRepo,
		channelRepo:            channelRepo,
		botName:                cfg.BotName,
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			_ = telegramConnectionRepo.DeleteExpired(context.Background())
		}
	}()

	return &s
}

func (s *telegramConnectionService) CreateConnectLink(ctx context.Context, userID uuid.UUID) (string, error) {
	if err := s.ensureNoTelegramChannel(ctx, userID); err != nil {
		return "", err
	}

	token, err := randomToken()
	if err != nil {
		return "", err
	}
	conn := &models.TelegramConnection{
		UserID:    userID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.telegramConnectionRepo.Create(ctx, conn); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://t.me/%s?start=%s", s.botName, token), nil
}

func (s *telegramConnectionService) ensureNoTelegramChannel(ctx context.Context, userID uuid.UUID) error {
	channels, err := s.channelRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, channel := range channels {
		if channel.ChannelType == models.ChannelTypeTelegram {
			return apperrors.ErrTelegramChannelAlreadyExists
		}
	}

	return nil
}

func (s *telegramConnectionService) CompleteConnection(ctx context.Context, token string, chatID int64) (*models.NotificationChannel, error) {
	conn, err := s.telegramConnectionRepo.Consume(ctx, hashToken(token))
	if err != nil {
		return nil, err
	}
	if err := s.ensureNoTelegramChannel(ctx, conn.UserID); err != nil {
		return nil, err
	}
	channel := &models.NotificationChannel{
		UserID:      conn.UserID,
		ChannelType: models.ChannelTypeTelegram,
		Destination: strconv.FormatInt(chatID, 10),
		IsActive:    true,
	}
	if err := s.channelRepo.Create(ctx, channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate telegram connection token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
