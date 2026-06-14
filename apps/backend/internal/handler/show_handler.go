package handler

import (
	"net/http"

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
