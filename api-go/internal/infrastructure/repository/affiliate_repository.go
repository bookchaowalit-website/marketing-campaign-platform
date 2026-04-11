package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// affiliateRepository implements domain.AffiliateRepository
type affiliateRepository struct {
	db *sql.DB
}

// NewAffiliateRepository creates a new affiliate repository
func NewAffiliateRepository(db *sql.DB) domain.AffiliateRepository {
	return &affiliateRepository{db: db}
}

// Create creates a new affiliate
func (r *affiliateRepository) Create(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `
		INSERT INTO affiliates (id, user_id, organization_id, affiliate_code, status, commission_rate, 
			commission_type, fixed_commission_amount, total_clicks, total_signups, total_conversions, 
			total_earnings, total_paid, balance, minimum_payout, paypal_email, bank_name, bank_account, 
			bank_routing, tax_id, website, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`

	now := time.Now()
	affiliate.CreatedAt = now
	affiliate.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		affiliate.ID,
		affiliate.UserID,
		affiliate.OrganizationID,
		affiliate.AffiliateCode,
		affiliate.Status,
		affiliate.CommissionRate,
		affiliate.CommissionType,
		affiliate.FixedCommissionAmount,
		affiliate.TotalClicks,
		affiliate.TotalSignups,
		affiliate.TotalConversions,
		affiliate.TotalEarnings,
		affiliate.TotalPaid,
		affiliate.Balance,
		affiliate.MinimumPayout,
		affiliate.PaypalEmail,
		affiliate.BankName,
		affiliate.BankAccount,
		affiliate.BankRouting,
		affiliate.TaxID,
		affiliate.Website,
		affiliate.Description,
		affiliate.CreatedAt,
		affiliate.UpdatedAt,
	)

	return err
}

// GetByID retrieves an affiliate by ID
func (r *affiliateRepository) GetByID(ctx context.Context, id string) (*domain.Affiliate, error) {
	query := `
		SELECT id, user_id, organization_id, affiliate_code, status, commission_rate, commission_type,
			fixed_commission_amount, total_clicks, total_signups, total_conversions, total_earnings,
			total_paid, balance, minimum_payout, paypal_email, bank_name, bank_account, bank_routing,
			tax_id, website, description, created_at, updated_at, approved_at, approved_by
		FROM affiliates WHERE id = $1
	`

	affiliate := &domain.Affiliate{}
	var approvedAt sql.NullTime
	var approvedBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&affiliate.ID,
		&affiliate.UserID,
		&affiliate.OrganizationID,
		&affiliate.AffiliateCode,
		&affiliate.Status,
		&affiliate.CommissionRate,
		&affiliate.CommissionType,
		&affiliate.FixedCommissionAmount,
		&affiliate.TotalClicks,
		&affiliate.TotalSignups,
		&affiliate.TotalConversions,
		&affiliate.TotalEarnings,
		&affiliate.TotalPaid,
		&affiliate.Balance,
		&affiliate.MinimumPayout,
		&affiliate.PaypalEmail,
		&affiliate.BankName,
		&affiliate.BankAccount,
		&affiliate.BankRouting,
		&affiliate.TaxID,
		&affiliate.Website,
		&affiliate.Description,
		&affiliate.CreatedAt,
		&affiliate.UpdatedAt,
		&approvedAt,
		&approvedBy,
	)

	if err != nil {
		return nil, err
	}

	if approvedAt.Valid {
		affiliate.ApprovedAt = &approvedAt.Time
	}
	if approvedBy.Valid {
		affiliate.ApprovedBy = approvedBy.String
	}

	return affiliate, nil
}

// GetByUserID retrieves an affiliate by user ID
func (r *affiliateRepository) GetByUserID(ctx context.Context, userID string) (*domain.Affiliate, error) {
	query := `
		SELECT id, user_id, organization_id, affiliate_code, status, commission_rate, commission_type,
			fixed_commission_amount, total_clicks, total_signups, total_conversions, total_earnings,
			total_paid, balance, minimum_payout, paypal_email, bank_name, bank_account, bank_routing,
			tax_id, website, description, created_at, updated_at, approved_at, approved_by
		FROM affiliates WHERE user_id = $1
	`

	affiliate := &domain.Affiliate{}
	var approvedAt sql.NullTime
	var approvedBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&affiliate.ID,
		&affiliate.UserID,
		&affiliate.OrganizationID,
		&affiliate.AffiliateCode,
		&affiliate.Status,
		&affiliate.CommissionRate,
		&affiliate.CommissionType,
		&affiliate.FixedCommissionAmount,
		&affiliate.TotalClicks,
		&affiliate.TotalSignups,
		&affiliate.TotalConversions,
		&affiliate.TotalEarnings,
		&affiliate.TotalPaid,
		&affiliate.Balance,
		&affiliate.MinimumPayout,
		&affiliate.PaypalEmail,
		&affiliate.BankName,
		&affiliate.BankAccount,
		&affiliate.BankRouting,
		&affiliate.TaxID,
		&affiliate.Website,
		&affiliate.Description,
		&affiliate.CreatedAt,
		&affiliate.UpdatedAt,
		&approvedAt,
		&approvedBy,
	)

	if err != nil {
		return nil, err
	}

	if approvedAt.Valid {
		affiliate.ApprovedAt = &approvedAt.Time
	}
	if approvedBy.Valid {
		affiliate.ApprovedBy = approvedBy.String
	}

	return affiliate, nil
}

// GetByCode retrieves an affiliate by affiliate code
func (r *affiliateRepository) GetByCode(ctx context.Context, code string) (*domain.Affiliate, error) {
	query := `
		SELECT id, user_id, organization_id, affiliate_code, status, commission_rate, commission_type,
			fixed_commission_amount, total_clicks, total_signups, total_conversions, total_earnings,
			total_paid, balance, minimum_payout, paypal_email, bank_name, bank_account, bank_routing,
			tax_id, website, description, created_at, updated_at, approved_at, approved_by
		FROM affiliates WHERE affiliate_code = $1
	`

	affiliate := &domain.Affiliate{}
	var approvedAt sql.NullTime
	var approvedBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&affiliate.ID,
		&affiliate.UserID,
		&affiliate.OrganizationID,
		&affiliate.AffiliateCode,
		&affiliate.Status,
		&affiliate.CommissionRate,
		&affiliate.CommissionType,
		&affiliate.FixedCommissionAmount,
		&affiliate.TotalClicks,
		&affiliate.TotalSignups,
		&affiliate.TotalConversions,
		&affiliate.TotalEarnings,
		&affiliate.TotalPaid,
		&affiliate.Balance,
		&affiliate.MinimumPayout,
		&affiliate.PaypalEmail,
		&affiliate.BankName,
		&affiliate.BankAccount,
		&affiliate.BankRouting,
		&affiliate.TaxID,
		&affiliate.Website,
		&affiliate.Description,
		&affiliate.CreatedAt,
		&affiliate.UpdatedAt,
		&approvedAt,
		&approvedBy,
	)

	if err != nil {
		return nil, err
	}

	if approvedAt.Valid {
		affiliate.ApprovedAt = &approvedAt.Time
	}
	if approvedBy.Valid {
		affiliate.ApprovedBy = approvedBy.String
	}

	return affiliate, nil
}

// GetAllByOrganizationID retrieves all affiliates for an organization with pagination
func (r *affiliateRepository) GetAllByOrganizationID(ctx context.Context, organizationID string, limit, offset int) ([]*domain.Affiliate, error) {
	query := `
		SELECT id, user_id, organization_id, affiliate_code, status, commission_rate, commission_type,
			fixed_commission_amount, total_clicks, total_signups, total_conversions, total_earnings,
			total_paid, balance, minimum_payout, paypal_email, bank_name, bank_account, bank_routing,
			tax_id, website, description, created_at, updated_at, approved_at, approved_by
		FROM affiliates WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	affiliates := make([]*domain.Affiliate, 0)
	for rows.Next() {
		affiliate := &domain.Affiliate{}
		var approvedAt sql.NullTime
		var approvedBy sql.NullString

		err := rows.Scan(
			&affiliate.ID,
			&affiliate.UserID,
			&affiliate.OrganizationID,
			&affiliate.AffiliateCode,
			&affiliate.Status,
			&affiliate.CommissionRate,
			&affiliate.CommissionType,
			&affiliate.FixedCommissionAmount,
			&affiliate.TotalClicks,
			&affiliate.TotalSignups,
			&affiliate.TotalConversions,
			&affiliate.TotalEarnings,
			&affiliate.TotalPaid,
			&affiliate.Balance,
			&affiliate.MinimumPayout,
			&affiliate.PaypalEmail,
			&affiliate.BankName,
			&affiliate.BankAccount,
			&affiliate.BankRouting,
			&affiliate.TaxID,
			&affiliate.Website,
			&affiliate.Description,
			&affiliate.CreatedAt,
			&affiliate.UpdatedAt,
			&approvedAt,
			&approvedBy,
		)
		if err != nil {
			return nil, err
		}

		if approvedAt.Valid {
			affiliate.ApprovedAt = &approvedAt.Time
		}
		if approvedBy.Valid {
			affiliate.ApprovedBy = approvedBy.String
		}

		affiliates = append(affiliates, affiliate)
	}

	return affiliates, rows.Err()
}

// GetByStatus retrieves affiliates by status for an organization
func (r *affiliateRepository) GetByStatus(ctx context.Context, organizationID, status string, limit, offset int) ([]*domain.Affiliate, error) {
	query := `
		SELECT id, user_id, organization_id, affiliate_code, status, commission_rate, commission_type,
			fixed_commission_amount, total_clicks, total_signups, total_conversions, total_earnings,
			total_paid, balance, minimum_payout, paypal_email, bank_name, bank_account, bank_routing,
			tax_id, website, description, created_at, updated_at, approved_at, approved_by
		FROM affiliates WHERE organization_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	affiliates := make([]*domain.Affiliate, 0)
	for rows.Next() {
		affiliate := &domain.Affiliate{}
		var approvedAt sql.NullTime
		var approvedBy sql.NullString

		err := rows.Scan(
			&affiliate.ID,
			&affiliate.UserID,
			&affiliate.OrganizationID,
			&affiliate.AffiliateCode,
			&affiliate.Status,
			&affiliate.CommissionRate,
			&affiliate.CommissionType,
			&affiliate.FixedCommissionAmount,
			&affiliate.TotalClicks,
			&affiliate.TotalSignups,
			&affiliate.TotalConversions,
			&affiliate.TotalEarnings,
			&affiliate.TotalPaid,
			&affiliate.Balance,
			&affiliate.MinimumPayout,
			&affiliate.PaypalEmail,
			&affiliate.BankName,
			&affiliate.BankAccount,
			&affiliate.BankRouting,
			&affiliate.TaxID,
			&affiliate.Website,
			&affiliate.Description,
			&affiliate.CreatedAt,
			&affiliate.UpdatedAt,
			&approvedAt,
			&approvedBy,
		)
		if err != nil {
			return nil, err
		}

		if approvedAt.Valid {
			affiliate.ApprovedAt = &approvedAt.Time
		}
		if approvedBy.Valid {
			affiliate.ApprovedBy = approvedBy.String
		}

		affiliates = append(affiliates, affiliate)
	}

	return affiliates, rows.Err()
}

// Update updates an existing affiliate
func (r *affiliateRepository) Update(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `
		UPDATE affiliates SET
			status = $2, commission_rate = $3, commission_type = $4, fixed_commission_amount = $5,
			total_clicks = $6, total_signups = $7, total_conversions = $8, total_earnings = $9,
			total_paid = $10, balance = $11, minimum_payout = $12, paypal_email = $13,
			bank_name = $14, bank_account = $15, bank_routing = $16, tax_id = $17,
			website = $18, description = $19, updated_at = $20
		WHERE id = $1
	`

	affiliate.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		affiliate.ID,
		affiliate.Status,
		affiliate.CommissionRate,
		affiliate.CommissionType,
		affiliate.FixedCommissionAmount,
		affiliate.TotalClicks,
		affiliate.TotalSignups,
		affiliate.TotalConversions,
		affiliate.TotalEarnings,
		affiliate.TotalPaid,
		affiliate.Balance,
		affiliate.MinimumPayout,
		affiliate.PaypalEmail,
		affiliate.BankName,
		affiliate.BankAccount,
		affiliate.BankRouting,
		affiliate.TaxID,
		affiliate.Website,
		affiliate.Description,
		affiliate.UpdatedAt,
	)

	return err
}

// Delete deletes an affiliate
func (r *affiliateRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM affiliates WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetStats retrieves affiliate statistics for an organization
func (r *affiliateRepository) GetStats(ctx context.Context, organizationID string) (*domain.AffiliateStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_affiliates,
			COUNT(*) FILTER (WHERE status = 'active') as active_affiliates,
			COALESCE(SUM(total_clicks), 0) as total_clicks,
			COALESCE(SUM(total_signups), 0) as total_signups,
			COALESCE(SUM(total_conversions), 0) as total_conversions,
			COALESCE(SUM(total_earnings), 0) as total_earnings,
			COALESCE(SUM(total_paid), 0) as total_paid,
			COALESCE(SUM(balance), 0) as pending_payouts
		FROM affiliates WHERE organization_id = $1
	`

	stats := &domain.AffiliateStats{}
	err := r.db.QueryRowContext(ctx, query, organizationID).Scan(
		&stats.TotalAffiliates,
		&stats.ActiveAffiliates,
		&stats.TotalClicks,
		&stats.TotalSignups,
		&stats.TotalConversions,
		&stats.TotalEarnings,
		&stats.TotalPaid,
		&stats.PendingPayouts,
	)

	if err != nil {
		return nil, err
	}

	// Calculate conversion rate
	if stats.TotalClicks > 0 {
		stats.ConversionRate = float64(stats.TotalConversions) / float64(stats.TotalClicks) * 100
	}

	// Get leaderboard
	leaderboard, err := r.GetLeaderboard(ctx, organizationID, 10)
	if err != nil {
		return nil, err
	}
	stats.TopAffiliates = leaderboard

	return stats, nil
}

// GetLeaderboard retrieves the top performing affiliates
func (r *affiliateRepository) GetLeaderboard(ctx context.Context, organizationID string, limit int) ([]domain.AffiliateLeaderboard, error) {
	query := `
		SELECT a.id, a.affiliate_code, COALESCE(u.first_name || ' ' || u.last_name, 'Unknown') as name,
			a.total_conversions, a.total_earnings
		FROM affiliates a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.organization_id = $1 AND a.status = 'active'
		ORDER BY a.total_earnings DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaderboard := make([]domain.AffiliateLeaderboard, 0)
	for rows.Next() {
		entry := domain.AffiliateLeaderboard{}
		err := rows.Scan(
			&entry.AffiliateID,
			&entry.AffiliateCode,
			&entry.Name,
			&entry.Conversions,
			&entry.Earnings,
		)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}

	return leaderboard, rows.Err()
}

// GetDashboard retrieves dashboard data for an affiliate
func (r *affiliateRepository) GetDashboard(ctx context.Context, affiliateID string) (*domain.AffiliateDashboard, error) {
	// Get affiliate
	affiliate, err := r.GetByID(ctx, affiliateID)
	if err != nil {
		return nil, err
	}

	dashboard := &domain.AffiliateDashboard{
		Affiliate: *affiliate,
	}

	// Get recent referrals
	referralQuery := `
		SELECT id, affiliate_id, affiliate_link_id, referred_user_id, referral_code, status,
			click_timestamp, signup_timestamp, conversion_timestamp, order_id, order_amount,
			commission_amount, commission_status, ip_address, user_agent, referrer_url, created_at, updated_at
		FROM referrals WHERE affiliate_id = $1
		ORDER BY created_at DESC
		LIMIT 10
	`

	rows, err := r.db.QueryContext(ctx, referralQuery, affiliateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dashboard.RecentReferrals = make([]domain.Referral, 0)
	for rows.Next() {
		referral := domain.Referral{}
		var signupTime, convTime sql.NullTime
		var referredUserID, linkID, orderID, ip, userAgent, referrerURL sql.NullString

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

		if signupTime.Valid {
			referral.SignupTimestamp = &signupTime.Time
		}
		if convTime.Valid {
			referral.ConversionTimestamp = &convTime.Time
		}
		if linkID.Valid {
			referral.AffiliateLinkID = linkID.String
		}
		if referredUserID.Valid {
			referral.ReferredUserID = referredUserID.String
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

		dashboard.RecentReferrals = append(dashboard.RecentReferrals, referral)
	}

	// Get recent commissions
	commissionQuery := `
		SELECT id, affiliate_id, referral_id, amount, status, type, description, reference_id,
			created_at, updated_at, approved_at, approved_by, paid_at, payout_id
		FROM commissions WHERE affiliate_id = $1
		ORDER BY created_at DESC
		LIMIT 10
	`

	commRows, err := r.db.QueryContext(ctx, commissionQuery, affiliateID)
	if err != nil {
		return nil, err
	}
	defer commRows.Close()

	dashboard.RecentCommissions = make([]domain.Commission, 0)
	for commRows.Next() {
		commission := domain.Commission{}
		var referralID, description, refID, approvedBy, payoutID sql.NullString
		var approvedAt, paidAt sql.NullTime

		err := commRows.Scan(
			&commission.ID,
			&commission.AffiliateID,
			&referralID,
			&commission.Amount,
			&commission.Status,
			&commission.Type,
			&description,
			&refID,
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
		if refID.Valid {
			commission.ReferenceID = refID.String
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

		dashboard.RecentCommissions = append(dashboard.RecentCommissions, commission)
	}

	// Get top links
	linkQuery := `
		SELECT id, affiliate_id, campaign_id, link_code, name, destination_url,
			click_count, signup_count, conversion_count, earnings, is_active, created_at, updated_at
		FROM affiliate_links WHERE affiliate_id = $1
		ORDER BY earnings DESC
		LIMIT 5
	`

	linkRows, err := r.db.QueryContext(ctx, linkQuery, affiliateID)
	if err != nil {
		return nil, err
	}
	defer linkRows.Close()

	dashboard.TopLinks = make([]domain.AffiliateLink, 0)
	for linkRows.Next() {
		link := domain.AffiliateLink{}
		var campaignID sql.NullString

		err := linkRows.Scan(
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

		dashboard.TopLinks = append(dashboard.TopLinks, link)
	}

	// Get monthly earnings history
	earningsQuery := `
		SELECT 
			DATE_TRUNC('month', created_at) as month,
			SUM(amount) as earnings,
			COUNT(*) as conversions
		FROM commissions
		WHERE affiliate_id = $1 AND status != 'cancelled' AND status != 'refunded'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY month DESC
		LIMIT 12
	`

	earnRows, err := r.db.QueryContext(ctx, earningsQuery, affiliateID)
	if err != nil {
		return nil, err
	}
	defer earnRows.Close()

	dashboard.EarningsHistory = make([]domain.MonthlyEarnings, 0)
	for earnRows.Next() {
		earning := domain.MonthlyEarnings{}
		var month time.Time
		err := earnRows.Scan(&month, &earning.Earnings, &earning.Conversions)
		if err != nil {
			return nil, err
		}
		earning.Month = month.Format("2006-01")
		dashboard.EarningsHistory = append(dashboard.EarningsHistory, earning)
	}

	return dashboard, nil
}