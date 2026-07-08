package handler

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/devrapture/pod-events/internal/config"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/notifications/telegram"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const telegramSecretHeader = "X-Telegram-Bot-Api-Secret-Token"

type TelegramWebHookHandler struct {
	cfg                       *config.Config
	notifier                  *telegram.Notifier
	telegramConnectionService services.TelegramConnectionServices
	logger                    *zap.Logger
}

type SendTestMessageRequest struct {
	ChatID int64 `json:"chat_id"`
}

func NewTelegramWebHookHandler(cfg *config.Config, notifier *telegram.Notifier, telegramConnectionService services.TelegramConnectionServices, logger *zap.Logger) *TelegramWebHookHandler {
	return &TelegramWebHookHandler{
		cfg:                       cfg,
		notifier:                  notifier,
		telegramConnectionService: telegramConnectionService,
		logger:                    logger,
	}
}

// CreateConnectLink creates a one-time Telegram connection link.
//
//	@Summary     Generate Telegram connect link
//	@Description Generate a one-time link to connect a Telegram chat to the user's PodEvents account
//	@Tags        Telegram
//	@Security    BearerAuth
//	@Produce     json
//	@Success     200 {object} response.APIResponse "Telegram connect link created"
//	@Failure     409 {object} response.APIResponse "already have a Telegram channel"
//	@Router      /telegram/generate-link [post]
func (h *TelegramWebHookHandler) CreateConnectLink(c *gin.Context) {
	userID, _ := c.Get("userID")
	link, err := h.telegramConnectionService.CreateConnectLink(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, apperrors.ErrTelegramChannelAlreadyExists) {
			response.ErrorResponse(c, http.StatusConflict, "You already have a Telegram channel. Delete it before creating a new Telegram connection link.")
			return
		}
		h.logger.Error("failed to create telegram connect link", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create telegram connect link")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Telegram connect link created", gin.H{"url": link}, nil)
}

// Handle processes incoming Telegram webhook updates.
//
//	@Summary     Telegram webhook
//	@Description Handle incoming updates from Telegram (bot commands, messages)
//	@Tags        Telegram
//	@Accept      json
//	@Param       update body telegram.Update true "Telegram update"
//	@Success     200
//	@Failure     401 "invalid telegram secret"
//	@Router      /webhooks/telegram [post]
func (h *TelegramWebHookHandler) Handle(c *gin.Context) {
	if !h.validSecret(c.GetHeader(telegramSecretHeader)) {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid telegram secret")
		return
	}

	var update telegram.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		h.logger.Warn("invalid telegram webhook payload", zap.Error(err))
		c.Status(http.StatusOK)
		return
	}

	if update.Message == nil {
		c.Status(http.StatusOK)
		return
	}

	// chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)
	if strings.HasPrefix(text, "/start ") {
		token := strings.TrimSpace(strings.TrimPrefix(text, "/start "))
		chatID := update.Message.Chat.ID

		name := strings.TrimSpace(strings.Join([]string{
			update.Message.From.FirstName,
			update.Message.From.LastName,
		}, " "))
		if name == "" {
			name = "there"
		}

		_, err := h.telegramConnectionService.CompleteConnection(
			c.Request.Context(),
			token,
			chatID,
		)
		if err != nil {
			if errors.Is(err, apperrors.ErrTelegramChannelAlreadyExists) {
				_ = h.notifier.SendToChatID(
					c.Request.Context(),
					"PodEvents is already connected to a Telegram channel. Delete the existing Telegram channel in the app before connecting a new one.",
					chatID,
				)
				c.Status(http.StatusOK)
				return
			}
			h.logger.Warn("failed to complete telegram connection", zap.Error(err))
			_ = h.notifier.SendToChatID(
				c.Request.Context(),
				"This PodEvents connection link is invalid or expired. Please generate a new one from the app.",
				chatID,
			)
			c.Status(http.StatusOK)
			return
		}

		appURL := h.cfg.FrontendURL
		message := fmt.Sprintf("Hi %s, PodEvents is now connected to this Telegram chat. You can now receive notifications in this chat. You can manage your connection in the app at %s/dashboard/channels", name, appURL)
		_ = h.notifier.SendToChatID(c.Request.Context(), message, chatID)
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}

func (h *TelegramWebHookHandler) validSecret(got string) bool {
	expected := h.cfg.TelegramWebhookSecret
	if expected == "" || got == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1
}
