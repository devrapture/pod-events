package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/devrapture/pod-events/internal/dto"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ShowHandler struct {
	showService services.ShowServices
	logger      *zap.Logger
}

func NewShowHandler(showService services.ShowServices, logger *zap.Logger) *ShowHandler {
	return &ShowHandler{
		showService: showService,
		logger:      logger,
	}
}

var spotifyShowIDPattern = regexp.MustCompile(`^[A-Za-z0-9]{22}$`)

// GetUserSavedShows returns a list of shows saved by a user on Spotify.
// GET /shows/saved
func (h *ShowHandler) GetUserSavedShows(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("q")
	show, err := h.showService.GetUserSavedShows(c.Request.Context(), userID.(uuid.UUID), query)
	if err != nil {
		h.logger.Error("failed to get saved shows", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get saved shows")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "show fetched successfully", show, nil)
}

// SearchShows searches for shows on Spotify".
// GET /shows/search
func (h *ShowHandler) SearchShows(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := strings.TrimSpace(c.Query("q"))

	limit := c.DefaultQuery("limit", "10")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	if limitInt <= 0 {
		response.ErrorResponse(c, http.StatusBadRequest, "limit must be > 0")
		return
	}

	offset := c.DefaultQuery("offset", "0")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	if offsetInt < 0 {
		response.ErrorResponse(c, http.StatusBadRequest, "offset must be >= 0")
		return
	}

	if query == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "query is required")
		return
	}
	show, err := h.showService.SearchShows(c.Request.Context(), userID.(uuid.UUID), query, limitInt, offsetInt)
	if err != nil {
		h.logger.Error("failed to search shows", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to search shows")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "show fetched successfully", show, nil)
}

// Subscribe subscribes the current user to a podcast show.
//
// POST /api/shows/:spotifyShowId/subscribe
func (h *ShowHandler) Subscribe(c *gin.Context) {
	userID, _ := c.Get("userID")
	spotifyShowID := strings.TrimSpace(c.Param("spotifyShowId"))
	if spotifyShowID == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "spotify show ID is required")
		return
	}

	if !spotifyShowIDPattern.MatchString(spotifyShowID) {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid spotify show ID")
		return
	}
	subscription, err := h.showService.Subscribe(c.Request.Context(), userID.(uuid.UUID), spotifyShowID)
	if err != nil {
		if errors.Is(err, apperrors.ErrSubscriptionAlreadyExists) {
			response.ErrorResponse(c, http.StatusConflict, "already subscribed to this show")
			return
		}
		h.logger.Error("failed to subscribe", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to subscribe")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "subscribed successfully", subscription, nil)
}

// Unsubscribe removes a subscription by its UUID.
//
// DELETE /api/subscriptions/:id
func (h *ShowHandler) Unsubscribe(c *gin.Context) {
	userID, _ := c.Get("userID")
	subscriptionIDStr := c.Param("id")
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid subscription ID")
		return
	}
	if err := h.showService.Unsubscribe(c.Request.Context(), userID.(uuid.UUID), subscriptionID); err != nil {
		if errors.Is(err, apperrors.ErrSubscriptionNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, "subscription not found")
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to unsubscribe")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "unsubscribed successfully", nil, nil)
}

// GetSubscriptions returns all of the user's active subscriptions.
//
// GET /api/subscriptions
func (h *ShowHandler) GetSubscriptions(c *gin.Context) {
	userID, _ := c.Get("userID")
	subscriptions, err := h.showService.GetSubscriptions(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("failed to get subscriptions", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get subscriptions")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "subscriptions fetched successfully", dto.ToSubscriptionResponses(subscriptions), nil)
}
