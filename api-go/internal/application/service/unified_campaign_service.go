package service

import (
	"context"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type UnifiedCampaignService struct {
	repo domain.UnifiedCampaignRepository
}

func NewUnifiedCampaignService(repo domain.UnifiedCampaignRepository) *UnifiedCampaignService {
	return &UnifiedCampaignService{repo: repo}
}

func (s *UnifiedCampaignService) CreateCampaign(ctx context.Context, c *domain.UnifiedCampaign) error {
	if !domain.ValidChannels[c.Channel] {
		return domain.ErrInvalidChannel
	}
	if !domain.ValidChannelCategories[c.ChannelCategory] {
		return domain.ErrInvalidChannelCategory
	}
	return s.repo.Create(c)
}

func (s *UnifiedCampaignService) GetCampaign(ctx context.Context, id string) (*domain.UnifiedCampaign, error) {
	return s.repo.GetByID(id)
}

func (s *UnifiedCampaignService) ListCampaigns(ctx context.Context, orgID, channel, status string, limit, offset int) ([]*domain.UnifiedCampaign, int, error) {
	return s.repo.List(orgID, channel, status, limit, offset)
}

func (s *UnifiedCampaignService) UpdateCampaign(ctx context.Context, c *domain.UnifiedCampaign) error {
	return s.repo.Update(c)
}

func (s *UnifiedCampaignService) DeleteCampaign(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}

func (s *UnifiedCampaignService) SubmitForApproval(ctx context.Context, id string) error {
	return s.repo.SubmitForApproval(id)
}

func (s *UnifiedCampaignService) ApproveCampaign(ctx context.Context, id, approverID string) error {
	return s.repo.Approve(id, approverID)
}

func (s *UnifiedCampaignService) RejectCampaign(ctx context.Context, id, reason string) error {
	return s.repo.Reject(id, reason)
}
