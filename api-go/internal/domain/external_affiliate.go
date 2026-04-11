package domain

import "time"

// Platform types
type Platform string

const (
	PlatformShopee Platform = "shopee"
	PlatformTikTok Platform = "tiktok"
	PlatformLazada Platform = "lazada"
	PlatformAmazon Platform = "amazon"
	PlatformCustom Platform = "custom"
)

// ExternalAffiliateAccount represents a connected external affiliate account
type ExternalAffiliateAccount struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	OrganizationID      string     `json:"organization_id"`
	Platform            Platform   `json:"platform"`
	PlatformAccountID   string     `json:"platform_account_id"`
	PlatformAccountName string     `json:"platform_account_name"`
	PlatformUsername    string     `json:"platform_username"`
	Status              string     `json:"status"` // active, pending, disconnected, expired, error

	// API Credentials
	APIKey          string     `json:"-"` // Never expose in JSON
	APISecret       string     `json:"-"`
	AccessToken     string     `json:"-"`
	RefreshToken    string     `json:"-"`
	TokenExpiresAt  *time.Time `json:"token_expires_at,omitempty"`

	// Sync Settings
	AutoSync         bool       `json:"auto_sync"`
	SyncIntervalMins int        `json:"sync_interval_minutes"`
	LastSyncAt       *time.Time `json:"last_sync_at,omitempty"`
	LastSyncStatus   string     `json:"last_sync_status"`
	LastSyncError    string     `json:"last_sync_error,omitempty"`

	// Stats
	TotalEarnings       float64 `json:"total_earnings"`
	TotalOrders         int     `json:"total_orders"`
	TotalClicks         int64   `json:"total_clicks"`
	TotalProducts       int     `json:"total_products"`
	PendingCommission   float64 `json:"pending_commission"`

	// Metadata
	ConnectionData string    `json:"-"` // JSON field for additional data
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ExternalAffiliateProduct represents a product available for promotion
type ExternalAffiliateProduct struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"account_id"`
	Platform       Platform  `json:"platform"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ProductURL     string    `json:"product_url"`
	ProductImageURL string   `json:"product_image_url"`
	Category       string    `json:"category"`
	Brand          string    `json:"brand"`

	// Pricing
	OriginalPrice float64 `json:"original_price"`
	SalePrice     float64 `json:"sale_price"`
	Currency      string  `json:"currency"`

	// Commission
	CommissionRate     float64 `json:"commission_rate"`
	CommissionType     string  `json:"commission_type"`
	EstimatedEarnings  float64 `json:"estimated_earnings"`

	// Stats
	TotalClicks   int64   `json:"total_clicks"`
	TotalOrders   int     `json:"total_orders"`
	TotalEarnings float64 `json:"total_earnings"`

	// Status
	ProductStatus string `json:"product_status"`

	// Tracking
	AffiliateLink string `json:"affiliate_link"`
	ShortLink     string `json:"short_link"`

	// Metadata
	ProductData   string     `json:"-"` // JSON
	LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ExternalAffiliateOrder represents an order tracked through an affiliate link
type ExternalAffiliateOrder struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	ProductID     string    `json:"product_id,omitempty"`

	// Order Identification
	Platform       Platform `json:"platform"`
	OrderID        string   `json:"order_id"`
	ParentOrderID  string   `json:"parent_order_id,omitempty"`

	// Order Details
	OrderStatus string     `json:"order_status"`
	OrderDate   *time.Time `json:"order_date,omitempty"`
	OrderItems  int        `json:"order_items"`

	// Amounts
	OrderAmount    float64 `json:"order_amount"`
	ProductAmount  float64 `json:"product_amount"`
	ShippingAmount float64 `json:"shipping_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	Currency       string  `json:"currency"`

	// Commission
	CommissionRate   float64 `json:"commission_rate"`
	CommissionAmount float64 `json:"commission_amount"`
	CommissionStatus string  `json:"commission_status"`
	CommissionPaidAt *time.Time `json:"commission_paid_at,omitempty"`

	// Customer
	CustomerID string `json:"customer_id,omitempty"`

	// Attribution
	TrafficSource string `json:"traffic_source,omitempty"`
	UTMSource     string `json:"utm_source,omitempty"`
	UTMMedium     string `json:"utm_medium,omitempty"`
	UTMCampaign   string `json:"utm_campaign,omitempty"`

	// Metadata
	OrderData string    `json:"-"` // JSON
	SyncedAt  time.Time `json:"synced_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExternalAffiliateEarnings represents daily earnings summary
type ExternalAffiliateEarnings struct {
	ID        string   `json:"id"`
	AccountID string   `json:"account_id"`
	Platform  Platform `json:"platform"`
	Date      time.Time `json:"date"`

	// Sales
	GrossSales float64 `json:"gross_sales"`
	NetSales   float64 `json:"net_sales"`
	OrdersCount int    `json:"orders_count"`
	ItemsSold   int    `json:"items_sold"`

	// Commission
	CommissionEarned   float64 `json:"commission_earned"`
	CommissionPending  float64 `json:"commission_pending"`
	CommissionConfirmed float64 `json:"commission_confirmed"`

	// Traffic
	Clicks      int64 `json:"clicks"`
	Impressions int64 `json:"impressions"`

	// Performance
	ConversionRate float64 `json:"conversion_rate"`
	EPC            float64 `json:"epc"` // Earnings Per Click

	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExternalAffiliateSyncLog represents a sync operation log
type ExternalAffiliateSyncLog struct {
	ID        string   `json:"id"`
	AccountID string   `json:"account_id"`
	Platform  Platform `json:"platform"`

	SyncType string `json:"sync_type"` // full, incremental, products, orders, earnings
	Status   string `json:"status"`    // started, success, partial, failed

	// Timing
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	DurationMs  int64      `json:"duration_ms,omitempty"`

	// Results
	RecordsProcessed int `json:"records_processed"`
	RecordsCreated   int `json:"records_created"`
	RecordsUpdated   int `json:"records_updated"`
	RecordsFailed    int `json:"records_failed"`

	// Error
	ErrorMessage string `json:"error_message,omitempty"`
	ErrorDetails string `json:"error_details,omitempty"` // JSON

	// API Info
	APICallsCount int  `json:"api_calls_count"`
	RateLimitHit  bool `json:"rate_limit_hit"`

	CreatedAt time.Time `json:"created_at"`
}

// ExternalAffiliateLink represents a custom tracking link
type ExternalAffiliateLink struct {
	ID        string   `json:"id"`
	AccountID string   `json:"account_id"`
	ProductID string   `json:"product_id,omitempty"`

	LinkCode    string `json:"link_code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// URLs
	OriginalURL   string `json:"original_url"`
	AffiliateURL  string `json:"affiliate_url"`
	ShortURL      string `json:"short_url"`

	// Tracking Parameters
	UTMSource   string `json:"utm_source,omitempty"`
	UTMMedium   string `json:"utm_medium,omitempty"`
	UTMCampaign string `json:"utm_campaign,omitempty"`
	UTMTerm     string `json:"utm_term,omitempty"`
	UTMContent  string `json:"utm_content,omitempty"`
	SubID       string `json:"sub_id,omitempty"`

	// Stats
	TotalClicks   int64   `json:"total_clicks"`
	UniqueClicks  int64   `json:"unique_clicks"`
	TotalOrders   int     `json:"total_orders"`
	TotalEarnings float64 `json:"total_earnings"`

	// Settings
	IsActive  bool       `json:"is_active"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UnifiedAffiliateDashboard represents combined data from all platforms
type UnifiedAffiliateDashboard struct {
	// Overall Stats
	TotalEarnings       float64 `json:"total_earnings"`
	TotalPendingCommission float64 `json:"total_pending_commission"`
	TotalOrders         int     `json:"total_orders"`
	TotalClicks         int64   `json:"total_clicks"`
	TotalProducts       int     `json:"total_products"`

	// By Platform
	PlatformStats []PlatformStats `json:"platform_stats"`

	// Recent Orders
	RecentOrders []ExternalAffiliateOrder `json:"recent_orders"`

	// Top Products
	TopProducts []ExternalAffiliateProduct `json:"top_products"`

	// Earnings Trend
	EarningsTrend []DailyEarnings `json:"earnings_trend"`

	// Connected Accounts
	ConnectedAccounts []ExternalAffiliateAccount `json:"connected_accounts"`
}

// PlatformStats represents stats for a specific platform
type PlatformStats struct {
	Platform           Platform `json:"platform"`
	AccountCount       int      `json:"account_count"`
	TotalEarnings      float64  `json:"total_earnings"`
	TotalOrders        int      `json:"total_orders"`
	TotalClicks        int64    `json:"total_clicks"`
	PendingCommission  float64  `json:"pending_commission"`
	AvgConversionRate  float64  `json:"avg_conversion_rate"`
	LastSyncAt         *time.Time `json:"last_sync_at,omitempty"`
}

// DailyEarnings represents earnings for a specific day
type DailyEarnings struct {
	Date              time.Time `json:"date"`
	TotalEarnings     float64   `json:"total_earnings"`
	OrdersCount       int       `json:"orders_count"`
	Clicks            int64     `json:"clicks"`
	ConversionRate    float64   `json:"conversion_rate"`
	ByPlatform        map[Platform]float64 `json:"by_platform"`
}

// PlatformConnectionRequest represents a request to connect a platform
type PlatformConnectionRequest struct {
	Platform          Platform `json:"platform"`
	APIKey            string   `json:"api_key,omitempty"`
	APISecret         string   `json:"api_secret,omitempty"`
	AccessToken       string   `json:"access_token,omitempty"`
	RefreshToken      string   `json:"refresh_token,omitempty"`
	PlatformAccountID string   `json:"platform_account_id,omitempty"`
	PlatformUsername  string   `json:"platform_username,omitempty"`
	AuthorizationCode string   `json:"authorization_code,omitempty"` // For OAuth flows
}

// PlatformSyncRequest represents a request to sync data from a platform
type PlatformSyncRequest struct {
	AccountID string `json:"account_id"`
	SyncType  string `json:"sync_type"` // full, incremental, products, orders, earnings
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}