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

func (h *CronJobHandler) CheckEpisodes(c *gin.Context) {
	result, err := h.episodeChecker.Run(c.Request.Context())
	if err != nil {
		h.logger.Error("cron job failed", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Cron job failed")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Cron job completed", result, nil)
}
