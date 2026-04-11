package dto

// ConnectPlatformRequest represents a request to connect a platform
type ConnectPlatformRequest struct {
	Platform          string `json:"platform" binding:"required"` // shopee, tiktok, lazada, amazon
	APIKey            string `json:"api_key,omitempty"`
	APISecret         string `json:"api_secret,omitempty"`
	AccessToken       string `json:"access_token,omitempty"`
	RefreshToken      string `json:"refresh_token,omitempty"`
	PlatformAccountID string `json:"platform_account_id,omitempty"`
	PlatformUsername  string `json:"platform_username,omitempty"`
	AuthorizationCode string `json:"authorization_code,omitempty"` // For OAuth flows
}

// SyncAccountRequest represents a request to sync an account
type SyncAccountRequest struct {
	SyncType string `json:"sync_type,omitempty"` // full, incremental, products, orders
}

// GenerateLinkRequest represents a request to generate an affiliate link
type GenerateLinkRequest struct {
	ProductURL string `json:"product_url" binding:"required"`
	CustomName string `json:"custom_name,omitempty"`
}

// ExternalAccountResponse represents a connected external account
type ExternalAccountResponse struct {
	ID                  string  `json:"id"`
	Platform            string  `json:"platform"`
	PlatformAccountID   string  `json:"platform_account_id"`
	PlatformAccountName string  `json:"platform_account_name"`
	PlatformUsername    string  `json:"platform_username"`
	Status              string  `json:"status"`
	AutoSync            bool    `json:"auto_sync"`
	SyncIntervalMins    int     `json:"sync_interval_minutes"`
	LastSyncAt          string  `json:"last_sync_at,omitempty"`
	LastSyncStatus      string  `json:"last_sync_status"`
	LastSyncError       string  `json:"last_sync_error,omitempty"`
	TotalEarnings       float64 `json:"total_earnings"`
	TotalOrders         int     `json:"total_orders"`
	TotalClicks         int64   `json:"total_clicks"`
	TotalProducts       int     `json:"total_products"`
	PendingCommission   float64 `json:"pending_commission"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}

// ExternalProductResponse represents an external affiliate product
type ExternalProductResponse struct {
	ID               string  `json:"id"`
	Platform         string  `json:"platform"`
	ProductID        string  `json:"product_id"`
	ProductName      string  `json:"product_name"`
	ProductURL       string  `json:"product_url"`
	ProductImageURL  string  `json:"product_image_url"`
	Category         string  `json:"category"`
	Brand            string  `json:"brand"`
	OriginalPrice    float64 `json:"original_price"`
	SalePrice        float64 `json:"sale_price"`
	Currency         string  `json:"currency"`
	CommissionRate   float64 `json:"commission_rate"`
	CommissionType   string  `json:"commission_type"`
	EstimatedEarnings float64 `json:"estimated_earnings"`
	TotalClicks      int64   `json:"total_clicks"`
	TotalOrders      int     `json:"total_orders"`
	TotalEarnings    float64 `json:"total_earnings"`
	ProductStatus    string  `json:"product_status"`
	AffiliateLink    string  `json:"affiliate_link"`
}

// ExternalOrderResponse represents an external affiliate order
type ExternalOrderResponse struct {
	ID               string  `json:"id"`
	Platform         string  `json:"platform"`
	OrderID          string  `json:"order_id"`
	OrderStatus      string  `json:"order_status"`
	OrderDate        string  `json:"order_date,omitempty"`
	OrderItems       int     `json:"order_items"`
	OrderAmount      float64 `json:"order_amount"`
	ProductAmount    float64 `json:"product_amount"`
	ShippingAmount   float64 `json:"shipping_amount"`
	DiscountAmount   float64 `json:"discount_amount"`
	Currency         string  `json:"currency"`
	CommissionRate   float64 `json:"commission_rate"`
	CommissionAmount float64 `json:"commission_amount"`
	CommissionStatus string  `json:"commission_status"`
	TrafficSource    string  `json:"traffic_source,omitempty"`
}

// PlatformStatsResponse represents stats for a platform
type PlatformStatsResponse struct {
	Platform          string  `json:"platform"`
	AccountCount      int     `json:"account_count"`
	TotalEarnings     float64 `json:"total_earnings"`
	TotalOrders       int     `json:"total_orders"`
	TotalClicks       int64   `json:"total_clicks"`
	PendingCommission float64 `json:"pending_commission"`
	AvgConversionRate float64 `json:"avg_conversion_rate"`
	LastSyncAt        string  `json:"last_sync_at,omitempty"`
}

// DailyEarningsResponse represents daily earnings
type DailyEarningsResponse struct {
	Date           string  `json:"date"`
	TotalEarnings  float64 `json:"total_earnings"`
	OrdersCount    int     `json:"orders_count"`
	Clicks         int64   `json:"clicks"`
	ConversionRate float64 `json:"conversion_rate"`
}

// SyncLogResponse represents a sync log entry
type SyncLogResponse struct {
	ID             string `json:"id"`
	SyncType       string `json:"sync_type"`
	Status         string `json:"status"`
	StartedAt      string `json:"started_at"`
	CompletedAt    string `json:"completed_at,omitempty"`
	RecordsCreated int    `json:"records_created"`
	RecordsUpdated int    `json:"records_updated"`
	RecordsFailed  int    `json:"records_failed"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

// UnifiedDashboardResponse represents the unified dashboard
type UnifiedDashboardResponse struct {
	TotalEarnings         float64                   `json:"total_earnings"`
	TotalPendingCommission float64                   `json:"total_pending_commission"`
	TotalOrders           int                       `json:"total_orders"`
	TotalClicks           int64                     `json:"total_clicks"`
	TotalProducts         int                       `json:"total_products"`
	PlatformStats         []PlatformStatsResponse   `json:"platform_stats"`
	RecentOrders          []ExternalOrderResponse   `json:"recent_orders"`
	TopProducts           []ExternalProductResponse `json:"top_products"`
	EarningsTrend         []DailyEarningsResponse   `json:"earnings_trend"`
	ConnectedAccounts     []ExternalAccountResponse `json:"connected_accounts"`
}