package service

import (
	"context"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// PlatformIntegration defines the interface for external affiliate platform integrations
type PlatformIntegration interface {
	// Platform identification
	GetPlatform() domain.Platform

	// Connection
	Connect(ctx context.Context, req *domain.PlatformConnectionRequest) (*domain.ExternalAffiliateAccount, error)
	Disconnect(ctx context.Context, accountID string) error
	ValidateCredentials(ctx context.Context, account *domain.ExternalAffiliateAccount) error
	RefreshTokens(ctx context.Context, account *domain.ExternalAffiliateAccount) error

	// Sync operations
	SyncProducts(ctx context.Context, account *domain.ExternalAffiliateAccount, limit, offset int) ([]*domain.ExternalAffiliateProduct, error)
	SyncOrders(ctx context.Context, account *domain.ExternalAffiliateAccount, startDate, endDate time.Time) ([]*domain.ExternalAffiliateOrder, error)
	SyncEarnings(ctx context.Context, account *domain.ExternalAffiliateAccount, startDate, endDate time.Time) ([]*domain.ExternalAffiliateEarnings, error)
	FullSync(ctx context.Context, account *domain.ExternalAffiliateAccount) error

	// Generate affiliate links
	GenerateAffiliateLink(ctx context.Context, account *domain.ExternalAffiliateAccount, productURL string) (string, error)
	GenerateShortLink(ctx context.Context, account *domain.ExternalAffiliateAccount, productURL, customName string) (string, error)

	// Get stats
	GetAccountStats(ctx context.Context, account *domain.ExternalAffiliateAccount) (*PlatformAccountStats, error)
}

// PlatformAccountStats represents account stats from a platform
type PlatformAccountStats struct {
	TotalEarnings     float64   `json:"total_earnings"`
	PendingCommission float64   `json:"pending_commission"`
	TotalOrders       int       `json:"total_orders"`
	TotalClicks       int64     `json:"total_clicks"`
	TotalProducts     int       `json:"total_products"`
	LastUpdated       time.Time `json:"last_updated"`
}

// PlatformAPIError represents an error from a platform API
type PlatformAPIError struct {
	Platform    domain.Platform
	Code        string
	Message     string
	StatusCode  int
	RateLimited bool
}

func (e *PlatformAPIError) Error() string {
	return e.Message
}

// PlatformProductData represents product data from a platform API
type PlatformProductData struct {
	ProductID         string  `json:"product_id"`
	ProductName       string  `json:"product_name"`
	ProductURL        string  `json:"product_url"`
	ProductImageURL   string  `json:"product_image_url"`
	Category          string  `json:"category"`
	Brand             string  `json:"brand"`
	OriginalPrice     float64 `json:"original_price"`
	SalePrice         float64 `json:"sale_price"`
	Currency          string  `json:"currency"`
	CommissionRate    float64 `json:"commission_rate"`
	CommissionType    string  `json:"commission_type"`
	EstimatedEarnings float64 `json:"estimated_earnings"`
	ProductStatus     string  `json:"product_status"`
	AffiliateLink     string  `json:"affiliate_link"`
}

// PlatformOrderData represents order data from a platform API
type PlatformOrderData struct {
	OrderID         string     `json:"order_id"`
	ParentOrderID   string     `json:"parent_order_id"`
	OrderStatus     string     `json:"order_status"`
	OrderDate       *time.Time `json:"order_date"`
	OrderItems      int        `json:"order_items"`
	OrderAmount     float64    `json:"order_amount"`
	ProductAmount   float64    `json:"product_amount"`
	ShippingAmount  float64    `json:"shipping_amount"`
	DiscountAmount  float64    `json:"discount_amount"`
	Currency        string     `json:"currency"`
	CommissionRate  float64    `json:"commission_rate"`
	CommissionAmount float64   `json:"commission_amount"`
	CommissionStatus string    `json:"commission_status"`
	CustomerID      string     `json:"customer_id"`
	ProductID       string     `json:"product_id"`
}