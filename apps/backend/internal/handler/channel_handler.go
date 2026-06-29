package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/devrapture/pod-events/internal/dto"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChannelHandler struct {
	channelService services.ChannelServices
	logger         *zap.Logger
}

func NewChannelHandler(channelService services.ChannelServices, logger *zap.Logger) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
		logger:         logger,
	}
}

func (h *ChannelHandler) CreateChannel(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	channel, err := h.channelService.Create(c.Request.Context(), userID.(uuid.UUID), req)
	if err != nil {
		if isChannelValidationError(err) {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		h.logger.Error("error creating channel", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create channel")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Channel created", channel, nil)
}

func (h *ChannelHandler) GetChannels(c *gin.Context) {
	userID, _ := c.Get("userID")
	channels, err := h.channelService.GetByUser(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch channels")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Channels fetched", channels, nil)
}

func (h *ChannelHandler) ToggleActive(c *gin.Context) {
	userID, _ := c.Get("userID")
	channelIDStr := c.Param("channelID")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	var req dto.ToggleActiveChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	if err := h.channelService.ToggleActive(c.Request.Context(), userID.(uuid.UUID), channelID, *req.IsActive); err != nil {
		if errors.Is(err, apperrors.ErrChannelIDNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, "Channel not found")
			return
		}
		h.logger.Error("error toggling channel", zap.Any("channelID", channelID), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update channel status")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Channel status updated", nil, nil)
}

func (h *ChannelHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	channelIDStr := c.Param("channelID")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid channel ID")
		return
	}
	if err := h.channelService.DeleteChannel(c.Request.Context(), userID.(uuid.UUID), channelID); err != nil {
		if errors.Is(err, apperrors.ErrChannelIDNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, "Channel not found")
			return
		}
		h.logger.Error("error deleting channel", zap.Any("channelID", channelID), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Error deleting channel")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Channel deleted", nil, nil)
}

func isChannelValidationError(err error) bool {
	msg := err.Error()

	return strings.Contains(msg, "invalid channel_type") ||
		strings.Contains(msg, "slack destination") ||
		strings.Contains(msg, "discord destination") ||
		strings.Contains(msg, "telegram destination") ||
		strings.Contains(msg, "whatsapp destination")
}
