package service

import (
	"context"
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

// TikTokIntegration implements PlatformIntegration for TikTok Shop
type TikTokIntegration struct {
	baseURL      string
	appID        string
	appSecret    string
	redirectURL  string

	accountRepo   domain.ExternalAffiliateAccountRepository
	productRepo   domain.ExternalAffiliateProductRepository
	orderRepo     domain.ExternalAffiliateOrderRepository
	earningsRepo  domain.ExternalAffiliateEarningsRepository
	syncLogRepo   domain.ExternalAffiliateSyncLogRepository
}

// NewTikTokIntegration creates a new TikTok integration
func NewTikTokIntegration(
	appID, appSecret, redirectURL string,
	accountRepo domain.ExternalAffiliateAccountRepository,
	productRepo domain.ExternalAffiliateProductRepository,
	orderRepo domain.ExternalAffiliateOrderRepository,
	earningsRepo domain.ExternalAffiliateEarningsRepository,
	syncLogRepo domain.ExternalAffiliateSyncLogRepository,
) *TikTokIntegration {
	return &TikTokIntegration{
		baseURL:      "https://open-api.tiktokglobal.com",
		appID:        appID,
		appSecret:    appSecret,
		redirectURL:  redirectURL,
		accountRepo:  accountRepo,
		productRepo:  productRepo,
		orderRepo:    orderRepo,
		earningsRepo: earningsRepo,
		syncLogRepo:  syncLogRepo,
	}
}

// GetPlatform returns the platform type
func (t *TikTokIntegration) GetPlatform() domain.Platform {
	return domain.PlatformTikTok
}

// Connect connects a TikTok affiliate account
func (t *TikTokIntegration) Connect(ctx context.Context, req *domain.PlatformConnectionRequest) (*domain.ExternalAffiliateAccount, error) {
	// For TikTok Shop, connection typically happens via OAuth
	// The authorization code is exchanged for access tokens
	
	if req.AuthorizationCode != "" {
		// Exchange authorization code for tokens
		tokens, err := t.exchangeCodeForTokens(ctx, req.AuthorizationCode)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
		}

		req.AccessToken = tokens.AccessToken
		req.RefreshToken = tokens.RefreshToken
	}

	account := &domain.ExternalAffiliateAccount{
		ID:                uuid.New().String(),
		Platform:          domain.PlatformTikTok,
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

	// Validate credentials
	if err := t.ValidateCredentials(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to validate TikTok credentials: %w", err)
	}

	account.Status = "active"
	return account, nil
}

// Disconnect disconnects a TikTok account
func (t *TikTokIntegration) Disconnect(ctx context.Context, accountID string) error {
	account, err := t.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	account.Status = "disconnected"
	account.AccessToken = ""
	account.RefreshToken = ""
	account.UpdatedAt = time.Now()

	return t.accountRepo.Update(ctx, account)
}

// ValidateCredentials validates TikTok API credentials
func (t *TikTokIntegration) ValidateCredentials(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	// Make a test API call to validate credentials
	_, err := t.makeAPICall(ctx, account, "/api/creator/info", nil, "GET")
	return err
}

// RefreshTokens refreshes TikTok access tokens
func (t *TikTokIntegration) RefreshTokens(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	if account.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	path := "/api/token/refresh"
	params := map[string]string{
		"refresh_token": account.RefreshToken,
		"app_id":        t.appID,
	}

	resp, err := t.makeAPICall(ctx, account, path, params, "POST")
	if err != nil {
		return err
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.Unmarshal(resp, &tokenResp); err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return t.accountRepo.UpdateTokens(ctx, account.ID, tokenResp.AccessToken, tokenResp.RefreshToken, &expiresAt)
}

// SyncProducts syncs products from TikTok Shop
func (t *TikTokIntegration) SyncProducts(ctx context.Context, account *domain.ExternalAffiliateAccount, limit, offset int) ([]*domain.ExternalAffiliateProduct, error) {
	path := "/api/affiliate/products"
	params := map[string]string{
		"page_number": strconv.Itoa(offset/limit + 1),
		"page_size":   strconv.Itoa(limit),
	}

	resp, err := t.makeAPICall(ctx, account, path, params, "GET")
	if err != nil {
		return nil, err
	}

	var productResp struct {
		Products []TikTokProduct `json:"products"`
		Total    int             `json:"total"`
	}

	if err := json.Unmarshal(resp, &productResp); err != nil {
		return nil, err
	}

	products := make([]*domain.ExternalAffiliateProduct, 0, len(productResp.Products))
	for _, tp := range productResp.Products {
		product := &domain.ExternalAffiliateProduct{
			ID:              uuid.New().String(),
			AccountID:       account.ID,
			Platform:        domain.PlatformTikTok,
			ProductID:       tp.ProductID,
			ProductName:     tp.Title,
			ProductURL:      tp.ProductURL,
			ProductImageURL: tp.MainImage,
			Category:        tp.Category,
			Brand:           tp.Brand,
			OriginalPrice:   tp.OriginalPrice,
			SalePrice:       tp.SalePrice,
			Currency:        "USD",
			CommissionRate:  tp.CommissionRate,
			CommissionType:  "percentage",
			ProductStatus:   tp.Status,
			AffiliateLink:   tp.AffiliateLink,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := t.productRepo.Upsert(ctx, product); err != nil {
			continue
		}
		products = append(products, product)
	}

	return products, nil
}

// SyncOrders syncs orders from TikTok Shop
func (t *TikTokIntegration) SyncOrders(ctx context.Context, account *domain.ExternalAffiliateAccount, startDate, endDate time.Time) ([]*domain.ExternalAffiliateOrder, error) {
	path := "/api/affiliate/orders"
	params := map[string]string{
		"start_time": startDate.Format(time.RFC3339),
		"end_time":   endDate.Format(time.RFC3339),
	}

	resp, err := t.makeAPICall(ctx, account, path, params, "GET")
	if err != nil {
		return nil, err
	}

	var orderResp struct {
		Orders []TikTokOrder `json:"orders"`
		Total  int           `json:"total"`
	}

	if err := json.Unmarshal(resp, &orderResp); err != nil {
		return nil, err
	}

	orders := make([]*domain.ExternalAffiliateOrder, 0, len(orderResp.Orders))
	for _, to := range orderResp.Orders {
		order := &domain.ExternalAffiliateOrder{
			ID:               uuid.New().String(),
			AccountID:        account.ID,
			Platform:         domain.PlatformTikTok,
			OrderID:          to.OrderID,
			OrderStatus:      to.OrderStatus,
			OrderDate:        to.CreateTime,
			OrderItems:       to.ItemCount,
			OrderAmount:      to.TotalAmount,
			Currency:         "USD",
			CommissionAmount: to.EstimatedCommission,
			CommissionStatus: to.CommissionStatus,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			SyncedAt:         time.Now(),
		}

		if err := t.orderRepo.Upsert(ctx, order); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// SyncEarnings syncs earnings from TikTok Shop
func (t *TikTokIntegration) SyncEarnings(ctx context.Context, account *domain.ExternalAffiliateAccount, startDate, endDate time.Time) ([]*domain.ExternalAffiliateEarnings, error) {
	path := "/api/affiliate/earnings"
	params := map[string]string{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}

	resp, err := t.makeAPICall(ctx, account, path, params, "GET")
	if err != nil {
		return nil, err
	}

	var earningsResp struct {
		DailyEarnings []TikTokDailyEarnings `json:"daily_earnings"`
	}

	if err := json.Unmarshal(resp, &earningsResp); err != nil {
		return nil, err
	}

	earnings := make([]*domain.ExternalAffiliateEarnings, 0, len(earningsResp.DailyEarnings))
	for _, te := range earningsResp.DailyEarnings {
		date, _ := time.Parse("2006-01-02", te.Date)
		earning := &domain.ExternalAffiliateEarnings{
			ID:               uuid.New().String(),
			AccountID:        account.ID,
			Platform:         domain.PlatformTikTok,
			Date:             date,
			GrossSales:       te.GrossSales,
			NetSales:         te.NetSales,
			OrdersCount:      te.Orders,
			CommissionEarned: te.Commission,
			Clicks:           int64(te.Clicks),
			ConversionRate:   te.ConversionRate,
			Currency:         "USD",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := t.earningsRepo.Upsert(ctx, earning); err != nil {
			continue
		}
		earnings = append(earnings, earning)
	}

	return earnings, nil
}

// FullSync performs a full sync for a TikTok account
func (t *TikTokIntegration) FullSync(ctx context.Context, account *domain.ExternalAffiliateAccount) error {
	now := time.Now()
	syncLog := &domain.ExternalAffiliateSyncLog{
		ID:        uuid.New().String(),
		AccountID: account.ID,
		Platform:  domain.PlatformTikTok,
		SyncType:  "full",
		Status:    "started",
		StartedAt: now,
	}

	if err := t.syncLogRepo.Create(ctx, syncLog); err != nil {
		return err
	}

	var recordsCreated, recordsUpdated, recordsFailed int

	// Sync products
	products, err := t.SyncProducts(ctx, account, 100, 0)
	if err != nil {
		recordsFailed++
	} else {
		recordsCreated += len(products)
	}

	// Sync orders (last 30 days)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	orders, err := t.SyncOrders(ctx, account, startDate, endDate)
	if err != nil {
		recordsFailed++
	} else {
		recordsCreated += len(orders)
	}

	// Sync earnings
	earnings, err := t.SyncEarnings(ctx, account, startDate, endDate)
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
	t.syncLogRepo.UpdateStatus(ctx, syncLog.ID, syncStatus, time.Now(), recordsCreated, recordsUpdated, recordsFailed, "")

	// Update account stats
	stats, _ := t.GetAccountStats(ctx, account)
	if stats != nil {
		t.accountRepo.UpdateStats(ctx, account.ID, stats.TotalEarnings, stats.TotalOrders, stats.TotalProducts, stats.TotalClicks, stats.PendingCommission)
	}

	// Update sync status
	t.accountRepo.UpdateSyncStatus(ctx, account.ID, syncStatus, "")

	return nil
}

// GenerateAffiliateLink generates an affiliate link for a product
func (t *TikTokIntegration) GenerateAffiliateLink(ctx context.Context, account *domain.ExternalAffiliateAccount, productURL string) (string, error) {
	path := "/api/affiliate/generate_link"
	params := map[string]string{
		"product_url": productURL,
	}

	resp, err := t.makeAPICall(ctx, account, path, params, "POST")
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
func (t *TikTokIntegration) GenerateShortLink(ctx context.Context, account *domain.ExternalAffiliateAccount, productURL, customName string) (string, error) {
	// Generate the affiliate link first
	affiliateLink, err := t.GenerateAffiliateLink(ctx, account, productURL)
	if err != nil {
		return "", err
	}

	return affiliateLink, nil
}

// GetAccountStats gets account statistics from TikTok
func (t *TikTokIntegration) GetAccountStats(ctx context.Context, account *domain.ExternalAffiliateAccount) (*PlatformAccountStats, error) {
	path := "/api/affiliate/stats"
	params := map[string]string{
		"period": "30d",
	}

	resp, err := t.makeAPICall(ctx, account, path, params, "GET")
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

// GetOAuthURL generates the OAuth authorization URL
func (t *TikTokIntegration) GetOAuthURL(state string) string {
	params := url.Values{}
	params.Set("app_id", t.appID)
	params.Set("state", state)
	params.Set("redirect_uri", t.redirectURL)
	params.Set("scope", "user.info.basic,affiliate.products,affiliate.orders")

	return "https://www.tiktok.com/auth/authorize/?" + params.Encode()
}

// exchangeCodeForTokens exchanges authorization code for tokens
func (t *TikTokIntegration) exchangeCodeForTokens(ctx context.Context, code string) (*TokenResponse, error) {
	path := "/api/token/create"
	params := map[string]string{
		"code":         code,
		"app_id":       t.appID,
		"app_secret":   t.appSecret,
		"redirect_uri": t.redirectURL,
	}

	resp, err := t.makeAPICallRaw(ctx, path, params, "POST")
	if err != nil {
		return nil, err
	}

	var tokenResp struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int    `json:"expires_in"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &tokenResp); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  tokenResp.Data.AccessToken,
		RefreshToken: tokenResp.Data.RefreshToken,
		ExpiresIn:    tokenResp.Data.ExpiresIn,
	}, nil
}

// makeAPICall makes an API call with the account's access token
func (t *TikTokIntegration) makeAPICall(ctx context.Context, account *domain.ExternalAffiliateAccount, path string, params map[string]string, method string) ([]byte, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["access_token"] = account.AccessToken

	return t.makeAPICallRaw(ctx, path, params, method)
}

// makeAPICallRaw makes a raw API call
func (t *TikTokIntegration) makeAPICallRaw(ctx context.Context, path string, params map[string]string, method string) ([]byte, error) {
	fullURL := t.baseURL + path

	// Build query string
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	var resp *http.Response
	var err error

	if method == "POST" {
		resp, err = http.PostForm(fullURL, query)
	} else {
		resp, err = http.Get(fullURL + "?" + query.Encode())
	}

	if err != nil {
		return nil, &PlatformAPIError{
			Platform:   domain.PlatformTikTok,
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
			Platform:    domain.PlatformTikTok,
			Message:     string(body),
			StatusCode:  resp.StatusCode,
			RateLimited: resp.StatusCode == 429,
		}
	}

	return body, nil
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

// TikTok API response structures
type TikTokProduct struct {
	ProductID      string  `json:"product_id"`
	Title          string  `json:"title"`
	MainImage      string  `json:"main_image"`
	ProductURL     string  `json:"product_url"`
	Category       string  `json:"category_name"`
	Brand          string  `json:"brand"`
	OriginalPrice  float64 `json:"original_price"`
	SalePrice      float64 `json:"sale_price"`
	CommissionRate float64 `json:"commission_rate"`
	Status         string  `json:"status"`
	AffiliateLink  string  `json:"affiliate_link"`
}

type TikTokOrder struct {
	OrderID           string     `json:"order_id"`
	OrderStatus       string     `json:"order_status"`
	CreateTime        *time.Time `json:"create_time"`
	ItemCount         int        `json:"item_count"`
	TotalAmount       float64    `json:"total_amount"`
	EstimatedCommission float64   `json:"estimated_commission"`
	CommissionStatus  string     `json:"commission_status"`
}

type TikTokDailyEarnings struct {
	Date           string  `json:"date"`
	GrossSales     float64 `json:"gross_sales"`
	NetSales       float64 `json:"net_sales"`
	Orders         int     `json:"orders"`
	Commission     float64 `json:"commission"`
	Clicks         int     `json:"clicks"`
	ConversionRate float64 `json:"conversion_rate"`
}