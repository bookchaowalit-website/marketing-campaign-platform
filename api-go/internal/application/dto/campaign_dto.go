package dto

// CreateCampaignRequest represents the request to create a campaign
type CreateCampaignRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"` // SEO, SEM, Social, Email, Display
	Budget      float64 `json:"budget" binding:"required,min=0"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

// UpdateCampaignRequest represents the request to update a campaign
type UpdateCampaignRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Budget      float64 `json:"budget"`
	Spent       float64 `json:"spent"`
	Impressions int64   `json:"impressions"`
	Clicks      int64   `json:"clicks"`
	Conversions int64   `json:"conversions"`
	EndDate     *string `json:"end_date,omitempty"`
}

// CampaignResponse represents the campaign response
type CampaignResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	Budget          float64 `json:"budget"`
	Spent           float64 `json:"spent"`
	Impressions     int64   `json:"impressions"`
	Clicks          int64   `json:"clicks"`
	Conversions     int64   `json:"conversions"`
	CTR             float64 `json:"ctr"`
	ConversionRate  float64 `json:"conversion_rate"`
	BudgetUtilization float64 `json:"budget_utilization"`
	StartDate       string  `json:"start_date"`
	EndDate         *string `json:"end_date,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// CreateABTestRequest represents the request to create an A/B test
type CreateABTestRequest struct {
	Name     string  `json:"name" binding:"required"`
	Variant  string  `json:"variant" binding:"required"` // A or B
	StartDate string `json:"start_date" binding:"required"`
}

// ABTestResponse represents the A/B test response
type ABTestResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Variant        string  `json:"variant"`
	Status         string  `json:"status"`
	Visitors       int64   `json:"visitors"`
	Conversions    int64   `json:"conversions"`
	ConversionRate float64 `json:"conversion_rate"`
	Confidence     float64 `json:"confidence"`
	StartDate      string  `json:"start_date"`
	EndDate        *string `json:"end_date,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// DashboardResponse represents the dashboard response
type DashboardResponse struct {
	TotalCampaigns  int64                `json:"total_campaigns"`
	ActiveCampaigns int64                `json:"active_campaigns"`
	TotalSpend      float64              `json:"total_spend"`
	TotalRevenue    float64              `json:"total_revenue"`
	ROI             float64              `json:"roi"`
	ConversionRate  float64              `json:"conversion_rate"`
	TopChannels     []ChannelPerformance `json:"top_channels"`
}

// ChannelResponse represents a channel performance
type ChannelPerformance struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
