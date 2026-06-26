package handler

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/notifications/telegram"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const telegramSecretHeader = "X-Telegram-Bot-Api-Secret-Token"

type TelegramWebHookHandler struct {
	cfg      *config.Config
	notifier *telegram.Notifier
	logger   *zap.Logger
}

type SendTestMessageRequest struct {
	ChatID int64 `json:"chat_id"`
}

func NewTelegramWebHookHandler(cfg *config.Config, notifier *telegram.Notifier, logger *zap.Logger) *TelegramWebHookHandler {
	return &TelegramWebHookHandler{
		cfg:      cfg,
		notifier: notifier,
		logger:   logger,
	}
}

func (h *TelegramWebHookHandler) SendTestMessage(c *gin.Context) {
	message := `🧪 Test notification
	
	Good news! Your Telegram notifications are working.
	
	From now on, whenever a podcast you're subscribed to releases a new episode, you'll receive an alert like this.
	
	Happy listening! 🎧`

	var request SendTestMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body. Please provide a valid chat_id.")
		c.Status(http.StatusOK)
		return
	}

	if err := h.notifier.SendToChatID(c.Request.Context(), message, request.ChatID); err != nil {
		h.logger.Error("failed to send telegram chat id", zap.Error(err), zap.Int64("chat_id", request.ChatID))
	}

	response.SuccessResponse(c, http.StatusOK, "test message sent successfully", nil, nil)
}

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

	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)
	name := strings.TrimSpace(strings.Join([]string{
		update.Message.From.FirstName,
		update.Message.From.LastName,
	}, " "))
	if name == "" {
		name = "there"
	}
	message := fmt.Sprintf("Hi %s, this is PodEvents. Your Telegram chat ID is %d.", name, chatID)
	switch text {
	case "/start":
		if err := h.notifier.SendToChatID(c.Request.Context(), message, chatID); err != nil {
			h.logger.Error("failed to send telegram chat id", zap.Error(err), zap.Int64("chat_id", chatID))
		}
	case "/chatid":
		if err := h.notifier.SendToChatID(c.Request.Context(), message, chatID); err != nil {
			h.logger.Error("failed to send telegram chat id", zap.Error(err), zap.Int64("chat_id", chatID))
		}
	default:
		// Ignore normal user messages.
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
