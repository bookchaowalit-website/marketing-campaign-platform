package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// referralRepository implements domain.ReferralRepository
type referralRepository struct {
	db *sql.DB
}

// NewReferralRepository creates a new referral repository
func NewReferralRepository(db *sql.DB) domain.ReferralRepository {
	return &referralRepository{db: db}
}

// Create creates a new referral
func (r *referralRepository) Create(ctx context.Context, referral *domain.Referral) error {
	query := `
		INSERT INTO referrals (id, affiliate_id, affiliate_link_id, referred_user_id, referral_code,
			status, click_timestamp, signup_timestamp, conversion_timestamp, order_id, order_amount,
			commission_amount, commission_status, ip_address, user_agent, referrer_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	now := time.Now()
	referral.CreatedAt = now
	referral.UpdatedAt = now
	referral.ClickTimestamp = now

	var linkID, referredUserID, signupTime, convTime, orderID, ip, userAgent, referrerURL interface{}
	if referral.AffiliateLinkID != "" {
		linkID = referral.AffiliateLinkID
	}
	if referral.ReferredUserID != "" {
		referredUserID = referral.ReferredUserID
	}
	if referral.SignupTimestamp != nil {
		signupTime = *referral.SignupTimestamp
	}
	if referral.ConversionTimestamp != nil {
		convTime = *referral.ConversionTimestamp
	}
	if referral.OrderID != "" {
		orderID = referral.OrderID
	}
	if referral.IPAddress != "" {
		ip = referral.IPAddress
	}
	if referral.UserAgent != "" {
		userAgent = referral.UserAgent
	}
	if referral.ReferrerURL != "" {
		referrerURL = referral.ReferrerURL
	}

	_, err := r.db.ExecContext(ctx, query,
		referral.ID,
		referral.AffiliateID,
		linkID,
		referredUserID,
		referral.ReferralCode,
		referral.Status,
		referral.ClickTimestamp,
		signupTime,
		convTime,
		orderID,
		referral.OrderAmount,
		referral.CommissionAmount,
		referral.CommissionStatus,
		ip,
		userAgent,
		referrerURL,
		referral.CreatedAt,
		referral.UpdatedAt,
	)

	return err
}

// GetByID retrieves a referral by ID
func (r *referralRepository) GetByID(ctx context.Context, id string) (*domain.Referral, error) {
	query := `
		SELECT id, affiliate_id, affiliate_link_id, referred_user_id, referral_code, status,
			click_timestamp, signup_timestamp, conversion_timestamp, order_id, order_amount,
			commission_amount, commission_status, ip_address, user_agent, referrer_url, created_at, updated_at
		FROM referrals WHERE id = $1
	`

	referral := &domain.Referral{}
	var linkID, referredUserID, orderID, ip, userAgent, referrerURL sql.NullString
	var signupTime, convTime sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&referral.ID,
		&referral.AffiliateID,
		&linkID,
		&referredUserID,
		&referral.ReferralCode,
		&referral.Status,
		&referral.ClickTimestamp,
		&signupTime,
		&convTime,
		&orderID,
		&referral.OrderAmount,
		&referral.CommissionAmount,
		&referral.CommissionStatus,
		&ip,
		&userAgent,
		&referrerURL,
		&referral.CreatedAt,
		&referral.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if linkID.Valid {
		referral.AffiliateLinkID = linkID.String
	}
	if referredUserID.Valid {
		referral.ReferredUserID = referredUserID.String
	}
	if signupTime.Valid {
		referral.SignupTimestamp = &signupTime.Time
	}
	if convTime.Valid {
		referral.ConversionTimestamp = &convTime.Time
	}
	if orderID.Valid {
		referral.OrderID = orderID.String
	}
	if ip.Valid {
		referral.IPAddress = ip.String
	}
	if userAgent.Valid {
		referral.UserAgent = userAgent.String
	}
	if referrerURL.Valid {
		referral.ReferrerURL = referrerURL.String
	}

	return referral, nil
}

// GetByReferralCode retrieves a referral by referral code
func (r *referralRepository) GetByReferralCode(ctx context.Context, code string) (*domain.Referral, error) {
	query := `
		SELECT id, affiliate_id, affiliate_link_id, referred_user_id, referral_code, status,
			click_timestamp, signup_timestamp, conversion_timestamp, order_id, order_amount,
			commission_amount, commission_status, ip_address, user_agent, referrer_url, created_at, updated_at
		FROM referrals WHERE referral_code = $1 ORDER BY created_at DESC LIMIT 1
	`

	referral := &domain.Referral{}
	var linkID, referredUserID, orderID, ip, userAgent, referrerURL sql.NullString
	var signupTime, convTime sql.NullTime

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&referral.ID,
		&referral.AffiliateID,
		&linkID,
		&referredUserID,
		&referral.ReferralCode,
		&referral.Status,
		&referral.ClickTimestamp,
		&signupTime,
		&convTime,
		&orderID,
		&referral.OrderAmount,
		&referral.CommissionAmount,
		&referral.CommissionStatus,
		&ip,
		&userAgent,
		&referrerURL,
		&referral.CreatedAt,
		&referral.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if linkID.Valid {
		referral.AffiliateLinkID = linkID.String
	}
	if referredUserID.Valid {
		referral.ReferredUserID = referredUserID.String
	}
	if signupTime.Valid {
		referral.SignupTimestamp = &signupTime.Time
	}
	if convTime.Valid {
		referral.ConversionTimestamp = &convTime.Time
	}
	if orderID.Valid {
		referral.OrderID = orderID.String
	}
	if ip.Valid {
		referral.IPAddress = ip.String
	}
	if userAgent.Valid {
		referral.UserAgent = userAgent.String
	}
	if referrerURL.Valid {
		referral.ReferrerURL = referrerURL.String
	}

	return referral, nil
}

// GetAllByAffiliateID retrieves all referrals for an affiliate with pagination
func (r *referralRepository) GetAllByAffiliateID(ctx context.Context, affiliateID string, limit, offset int) ([]*domain.Referral, error) {
	query := `
		SELECT id, affiliate_id, affiliate_link_id, referred_user_id, referral_code, status,
			click_timestamp, signup_timestamp, conversion_timestamp, order_id, order_amount,
			commission_amount, commission_status, ip_address, user_agent, referrer_url, created_at, updated_at
		FROM referrals WHERE affiliate_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	referrals := make([]*domain.Referral, 0)
	for rows.Next() {
		referral := &domain.Referral{}
		var linkID, referredUserID, orderID, ip, userAgent, referrerURL sql.NullString
		var signupTime, convTime sql.NullTime

		err := rows.Scan(
			&referral.ID,
			&referral.AffiliateID,
			&linkID,
			&referredUserID,
			&referral.ReferralCode,
			&referral.Status,
			&referral.ClickTimestamp,
			&signupTime,
			&convTime,
			&orderID,
			&referral.OrderAmount,
			&referral.CommissionAmount,
			&referral.CommissionStatus,
			&ip,
			&userAgent,
			&referrerURL,
			&referral.CreatedAt,
			&referral.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if linkID.Valid {
			referral.AffiliateLinkID = linkID.String
		}
		if referredUserID.Valid {
			referral.ReferredUserID = referredUserID.String
		}
		if signupTime.Valid {
			referral.SignupTimestamp = &signupTime.Time
		}
		if convTime.Valid {
			referral.ConversionTimestamp = &convTime.Time
		}
		if orderID.Valid {
			referral.OrderID = orderID.String
		}
		if ip.Valid {
			referral.IPAddress = ip.String
		}
		if userAgent.Valid {
			referral.UserAgent = userAgent.String
		}
		if referrerURL.Valid {
			referral.ReferrerURL = referrerURL.String
		}

		referrals = append(referrals, referral)
	}

	return referrals, rows.Err()
}

// GetByStatus retrieves referrals by status for an affiliate
func (r *referralRepository) GetByStatus(ctx context.Context, affiliateID, status string, limit, offset int) ([]*domain.Referral, error) {
	query := `
		SELECT id, affiliate_id, affiliate_link_id, referred_user_id, referral_code, status,
			click_timestamp, signup_timestamp, conversion_timestamp, order_id, order_amount,
			commission_amount, commission_status, ip_address, user_agent, referrer_url, created_at, updated_at
		FROM referrals WHERE affiliate_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	referrals := make([]*domain.Referral, 0)
	for rows.Next() {
		referral := &domain.Referral{}
		var linkID, referredUserID, orderID, ip, userAgent, referrerURL sql.NullString
		var signupTime, convTime sql.NullTime

		err := rows.Scan(
			&referral.ID,
			&referral.AffiliateID,
			&linkID,
			&referredUserID,
			&referral.ReferralCode,
			&referral.Status,
			&referral.ClickTimestamp,
			&signupTime,
			&convTime,
			&orderID,
			&referral.OrderAmount,
			&referral.CommissionAmount,
			&referral.CommissionStatus,
			&ip,
			&userAgent,
			&referrerURL,
			&referral.CreatedAt,
			&referral.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if linkID.Valid {
			referral.AffiliateLinkID = linkID.String
		}
		if referredUserID.Valid {
			referral.ReferredUserID = referredUserID.String
		}
		if signupTime.Valid {
			referral.SignupTimestamp = &signupTime.Time
		}
		if convTime.Valid {
			referral.ConversionTimestamp = &convTime.Time
		}
		if orderID.Valid {
			referral.OrderID = orderID.String
		}
		if ip.Valid {
			referral.IPAddress = ip.String
		}
		if userAgent.Valid {
			referral.UserAgent = userAgent.String
		}
		if referrerURL.Valid {
			referral.ReferrerURL = referrerURL.String
		}

		referrals = append(referrals, referral)
	}

	return referrals, rows.Err()
}

// Update updates an existing referral
func (r *referralRepository) Update(ctx context.Context, referral *domain.Referral) error {
	query := `
		UPDATE referrals SET
			status = $2, signup_timestamp = $3, conversion_timestamp = $4, order_id = $5,
			order_amount = $6, commission_amount = $7, commission_status = $8, referred_user_id = $9,
			updated_at = $10
		WHERE id = $1
	`

	referral.UpdatedAt = time.Now()

	var signupTime, convTime, orderID, referredUserID interface{}
	if referral.SignupTimestamp != nil {
		signupTime = *referral.SignupTimestamp
	}
	if referral.ConversionTimestamp != nil {
		convTime = *referral.ConversionTimestamp
	}
	if referral.OrderID != "" {
		orderID = referral.OrderID
	}
	if referral.ReferredUserID != "" {
		referredUserID = referral.ReferredUserID
	}

	_, err := r.db.ExecContext(ctx, query,
		referral.ID,
		referral.Status,
		signupTime,
		convTime,
		orderID,
		referral.OrderAmount,
		referral.CommissionAmount,
		referral.CommissionStatus,
		referredUserID,
		referral.UpdatedAt,
	)

	return err
}

// UpdateStatus updates the status of a referral
func (r *referralRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE referrals SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status, time.Now())
	return err
}

// GetRecentByAffiliateID retrieves recent referrals for an affiliate
func (r *referralRepository) GetRecentByAffiliateID(ctx context.Context, affiliateID string, limit int) ([]*domain.Referral, error) {
	return r.GetAllByAffiliateID(ctx, affiliateID, limit, 0)
}