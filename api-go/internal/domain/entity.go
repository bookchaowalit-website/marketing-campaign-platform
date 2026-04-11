package domain

import "time"

// Campaign represents a marketing campaign entity
type Campaign struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`    // SEO, SEM, Social, Email, Display
	Status         string     `json:"status"`  // active, paused, completed, draft
	Budget         float64    `json:"budget"`
	Spent          float64    `json:"spent"`
	Impressions    int64      `json:"impressions"`
	Clicks         int64      `json:"clicks"`
	Conversions    int64      `json:"conversions"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CTR returns the click-through rate
func (c *Campaign) CTR() float64 {
	if c.Impressions == 0 {
		return 0
	}
	return float64(c.Clicks) / float64(c.Impressions) * 100
}

// ConversionRate returns the conversion rate
func (c *Campaign) ConversionRate() float64 {
	if c.Clicks == 0 {
		return 0
	}
	return float64(c.Conversions) / float64(c.Clicks) * 100
}

// BudgetUtilization returns the percentage of budget spent
func (c *Campaign) BudgetUtilization() float64 {
	if c.Budget == 0 {
		return 0
	}
	return c.Spent / c.Budget * 100
}

// Dashboard represents a marketing dashboard entity
type Dashboard struct {
	TotalCampaigns  int64                `json:"total_campaigns"`
	ActiveCampaigns int64                `json:"active_campaigns"`
	TotalSpend      float64              `json:"total_spend"`
	TotalRevenue    float64              `json:"total_revenue"`
	ROI             float64              `json:"roi"`
	ConversionRate  float64              `json:"conversion_rate"`
	TopChannels     []ChannelPerformance `json:"top_channels"`
	RecentCampaigns []Campaign           `json:"recent_campaigns"`
}

// ChannelPerformance represents performance by channel
type ChannelPerformance struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
