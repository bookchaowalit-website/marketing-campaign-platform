package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// externalAffiliateAccountRepository implements domain.ExternalAffiliateAccountRepository
type externalAffiliateAccountRepository struct {
	db *sql.DB
}

// NewExternalAffiliateAccountRepository creates a new external affiliate account repository
func NewExternalAffiliateAccountRepository(db *sql.DB) domain.ExternalAffiliateAccountRepository {
	return &externalAffiliateAccountRepository{db: db}
}

// Create creates a new external affiliate account
func (r *externalAffiliateAccountRepository) Create(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	query := `
		INSERT INTO external_affiliate_accounts (
			id, user_id, organization_id, platform, platform_account_id, platform_account_name,
			platform_username, status, api_key, api_secret, access_token, refresh_token,
			token_expires_at, auto_sync, sync_interval_minutes, total_earnings, total_orders,
			total_clicks, total_products, pending_commission, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	var tokenExpiresAt interface{}
	if account.TokenExpiresAt != nil {
		tokenExpiresAt = *account.TokenExpiresAt
	}

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.OrganizationID,
		account.Platform,
		account.PlatformAccountID,
		account.PlatformAccountName,
		account.PlatformUsername,
		account.Status,
		account.APIKey,
		account.APISecret,
		account.AccessToken,
		account.RefreshToken,
		tokenExpiresAt,
		account.AutoSync,
		account.SyncIntervalMins,
		account.TotalEarnings,
		account.TotalOrders,
		account.TotalClicks,
		account.TotalProducts,
		account.PendingCommission,
		account.CreatedAt,
		account.UpdatedAt,
	)

	return err
}

// GetByID retrieves an external affiliate account by ID
func (r *externalAffiliateAccountRepository) GetByID(ctx context.Context, id string) (*domain.ExternalAffiliateAccount, error) {
	query := `
		SELECT id, user_id, organization_id, platform, platform_account_id, platform_account_name,
			platform_username, status, api_key, api_secret, access_token, refresh_token,
			token_expires_at, auto_sync, sync_interval_minutes, last_sync_at, last_sync_status,
			last_sync_error, total_earnings, total_orders, total_clicks, total_products,
			pending_commission, created_at, updated_at
		FROM external_affiliate_accounts WHERE id = $1
	`

	account := &domain.ExternalAffiliateAccount{}
	var lastSyncAt, tokenExpiresAt sql.NullTime
	var lastSyncError sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.OrganizationID,
		&account.Platform,
		&account.PlatformAccountID,
		&account.PlatformAccountName,
		&account.PlatformUsername,
		&account.Status,
		&account.APIKey,
		&account.APISecret,
		&account.AccessToken,
		&account.RefreshToken,
		&tokenExpiresAt,
		&account.AutoSync,
		&account.SyncIntervalMins,
		&lastSyncAt,
		&account.LastSyncStatus,
		&lastSyncError,
		&account.TotalEarnings,
		&account.TotalOrders,
		&account.TotalClicks,
		&account.TotalProducts,
		&account.PendingCommission,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if lastSyncAt.Valid {
		account.LastSyncAt = &lastSyncAt.Time
	}
	if tokenExpiresAt.Valid {
		account.TokenExpiresAt = &tokenExpiresAt.Time
	}
	if lastSyncError.Valid {
		account.LastSyncError = lastSyncError.String
	}

	return account, nil
}

// GetByUserIDAndPlatform retrieves an external affiliate account by user ID and platform
func (r *externalAffiliateAccountRepository) GetByUserIDAndPlatform(ctx context.Context, userID string, platform domain.Platform) (*domain.ExternalAffiliateAccount, error) {
	query := `
		SELECT id, user_id, organization_id, platform, platform_account_id, platform_account_name,
			platform_username, status, api_key, api_secret, access_token, refresh_token,
			token_expires_at, auto_sync, sync_interval_minutes, last_sync_at, last_sync_status,
			last_sync_error, total_earnings, total_orders, total_clicks, total_products,
			pending_commission, created_at, updated_at
		FROM external_affiliate_accounts WHERE user_id = $1 AND platform = $2
	`

	account := &domain.ExternalAffiliateAccount{}
	var lastSyncAt, tokenExpiresAt sql.NullTime
	var lastSyncError sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID, platform).Scan(
		&account.ID,
		&account.UserID,
		&account.OrganizationID,
		&account.Platform,
		&account.PlatformAccountID,
		&account.PlatformAccountName,
		&account.PlatformUsername,
		&account.Status,
		&account.APIKey,
		&account.APISecret,
		&account.AccessToken,
		&account.RefreshToken,
		&tokenExpiresAt,
		&account.AutoSync,
		&account.SyncIntervalMins,
		&lastSyncAt,
		&account.LastSyncStatus,
		&lastSyncError,
		&account.TotalEarnings,
		&account.TotalOrders,
		&account.TotalClicks,
		&account.TotalProducts,
		&account.PendingCommission,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if lastSyncAt.Valid {
		account.LastSyncAt = &lastSyncAt.Time
	}
	if tokenExpiresAt.Valid {
		account.TokenExpiresAt = &tokenExpiresAt.Time
	}
	if lastSyncError.Valid {
		account.LastSyncError = lastSyncError.String
	}

	return account, nil
}

// GetAllByUserID retrieves all external affiliate accounts for a user
func (r *externalAffiliateAccountRepository) GetAllByUserID(ctx context.Context, userID string) ([]*domain.ExternalAffiliateAccount, error) {
	query := `
		SELECT id, user_id, organization_id, platform, platform_account_id, platform_account_name,
			platform_username, status, auto_sync, sync_interval_minutes, last_sync_at, last_sync_status,
			last_sync_error, total_earnings, total_orders, total_clicks, total_products,
			pending_commission, created_at, updated_at
		FROM external_affiliate_accounts WHERE user_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*domain.ExternalAffiliateAccount, 0)
	for rows.Next() {
		account := &domain.ExternalAffiliateAccount{}
		var lastSyncAt sql.NullTime
		var lastSyncError sql.NullString

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.OrganizationID,
			&account.Platform,
			&account.PlatformAccountID,
			&account.PlatformAccountName,
			&account.PlatformUsername,
			&account.Status,
			&account.AutoSync,
			&account.SyncIntervalMins,
			&lastSyncAt,
			&account.LastSyncStatus,
			&lastSyncError,
			&account.TotalEarnings,
			&account.TotalOrders,
			&account.TotalClicks,
			&account.TotalProducts,
			&account.PendingCommission,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastSyncAt.Valid {
			account.LastSyncAt = &lastSyncAt.Time
		}
		if lastSyncError.Valid {
			account.LastSyncError = lastSyncError.String
		}

		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

// GetAllByOrganizationID retrieves all external affiliate accounts for an organization
func (r *externalAffiliateAccountRepository) GetAllByOrganizationID(ctx context.Context, organizationID string) ([]*domain.ExternalAffiliateAccount, error) {
	query := `
		SELECT id, user_id, organization_id, platform, platform_account_id, platform_account_name,
			platform_username, status, auto_sync, sync_interval_minutes, last_sync_at, last_sync_status,
			last_sync_error, total_earnings, total_orders, total_clicks, total_products,
			pending_commission, created_at, updated_at
		FROM external_affiliate_accounts WHERE organization_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*domain.ExternalAffiliateAccount, 0)
	for rows.Next() {
		account := &domain.ExternalAffiliateAccount{}
		var lastSyncAt sql.NullTime
		var lastSyncError sql.NullString

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.OrganizationID,
			&account.Platform,
			&account.PlatformAccountID,
			&account.PlatformAccountName,
			&account.PlatformUsername,
			&account.Status,
			&account.AutoSync,
			&account.SyncIntervalMins,
			&lastSyncAt,
			&account.LastSyncStatus,
			&lastSyncError,
			&account.TotalEarnings,
			&account.TotalOrders,
			&account.TotalClicks,
			&account.TotalProducts,
			&account.PendingCommission,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastSyncAt.Valid {
			account.LastSyncAt = &lastSyncAt.Time
		}
		if lastSyncError.Valid {
			account.LastSyncError = lastSyncError.String
		}

		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

// Update updates an external affiliate account
func (r *externalAffiliateAccountRepository) Update(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	query := `
		UPDATE external_affiliate_accounts SET
			platform_account_name = $2, platform_username = $3, status = $4,
			auto_sync = $5, sync_interval_minutes = $6, updated_at = $7
		WHERE id = $1
	`

	account.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.PlatformAccountName,
		account.PlatformUsername,
		account.Status,
		account.AutoSync,
		account.SyncIntervalMins,
		account.UpdatedAt,
	)

	return err
}

// UpdateTokens updates the access tokens for an account
func (r *externalAffiliateAccountRepository) UpdateTokens(ctx context.Context, id, accessToken, refreshToken string, expiresAt *time.Time) error {
	query := `
		UPDATE external_affiliate_accounts SET
			access_token = $2, refresh_token = $3, token_expires_at = $4, status = 'active', updated_at = $5
		WHERE id = $1
	`

	var expiresAtVal interface{}
	if expiresAt != nil {
		expiresAtVal = *expiresAt
	}

	_, err := r.db.ExecContext(ctx, query, id, accessToken, refreshToken, expiresAtVal, time.Now())
	return err
}

// UpdateSyncStatus updates the sync status for an account
func (r *externalAffiliateAccountRepository) UpdateSyncStatus(ctx context.Context, id, status, errorMsg string) error {
	query := `
		UPDATE external_affiliate_accounts SET
			last_sync_at = $2, last_sync_status = $3, last_sync_error = $4, updated_at = $5
		WHERE id = $1
	`

	var errorMsgVal interface{}
	if errorMsg != "" {
		errorMsgVal = errorMsg
	}

	_, err := r.db.ExecContext(ctx, query, id, time.Now(), status, errorMsgVal, time.Now())
	return err
}

// UpdateStats updates the cached stats for an account
func (r *externalAffiliateAccountRepository) UpdateStats(ctx context.Context, id string, earnings float64, orders, products int, clicks int64, pendingCommission float64) error {
	query := `
		UPDATE external_affiliate_accounts SET
			total_earnings = $2, total_orders = $3, total_clicks = $4, total_products = $5,
			pending_commission = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, earnings, orders, clicks, products, pendingCommission, time.Now())
	return err
}

// Delete deletes an external affiliate account
func (r *externalAffiliateAccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM external_affiliate_accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}