package handler

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/devrapture/pod-events/internal/cron"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const cronJobTimeout = 25 * time.Second

type episodeRunner interface {
	Run(ctx context.Context) (*cron.CheckResult, error)
}

type CronJobHandler struct {
	logger         *zap.Logger
	episodeChecker episodeRunner
	timeout        time.Duration
	runMu          sync.Mutex
}

func NewCronJobHandler(logger *zap.Logger, episodeChecker *cron.EpisodeChecker) *CronJobHandler {
	return &CronJobHandler{
		logger:         logger,
		episodeChecker: episodeChecker,
		timeout:        cronJobTimeout,
	}
}

// CheckEpisodes triggers a manual check for new podcast episodes across all tracked shows.
//
//	@Summary     Trigger episode check
//	@Description Manually trigger the cron job that polls Spotify for new episodes and sends notifications to subscribers
//	@Tags        Cron
//	@Produce     json
//	@Param       X-Cron-Secret header string true "Cron secret for authorization"
//	@Success     200 {object} response.APIResponse{data=cron.CheckResult} "Cron job completed"
//	@Success     202 {object} response.APIResponse "Cron job already running"
//	@Failure     401 {object} response.APIResponse "Unauthorized"
//	@Failure     500 {object} response.APIResponse "Cron job failed"
//	@Failure     504 {object} response.APIResponse "Cron job timed out"
//	@Router       /cron/check-episodes [post]
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
	defer h.runMu.Unlock()

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	result, err := h.episodeChecker.Run(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			h.logger.Error("cron job timed out", zap.Error(err))
			response.ErrorResponse(c, http.StatusGatewayTimeout, "Cron job timed out")
			return
		}

		h.logger.Error("cron job failed", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Cron job failed")
		return
	}

	h.logger.Info("cron job completed", zap.Any("result", result))

	response.SuccessResponse(
		c,
		http.StatusOK,
		"Cron job completed",
		result,
		nil,
	)
}
