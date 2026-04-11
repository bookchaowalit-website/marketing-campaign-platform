package domain

import (
	"encoding/json"
	"time"
)

// ===========================
// AD PLATFORM
// ===========================

type AdPlatform struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Platform       string    `json:"platform"` // google_ads, facebook_ads, tiktok_ads, linkedin_ads, twitter_ads
	AccountID      string    `json:"account_id"`
	Status         string    `json:"status"`
	ConnectedAt    time.Time `json:"connected_at"`
	LastSyncAt     *time.Time `json:"last_sync_at,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ===========================
// AD CAMPAIGN
// ===========================

type AdCampaign struct {
	ID                string                 `json:"id"`
	OrganizationID    string                 `json:"organization_id"`
	UnifiedCampaignID *string                `json:"unified_campaign_id,omitempty"`
	PlatformID        string                 `json:"platform_id"`
	PlatformCampaignID *string               `json:"platform_campaign_id,omitempty"`
	Name              string                 `json:"name"`
	Status            string                 `json:"status"`
	Budget            float64                `json:"budget"`
	Spent             float64                `json:"spent"`
	Impressions       int64                  `json:"impressions"`
	Clicks            int64                  `json:"clicks"`
	Conversions       int64                  `json:"conversions"`
	CPC               float64                `json:"cpc"`
	CPM               float64                `json:"cpm"`
	CTR               float64                `json:"ctr"`
	StartDate         *time.Time             `json:"start_date,omitempty"`
	EndDate           *time.Time             `json:"end_date,omitempty"`
	Targeting         map[string]interface{} `json:"targeting,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// ===========================
// AD GROUP
// ===========================

type AdGroup struct {
	ID                 string     `json:"id"`
	CampaignID         string     `json:"campaign_id"`
	PlatformAdGroupID  *string    `json:"platform_ad_group_id,omitempty"`
	Name               string     `json:"name"`
	Status             string     `json:"status"`
	Budget             float64    `json:"budget"`
	Spent              float64    `json:"spent"`
	Impressions        int64      `json:"impressions"`
	Clicks             int64      `json:"clicks"`
	Conversions        int64      `json:"conversions"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ===========================
// AD CREATIVE
// ===========================

type AdCreative struct {
	ID                string     `json:"id"`
	AdGroupID         string     `json:"ad_group_id"`
	Type              string     `json:"type"` // text, image, video, carousel
	Headline          string     `json:"headline,omitempty"`
	Description       string     `json:"description,omitempty"`
	ImageURL          string     `json:"image_url,omitempty"`
	VideoURL          string     `json:"video_url,omitempty"`
	LandingPageURL    string     `json:"landing_page_url,omitempty"`
	Status            string     `json:"status"`
	ApprovalStatus    string     `json:"approval_status"`
	PlatformCreativeID *string   `json:"platform_creative_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// ===========================
// AD KEYWORD
// ===========================

type AdKeyword struct {
	ID           string     `json:"id"`
	AdGroupID    string     `json:"ad_group_id"`
	Keyword      string     `json:"keyword"`
	MatchType    string     `json:"match_type"` // exact, phrase, broad, negative
	BidAmount    float64    `json:"bid_amount"`
	MaxCPC       float64    `json:"max_cpc,omitempty"`
	Impressions  int64      `json:"impressions"`
	Clicks       int64      `json:"clicks"`
	Conversions  int64      `json:"conversions"`
	Cost         float64    `json:"cost"`
	AvgPosition  float64    `json:"avg_position,omitempty"`
	QualityScore float64    `json:"quality_score,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ===========================
// AD PERFORMANCE DAILY
// ===========================

type AdPerformanceDaily struct {
	ID          string    `json:"id"`
	AdID        string    `json:"ad_id"`
	Date        time.Time `json:"date"`
	Impressions int64     `json:"impressions"`
	Clicks      int64     `json:"clicks"`
	Conversions int64     `json:"conversions"`
	Cost        float64   `json:"cost"`
	CTR         float64   `json:"ctr"`
	CPC         float64   `json:"cpc"`
	CPA         float64   `json:"cpa"`
	ROAS        float64   `json:"roas"`
	CreatedAt   time.Time `json:"created_at"`
}

// ===========================
// REPOSITORY INTERFACES
// ===========================

type AdPlatformRepository interface {
	Create(p *AdPlatform) error
	GetByID(id string) (*AdPlatform, error)
	ListByOrg(orgID string) ([]*AdPlatform, error)
	Update(p *AdPlatform) error
	Delete(id string) error
	UpdateSyncStatus(id string, status string, errMsg string) error
}

type AdCampaignRepository interface {
	Create(c *AdCampaign) error
	GetByID(id string) (*AdCampaign, error)
	ListByOrg(orgID, status string, limit, offset int) ([]*AdCampaign, int, error)
	ListByPlatform(platformID string) ([]*AdCampaign, error)
	Update(c *AdCampaign) error
	Delete(id string) error
	UpdateStatus(id, status string) error
}

type AdGroupRepository interface {
	Create(g *AdGroup) error
	GetByID(id string) (*AdGroup, error)
	ListByCampaign(campaignID string) ([]*AdGroup, error)
	Update(g *AdGroup) error
	Delete(id string) error
}

type AdCreativeRepository interface {
	Create(c *AdCreative) error
	GetByID(id string) (*AdCreative, error)
	ListByAdGroup(adGroupID string) ([]*AdCreative, error)
	Update(c *AdCreative) error
	Delete(id string) error
}

type AdKeywordRepository interface {
	Create(k *AdKeyword) error
	GetByID(id string) (*AdKeyword, error)
	ListByAdGroup(adGroupID string) ([]*AdKeyword, error)
	Update(k *AdKeyword) error
	Delete(id string) error
}

type AdPerformanceRepository interface {
	Create(p *AdPerformanceDaily) error
	GetByAdID(adID string, startDate, endDate time.Time) ([]*AdPerformanceDaily, error)
	Upsert(p *AdPerformanceDaily) error
}

// ===========================
// VALIDATION CONSTANTS
// ===========================

var ValidAdPlatforms = map[string]bool{
	"google_ads":     true,
	"facebook_ads":   true,
	"tiktok_ads":     true,
	"linkedin_ads":   true,
	"twitter_ads":    true,
}

var ValidAdCreativeTypes = map[string]bool{
	"text":     true,
	"image":    true,
	"video":    true,
	"carousel": true,
}

var ValidKeywordMatchTypes = map[string]bool{
	"exact":   true,
	"phrase":  true,
	"broad":   true,
	"negative": true,
}

// Helper method for JSON marshaling of targeting
func (c *AdCampaign) TargetingJSON() json.RawMessage {
	if c.Targeting == nil {
		return json.RawMessage("{}")
	}
	b, _ := json.Marshal(c.Targeting)
	return json.RawMessage(b)
}
