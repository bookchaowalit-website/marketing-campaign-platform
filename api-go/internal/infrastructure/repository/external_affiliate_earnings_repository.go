package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// externalAffiliateEarningsRepository implements domain.ExternalAffiliateEarningsRepository
type externalAffiliateEarningsRepository struct {
	db *sql.DB
}

// NewExternalAffiliateEarningsRepository creates a new earnings repository
func NewExternalAffiliateEarningsRepository(db *sql.DB) domain.ExternalAffiliateEarningsRepository {
	return &externalAffiliateEarningsRepository{db: db}
}

// Create creates a new earnings record
func (r *externalAffiliateEarningsRepository) Create(ctx context.Context, earnings *domain.ExternalAffiliateEarnings) error {
	query := `
		INSERT INTO external_affiliate_earnings (
			id, account_id, platform, date, gross_sales, net_sales, orders_count, items_sold,
			commission_earned, commission_pending, commission_confirmed, clicks, impressions,
			conversion_rate, epc, currency, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	now := time.Now()
	earnings.CreatedAt = now
	earnings.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		earnings.ID, earnings.AccountID, earnings.Platform, earnings.Date,
		earnings.GrossSales, earnings.NetSales, earnings.OrdersCount, earnings.ItemsSold,
		earnings.CommissionEarned, earnings.CommissionPending, earnings.CommissionConfirmed,
		earnings.Clicks, earnings.Impressions, earnings.ConversionRate, earnings.EPC,
		earnings.Currency, earnings.CreatedAt, earnings.UpdatedAt,
	)

	return err
}

// GetByAccountAndDate retrieves earnings by account and date
func (r *externalAffiliateEarningsRepository) GetByAccountAndDate(ctx context.Context, accountID string, date time.Time) (*domain.ExternalAffiliateEarnings, error) {
	query := `
		SELECT id, account_id, platform, date, gross_sales, net_sales, orders_count, items_sold,
			commission_earned, commission_pending, commission_confirmed, clicks, impressions,
			conversion_rate, epc, currency, created_at, updated_at
		FROM external_affiliate_earnings WHERE account_id = $1 AND date = $2
	`

	earnings := &domain.ExternalAffiliateEarnings{}
	err := r.db.QueryRowContext(ctx, query, accountID, date).Scan(
		&earnings.ID, &earnings.AccountID, &earnings.Platform, &earnings.Date,
		&earnings.GrossSales, &earnings.NetSales, &earnings.OrdersCount, &earnings.ItemsSold,
		&earnings.CommissionEarned, &earnings.CommissionPending, &earnings.CommissionConfirmed,
		&earnings.Clicks, &earnings.Impressions, &earnings.ConversionRate, &earnings.EPC,
		&earnings.Currency, &earnings.CreatedAt, &earnings.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return earnings, nil
}

// GetByDateRange retrieves earnings within a date range for an account
func (r *externalAffiliateEarningsRepository) GetByDateRange(ctx context.Context, accountID string, start, end time.Time) ([]*domain.ExternalAffiliateEarnings, error) {
	query := `
		SELECT id, account_id, platform, date, gross_sales, net_sales, orders_count, items_sold,
			commission_earned, commission_pending, commission_confirmed, clicks, impressions,
			conversion_rate, epc, currency, created_at, updated_at
		FROM external_affiliate_earnings WHERE account_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEarnings(rows)
}

// GetByUserIDAndDateRange retrieves earnings within a date range for a user
func (r *externalAffiliateEarningsRepository) GetByUserIDAndDateRange(ctx context.Context, userID string, start, end time.Time) ([]*domain.ExternalAffiliateEarnings, error) {
	query := `
		SELECT e.id, e.account_id, e.platform, e.date, e.gross_sales, e.net_sales, e.orders_count, e.items_sold,
			e.commission_earned, e.commission_pending, e.commission_confirmed, e.clicks, e.impressions,
			e.conversion_rate, e.epc, e.currency, e.created_at, e.updated_at
		FROM external_affiliate_earnings e
		JOIN external_affiliate_accounts a ON e.account_id = a.id
		WHERE a.user_id = $1 AND e.date >= $2 AND e.date <= $3
		ORDER BY e.date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEarnings(rows)
}

// Update updates earnings
func (r *externalAffiliateEarningsRepository) Update(ctx context.Context, earnings *domain.ExternalAffiliateEarnings) error {
	query := `
		UPDATE external_affiliate_earnings SET
			gross_sales = $2, net_sales = $3, orders_count = $4, items_sold = $5,
			commission_earned = $6, commission_pending = $7, commission_confirmed = $8,
			clicks = $9, impressions = $10, conversion_rate = $11, epc = $12, updated_at = $13
		WHERE id = $1
	`

	earnings.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		earnings.ID, earnings.GrossSales, earnings.NetSales, earnings.OrdersCount, earnings.ItemsSold,
		earnings.CommissionEarned, earnings.CommissionPending, earnings.CommissionConfirmed,
		earnings.Clicks, earnings.Impressions, earnings.ConversionRate, earnings.EPC, earnings.UpdatedAt,
	)

	return err
}

// Upsert creates or updates earnings
func (r *externalAffiliateEarningsRepository) Upsert(ctx context.Context, earnings *domain.ExternalAffiliateEarnings) error {
	existing, err := r.GetByAccountAndDate(ctx, earnings.AccountID, earnings.Date)
	if err == nil {
		earnings.ID = existing.ID
		earnings.CreatedAt = existing.CreatedAt
		return r.Update(ctx, earnings)
	}

	return r.Create(ctx, earnings)
}

// Helper function to scan earnings
func (r *externalAffiliateEarningsRepository) scanEarnings(rows *sql.Rows) ([]*domain.ExternalAffiliateEarnings, error) {
	earnings := make([]*domain.ExternalAffiliateEarnings, 0)
	for rows.Next() {
		e := &domain.ExternalAffiliateEarnings{}
		err := rows.Scan(
			&e.ID, &e.AccountID, &e.Platform, &e.Date,
			&e.GrossSales, &e.NetSales, &e.OrdersCount, &e.ItemsSold,
			&e.CommissionEarned, &e.CommissionPending, &e.CommissionConfirmed,
			&e.Clicks, &e.Impressions, &e.ConversionRate, &e.EPC,
			&e.Currency, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		earnings = append(earnings, e)
	}

	return earnings, rows.Err()
}