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

// BudgetPeriod represents a budget planning period
type BudgetPeriod struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	PeriodType     string    `json:"period_type"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	TotalBudget    float64   `json:"total_budget"`
	TotalAllocated float64   `json:"total_allocated"`
	TotalSpent     float64   `json:"total_spent"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes,omitempty"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BudgetAllocation represents budget allocated to a channel
type BudgetAllocation struct {
	ID              string  `json:"id"`
	PeriodID        string  `json:"period_id"`
	Channel         string  `json:"channel"`
	ChannelCategory string  `json:"channel_category"`
	AllocatedAmount float64 `json:"allocated_amount"`
	SpentAmount     float64 `json:"spent_amount"`
	ForecastAmount  float64 `json:"forecast_amount,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BudgetTransaction tracks actual spend
type BudgetTransaction struct {
	ID            string    `json:"id"`
	OrganizationID string   `json:"organization_id"`
	CampaignID    *string   `json:"campaign_id,omitempty"`
	AllocationID  *string   `json:"allocation_id,omitempty"`
	Channel       string    `json:"channel"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	TransactionDate time.Time `json:"transaction_date"`
	Description   string    `json:"description,omitempty"`
	InvoiceNumber string    `json:"invoice_number,omitempty"`
	Vendor        string    `json:"vendor,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CalendarEvent represents a marketing calendar event
type CalendarEvent struct {
	ID          string    `json:"id"`
	OrganizationID string `json:"organization_id"`
	CampaignID  *string   `json:"campaign_id,omitempty"`
	EventType   string    `json:"event_type"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time,omitempty"`
	AllDay      bool      `json:"all_day"`
	Color       string    `json:"color"`
	Channel     string    `json:"channel,omitempty"`
	Status      string    `json:"status"`
	Location    string    `json:"location,omitempty"`
	LocationURL string    `json:"location_url,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ApprovalWorkflow defines an approval process
type ApprovalWorkflow struct {
	ID              string   `json:"id"`
	OrganizationID  string   `json:"organization_id"`
	Name            string   `json:"name"`
	Description     string   `json:"description,omitempty"`
	TriggerType     string   `json:"trigger_type"`
	Channels        []string `json:"channels,omitempty"`
	MinApprovers    int      `json:"min_approvers"`
	ApprovalOrder   string   `json:"approval_order"`
	AutoApproveAfterHours int `json:"auto_approve_after_hours"`
	Status          string   `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ApprovalRequest represents a pending approval
type ApprovalRequest struct {
	ID           string    `json:"id"`
	OrganizationID string `json:"organization_id"`
	WorkflowID   *string   `json:"workflow_id,omitempty"`
	EntityType   string    `json:"entity_type"`
	EntityID     string    `json:"entity_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	SubmittedBy  string    `json:"submitted_by"`
	SubmittedAt  time.Time `json:"submitted_at"`
	Status       string    `json:"status"`
	Deadline     time.Time `json:"deadline,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	Comments     string    `json:"comments,omitempty"`
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
