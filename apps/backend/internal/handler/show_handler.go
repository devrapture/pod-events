package handler

import (
	"net/http"
	"strconv"

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
	query := c.Query("q")

	limit := c.DefaultQuery("limit", "20")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 20
	}

	offset := c.DefaultQuery("offset", "0")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
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
