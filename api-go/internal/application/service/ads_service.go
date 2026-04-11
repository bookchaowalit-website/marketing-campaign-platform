package service

import (
	"context"
	"errors"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

var (
	ErrInvalidPlatform   = errors.New("invalid ad platform")
	ErrInvalidCreativType = errors.New("invalid creative type")
	ErrInvalidMatchType  = errors.New("invalid keyword match type")
)

type PPCAdsService struct {
	platformRepo   domain.AdPlatformRepository
	campaignRepo   domain.AdCampaignRepository
	groupRepo      domain.AdGroupRepository
	creativeRepo   domain.AdCreativeRepository
	keywordRepo    domain.AdKeywordRepository
	performanceRepo domain.AdPerformanceRepository
}

func NewPPCAdsService(
	platformRepo domain.AdPlatformRepository,
	campaignRepo domain.AdCampaignRepository,
	groupRepo domain.AdGroupRepository,
	creativeRepo domain.AdCreativeRepository,
	keywordRepo domain.AdKeywordRepository,
	performanceRepo domain.AdPerformanceRepository,
) *PPCAdsService {
	return &PPCAdsService{
		platformRepo:    platformRepo,
		campaignRepo:    campaignRepo,
		groupRepo:       groupRepo,
		creativeRepo:    creativeRepo,
		keywordRepo:     keywordRepo,
		performanceRepo: performanceRepo,
	}
}

// ===========================
// PLATFORM OPERATIONS
// ===========================

func (s *PPCAdsService) ConnectPlatform(ctx context.Context, p *domain.AdPlatform) error {
	if !domain.ValidAdPlatforms[p.Platform] {
		return ErrInvalidPlatform
	}
	p.Status = "connected"
	return s.platformRepo.Create(p)
}

func (s *PPCAdsService) ListPlatforms(ctx context.Context, orgID string) ([]*domain.AdPlatform, error) {
	return s.platformRepo.ListByOrg(orgID)
}

func (s *PPCAdsService) GetPlatform(ctx context.Context, id string) (*domain.AdPlatform, error) {
	return s.platformRepo.GetByID(id)
}

// ===========================
// CAMPAIGN OPERATIONS
// ===========================

func (s *PPCAdsService) CreateCampaign(ctx context.Context, c *domain.AdCampaign) error {
	return s.campaignRepo.Create(c)
}

func (s *PPCAdsService) GetCampaign(ctx context.Context, id string) (*domain.AdCampaign, error) {
	return s.campaignRepo.GetByID(id)
}

func (s *PPCAdsService) ListCampaigns(ctx context.Context, orgID, status string, limit, offset int) ([]*domain.AdCampaign, int, error) {
	return s.campaignRepo.ListByOrg(orgID, status, limit, offset)
}

func (s *PPCAdsService) PauseCampaign(ctx context.Context, id string) error {
	return s.campaignRepo.UpdateStatus(id, "paused")
}

func (s *PPCAdsService) ActivateCampaign(ctx context.Context, id string) error {
	return s.campaignRepo.UpdateStatus(id, "active")
}

func (s *PPCAdsService) SyncCampaign(ctx context.Context, id string) error {
	// In production, this would call the platform API to sync data
	// For now, update the last sync timestamp at platform level
	campaign, err := s.campaignRepo.GetByID(id)
	if err != nil {
		return err
	}
	// Trigger platform sync
	return s.platformRepo.UpdateSyncStatus(campaign.PlatformID, "connected", "")
}

func (s *PPCAdsService) GetCampaignPerformance(ctx context.Context, campaignID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	_, err := s.campaignRepo.GetByID(campaignID)
	if err != nil {
		return nil, err
	}

	// In production, aggregate from ad_performance_daily
	// For now, return campaign stats
	campaign, err := s.campaignRepo.GetByID(campaignID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"campaign":    campaign,
		"start_date":  startDate.Format("2006-01-02"),
		"end_date":    endDate.Format("2006-01-02"),
		"impressions": campaign.Impressions,
		"clicks":      campaign.Clicks,
		"conversions": campaign.Conversions,
		"spent":       campaign.Spent,
		"cpc":         campaign.CPC,
		"cpm":         campaign.CPM,
		"ctr":         campaign.CTR,
	}, nil
}

// ===========================
// AD GROUP OPERATIONS
// ===========================

func (s *PPCAdsService) CreateAdGroup(ctx context.Context, g *domain.AdGroup) error {
	return s.groupRepo.Create(g)
}

func (s *PPCAdsService) GetAdGroup(ctx context.Context, id string) (*domain.AdGroup, error) {
	return s.groupRepo.GetByID(id)
}

func (s *PPCAdsService) ListAdGroups(ctx context.Context, campaignID string) ([]*domain.AdGroup, error) {
	return s.groupRepo.ListByCampaign(campaignID)
}

// ===========================
// AD CREATIVE OPERATIONS
// ===========================

func (s *PPCAdsService) CreateAdCreative(ctx context.Context, c *domain.AdCreative) error {
	if !domain.ValidAdCreativeTypes[c.Type] {
		return ErrInvalidCreativType
	}
	return s.creativeRepo.Create(c)
}

func (s *PPCAdsService) GetAdCreative(ctx context.Context, id string) (*domain.AdCreative, error) {
	return s.creativeRepo.GetByID(id)
}

func (s *PPCAdsService) ListAdCreatives(ctx context.Context, adGroupID string) ([]*domain.AdCreative, error) {
	return s.creativeRepo.ListByAdGroup(adGroupID)
}

// ===========================
// KEYWORD OPERATIONS
// ===========================

func (s *PPCAdsService) CreateKeyword(ctx context.Context, k *domain.AdKeyword) error {
	if !domain.ValidKeywordMatchTypes[k.MatchType] {
		return ErrInvalidMatchType
	}
	return s.keywordRepo.Create(k)
}

func (s *PPCAdsService) GetKeyword(ctx context.Context, id string) (*domain.AdKeyword, error) {
	return s.keywordRepo.GetByID(id)
}

func (s *PPCAdsService) ListKeywords(ctx context.Context, adGroupID string) ([]*domain.AdKeyword, error) {
	return s.keywordRepo.ListByAdGroup(adGroupID)
}

// ===========================
// PERFORMANCE OPERATIONS
// ===========================

func (s *PPCAdsService) RecordPerformance(ctx context.Context, p *domain.AdPerformanceDaily) error {
	return s.performanceRepo.Upsert(p)
}

func (s *PPCAdsService) GetPerformance(ctx context.Context, adID string, startDate, endDate time.Time) ([]*domain.AdPerformanceDaily, error) {
	return s.performanceRepo.GetByAdID(adID, startDate, endDate)
}
