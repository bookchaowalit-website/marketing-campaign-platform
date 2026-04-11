package service

import (
	"context"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// CampaignService handles business logic for campaigns
type CampaignService struct {
	repo domain.CampaignRepository
}

// NewCampaignService creates a new campaign service
func NewCampaignService(repo domain.CampaignRepository) *CampaignService {
	return &CampaignService{repo: repo}
}

// CreateCampaign creates a new campaign
func (s *CampaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	return s.repo.Create(ctx, campaign)
}

// GetCampaign retrieves a campaign by ID
func (s *CampaignService) GetCampaign(ctx context.Context, id string) (*domain.Campaign, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllCampaigns retrieves all campaigns for an organization with pagination
func (s *CampaignService) GetAllCampaigns(ctx context.Context, organizationID string, limit, offset int) ([]*domain.Campaign, error) {
	return s.repo.GetAllByOrganizationID(ctx, organizationID, limit, offset)
}

// UpdateCampaign updates an existing campaign
func (s *CampaignService) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	campaign.UpdatedAt = time.Now()
	return s.repo.Update(ctx, campaign)
}

// DeleteCampaign deletes a campaign
func (s *CampaignService) DeleteCampaign(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetActiveCampaigns retrieves all active campaigns for an organization
func (s *CampaignService) GetActiveCampaigns(ctx context.Context, organizationID string, limit, offset int) ([]*domain.Campaign, error) {
	return s.repo.GetByStatus(ctx, organizationID, "active", limit, offset)
}

// PauseCampaign pauses an active campaign
func (s *CampaignService) PauseCampaign(ctx context.Context, id string) error {
	campaign, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	campaign.Status = "paused"
	campaign.UpdatedAt = time.Now()
	return s.repo.Update(ctx, campaign)
}

// ActivateCampaign activates a paused campaign
func (s *CampaignService) ActivateCampaign(ctx context.Context, id string) error {
	campaign, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	campaign.Status = "active"
	campaign.UpdatedAt = time.Now()
	return s.repo.Update(ctx, campaign)
}
