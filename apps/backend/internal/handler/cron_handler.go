package handler

import (
	"net/http"

	"github.com/devrapture/pod-events/internal/cron"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CronJobHandler struct {
	logger         *zap.Logger
	episodeChecker *cron.EpisodeChecker
}

func NewCronJobHandler(logger *zap.Logger, episodeChecker *cron.EpisodeChecker) *CronJobHandler {
	return &CronJobHandler{
		logger:         logger,
		episodeChecker: episodeChecker,
	}
}

// CheckEpisodes triggers a manual check for new podcast episodes across all tracked shows.
//
//	@Summary     Trigger episode check
//	@Description Manually trigger the cron job that polls Spotify for new episodes and sends notifications to subscribers
//	@Tags        Cron
//	@Produce     json
//	@Param       X-Cron-Secret header string true "Cron secret for authorization"
//	@Success     200 {object} response.APIResponse "Cron job completed"
//	@Failure     401 {object} response.APIResponse "Unauthorized"
//	@Failure     500 {object} response.APIResponse "Cron job failed"
//	@Router      /../cron/check-episodes [post]
func (h *CronJobHandler) CheckEpisodes(c *gin.Context) {
	result, err := h.episodeChecker.Run(c.Request.Context())
	if err != nil {
		h.logger.Error("cron job failed", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Cron job failed")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Cron job completed", result, nil)
}
