package service

import (
	"context"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// ExternalAffiliateService handles business logic for external affiliate integrations
type ExternalAffiliateService struct {
	accountRepo   domain.ExternalAffiliateAccountRepository
	productRepo   domain.ExternalAffiliateProductRepository
	orderRepo     domain.ExternalAffiliateOrderRepository
	earningsRepo  domain.ExternalAffiliateEarningsRepository
	syncLogRepo   domain.ExternalAffiliateSyncLogRepository

	integrations map[domain.Platform]PlatformIntegration
}

// NewExternalAffiliateService creates a new external affiliate service
func NewExternalAffiliateService(
	accountRepo domain.ExternalAffiliateAccountRepository,
	productRepo domain.ExternalAffiliateProductRepository,
	orderRepo domain.ExternalAffiliateOrderRepository,
	earningsRepo domain.ExternalAffiliateEarningsRepository,
	syncLogRepo domain.ExternalAffiliateSyncLogRepository,
) *ExternalAffiliateService {
	return &ExternalAffiliateService{
		accountRepo:  accountRepo,
		productRepo:  productRepo,
		orderRepo:    orderRepo,
		earningsRepo: earningsRepo,
		syncLogRepo:  syncLogRepo,
		integrations: make(map[domain.Platform]PlatformIntegration),
	}
}

// RegisterIntegration registers a platform integration
func (s *ExternalAffiliateService) RegisterIntegration(platform domain.Platform, integration PlatformIntegration) {
	s.integrations[platform] = integration
}

// ConnectPlatform connects a platform account
func (s *ExternalAffiliateService) ConnectPlatform(ctx context.Context, userID, organizationID string, req *domain.PlatformConnectionRequest) (*domain.ExternalAffiliateAccount, error) {
	integration, ok := s.integrations[req.Platform]
	if !ok {
		return nil, &PlatformAPIError{
			Platform: req.Platform,
			Message:  "platform not supported",
		}
	}

	account, err := integration.Connect(ctx, req)
	if err != nil {
		return nil, err
	}

	account.UserID = userID
	account.OrganizationID = organizationID

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// DisconnectPlatform disconnects a platform account
func (s *ExternalAffiliateService) DisconnectPlatform(ctx context.Context, accountID string) error {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	integration, ok := s.integrations[account.Platform]
	if ok {
		if err := integration.Disconnect(ctx, accountID); err != nil {
			// Log but don't fail
		}
	}

	return s.accountRepo.Delete(ctx, accountID)
}

// GetConnectedAccounts retrieves all connected accounts for a user
func (s *ExternalAffiliateService) GetConnectedAccounts(ctx context.Context, userID string) ([]*domain.ExternalAffiliateAccount, error) {
	return s.accountRepo.GetAllByUserID(ctx, userID)
}

// GetAccount retrieves an account by ID
func (s *ExternalAffiliateService) GetAccount(ctx context.Context, accountID string) (*domain.ExternalAffiliateAccount, error) {
	return s.accountRepo.GetByID(ctx, accountID)
}

// SyncAccount triggers a sync for an account
func (s *ExternalAffiliateService) SyncAccount(ctx context.Context, accountID string, syncType string) error {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	integration, ok := s.integrations[account.Platform]
	if !ok {
		return &PlatformAPIError{
			Platform: account.Platform,
			Message:  "platform integration not available",
		}
	}

	// Update sync status to indicate sync started
	s.accountRepo.UpdateSyncStatus(ctx, accountID, "pending", "")

	// Perform the sync
	if err := integration.FullSync(ctx, account); err != nil {
		s.accountRepo.UpdateSyncStatus(ctx, accountID, "failed", err.Error())
		return err
	}

	return nil
}

// GetUnifiedDashboard retrieves a unified dashboard for all connected platforms
func (s *ExternalAffiliateService) GetUnifiedDashboard(ctx context.Context, userID string) (*domain.UnifiedAffiliateDashboard, error) {
	// Get all connected accounts
	accounts, err := s.accountRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dashboard := &domain.UnifiedAffiliateDashboard{
		PlatformStats:     make([]domain.PlatformStats, 0),
		RecentOrders:      make([]domain.ExternalAffiliateOrder, 0),
		TopProducts:       make([]domain.ExternalAffiliateProduct, 0),
		EarningsTrend:     make([]domain.DailyEarnings, 0),
		ConnectedAccounts: make([]domain.ExternalAffiliateAccount, 0),
	}

	// Aggregate stats by platform
	platformStatsMap := make(map[domain.Platform]*domain.PlatformStats)

	for _, account := range accounts {
		platform := account.Platform
		if _, exists := platformStatsMap[platform]; !exists {
			platformStatsMap[platform] = &domain.PlatformStats{
				Platform:     platform,
				AccountCount: 0,
			}
		}

		stats := platformStatsMap[platform]
		stats.AccountCount++
		stats.TotalEarnings += account.TotalEarnings
		stats.TotalOrders += account.TotalOrders
		stats.TotalClicks += account.TotalClicks
		stats.PendingCommission += account.PendingCommission
		stats.LastSyncAt = account.LastSyncAt

		dashboard.TotalEarnings += account.TotalEarnings
		dashboard.TotalOrders += account.TotalOrders
		dashboard.TotalClicks += account.TotalClicks
		dashboard.TotalProducts += account.TotalProducts
		dashboard.TotalPendingCommission += account.PendingCommission

		dashboard.ConnectedAccounts = append(dashboard.ConnectedAccounts, *account)
	}

	// Convert map to slice
	for _, stats := range platformStatsMap {
		// Calculate conversion rate
		if stats.TotalClicks > 0 {
			stats.AvgConversionRate = float64(stats.TotalOrders) / float64(stats.TotalClicks) * 100
		}
		dashboard.PlatformStats = append(dashboard.PlatformStats, *stats)
	}

	// Get recent orders
	recentOrders, err := s.orderRepo.GetRecentByUserID(ctx, userID, 10)
	if err == nil {
		for _, order := range recentOrders {
			dashboard.RecentOrders = append(dashboard.RecentOrders, *order)
		}
	}

	// Get top products
	topProducts, err := s.productRepo.GetTopProducts(ctx, userID, 10)
	if err == nil {
		for _, product := range topProducts {
			dashboard.TopProducts = append(dashboard.TopProducts, *product)
		}
	}

	// Get earnings trend (last 30 days)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	earningsTrend, err := s.earningsRepo.GetByUserIDAndDateRange(ctx, userID, startDate, endDate)
	if err == nil {
		// Group by date
		dateEarnings := make(map[time.Time]*domain.DailyEarnings)
		for _, e := range earningsTrend {
			date := e.Date
			if _, exists := dateEarnings[date]; !exists {
				dateEarnings[date] = &domain.DailyEarnings{
					Date:       date,
					ByPlatform: make(map[domain.Platform]float64),
				}
			}
			de := dateEarnings[date]
			de.TotalEarnings += e.CommissionEarned
			de.OrdersCount += e.OrdersCount
			de.Clicks += e.Clicks
			de.ByPlatform[e.Platform] += e.CommissionEarned
		}

		// Convert to slice
		for _, de := range dateEarnings {
			if de.Clicks > 0 {
				de.ConversionRate = float64(de.OrdersCount) / float64(de.Clicks) * 100
			}
			dashboard.EarningsTrend = append(dashboard.EarningsTrend, *de)
		}
	}

	return dashboard, nil
}

// GetPlatformStats retrieves stats for a specific platform
func (s *ExternalAffiliateService) GetPlatformStats(ctx context.Context, userID string, platform domain.Platform) (*domain.PlatformStats, error) {
	accounts, err := s.accountRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &domain.PlatformStats{
		Platform: platform,
	}

	for _, account := range accounts {
		if account.Platform == platform {
			stats.AccountCount++
			stats.TotalEarnings += account.TotalEarnings
			stats.TotalOrders += account.TotalOrders
			stats.TotalClicks += account.TotalClicks
			stats.PendingCommission += account.PendingCommission

			if stats.LastSyncAt == nil || (account.LastSyncAt != nil && account.LastSyncAt.After(*stats.LastSyncAt)) {
				stats.LastSyncAt = account.LastSyncAt
			}
		}
	}

	if stats.TotalClicks > 0 {
		stats.AvgConversionRate = float64(stats.TotalOrders) / float64(stats.TotalClicks) * 100
	}

	return stats, nil
}

// GetProducts retrieves products for an account
func (s *ExternalAffiliateService) GetProducts(ctx context.Context, accountID string, limit, offset int) ([]*domain.ExternalAffiliateProduct, error) {
	return s.productRepo.GetAllByAccountID(ctx, accountID, limit, offset)
}

// GetOrders retrieves orders for an account
func (s *ExternalAffiliateService) GetOrders(ctx context.Context, accountID string, limit, offset int) ([]*domain.ExternalAffiliateOrder, error) {
	return s.orderRepo.GetAllByAccountID(ctx, accountID, limit, offset)
}

// GetSyncLogs retrieves sync logs for an account
func (s *ExternalAffiliateService) GetSyncLogs(ctx context.Context, accountID string, limit int) ([]*domain.ExternalAffiliateSyncLog, error) {
	return s.syncLogRepo.GetRecentByAccountID(ctx, accountID, limit)
}

// GenerateAffiliateLink generates an affiliate link for a product
func (s *ExternalAffiliateService) GenerateAffiliateLink(ctx context.Context, accountID, productURL string) (string, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return "", err
	}

	integration, ok := s.integrations[account.Platform]
	if !ok {
		return "", &PlatformAPIError{
			Platform: account.Platform,
			Message:  "platform integration not available",
		}
	}

	return integration.GenerateAffiliateLink(ctx, account, productURL)
}

// GetOAuthURL gets the OAuth authorization URL for a platform
func (s *ExternalAffiliateService) GetOAuthURL(ctx context.Context, platform domain.Platform, state string) (string, error) {
	integration, ok := s.integrations[platform]
	if !ok {
		return "", &PlatformAPIError{
			Platform: platform,
			Message:  "platform not supported",
		}
	}

	// Check if the integration has OAuth support
	switch p := integration.(type) {
	case *TikTokIntegration:
		return p.GetOAuthURL(state), nil
	default:
		return "", &PlatformAPIError{
			Platform: platform,
			Message:  "OAuth not supported for this platform",
		}
	}
}

// HandleOAuthCallback handles the OAuth callback from a platform
func (s *ExternalAffiliateService) HandleOAuthCallback(ctx context.Context, userID, organizationID string, platform domain.Platform, code string) (*domain.ExternalAffiliateAccount, error) {
	return s.ConnectPlatform(ctx, userID, organizationID, &domain.PlatformConnectionRequest{
		Platform:          platform,
		AuthorizationCode: code,
	})
}