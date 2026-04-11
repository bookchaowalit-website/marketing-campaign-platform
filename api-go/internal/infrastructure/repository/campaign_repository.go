package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// campaignRepository implements domain.CampaignRepository
type campaignRepository struct {
	db *sql.DB
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(db *sql.DB) domain.CampaignRepository {
	return &campaignRepository{db: db}
}

// Create creates a new campaign
func (r *campaignRepository) Create(ctx context.Context, campaign *domain.Campaign) error {
	query := `
		INSERT INTO campaigns (id, organization_id, name, type, status, budget, spent, impressions, clicks, conversions, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	now := time.Now()
	campaign.CreatedAt = now
	campaign.UpdatedAt = now

	var endDate interface{}
	if campaign.EndDate != nil {
		endDate = *campaign.EndDate
	}

	_, err := r.db.ExecContext(ctx, query,
		campaign.ID,
		campaign.OrganizationID,
		campaign.Name,
		campaign.Type,
		campaign.Status,
		campaign.Budget,
		campaign.Spent,
		campaign.Impressions,
		campaign.Clicks,
		campaign.Conversions,
		campaign.StartDate,
		endDate,
		campaign.CreatedAt,
		campaign.UpdatedAt,
	)

	return err
}

// GetByID retrieves a campaign by ID
func (r *campaignRepository) GetByID(ctx context.Context, id string) (*domain.Campaign, error) {
	query := `
		SELECT id, organization_id, name, type, status, budget, spent, impressions, clicks, conversions, start_date, end_date, created_at, updated_at
		FROM campaigns
		WHERE id = $1
	`

	campaign := &domain.Campaign{}
	var endDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&campaign.ID,
		&campaign.OrganizationID,
		&campaign.Name,
		&campaign.Type,
		&campaign.Status,
		&campaign.Budget,
		&campaign.Spent,
		&campaign.Impressions,
		&campaign.Clicks,
		&campaign.Conversions,
		&campaign.StartDate,
		&endDate,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if endDate.Valid {
		campaign.EndDate = &endDate.Time
	}

	return campaign, nil
}

// GetAllByOrganizationID retrieves all campaigns for an organization with pagination
func (r *campaignRepository) GetAllByOrganizationID(ctx context.Context, organizationID string, limit, offset int) ([]*domain.Campaign, error) {
	query := `
		SELECT id, organization_id, name, type, status, budget, spent, impressions, clicks, conversions, start_date, end_date, created_at, updated_at
		FROM campaigns
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := make([]*domain.Campaign, 0)
	for rows.Next() {
		campaign := &domain.Campaign{}
		var endDate sql.NullTime

		err := rows.Scan(
			&campaign.ID,
			&campaign.OrganizationID,
			&campaign.Name,
			&campaign.Type,
			&campaign.Status,
			&campaign.Budget,
			&campaign.Spent,
			&campaign.Impressions,
			&campaign.Clicks,
			&campaign.Conversions,
			&campaign.StartDate,
			&endDate,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			campaign.EndDate = &endDate.Time
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, rows.Err()
}

// Update updates an existing campaign
func (r *campaignRepository) Update(ctx context.Context, campaign *domain.Campaign) error {
	query := `
		UPDATE campaigns
		SET name = $2, type = $3, status = $4, budget = $5, spent = $6,
		    impressions = $7, clicks = $8, conversions = $9, start_date = $10,
		    end_date = $11, updated_at = $12
		WHERE id = $1
	`

	campaign.UpdatedAt = time.Now()

	var endDate interface{}
	if campaign.EndDate != nil {
		endDate = *campaign.EndDate
	}

	_, err := r.db.ExecContext(ctx, query,
		campaign.ID,
		campaign.Name,
		campaign.Type,
		campaign.Status,
		campaign.Budget,
		campaign.Spent,
		campaign.Impressions,
		campaign.Clicks,
		campaign.Conversions,
		campaign.StartDate,
		endDate,
		campaign.UpdatedAt,
	)

	return err
}

// Delete deletes a campaign
func (r *campaignRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM campaigns WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByStatus retrieves campaigns by status for an organization
func (r *campaignRepository) GetByStatus(ctx context.Context, organizationID, status string, limit, offset int) ([]*domain.Campaign, error) {
	query := `
		SELECT id, organization_id, name, type, status, budget, spent, impressions, clicks, conversions, start_date, end_date, created_at, updated_at
		FROM campaigns
		WHERE organization_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := make([]*domain.Campaign, 0)
	for rows.Next() {
		campaign := &domain.Campaign{}
		var endDate sql.NullTime

		err := rows.Scan(
			&campaign.ID,
			&campaign.OrganizationID,
			&campaign.Name,
			&campaign.Type,
			&campaign.Status,
			&campaign.Budget,
			&campaign.Spent,
			&campaign.Impressions,
			&campaign.Clicks,
			&campaign.Conversions,
			&campaign.StartDate,
			&endDate,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			campaign.EndDate = &endDate.Time
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, rows.Err()
}

// GetByType retrieves campaigns by type for an organization
func (r *campaignRepository) GetByType(ctx context.Context, organizationID, campaignType string, limit, offset int) ([]*domain.Campaign, error) {
	query := `
		SELECT id, organization_id, name, type, status, budget, spent, impressions, clicks, conversions, start_date, end_date, created_at, updated_at
		FROM campaigns
		WHERE organization_id = $1 AND type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, campaignType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := make([]*domain.Campaign, 0)
	for rows.Next() {
		campaign := &domain.Campaign{}
		var endDate sql.NullTime

		err := rows.Scan(
			&campaign.ID,
			&campaign.OrganizationID,
			&campaign.Name,
			&campaign.Type,
			&campaign.Status,
			&campaign.Budget,
			&campaign.Spent,
			&campaign.Impressions,
			&campaign.Clicks,
			&campaign.Conversions,
			&campaign.StartDate,
			&endDate,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			campaign.EndDate = &endDate.Time
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, rows.Err()
}

// Helper function to generate UUID
func generateUUID() string {
	return uuid.New().String()
}
