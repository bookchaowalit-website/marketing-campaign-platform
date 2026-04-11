package repository

import (
	"context"
	"database/sql"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// dashboardRepository implements domain.DashboardRepository
type dashboardRepository struct {
	db *sql.DB
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(db *sql.DB) domain.DashboardRepository {
	return &dashboardRepository{
		db: db,
	}
}

// Default organization ID for demo purposes
const defaultOrganizationID = "00000000-0000-0000-0000-000000000001"

// GetOverview retrieves the dashboard overview for an organization
func (r *dashboardRepository) GetOverview(ctx context.Context, organizationID string) (*domain.Dashboard, error) {
	if organizationID == "" {
		organizationID = defaultOrganizationID
	}

	// Get campaign statistics for organization
	var totalCampaigns, activeCampaigns int64
	var totalSpend, totalRevenue float64

	// Count total campaigns
	err := r.db.QueryRowContext(ctx, 
		"SELECT COUNT(*) FROM campaigns WHERE organization_id = $1", 
		organizationID,
	).Scan(&totalCampaigns)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Count active campaigns
	err = r.db.QueryRowContext(ctx, 
		"SELECT COUNT(*) FROM campaigns WHERE organization_id = $1 AND status = 'active'", 
		organizationID,
	).Scan(&activeCampaigns)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get total spend
	err = r.db.QueryRowContext(ctx, 
		"SELECT COALESCE(SUM(spent), 0) FROM campaigns WHERE organization_id = $1", 
		organizationID,
	).Scan(&totalSpend)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Calculate revenue (mock: 2.5x spend for demo)
	totalRevenue = totalSpend * 2.5

	// Calculate ROI
	var roi float64
	if totalSpend > 0 {
		roi = ((totalRevenue - totalSpend) / totalSpend) * 100
	}

	// Get conversion rate
	var totalClicks, totalConversions int64
	err = r.db.QueryRowContext(ctx, 
		"SELECT COALESCE(SUM(clicks), 0), COALESCE(SUM(conversions), 0) FROM campaigns WHERE organization_id = $1", 
		organizationID,
	).Scan(&totalClicks, &totalConversions)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var conversionRate float64
	if totalClicks > 0 {
		conversionRate = float64(totalConversions) / float64(totalClicks) * 100
	}

	// Top channels (mock data for now)
	topChannels := []domain.ChannelPerformance{
		{Name: "Google Ads", Value: 45},
		{Name: "Facebook", Value: 30},
		{Name: "Organic Search", Value: 15},
		{Name: "Email", Value: 10},
	}

	return &domain.Dashboard{
		TotalCampaigns:  totalCampaigns,
		ActiveCampaigns: activeCampaigns,
		TotalSpend:      totalSpend,
		TotalRevenue:    totalRevenue,
		ROI:             roi,
		ConversionRate:  conversionRate,
		TopChannels:     topChannels,
	}, nil
}
