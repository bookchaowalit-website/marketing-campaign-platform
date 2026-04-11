package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// ShopeeIntegration implements PlatformIntegration for Shopee
type ShopeeIntegration struct {
	baseURL      string
	partnerID    string
	partnerKey   string
	redirectURL  string

	accountRepo   domain.ExternalAffiliateAccountRepository
	productRepo   domain.ExternalAffiliateProductRepository
	orderRepo     domain.ExternalAffiliateOrderRepository
	earningsRepo  domain.ExternalAffiliateEarningsRepository
	syncLogRepo   domain.ExternalAffiliateSyncLogRepository
}

// NewShopeeIntegration creates a new Shopee integration
func NewShopeeIntegration(
	partnerID, partnerKey, redirectURL string,
	accountRepo domain.ExternalAffiliateAccountRepository,
	productRepo domain.ExternalAffiliateProductRepository,
	orderRepo domain.ExternalAffiliateOrderRepository,
	earningsRepo domain.ExternalAffiliateEarningsRepository,
	syncLogRepo domain.ExternalAffiliateSyncLogRepository,
) *ShopeeIntegration {
	return &ShopeeIntegration{
		baseURL:      "https://partner.shopeemobile.com",
		partnerID:    partnerID,
		partnerKey:   partnerKey,
		redirectURL:  redirectURL,
		accountRepo:  accountRepo,
		productRepo:  productRepo,
		orderRepo:    orderRepo,
		earningsRepo: earningsRepo,
		syncLogRepo:  syncLogRepo,
	}
}

// GetPlatform returns the platform type
func (s *ShopeeIntegration) GetPlatform() domain.Platform {
	return domain.PlatformShopee
}

// Connect connects a Shopee affiliate account
func (s *ShopeeIntegration) Connect(ctx context.Context, req *domain.PlatformConnectionRequest) (*domain.ExternalAffiliateAccount, error) {
	// For Shopee, we can connect via API keys or OAuth
	account := &domain.ExternalAffiliateAccount{
		ID:                uuid.New().String(),
		Platform:          domain.PlatformShopee,
		PlatformAccountID: req.PlatformAccountID,
		PlatformUsername:  req.PlatformUsername,
		APIKey:            req.APIKey,
		APISecret:         req.APISecret,
		AccessToken:       req.AccessToken,
		RefreshToken:      req.RefreshToken,
		Status:            "pending",
		AutoSync:          true,
		SyncIntervalMins:  60,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Validate credentials by making a test API call
	if err := s.ValidateCredentials(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to validate Shopee credentials: %w", err)
	}

	account.Status = "active"
	return account, nil
}

// Disconnect disconnects a Shopee account
func (s *ShopeeIntegration) Disconnect(ctx context.Context, accountID string) error {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	account.Status = "disconnected"
	account.AccessToken = ""
	account.RefreshToken = ""
	account.UpdatedAt = time.Now()

	return s.accountRepo.Update(ctx, account)
}

// ValidateCredentials validates Shopee API credentials
func (s *ShopeeIntegration) ValidateCredentials(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	// Make a test API call to validate credentials
	// In production, this would call Shopee's GetShopInfo or similar endpoint
	_, err := s.makeAPICall(ctx, account, "/api/v2/shop/get_shop_info", nil, "GET")
	return err
}

// RefreshTokens refreshes Shopee access tokens
func (s *ShopeeIntegration) RefreshTokens(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	if account.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	path := "/api/v2/auth/access_token/get"
	params := map[string]string{
		"refresh_token": account.RefreshToken,
	}

	resp, err := s.makeAPICall(ctx, account, path, params, "POST")
	if err != nil {
		return err
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expire_in"`
	}

	if err := json.Unmarshal(resp, &tokenResp); err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return s.accountRepo.UpdateTokens(ctx, account.ID, tokenResp.AccessToken, tokenResp.RefreshToken, &expiresAt)
}

// SyncProducts syncs products from Shopee
func (s *ShopeeIntegration) SyncProducts(ctx context.Context, account *domain.ExternalAffiliateAccount, limit, offset int) ([]*domain.ExternalAffiliateProduct, error) {
	path := "/api/v2/affiliate/product/get_product_list"
	params := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
	}

	resp, err := s.makeAPICall(ctx, account, path, params, "GET")
	if err != nil {
		return nil, err
	}

	var productResp struct {
		Products []ShopeeProduct `json:"products"`
		Total    int             `json:"total"`
	}

	if err := json.Unmarshal(resp, &productResp); err != nil {
		return nil, err
	}

	products := make([]*domain.ExternalAffiliateProduct, 0, len(productResp.Products))
	for _, sp := range productResp.Products {
		product := &domain.ExternalAffiliateProduct{
			ID:               uuid.New().String(),
			AccountID:        account.ID,
			Platform:         domain.PlatformShopee,
			ProductID:        strconv.FormatInt(sp.ItemID, 10),
			ProductName:      sp.Name,
			ProductURL:       sp.ProductURL,
			ProductImageURL:  sp.ImageURL,
			Category:         sp.Category,
			OriginalPrice:    float64(sp.OriginalPrice) / 100,
			SalePrice:        float64(sp.SalePrice) / 100,
			Currency:         "THB", // Adjust based on region
			CommissionRate:   float64(sp.CommissionRate) / 10000,
			CommissionType:   "percentage",
			ProductStatus:    sp.Status,
			AffiliateLink:    sp.AffiliateLink,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := s.productRepo.Upsert(ctx, product); err != nil {
			continue
		}
		products = append(products, product)
	}

	return products, nil
}

// SyncOrders syncs orders from Shopee
func (s *ShopeeIntegration) SyncOrders(ctx context.Context, account *domain.ExternalAffiliateAccount, startDate, endDate time.Time) ([]*domain.ExternalAffiliateOrder, error) {
	path := "/api/v2/affiliate/order/get_order_list"
	params := map[string]string{
		"start_time": strconv.FormatInt(startDate.Unix(), 10),
		"end_time":   strconv.FormatInt(endDate.Unix(), 10),
	}

	resp, err := s.makeAPICall(ctx, account, path, params, "GET")
	if err != nil {
		return nil, err
	}

	var orderResp struct {
		Orders []ShopeeOrder `json:"orders"`
		Total  int           `json:"total"`
	}

	if err := json.Unmarshal(resp, &orderResp); err != nil {
		return nil, err
	}

	orders := make([]*domain.ExternalAffiliateOrder, 0, len(orderResp.Orders))
	for _, so := range orderResp.Orders {
		order := &domain.ExternalAffiliateOrder{
			ID:               uuid.New().String(),
			AccountID:        account.ID,
			Platform:         domain.PlatformShopee,
			OrderID:          strconv.FormatInt(so.OrderID, 10),
			OrderStatus:      so.OrderStatus,
			OrderDate:        so.OrderTime,
			OrderItems:       so.ItemCount,
			OrderAmount:      float64(so.TotalPrice) / 100,
			Currency:         "THB",
			CommissionAmount: float64(so.Commission) / 100,
			CommissionStatus: so.CommissionStatus,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			SyncedAt:         time.Now(),
		}

		if err := s.orderRepo.Upsert(ctx, order); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// SyncEarnings syncs earnings from Shopee
func (s *ShopeeIntegration) SyncEarnings(ctx context.Context, account *domain.ExternalAffiliateAccount, startDate, endDate time.Time) ([]*domain.ExternalAffiliateEarnings, error) {
	path := "/api/v2/affiliate/report/get_earnings"
	params := map[string]string{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}

	resp, err := s.makeAPICall(ctx, account, path, params, "GET")
	if err != nil {
		return nil, err
	}

	var earningsResp struct {
		DailyEarnings []ShopeeDailyEarnings `json:"daily_earnings"`
	}

	if err := json.Unmarshal(resp, &earningsResp); err != nil {
		return nil, err
	}

	earnings := make([]*domain.ExternalAffiliateEarnings, 0, len(earningsResp.DailyEarnings))
	for _, se := range earningsResp.DailyEarnings {
		date, _ := time.Parse("2006-01-02", se.Date)
		earning := &domain.ExternalAffiliateEarnings{
			ID:                uuid.New().String(),
			AccountID:         account.ID,
			Platform:          domain.PlatformShopee,
			Date:              date,
			GrossSales:        float64(se.GrossSales) / 100,
			NetSales:          float64(se.NetSales) / 100,
			OrdersCount:       se.Orders,
			CommissionEarned:  float64(se.Commission) / 100,
			CommissionPending: float64(se.PendingCommission) / 100,
			Clicks:            int64(se.Clicks),
			ConversionRate:    float64(se.ConversionRate) / 100,
			Currency:          "THB",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := s.earningsRepo.Upsert(ctx, earning); err != nil {
			continue
		}
		earnings = append(earnings, earning)
	}

	return earnings, nil
}

// FullSync performs a full sync for a Shopee account
func (s *ShopeeIntegration) FullSync(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	now := time.Now()
	syncLog := &domain.ExternalAffiliateSyncLog{
		ID:        uuid.New().String(),
		AccountID: account.ID,
		Platform:  domain.PlatformShopee,
		SyncType:  "full",
		Status:    "started",
		StartedAt: now,
	}

	if err := s.syncLogRepo.Create(ctx, syncLog); err != nil {
		return err
	}

	var recordsCreated, recordsUpdated, recordsFailed int

	// Sync products
	products, err := s.SyncProducts(ctx, account, 100, 0)
	if err != nil {
		recordsFailed++
	} else {
		recordsCreated += len(products)
	}

	// Sync orders (last 30 days)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	orders, err := s.SyncOrders(ctx, account, startDate, endDate)
	if err != nil {
		recordsFailed++
	} else {
		recordsCreated += len(orders)
	}

	// Sync earnings
	earnings, err := s.SyncEarnings(ctx, account, startDate, endDate)
	if err != nil {
		recordsFailed++
	} else {
		recordsCreated += len(earnings)
	}

	// Update sync log
	syncStatus := "success"
	if recordsFailed > 0 {
		syncStatus = "partial"
	}
	s.syncLogRepo.UpdateStatus(ctx, syncLog.ID, syncStatus, time.Now(), recordsCreated, recordsUpdated, recordsFailed, "")

	// Update account stats
	stats, _ := s.GetAccountStats(ctx, account)
	if stats != nil {
		s.accountRepo.UpdateStats(ctx, account.ID, stats.TotalEarnings, stats.TotalOrders, stats.TotalProducts, stats.TotalClicks, stats.PendingCommission)
	}

	// Update sync status
	s.accountRepo.UpdateSyncStatus(ctx, account.ID, syncStatus, "")

	return nil
}

// GenerateAffiliateLink generates an affiliate link for a product
func (s *ShopeeIntegration) GenerateAffiliateLink(ctx context.Context, account *domain.ExternalAffiliateAccount, productURL string) (string, error) {
	path := "/api/v2/affiliate/product/generate_link"
	params := map[string]string{
		"product_url": productURL,
	}

	resp, err := s.makeAPICall(ctx, account, path, params, "POST")
	if err != nil {
		return "", err
	}

	var linkResp struct {
		AffiliateLink string `json:"affiliate_link"`
	}

	if err := json.Unmarshal(resp, &linkResp); err != nil {
		return "", err
	}

	return linkResp.AffiliateLink, nil
}

// GenerateShortLink generates a short affiliate link
func (s *ShopeeIntegration) GenerateShortLink(ctx context.Context, account *domain.ExternalAffiliateAccount, productURL, customName string) (string, error) {
	// First generate the affiliate link
	affiliateLink, err := s.GenerateAffiliateLink(ctx, account, productURL)
	if err != nil {
		return "", err
	}

	// In production, you would use a URL shortener service
	// For now, we'll just return the affiliate link
	return affiliateLink, nil
}

// GetAccountStats gets account statistics from Shopee
func (s *ShopeeIntegration) GetAccountStats(ctx context.Context, account *domain.ExternalAffiliateAccount) (*PlatformAccountStats, error) {
	path := "/api/v2/affiliate/report/get_overview"
	params := map[string]string{
		"period": "30d",
	}

	resp, err := s.makeAPICall(ctx, account, path, params, "GET")
	if err != nil {
		return nil, err
	}

	var statsResp struct {
		TotalEarnings     float64 `json:"total_earnings"`
		PendingCommission float64 `json:"pending_commission"`
		TotalOrders       int     `json:"total_orders"`
		TotalClicks       int64   `json:"total_clicks"`
		TotalProducts     int     `json:"total_products"`
	}

	if err := json.Unmarshal(resp, &statsResp); err != nil {
		return nil, err
	}

	return &PlatformAccountStats{
		TotalEarnings:     statsResp.TotalEarnings,
		PendingCommission: statsResp.PendingCommission,
		TotalOrders:       statsResp.TotalOrders,
		TotalClicks:       statsResp.TotalClicks,
		TotalProducts:     statsResp.TotalProducts,
		LastUpdated:       time.Now(),
	}, nil
}

// Helper function to make API calls
func (s *ShopeeIntegration) makeAPICall(ctx context.Context, account *domain.ExternalAffiliateAccount, path string, params map[string]string, method string) ([]byte, error) {
	// Build URL with signature
	timestamp := time.Now().Unix()
	baseURL := s.baseURL + path

	// Add common parameters
	if params == nil {
		params = make(map[string]string)
	}
	params["partner_id"] = s.partnerID
	params["timestamp"] = strconv.FormatInt(timestamp, 10)
	if account.AccessToken != "" {
		params["access_token"] = account.AccessToken
	}
	if account.PlatformAccountID != "" {
		params["shop_id"] = account.PlatformAccountID
	}

	// Generate signature
	sign := s.generateSignature(path, timestamp, account.AccessToken, account.PlatformAccountID)

	// Build query string
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	query.Set("sign", sign)

	fullURL := baseURL + "?" + query.Encode()

	// Make HTTP request
	var resp *http.Response
	var err error

	if method == "POST" {
		resp, err = http.Post(fullURL, "application/json", nil)
	} else {
		resp, err = http.Get(fullURL)
	}

	if err != nil {
		return nil, &PlatformAPIError{
			Platform:   domain.PlatformShopee,
			Message:    err.Error(),
			StatusCode: 0,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if resp.StatusCode >= 400 {
		return nil, &PlatformAPIError{
			Platform:    domain.PlatformShopee,
			Message:     string(body),
			StatusCode:  resp.StatusCode,
			RateLimited: resp.StatusCode == 429,
		}
	}

	return body, nil
}

// generateSignature generates a Shopee API signature
func (s *ShopeeIntegration) generateSignature(path string, timestamp int64, accessToken, shopID string) string {
	baseString := fmt.Sprintf("%s%s%d", s.partnerID, path, timestamp)
	if accessToken != "" {
		baseString += accessToken
	}
	if shopID != "" {
		baseString += shopID
	}

	h := hmac.New(sha256.New, []byte(s.partnerKey))
	h.Write([]byte(baseString))
	return hex.EncodeToString(h.Sum(nil))
}

// Shopee API response structures
type ShopeeProduct struct {
	ItemID         int64   `json:"item_id"`
	Name           string  `json:"name"`
	ImageURL       string  `json:"image"`
	ProductURL     string  `json:"product_link"`
	Category       string  `json:"category_name"`
	OriginalPrice  int64   `json:"price_before_discount"`
	SalePrice      int64   `json:"price"`
	CommissionRate int     `json:"commission_rate"`
	Status         string  `json:"status"`
	AffiliateLink  string  `json:"affiliate_link"`
}

type ShopeeOrder struct {
	OrderID          int64      `json:"order_id"`
	OrderStatus      string     `json:"order_status"`
	OrderTime        *time.Time `json:"order_time"`
	ItemCount        int        `json:"item_count"`
	TotalPrice       int64      `json:"total_price"`
	Commission       int64      `json:"commission"`
	CommissionStatus string     `json:"commission_status"`
}

type ShopeeDailyEarnings struct {
	Date              string `json:"date"`
	GrossSales        int64  `json:"gross_sales"`
	NetSales          int64  `json:"net_sales"`
	Orders            int    `json:"orders"`
	Commission        int64  `json:"commission"`
	PendingCommission int64  `json:"pending_commission"`
	Clicks            int    `json:"clicks"`
	ConversionRate    int    `json:"conversion_rate"`
}