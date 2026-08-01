package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/devrapture/pod-events/internal/cron"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CronJobHandler struct {
	logger         *zap.Logger
	episodeChecker *cron.EpisodeChecker
	runMu          sync.Mutex
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
//	@Success     202 {object} response.APIResponse "Cron job started or already running"
//	@Failure     401 {object} response.APIResponse "Unauthorized"
//	@Router      /../cron/check-episodes [post]
func (h *CronJobHandler) CheckEpisodes(c *gin.Context) {
	// TryLock returns false if another run already has the lock.
	if !h.runMu.TryLock() {
		response.SuccessResponse(
			c,
			http.StatusAccepted,
			"Cron job already running",
			nil,
			nil,
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

	go func() {
		defer h.runMu.Unlock() // Allow the next run when this one finishes.
		defer cancel()         // Clean up the timeout context.

		result, err := h.episodeChecker.Run(ctx)
		if err != nil {
			h.logger.Error("cron job failed", zap.Error(err))
			return
		}

		h.logger.Info("cron job completed", zap.Any("result", result))
	}()

	response.SuccessResponse(
		c,
		http.StatusAccepted,
		"Cron job started",
		nil,
		nil,
	)
}
