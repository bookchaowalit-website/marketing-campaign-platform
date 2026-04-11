package domain

import (
	"context"
	"time"
)

// CampaignRepository defines the interface for campaign repository
type CampaignRepository interface {
	Create(ctx context.Context, campaign *Campaign) error
	GetByID(ctx context.Context, id string) (*Campaign, error)
	GetAllByOrganizationID(ctx context.Context, organizationID string, limit, offset int) ([]*Campaign, error)
	Update(ctx context.Context, campaign *Campaign) error
	Delete(ctx context.Context, id string) error
	GetByStatus(ctx context.Context, organizationID, status string, limit, offset int) ([]*Campaign, error)
	GetByType(ctx context.Context, organizationID, campaignType string, limit, offset int) ([]*Campaign, error)
}

// DashboardRepository defines the interface for dashboard repository
type DashboardRepository interface {
	GetOverview(ctx context.Context, organizationID string) (*Dashboard, error)
}

// AffiliateRepository defines the interface for affiliate repository
type AffiliateRepository interface {
	Create(ctx context.Context, affiliate *Affiliate) error
	GetByID(ctx context.Context, id string) (*Affiliate, error)
	GetByUserID(ctx context.Context, userID string) (*Affiliate, error)
	GetByCode(ctx context.Context, code string) (*Affiliate, error)
	GetAllByOrganizationID(ctx context.Context, organizationID string, limit, offset int) ([]*Affiliate, error)
	GetByStatus(ctx context.Context, organizationID, status string, limit, offset int) ([]*Affiliate, error)
	Update(ctx context.Context, affiliate *Affiliate) error
	Delete(ctx context.Context, id string) error
	GetStats(ctx context.Context, organizationID string) (*AffiliateStats, error)
	GetLeaderboard(ctx context.Context, organizationID string, limit int) ([]AffiliateLeaderboard, error)
	GetDashboard(ctx context.Context, affiliateID string) (*AffiliateDashboard, error)
}

// AffiliateLinkRepository defines the interface for affiliate link repository
type AffiliateLinkRepository interface {
	Create(ctx context.Context, link *AffiliateLink) error
	GetByID(ctx context.Context, id string) (*AffiliateLink, error)
	GetByCode(ctx context.Context, code string) (*AffiliateLink, error)
	GetAllByAffiliateID(ctx context.Context, affiliateID string, limit, offset int) ([]*AffiliateLink, error)
	Update(ctx context.Context, link *AffiliateLink) error
	Delete(ctx context.Context, id string) error
	IncrementClick(ctx context.Context, linkID string) error
	IncrementSignup(ctx context.Context, linkID string) error
	IncrementConversion(ctx context.Context, linkID string, earnings float64) error
}

// ReferralRepository defines the interface for referral repository
type ReferralRepository interface {
	Create(ctx context.Context, referral *Referral) error
	GetByID(ctx context.Context, id string) (*Referral, error)
	GetByReferralCode(ctx context.Context, code string) (*Referral, error)
	GetAllByAffiliateID(ctx context.Context, affiliateID string, limit, offset int) ([]*Referral, error)
	GetByStatus(ctx context.Context, affiliateID, status string, limit, offset int) ([]*Referral, error)
	Update(ctx context.Context, referral *Referral) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetRecentByAffiliateID(ctx context.Context, affiliateID string, limit int) ([]*Referral, error)
}

// CommissionRepository defines the interface for commission repository
type CommissionRepository interface {
	Create(ctx context.Context, commission *Commission) error
	GetByID(ctx context.Context, id string) (*Commission, error)
	GetAllByAffiliateID(ctx context.Context, affiliateID string, limit, offset int) ([]*Commission, error)
	GetByStatus(ctx context.Context, affiliateID, status string, limit, offset int) ([]*Commission, error)
	GetPendingByAffiliateID(ctx context.Context, affiliateID string) ([]*Commission, error)
	Update(ctx context.Context, commission *Commission) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Approve(ctx context.Context, id string, approvedBy string) error
	GetTotalPendingByAffiliateID(ctx context.Context, affiliateID string) (float64, error)
	GetRecentByAffiliateID(ctx context.Context, affiliateID string, limit int) ([]*Commission, error)
}

// AffiliatePayoutRepository defines the interface for affiliate payout repository
type AffiliatePayoutRepository interface {
	Create(ctx context.Context, payout *AffiliatePayout) error
	GetByID(ctx context.Context, id string) (*AffiliatePayout, error)
	GetAllByAffiliateID(ctx context.Context, affiliateID string, limit, offset int) ([]*AffiliatePayout, error)
	GetByStatus(ctx context.Context, affiliateID, status string, limit, offset int) ([]*AffiliatePayout, error)
	Update(ctx context.Context, payout *AffiliatePayout) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetPendingPayouts(ctx context.Context, organizationID string) ([]*AffiliatePayout, error)
}

// AffiliateSettingsRepository defines the interface for affiliate settings repository
type AffiliateSettingsRepository interface {
	Create(ctx context.Context, settings *AffiliateSettings) error
	GetByID(ctx context.Context, id string) (*AffiliateSettings, error)
	GetByOrganizationID(ctx context.Context, organizationID string) (*AffiliateSettings, error)
	Update(ctx context.Context, settings *AffiliateSettings) error
}

// ExternalAffiliateAccountRepository defines the interface for external affiliate accounts
type ExternalAffiliateAccountRepository interface {
	Create(ctx context.Context, account *ExternalAffiliateAccount) error
	GetByID(ctx context.Context, id string) (*ExternalAffiliateAccount, error)
	GetByUserIDAndPlatform(ctx context.Context, userID string, platform Platform) (*ExternalAffiliateAccount, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*ExternalAffiliateAccount, error)
	GetAllByOrganizationID(ctx context.Context, organizationID string) ([]*ExternalAffiliateAccount, error)
	Update(ctx context.Context, account *ExternalAffiliateAccount) error
	UpdateTokens(ctx context.Context, id, accessToken, refreshToken string, expiresAt *time.Time) error
	UpdateSyncStatus(ctx context.Context, id, status, errorMsg string) error
	UpdateStats(ctx context.Context, id string, earnings float64, orders, products int, clicks int64, pendingCommission float64) error
	Delete(ctx context.Context, id string) error
}

// ExternalAffiliateProductRepository defines the interface for external affiliate products
type ExternalAffiliateProductRepository interface {
	Create(ctx context.Context, product *ExternalAffiliateProduct) error
	GetByID(ctx context.Context, id string) (*ExternalAffiliateProduct, error)
	GetByAccountAndProductID(ctx context.Context, accountID, productID string) (*ExternalAffiliateProduct, error)
	GetAllByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*ExternalAffiliateProduct, error)
	GetTopProducts(ctx context.Context, userID string, limit int) ([]*ExternalAffiliateProduct, error)
	Update(ctx context.Context, product *ExternalAffiliateProduct) error
	UpdateStats(ctx context.Context, id string, clicks int64, orders int, earnings float64) error
	Upsert(ctx context.Context, product *ExternalAffiliateProduct) error
	Delete(ctx context.Context, id string) error
}

// ExternalAffiliateOrderRepository defines the interface for external affiliate orders
type ExternalAffiliateOrderRepository interface {
	Create(ctx context.Context, order *ExternalAffiliateOrder) error
	GetByID(ctx context.Context, id string) (*ExternalAffiliateOrder, error)
	GetByAccountAndOrderID(ctx context.Context, accountID, orderID string) (*ExternalAffiliateOrder, error)
	GetAllByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*ExternalAffiliateOrder, error)
	GetRecentByUserID(ctx context.Context, userID string, limit int) ([]*ExternalAffiliateOrder, error)
	GetByDateRange(ctx context.Context, accountID string, start, end time.Time) ([]*ExternalAffiliateOrder, error)
	Update(ctx context.Context, order *ExternalAffiliateOrder) error
	Upsert(ctx context.Context, order *ExternalAffiliateOrder) error
	Delete(ctx context.Context, id string) error
}

// ExternalAffiliateEarningsRepository defines the interface for external affiliate earnings
type ExternalAffiliateEarningsRepository interface {
	Create(ctx context.Context, earnings *ExternalAffiliateEarnings) error
	GetByAccountAndDate(ctx context.Context, accountID string, date time.Time) (*ExternalAffiliateEarnings, error)
	GetByDateRange(ctx context.Context, accountID string, start, end time.Time) ([]*ExternalAffiliateEarnings, error)
	GetByUserIDAndDateRange(ctx context.Context, userID string, start, end time.Time) ([]*ExternalAffiliateEarnings, error)
	Update(ctx context.Context, earnings *ExternalAffiliateEarnings) error
	Upsert(ctx context.Context, earnings *ExternalAffiliateEarnings) error
}

// ExternalAffiliateSyncLogRepository defines the interface for sync logs
type ExternalAffiliateSyncLogRepository interface {
	Create(ctx context.Context, log *ExternalAffiliateSyncLog) error
	GetByID(ctx context.Context, id string) (*ExternalAffiliateSyncLog, error)
	GetRecentByAccountID(ctx context.Context, accountID string, limit int) ([]*ExternalAffiliateSyncLog, error)
	UpdateStatus(ctx context.Context, id string, status string, completedAt time.Time, recordsCreated, recordsUpdated, recordsFailed int, errorMsg string) error
}

// ExternalAffiliateLinkRepository defines the interface for external affiliate links
type ExternalAffiliateLinkRepository interface {
	Create(ctx context.Context, link *ExternalAffiliateLink) error
	GetByID(ctx context.Context, id string) (*ExternalAffiliateLink, error)
	GetByCode(ctx context.Context, code string) (*ExternalAffiliateLink, error)
	GetAllByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*ExternalAffiliateLink, error)
	Update(ctx context.Context, link *ExternalAffiliateLink) error
	UpdateStats(ctx context.Context, id string, clicks, uniqueClicks int64, orders int, earnings float64) error
	Delete(ctx context.Context, id string) error
}
