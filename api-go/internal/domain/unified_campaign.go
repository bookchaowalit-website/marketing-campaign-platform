package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidChannel       = errors.New("invalid channel")
	ErrInvalidChannelCategory = errors.New("invalid channel category")
)

// UnifiedCampaign represents a campaign across ALL channels (online + offline)
type UnifiedCampaign struct {
	ID               string    `json:"id"`
	OrganizationID   string    `json:"organization_id"`
	GroupID          *string   `json:"group_id,omitempty"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	ChannelCategory  string    `json:"channel_category"` // online, offline, hybrid
	Channel          string    `json:"channel"`
	Budget           float64   `json:"budget"`
	Spent            float64   `json:"spent"`
	CPA              float64   `json:"cost_per_acquisition,omitempty"`
	ROAS             float64   `json:"return_on_ad_spend,omitempty"`
	TargetAudience   string    `json:"target_audience,omitempty"`
	Geography        []string  `json:"geography,omitempty"`
	Status           string    `json:"status"`
	StartDate        time.Time `json:"start_date,omitempty"`
	EndDate          time.Time `json:"end_date,omitempty"`
	Timezone         string    `json:"timezone"`
	UTMSource        string    `json:"utm_source,omitempty"`
	UTMMedium        string    `json:"utm_medium,omitempty"`
	UTMCampaign      string    `json:"utm_campaign,omitempty"`
	TrackingURL      string    `json:"tracking_url,omitempty"`
	CreativeBrief    string    `json:"creative_brief,omitempty"`
	ApprovalStatus   string    `json:"approval_status"`
	ApprovedBy       *string   `json:"approved_by,omitempty"`
	ApprovedAt       time.Time `json:"approved_at,omitempty"`
	RejectionReason  string    `json:"rejection_reason,omitempty"`
	Impressions      int64     `json:"impressions"`
	Clicks           int64     `json:"clicks"`
	Conversions      int64     `json:"conversions"`
	Revenue          float64   `json:"revenue"`
	Leads            int       `json:"leads"`
	OwnerID          string    `json:"owner_id"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CampaignGroup groups related cross-channel campaigns
type CampaignGroup struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	Objective      string    `json:"objective"`
	Status         string    `json:"status"`
	StartDate      time.Time `json:"start_date,omitempty"`
	EndDate        time.Time `json:"end_date,omitempty"`
	TotalBudget    float64   `json:"total_budget"`
	TotalSpent     float64   `json:"total_spent"`
	OwnerID        string    `json:"owner_id"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Channel definitions for validation
var ValidChannels = map[string]bool{
	// Online
	"google_ads": true, "facebook_ads": true, "instagram_ads": true, "tiktok_ads": true,
	"linkedin_ads": true, "twitter_ads": true, "seo": true, "email": true, "sms": true,
	"push_notifications": true, "video_youtube": true, "video_vimeo": true, "webinar": true,
	"content_marketing": true, "affiliate": true, "influencer": true, "display_ads": true,
	"native_ads": true, "podcast": true,
	// Offline
	"tv": true, "radio": true, "print_magazine": true, "print_newspaper": true,
	"billboard": true, "transit_ad": true, "direct_mail": true, "event": true,
	"trade_show": true, "conference": true, "sponsorship": true, "telemarketing": true,
	"pos_materials": true, "in_store_promotion": true, "brochure": true, "flyer": true,
	"guerrilla": true,
	// Hybrid
	"cross_channel": true, "omnichannel": true,
}

var ValidChannelCategories = map[string]bool{
	"online":  true,
	"offline": true,
	"hybrid":  true,
}

// UnifiedCampaignRepository defines the interface for unified campaigns
type UnifiedCampaignRepository interface {
	Create(c *UnifiedCampaign) error
	GetByID(id string) (*UnifiedCampaign, error)
	List(orgID, channel, status string, limit, offset int) ([]*UnifiedCampaign, int, error)
	Update(c *UnifiedCampaign) error
	Delete(id string) error
	SubmitForApproval(id string) error
	Approve(id, approverID string) error
	Reject(id, reason string) error
}
