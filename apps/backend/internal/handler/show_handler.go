package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/devrapture/pod-events/internal/dto"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/internal/spotify"
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

// Subscribe subscribes the current user to one or more podcast shows.
//
//	@Summary     Subscribe to shows
//	@Description Subscribe the current user to one or more podcast shows by Spotify show ID
//	@Tags        Subscriptions
//	@Security    BearerAuth
//	@Param       request body dto.SubscribeShowsRequest true "Spotify show IDs"
//	@Success     200 {object} response.APIResponse "subscribed successfully"
//	@Failure     400 {object} response.APIResponse "invalid or too many Spotify show IDs"
//	@Failure     422 {object} response.APIResponse "invalid request"
//	@Failure     404 {object} response.APIResponse "podcast show not found"
//	@Failure     409 {object} response.APIResponse "already subscribed"
//	@Failure     424 {object} response.APIResponse "Spotify authorization required"
//	@Failure     429 {object} response.APIResponse "Spotify rate limit exceeded"
//	@Failure     503 {object} response.APIResponse "Spotify unavailable"
//	@Failure     500 {object} response.APIResponse "internal server error"
//	@Router      /shows/subscribe [post]
func (h *ShowHandler) Subscribe(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.SubscribeShowsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	_, err := h.showService.Subscribe(c.Request.Context(), userID.(uuid.UUID), req.SpotifyShowIDs)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidSpotifyShowIDs):
			response.ErrorResponse(c, http.StatusBadRequest, "provide at least one valid Spotify show ID")
			return
		case errors.Is(err, apperrors.ErrTooManySpotifyShowIDs):
			response.ErrorResponse(c, http.StatusBadRequest, "you can subscribe to up to 50 shows at once")
			return
		case errors.Is(err, apperrors.ErrPodcastShowNotFound):
			response.ErrorResponse(c, http.StatusNotFound, "podcast show not found")
			return
		case errors.Is(err, apperrors.ErrSubscriptionAlreadyExists):
			response.ErrorResponse(c, http.StatusConflict, "you are already subscribed to one or more selected shows")
			return
		case errors.Is(err, apperrors.ErrorSpotifyTokenNotFound), errors.Is(err, apperrors.ErrSpotifyAuthorizationRequired):
			response.ErrorResponse(c, http.StatusFailedDependency, "your Spotify connection has expired; reconnect Spotify and try again")
			return
		}

		var rateLimitErr *spotify.RateLimitError
		if errors.As(err, &rateLimitErr) {
			c.Header("Retry-After", strconv.Itoa(rateLimitErr.RetryAfter))
			response.ErrorResponse(c, http.StatusTooManyRequests, fmt.Sprintf("Spotify is rate limiting requests; try again in %d seconds", rateLimitErr.RetryAfter))
			return
		}

		if errors.Is(err, apperrors.ErrSpotifyUnavailable) {
			response.ErrorResponse(c, http.StatusServiceUnavailable, "Spotify could not process the request; try again later")
			return
		}

		h.logger.Error("failed to subscribe", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "subscriptions could not be saved because of an internal server error; try again later")
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
