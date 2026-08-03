package handler

import (
	"net/http"

	"github.com/devrapture/pod-events/internal/dto"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DashboardShowHandler struct {
	dashboardShowService services.DashboardSummaryService
	logger               *zap.Logger
}

func NewDashboardShowHandler(dashboardShowService services.DashboardSummaryService, logger *zap.Logger) *DashboardShowHandler {
	return &DashboardShowHandler{
		dashboardShowService: dashboardShowService,
		logger:               logger,
	}
}

// GetDashboardSummary returns the dashboard overview summary.
//
//	@Summary     Get dashboard summary
//	@Description Get the dashboard overview with setup progress and statistics
//	@Tags        Dashboard
//	@Security    BearerAuth
//	@Produce     json
//	@Success     200 {object} response.APIResponse{data=dto.DashboardSummaryDTO} "dashboard summary fetched successfully"
//	@Failure     500 {object} response.APIResponse "failed to get dashboard summary"
//	@Router      /dashboard/summary [get]
func (h *DashboardShowHandler) GetDashboardSummary(c *gin.Context) {
	userID, _ := c.Get("userID")
	var summary *dto.DashboardSummaryDTO
	summary, err := h.dashboardShowService.GetDashboardSummary(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("failed to get dashboard summary", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get dashboard summary")
		return
	}
	response.SuccessResponse(c, http.StatusOK, "dashboard summary fetched successfully", summary, nil)
}
