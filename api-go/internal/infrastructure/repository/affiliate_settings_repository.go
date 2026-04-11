package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// affiliateSettingsRepository implements domain.AffiliateSettingsRepository
type affiliateSettingsRepository struct {
	db *sql.DB
}

// NewAffiliateSettingsRepository creates a new affiliate settings repository
func NewAffiliateSettingsRepository(db *sql.DB) domain.AffiliateSettingsRepository {
	return &affiliateSettingsRepository{db: db}
}

// Create creates new affiliate settings
func (r *affiliateSettingsRepository) Create(ctx context.Context, settings *domain.AffiliateSettings) error {
	query := `
		INSERT INTO affiliate_settings (id, organization_id, program_name, program_description,
			default_commission_rate, default_commission_type, default_fixed_commission,
			cookie_duration_days, minimum_payout_amount, payout_schedule, require_approval,
			allow_custom_commission_rates, allow_affiliate_links, terms_and_conditions,
			welcome_message, email_notifications, notify_on_signup, notify_on_conversion,
			notify_on_payout, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`

	now := time.Now()
	settings.CreatedAt = now
	settings.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		settings.ID,
		settings.OrganizationID,
		settings.ProgramName,
		settings.ProgramDescription,
		settings.DefaultCommissionRate,
		settings.DefaultCommissionType,
		settings.DefaultFixedCommission,
		settings.CookieDurationDays,
		settings.MinimumPayoutAmount,
		settings.PayoutSchedule,
		settings.RequireApproval,
		settings.AllowCustomCommissionRates,
		settings.AllowAffiliateLinks,
		settings.TermsAndConditions,
		settings.WelcomeMessage,
		settings.EmailNotifications,
		settings.NotifyOnSignup,
		settings.NotifyOnConversion,
		settings.NotifyOnPayout,
		settings.CreatedAt,
		settings.UpdatedAt,
	)

	return err
}

// GetByID retrieves affiliate settings by ID
func (r *affiliateSettingsRepository) GetByID(ctx context.Context, id string) (*domain.AffiliateSettings, error) {
	query := `
		SELECT id, organization_id, program_name, program_description,
			default_commission_rate, default_commission_type, default_fixed_commission,
			cookie_duration_days, minimum_payout_amount, payout_schedule, require_approval,
			allow_custom_commission_rates, allow_affiliate_links, terms_and_conditions,
			welcome_message, email_notifications, notify_on_signup, notify_on_conversion,
			notify_on_payout, created_at, updated_at
		FROM affiliate_settings WHERE id = $1
	`

	settings := &domain.AffiliateSettings{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&settings.ID,
		&settings.OrganizationID,
		&settings.ProgramName,
		&settings.ProgramDescription,
		&settings.DefaultCommissionRate,
		&settings.DefaultCommissionType,
		&settings.DefaultFixedCommission,
		&settings.CookieDurationDays,
		&settings.MinimumPayoutAmount,
		&settings.PayoutSchedule,
		&settings.RequireApproval,
		&settings.AllowCustomCommissionRates,
		&settings.AllowAffiliateLinks,
		&settings.TermsAndConditions,
		&settings.WelcomeMessage,
		&settings.EmailNotifications,
		&settings.NotifyOnSignup,
		&settings.NotifyOnConversion,
		&settings.NotifyOnPayout,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return settings, nil
}

// GetByOrganizationID retrieves affiliate settings by organization ID
func (r *affiliateSettingsRepository) GetByOrganizationID(ctx context.Context, organizationID string) (*domain.AffiliateSettings, error) {
	query := `
		SELECT id, organization_id, program_name, program_description,
			default_commission_rate, default_commission_type, default_fixed_commission,
			cookie_duration_days, minimum_payout_amount, payout_schedule, require_approval,
			allow_custom_commission_rates, allow_affiliate_links, terms_and_conditions,
			welcome_message, email_notifications, notify_on_signup, notify_on_conversion,
			notify_on_payout, created_at, updated_at
		FROM affiliate_settings WHERE organization_id = $1
	`

	settings := &domain.AffiliateSettings{}
	err := r.db.QueryRowContext(ctx, query, organizationID).Scan(
		&settings.ID,
		&settings.OrganizationID,
		&settings.ProgramName,
		&settings.ProgramDescription,
		&settings.DefaultCommissionRate,
		&settings.DefaultCommissionType,
		&settings.DefaultFixedCommission,
		&settings.CookieDurationDays,
		&settings.MinimumPayoutAmount,
		&settings.PayoutSchedule,
		&settings.RequireApproval,
		&settings.AllowCustomCommissionRates,
		&settings.AllowAffiliateLinks,
		&settings.TermsAndConditions,
		&settings.WelcomeMessage,
		&settings.EmailNotifications,
		&settings.NotifyOnSignup,
		&settings.NotifyOnConversion,
		&settings.NotifyOnPayout,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return settings, nil
}

// Update updates existing affiliate settings
func (r *affiliateSettingsRepository) Update(ctx context.Context, settings *domain.AffiliateSettings) error {
	query := `
		UPDATE affiliate_settings SET
			program_name = $2, program_description = $3, default_commission_rate = $4,
			default_commission_type = $5, default_fixed_commission = $6, cookie_duration_days = $7,
			minimum_payout_amount = $8, payout_schedule = $9, require_approval = $10,
			allow_custom_commission_rates = $11, allow_affiliate_links = $12, terms_and_conditions = $13,
			welcome_message = $14, email_notifications = $15, notify_on_signup = $16,
			notify_on_conversion = $17, notify_on_payout = $18, updated_at = $19
		WHERE id = $1
	`

	settings.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		settings.ID,
		settings.ProgramName,
		settings.ProgramDescription,
		settings.DefaultCommissionRate,
		settings.DefaultCommissionType,
		settings.DefaultFixedCommission,
		settings.CookieDurationDays,
		settings.MinimumPayoutAmount,
		settings.PayoutSchedule,
		settings.RequireApproval,
		settings.AllowCustomCommissionRates,
		settings.AllowAffiliateLinks,
		settings.TermsAndConditions,
		settings.WelcomeMessage,
		settings.EmailNotifications,
		settings.NotifyOnSignup,
		settings.NotifyOnConversion,
		settings.NotifyOnPayout,
		settings.UpdatedAt,
	)

	return err
}