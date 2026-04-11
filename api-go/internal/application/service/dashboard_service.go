package service

import (
	"context"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// DashboardService handles business logic for dashboard data
type DashboardService struct {
	repo domain.DashboardRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(repo domain.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

// GetOverview retrieves the dashboard overview for an organization
func (s *DashboardService) GetOverview(ctx context.Context, organizationID string) (*domain.Dashboard, error) {
	return s.repo.GetOverview(ctx, organizationID)
}
