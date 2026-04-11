package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// affiliateLinkRepository implements domain.AffiliateLinkRepository
type affiliateLinkRepository struct {
	db *sql.DB
}

// NewAffiliateLinkRepository creates a new affiliate link repository
func NewAffiliateLinkRepository(db *sql.DB) domain.AffiliateLinkRepository {
	return &affiliateLinkRepository{db: db}
}

// Create creates a new affiliate link
func (r *affiliateLinkRepository) Create(ctx context.Context, link *domain.AffiliateLink) error {
	query := `
		INSERT INTO affiliate_links (id, affiliate_id, campaign_id, link_code, name, destination_url,
			click_count, signup_count, conversion_count, earnings, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now

	var campaignID interface{}
	if link.CampaignID != "" {
		campaignID = link.CampaignID
	}

	_, err := r.db.ExecContext(ctx, query,
		link.ID,
		link.AffiliateID,
		campaignID,
		link.LinkCode,
		link.Name,
		link.DestinationURL,
		link.ClickCount,
		link.SignupCount,
		link.ConversionCount,
		link.Earnings,
		link.IsActive,
		link.CreatedAt,
		link.UpdatedAt,
	)

	return err
}

// GetByID retrieves an affiliate link by ID
func (r *affiliateLinkRepository) GetByID(ctx context.Context, id string) (*domain.AffiliateLink, error) {
	query := `
		SELECT id, affiliate_id, campaign_id, link_code, name, destination_url,
			click_count, signup_count, conversion_count, earnings, is_active, created_at, updated_at
		FROM affiliate_links WHERE id = $1
	`

	link := &domain.AffiliateLink{}
	var campaignID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&link.ID,
		&link.AffiliateID,
		&campaignID,
		&link.LinkCode,
		&link.Name,
		&link.DestinationURL,
		&link.ClickCount,
		&link.SignupCount,
		&link.ConversionCount,
		&link.Earnings,
		&link.IsActive,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if campaignID.Valid {
		link.CampaignID = campaignID.String
	}

	return link, nil
}

// GetByCode retrieves an affiliate link by link code
func (r *affiliateLinkRepository) GetByCode(ctx context.Context, code string) (*domain.AffiliateLink, error) {
	query := `
		SELECT id, affiliate_id, campaign_id, link_code, name, destination_url,
			click_count, signup_count, conversion_count, earnings, is_active, created_at, updated_at
		FROM affiliate_links WHERE link_code = $1
	`

	link := &domain.AffiliateLink{}
	var campaignID sql.NullString

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&link.ID,
		&link.AffiliateID,
		&campaignID,
		&link.LinkCode,
		&link.Name,
		&link.DestinationURL,
		&link.ClickCount,
		&link.SignupCount,
		&link.ConversionCount,
		&link.Earnings,
		&link.IsActive,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if campaignID.Valid {
		link.CampaignID = campaignID.String
	}

	return link, nil
}

// GetAllByAffiliateID retrieves all affiliate links for an affiliate with pagination
func (r *affiliateLinkRepository) GetAllByAffiliateID(ctx context.Context, affiliateID string, limit, offset int) ([]*domain.AffiliateLink, error) {
	query := `
		SELECT id, affiliate_id, campaign_id, link_code, name, destination_url,
			click_count, signup_count, conversion_count, earnings, is_active, created_at, updated_at
		FROM affiliate_links WHERE affiliate_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]*domain.AffiliateLink, 0)
	for rows.Next() {
		link := &domain.AffiliateLink{}
		var campaignID sql.NullString

		err := rows.Scan(
			&link.ID,
			&link.AffiliateID,
			&campaignID,
			&link.LinkCode,
			&link.Name,
			&link.DestinationURL,
			&link.ClickCount,
			&link.SignupCount,
			&link.ConversionCount,
			&link.Earnings,
			&link.IsActive,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if campaignID.Valid {
			link.CampaignID = campaignID.String
		}

		links = append(links, link)
	}

	return links, rows.Err()
}

// Update updates an existing affiliate link
func (r *affiliateLinkRepository) Update(ctx context.Context, link *domain.AffiliateLink) error {
	query := `
		UPDATE affiliate_links SET
			name = $2, destination_url = $3, is_active = $4, updated_at = $5
		WHERE id = $1
	`

	link.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		link.ID,
		link.Name,
		link.DestinationURL,
		link.IsActive,
		link.UpdatedAt,
	)

	return err
}

// Delete deletes an affiliate link
func (r *affiliateLinkRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM affiliate_links WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// IncrementClick increments the click count for a link
func (r *affiliateLinkRepository) IncrementClick(ctx context.Context, linkID string) error {
	query := `
		UPDATE affiliate_links SET click_count = click_count + 1, updated_at = $2 WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, linkID, time.Now())
	return err
}

// IncrementSignup increments the signup count for a link
func (r *affiliateLinkRepository) IncrementSignup(ctx context.Context, linkID string) error {
	query := `
		UPDATE affiliate_links SET signup_count = signup_count + 1, updated_at = $2 WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, linkID, time.Now())
	return err
}

// IncrementConversion increments the conversion count and earnings for a link
func (r *affiliateLinkRepository) IncrementConversion(ctx context.Context, linkID string, earnings float64) error {
	query := `
		UPDATE affiliate_links SET 
			conversion_count = conversion_count + 1, 
			earnings = earnings + $2,
			updated_at = $3 
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, linkID, earnings, time.Now())
	return err
}