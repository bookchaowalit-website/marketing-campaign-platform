package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// externalAffiliateProductRepository implements domain.ExternalAffiliateProductRepository
type externalAffiliateProductRepository struct {
	db *sql.DB
}

// NewExternalAffiliateProductRepository creates a new external affiliate product repository
func NewExternalAffiliateProductRepository(db *sql.DB) domain.ExternalAffiliateProductRepository {
	return &externalAffiliateProductRepository{db: db}
}

// Create creates a new product
func (r *externalAffiliateProductRepository) Create(ctx context.Context, product *domain.ExternalAffiliateProduct) error {
	query := `
		INSERT INTO external_affiliate_products (
			id, account_id, platform, product_id, product_name, product_url, product_image_url,
			category, brand, original_price, sale_price, currency, commission_rate, commission_type,
			estimated_earnings, product_status, affiliate_link, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.AccountID, product.Platform, product.ProductID, product.ProductName,
		product.ProductURL, product.ProductImageURL, product.Category, product.Brand,
		product.OriginalPrice, product.SalePrice, product.Currency, product.CommissionRate,
		product.CommissionType, product.EstimatedEarnings, product.ProductStatus, product.AffiliateLink,
		product.CreatedAt, product.UpdatedAt,
	)

	return err
}

// GetByID retrieves a product by ID
func (r *externalAffiliateProductRepository) GetByID(ctx context.Context, id string) (*domain.ExternalAffiliateProduct, error) {
	query := `
		SELECT id, account_id, platform, product_id, product_name, product_url, product_image_url,
			category, brand, original_price, sale_price, currency, commission_rate, commission_type,
			estimated_earnings, total_clicks, total_orders, total_earnings, product_status, affiliate_link,
			short_link, last_synced_at, created_at, updated_at
		FROM external_affiliate_products WHERE id = $1
	`

	product := &domain.ExternalAffiliateProduct{}
	var lastSyncedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.AccountID, &product.Platform, &product.ProductID, &product.ProductName,
		&product.ProductURL, &product.ProductImageURL, &product.Category, &product.Brand,
		&product.OriginalPrice, &product.SalePrice, &product.Currency, &product.CommissionRate,
		&product.CommissionType, &product.EstimatedEarnings, &product.TotalClicks, &product.TotalOrders,
		&product.TotalEarnings, &product.ProductStatus, &product.AffiliateLink, &product.ShortLink,
		&lastSyncedAt, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if lastSyncedAt.Valid {
		product.LastSyncedAt = &lastSyncedAt.Time
	}

	return product, nil
}

// GetByAccountAndProductID retrieves a product by account and platform product ID
func (r *externalAffiliateProductRepository) GetByAccountAndProductID(ctx context.Context, accountID, productID string) (*domain.ExternalAffiliateProduct, error) {
	query := `
		SELECT id, account_id, platform, product_id, product_name, product_url, product_image_url,
			category, brand, original_price, sale_price, currency, commission_rate, commission_type,
			estimated_earnings, total_clicks, total_orders, total_earnings, product_status, affiliate_link,
			short_link, last_synced_at, created_at, updated_at
		FROM external_affiliate_products WHERE account_id = $1 AND product_id = $2
	`

	product := &domain.ExternalAffiliateProduct{}
	var lastSyncedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, accountID, productID).Scan(
		&product.ID, &product.AccountID, &product.Platform, &product.ProductID, &product.ProductName,
		&product.ProductURL, &product.ProductImageURL, &product.Category, &product.Brand,
		&product.OriginalPrice, &product.SalePrice, &product.Currency, &product.CommissionRate,
		&product.CommissionType, &product.EstimatedEarnings, &product.TotalClicks, &product.TotalOrders,
		&product.TotalEarnings, &product.ProductStatus, &product.AffiliateLink, &product.ShortLink,
		&lastSyncedAt, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if lastSyncedAt.Valid {
		product.LastSyncedAt = &lastSyncedAt.Time
	}

	return product, nil
}

// GetAllByAccountID retrieves all products for an account
func (r *externalAffiliateProductRepository) GetAllByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*domain.ExternalAffiliateProduct, error) {
	query := `
		SELECT id, account_id, platform, product_id, product_name, product_url, product_image_url,
			category, brand, original_price, sale_price, currency, commission_rate, commission_type,
			estimated_earnings, total_clicks, total_orders, total_earnings, product_status, affiliate_link,
			created_at, updated_at
		FROM external_affiliate_products WHERE account_id = $1
		ORDER BY total_earnings DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*domain.ExternalAffiliateProduct, 0)
	for rows.Next() {
		product := &domain.ExternalAffiliateProduct{}
		err := rows.Scan(
			&product.ID, &product.AccountID, &product.Platform, &product.ProductID, &product.ProductName,
			&product.ProductURL, &product.ProductImageURL, &product.Category, &product.Brand,
			&product.OriginalPrice, &product.SalePrice, &product.Currency, &product.CommissionRate,
			&product.CommissionType, &product.EstimatedEarnings, &product.TotalClicks, &product.TotalOrders,
			&product.TotalEarnings, &product.ProductStatus, &product.AffiliateLink,
			&product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

// GetTopProducts retrieves top performing products for a user
func (r *externalAffiliateProductRepository) GetTopProducts(ctx context.Context, userID string, limit int) ([]*domain.ExternalAffiliateProduct, error) {
	query := `
		SELECT p.id, p.account_id, p.platform, p.product_id, p.product_name, p.product_url, p.product_image_url,
			p.category, p.brand, p.original_price, p.sale_price, p.currency, p.commission_rate, p.commission_type,
			p.estimated_earnings, p.total_clicks, p.total_orders, p.total_earnings, p.product_status, p.affiliate_link,
			p.created_at, p.updated_at
		FROM external_affiliate_products p
		JOIN external_affiliate_accounts a ON p.account_id = a.id
		WHERE a.user_id = $1
		ORDER BY p.total_earnings DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*domain.ExternalAffiliateProduct, 0)
	for rows.Next() {
		product := &domain.ExternalAffiliateProduct{}
		err := rows.Scan(
			&product.ID, &product.AccountID, &product.Platform, &product.ProductID, &product.ProductName,
			&product.ProductURL, &product.ProductImageURL, &product.Category, &product.Brand,
			&product.OriginalPrice, &product.SalePrice, &product.Currency, &product.CommissionRate,
			&product.CommissionType, &product.EstimatedEarnings, &product.TotalClicks, &product.TotalOrders,
			&product.TotalEarnings, &product.ProductStatus, &product.AffiliateLink,
			&product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

// Update updates a product
func (r *externalAffiliateProductRepository) Update(ctx context.Context, product *domain.ExternalAffiliateProduct) error {
	query := `
		UPDATE external_affiliate_products SET
			product_name = $2, product_url = $3, product_image_url = $4, category = $5, brand = $6,
			original_price = $7, sale_price = $8, commission_rate = $9, estimated_earnings = $10,
			product_status = $11, affiliate_link = $12, last_synced_at = $13, updated_at = $14
		WHERE id = $1
	`

	product.UpdatedAt = time.Now()
	now := time.Now()

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.ProductName, product.ProductURL, product.ProductImageURL, product.Category,
		product.Brand, product.OriginalPrice, product.SalePrice, product.CommissionRate, product.EstimatedEarnings,
		product.ProductStatus, product.AffiliateLink, now, product.UpdatedAt,
	)

	return err
}

// UpdateStats updates product stats
func (r *externalAffiliateProductRepository) UpdateStats(ctx context.Context, id string, clicks int64, orders int, earnings float64) error {
	query := `
		UPDATE external_affiliate_products SET
			total_clicks = $2, total_orders = $3, total_earnings = $4, updated_at = $5
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, clicks, orders, earnings, time.Now())
	return err
}

// Upsert creates or updates a product
func (r *externalAffiliateProductRepository) Upsert(ctx context.Context, product *domain.ExternalAffiliateProduct) error {
	existing, err := r.GetByAccountAndProductID(ctx, product.AccountID, product.ProductID)
	if err == nil {
		product.ID = existing.ID
		product.CreatedAt = existing.CreatedAt
		return r.Update(ctx, product)
	}

	return r.Create(ctx, product)
}

// Delete deletes a product
func (r *externalAffiliateProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM external_affiliate_products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}