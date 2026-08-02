package services

import (
	"context"

	"github.com/devrapture/pod-events/internal/dto"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/google/uuid"
)

type DashboardSummaryService interface {
	GetDashboardSummary(ctx context.Context, userID uuid.UUID) (*dto.DashboardSummaryDTO, error)
}

type dashboardSummaryService struct {
	dashboardSummaryRepository repositories.DashboardSummaryRepository
}

func NewDashboardSummaryService(dashboardSummaryRepository repositories.DashboardSummaryRepository) DashboardSummaryService {
	return &dashboardSummaryService{
		dashboardSummaryRepository: dashboardSummaryRepository,
	}
}

func (s *dashboardSummaryService) GetDashboardSummary(ctx context.Context, userID uuid.UUID) (*dto.DashboardSummaryDTO, error) {
	return s.dashboardSummaryRepository.GetDashboardSummary(ctx, userID)
}
