package handler

import (
	"errors"
	"net/http"
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

// GetUserSavedShows returns a list of shows saved by a user on Spotify.
//
//	@Summary     Get saved shows
//	@Description Get the current user's saved/podcasts from Spotify
//	@Tags        Shows
//	@Security    BearerAuth
//	@Param       q query string false "Search query to filter saved shows"
//	@Success     200 {object} response.APIResponse{data=[]dto.SavedShowResponse} "show fetched successfully"
//	@Router      /shows/saved [get]
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

// SearchShows searches for shows on Spotify.
//
//	@Summary     Search shows
//	@Description Search for podcast shows on Spotify
//	@Tags        Shows
//	@Security    BearerAuth
//	@Param       q      query string true  "Search query"
//	@Param       limit  query int    false "Maximum results (default 10)"
//	@Param       offset query int    false "Result offset (default 0)"
//	@Success     200 {object} response.APIResponse{data=[]dto.SavedShowResponse} "show fetched successfully"
//	@Failure     400 {object} response.APIResponse "invalid parameters"
//	@Router      /shows/search [get]
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
//	@Summary     Subscribe to a show
//	@Description Subscribe the current user to a podcast show by Spotify show ID
//	@Tags        Subscriptions
//	@Security    BearerAuth
//	@Param       spotifyShowId path string true "Spotify show ID (22 characters)"
//	@Success     200 {object} response.APIResponse{data=dto.SubscriptionResponse} "subscribed successfully"
//	@Failure     400 {object} response.APIResponse "invalid spotify show ID"
//	@Failure     404 {object} response.APIResponse "podcast show not found"
//	@Failure     409 {object} response.APIResponse "already subscribed"
//	@Router      /shows/{spotifyShowId}/subscribe [post]
func (h *ShowHandler) Subscribe(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.SubscribeShowsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	_, err := h.showService.Subscribe(c.Request.Context(), userID.(uuid.UUID), req.SpotifyShowIDs)
	if err != nil {

		if errors.Is(err, apperrors.ErrPodcastShowNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, "podcast show not found")
			return
		}

		if errors.Is(err, apperrors.ErrSubscriptionAlreadyExists) {
			response.ErrorResponse(c, http.StatusConflict, "already subscribed to this show")
			return
		}

		h.logger.Error("failed to subscribe", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to subscribe")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "subscribed successfully", nil, nil)
}

// Unsubscribe removes a subscription by its UUID.
//
//	@Summary     Unsubscribe from a show
//	@Description Remove a subscription by its UUID
//	@Tags        Subscriptions
//	@Security    BearerAuth
//	@Param       id path string true "Subscription UUID"
//	@Success     200 {object} response.APIResponse "unsubscribed successfully"
//	@Failure     400 {object} response.APIResponse "invalid subscription ID"
//	@Failure     404 {object} response.APIResponse "subscription not found"
//	@Router      /subscriptions/{id} [delete]
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
//	@Summary     Get subscriptions
//	@Description Get all of the current user's active subscriptions
//	@Tags        Subscriptions
//	@Security    BearerAuth
//	@Success     200 {object} response.APIResponse{data=[]dto.SubscriptionResponse} "subscriptions fetched successfully"
//	@Router      /subscriptions [get]
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
