package repository

import (
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

type unifiedCampaignRepository struct {
	db *sql.DB
}

// NewUnifiedCampaignRepository creates a new repository
func NewUnifiedCampaignRepository(db *sql.DB) domain.UnifiedCampaignRepository {
	return &unifiedCampaignRepository{db: db}
}

func (r *unifiedCampaignRepository) Create(c *domain.UnifiedCampaign) error {
	query := `
		INSERT INTO unified_campaignes (id, organization_id, group_id, name, description,
			channel_category, channel, budget, spent, target_audience, status,
			start_date, end_date, timezone, owner_id, created_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
	`
	now := time.Now()
	c.ID = uuid.New().String()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.Exec(query, c.ID, c.OrganizationID, c.GroupID, c.Name, c.Description,
		c.ChannelCategory, c.Channel, c.Budget, c.Spent, c.TargetAudience, c.Status,
		c.StartDate, c.EndDate, c.Timezone, c.OwnerID, c.CreatedBy, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *unifiedCampaignRepository) GetByID(id string) (*domain.UnifiedCampaign, error) {
	query := `
		SELECT id, organization_id, group_id, name, description, channel_category, channel,
			budget, spent, cpa, roas, target_audience, status, start_date, end_date,
			timezone, utm_source, utm_medium, utm_campaign, tracking_url, creative_brief,
			approval_status, approved_by, approved_at, rejection_reason,
			impressions, clicks, conversions, revenue, leads,
			owner_id, created_by, created_at, updated_at
		FROM unified_campaignes WHERE id = $1
	`
	c := &domain.UnifiedCampaign{}
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.OrganizationID, &c.GroupID, &c.Name, &c.Description,
		&c.ChannelCategory, &c.Channel, &c.Budget, &c.Spent, &c.CPA, &c.ROAS,
		&c.TargetAudience, &c.Status, &c.StartDate, &c.EndDate, &c.Timezone,
		&c.UTMSource, &c.UTMMedium, &c.UTMCampaign, &c.TrackingURL, &c.CreativeBrief,
		&c.ApprovalStatus, &c.ApprovedBy, &c.ApprovedAt, &c.RejectionReason,
		&c.Impressions, &c.Clicks, &c.Conversions, &c.Revenue, &c.Leads,
		&c.OwnerID, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *unifiedCampaignRepository) List(orgID, channel, status string, limit, offset int) ([]*domain.UnifiedCampaign, int, error) {
	countQuery := `SELECT COUNT(*) FROM unified_campaignes WHERE organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if channel != "" {
		countQuery += ` AND channel = $` + itoa(argIdx)
		args = append(args, channel)
		argIdx++
	}
	if status != "" {
		countQuery += ` AND status = $` + itoa(argIdx)
		args = append(args, status)
		argIdx++
	}

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, organization_id, group_id, name, description, channel_category, channel,
			budget, spent, cpa, roas, target_audience, status, start_date, end_date,
			timezone, utm_source, owner_id, created_by, created_at, updated_at
		FROM unified_campaignes WHERE organization_id = $1`
	if channel != "" {
		query += ` AND channel = $` + itoa(argIdx)
		args = append(args, channel)
		argIdx++
	}
	if status != "" {
		query += ` AND status = $` + itoa(argIdx)
		args = append(args, status)
		argIdx++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var campaigns []*domain.UnifiedCampaign
	for rows.Next() {
		c := &domain.UnifiedCampaign{}
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.GroupID, &c.Name, &c.Description,
			&c.ChannelCategory, &c.Channel, &c.Budget, &c.Spent, &c.CPA, &c.ROAS,
			&c.TargetAudience, &c.Status, &c.StartDate, &c.EndDate, &c.Timezone,
			&c.UTMSource, &c.OwnerID, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, total, nil
}

func (r *unifiedCampaignRepository) Update(c *domain.UnifiedCampaign) error {
	query := `
		UPDATE unified_campaignes SET name=$1, description=$2, budget=$3, spent=$4,
			target_audience=$5, status=$6, start_date=$7, end_date=$8, tracking_url=$9,
			creative_brief=$10, updated_at=$11 WHERE id=$12
	`
	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, c.Name, c.Description, c.Budget, c.Spent,
		c.TargetAudience, c.Status, c.StartDate, c.EndDate, c.TrackingURL,
		c.CreativeBrief, c.UpdatedAt, c.ID)
	return err
}

func (r *unifiedCampaignRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM unified_campaignes WHERE id = $1`, id)
	return err
}

func (r *unifiedCampaignRepository) SubmitForApproval(id string) error {
	_, err := r.db.Exec(`UPDATE unified_campaignes SET approval_status='pending', updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *unifiedCampaignRepository) Approve(id, approverID string) error {
	now := time.Now()
	_, err := r.db.Exec(`UPDATE unified_campaignes SET approval_status='approved', approved_by=$1, approved_at=$2, updated_at=NOW() WHERE id=$3`, approverID, now, id)
	return err
}

func (r *unifiedCampaignRepository) Reject(id, reason string) error {
	_, err := r.db.Exec(`UPDATE unified_campaignes SET approval_status='rejected', rejection_reason=$1, updated_at=NOW() WHERE id=$2`, reason, id)
	return err
}

func itoa(i int) string {
	return string(rune('0' + i))
}
