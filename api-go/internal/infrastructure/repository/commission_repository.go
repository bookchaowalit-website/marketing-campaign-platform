package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// commissionRepository implements domain.CommissionRepository
type commissionRepository struct {
	db *sql.DB
}

// NewCommissionRepository creates a new commission repository
func NewCommissionRepository(db *sql.DB) domain.CommissionRepository {
	return &commissionRepository{db: db}
}

// Create creates a new commission
func (r *commissionRepository) Create(ctx context.Context, commission *domain.Commission) error {
	query := `
		INSERT INTO commissions (id, affiliate_id, referral_id, amount, status, type, description,
			reference_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()
	commission.CreatedAt = now
	commission.UpdatedAt = now

	var referralID, description, referenceID interface{}
	if commission.ReferralID != "" {
		referralID = commission.ReferralID
	}
	if commission.Description != "" {
		description = commission.Description
	}
	if commission.ReferenceID != "" {
		referenceID = commission.ReferenceID
	}

	_, err := r.db.ExecContext(ctx, query,
		commission.ID,
		commission.AffiliateID,
		referralID,
		commission.Amount,
		commission.Status,
		commission.Type,
		description,
		referenceID,
		commission.CreatedAt,
		commission.UpdatedAt,
	)

	return err
}

// GetByID retrieves a commission by ID
func (r *commissionRepository) GetByID(ctx context.Context, id string) (*domain.Commission, error) {
	query := `
		SELECT id, affiliate_id, referral_id, amount, status, type, description, reference_id,
			created_at, updated_at, approved_at, approved_by, paid_at, payout_id
		FROM commissions WHERE id = $1
	`

	commission := &domain.Commission{}
	var referralID, description, referenceID, approvedBy, payoutID sql.NullString
	var approvedAt, paidAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&commission.ID,
		&commission.AffiliateID,
		&referralID,
		&commission.Amount,
		&commission.Status,
		&commission.Type,
		&description,
		&referenceID,
		&commission.CreatedAt,
		&commission.UpdatedAt,
		&approvedAt,
		&approvedBy,
		&paidAt,
		&payoutID,
	)

	if err != nil {
		return nil, err
	}

	if referralID.Valid {
		commission.ReferralID = referralID.String
	}
	if description.Valid {
		commission.Description = description.String
	}
	if referenceID.Valid {
		commission.ReferenceID = referenceID.String
	}
	if approvedAt.Valid {
		commission.ApprovedAt = &approvedAt.Time
	}
	if approvedBy.Valid {
		commission.ApprovedBy = approvedBy.String
	}
	if paidAt.Valid {
		commission.PaidAt = &paidAt.Time
	}
	if payoutID.Valid {
		commission.PayoutID = payoutID.String
	}

	return commission, nil
}

// GetAllByAffiliateID retrieves all commissions for an affiliate with pagination
func (r *commissionRepository) GetAllByAffiliateID(ctx context.Context, affiliateID string, limit, offset int) ([]*domain.Commission, error) {
	query := `
		SELECT id, affiliate_id, referral_id, amount, status, type, description, reference_id,
			created_at, updated_at, approved_at, approved_by, paid_at, payout_id
		FROM commissions WHERE affiliate_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commissions := make([]*domain.Commission, 0)
	for rows.Next() {
		commission := &domain.Commission{}
		var referralID, description, referenceID, approvedBy, payoutID sql.NullString
		var approvedAt, paidAt sql.NullTime

		err := rows.Scan(
			&commission.ID,
			&commission.AffiliateID,
			&referralID,
			&commission.Amount,
			&commission.Status,
			&commission.Type,
			&description,
			&referenceID,
			&commission.CreatedAt,
			&commission.UpdatedAt,
			&approvedAt,
			&approvedBy,
			&paidAt,
			&payoutID,
		)
		if err != nil {
			return nil, err
		}

		if referralID.Valid {
			commission.ReferralID = referralID.String
		}
		if description.Valid {
			commission.Description = description.String
		}
		if referenceID.Valid {
			commission.ReferenceID = referenceID.String
		}
		if approvedAt.Valid {
			commission.ApprovedAt = &approvedAt.Time
		}
		if approvedBy.Valid {
			commission.ApprovedBy = approvedBy.String
		}
		if paidAt.Valid {
			commission.PaidAt = &paidAt.Time
		}
		if payoutID.Valid {
			commission.PayoutID = payoutID.String
		}

		commissions = append(commissions, commission)
	}

	return commissions, rows.Err()
}

// GetByStatus retrieves commissions by status for an affiliate
func (r *commissionRepository) GetByStatus(ctx context.Context, affiliateID, status string, limit, offset int) ([]*domain.Commission, error) {
	query := `
		SELECT id, affiliate_id, referral_id, amount, status, type, description, reference_id,
			created_at, updated_at, approved_at, approved_by, paid_at, payout_id
		FROM commissions WHERE affiliate_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commissions := make([]*domain.Commission, 0)
	for rows.Next() {
		commission := &domain.Commission{}
		var referralID, description, referenceID, approvedBy, payoutID sql.NullString
		var approvedAt, paidAt sql.NullTime

		err := rows.Scan(
			&commission.ID,
			&commission.AffiliateID,
			&referralID,
			&commission.Amount,
			&commission.Status,
			&commission.Type,
			&description,
			&referenceID,
			&commission.CreatedAt,
			&commission.UpdatedAt,
			&approvedAt,
			&approvedBy,
			&paidAt,
			&payoutID,
		)
		if err != nil {
			return nil, err
		}

		if referralID.Valid {
			commission.ReferralID = referralID.String
		}
		if description.Valid {
			commission.Description = description.String
		}
		if referenceID.Valid {
			commission.ReferenceID = referenceID.String
		}
		if approvedAt.Valid {
			commission.ApprovedAt = &approvedAt.Time
		}
		if approvedBy.Valid {
			commission.ApprovedBy = approvedBy.String
		}
		if paidAt.Valid {
			commission.PaidAt = &paidAt.Time
		}
		if payoutID.Valid {
			commission.PayoutID = payoutID.String
		}

		commissions = append(commissions, commission)
	}

	return commissions, rows.Err()
}

// GetPendingByAffiliateID retrieves all pending commissions for an affiliate
func (r *commissionRepository) GetPendingByAffiliateID(ctx context.Context, affiliateID string) ([]*domain.Commission, error) {
	return r.GetByStatus(ctx, affiliateID, "pending", 1000, 0)
}

// Update updates an existing commission
func (r *commissionRepository) Update(ctx context.Context, commission *domain.Commission) error {
	query := `
		UPDATE commissions SET
			amount = $2, status = $3, type = $4, description = $5, reference_id = $6, updated_at = $7
		WHERE id = $1
	`

	commission.UpdatedAt = time.Now()

	var description, referenceID interface{}
	if commission.Description != "" {
		description = commission.Description
	}
	if commission.ReferenceID != "" {
		referenceID = commission.ReferenceID
	}

	_, err := r.db.ExecContext(ctx, query,
		commission.ID,
		commission.Amount,
		commission.Status,
		commission.Type,
		description,
		referenceID,
		commission.UpdatedAt,
	)

	return err
}

// UpdateStatus updates the status of a commission
func (r *commissionRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE commissions SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status, time.Now())
	return err
}

// Approve approves a commission
func (r *commissionRepository) Approve(ctx context.Context, id string, approvedBy string) error {
	query := `
		UPDATE commissions SET 
			status = 'approved', 
			approved_at = $3, 
			approved_by = $2, 
			updated_at = $3 
		WHERE id = $1
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, id, approvedBy, now)
	return err
}

// GetTotalPendingByAffiliateID gets the total pending commission amount for an affiliate
func (r *commissionRepository) GetTotalPendingByAffiliateID(ctx context.Context, affiliateID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) FROM commissions 
		WHERE affiliate_id = $1 AND status = 'pending'
	`

	var total float64
	err := r.db.QueryRowContext(ctx, query, affiliateID).Scan(&total)
	return total, err
}

// GetRecentByAffiliateID retrieves recent commissions for an affiliate
func (r *commissionRepository) GetRecentByAffiliateID(ctx context.Context, affiliateID string, limit int) ([]*domain.Commission, error) {
	return r.GetAllByAffiliateID(ctx, affiliateID, limit, 0)
}