package handler

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/devrapture/pod-events/internal/config"
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

func (h *TelegramWebHookHandler) CreateConnectLink(c *gin.Context) {
	userID, _ := c.Get("userID")
	link, err := h.telegramConnectionService.CreateConnectLink(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("failed to create telegram connect link", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create telegram connect link")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Telegram connect link created", gin.H{"url": link}, nil)
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
			h.logger.Warn("failed to complete telegram connection", zap.Error(err))
			_ = h.notifier.SendToChatID(
				c.Request.Context(),
				"This PodEvents connection link is invalid or expired. Please generate a new one from the app.",
				chatID,
			)
			c.Status(http.StatusOK)
			return
		}

		_ = h.notifier.SendToChatID(
			c.Request.Context(),
			fmt.Sprintf("Hi %s, PodEvents is now connected to this Telegram chat.", name),
			chatID,
		)

		c.Status(http.StatusOK)
		return
	}
	// name := strings.TrimSpace(strings.Join([]string{
	// 	update.Message.From.FirstName,
	// 	update.Message.From.LastName,
	// }, " "))
	// if name == "" {
	// 	name = "there"
	// }
	// message := fmt.Sprintf("Hi %s, this is PodEvents. Your Telegram chat ID is %d.", name, chatID)
	// switch text {
	// case cmdStart:
	// 	if err := h.notifier.SendToChatID(c.Request.Context(), message, chatID); err != nil {
	// 		h.logger.Error("failed to send telegram chat id", zap.Error(err), zap.Int64("chat_id", chatID))
	// 	}
	// case "/chatid":
	// 	if err := h.notifier.SendToChatID(c.Request.Context(), message, chatID); err != nil {
	// 		h.logger.Error("failed to send telegram chat id", zap.Error(err), zap.Int64("chat_id", chatID))
	// 	}
	// default:
	// 	// Ignore normal user messages.
	// }

	c.Status(http.StatusOK)
}

func (h *TelegramWebHookHandler) validSecret(got string) bool {
	expected := h.cfg.TelegramWebhookSecret
	if expected == "" || got == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1
}
