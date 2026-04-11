package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// externalAffiliateOrderRepository implements domain.ExternalAffiliateOrderRepository
type externalAffiliateOrderRepository struct {
	db *sql.DB
}

// NewExternalAffiliateOrderRepository creates a new external affiliate order repository
func NewExternalAffiliateOrderRepository(db *sql.DB) domain.ExternalAffiliateOrderRepository {
	return &externalAffiliateOrderRepository{db: db}
}

// Create creates a new order
func (r *externalAffiliateOrderRepository) Create(ctx context.Context, order *domain.ExternalAffiliateOrder) error {
	query := `
		INSERT INTO external_affiliate_orders (
			id, account_id, product_id, platform, order_id, parent_order_id, order_status, order_date,
			order_items, order_amount, product_amount, shipping_amount, discount_amount, currency,
			commission_rate, commission_amount, commission_status, customer_id, traffic_source,
			utm_source, utm_medium, utm_campaign, synced_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	`

	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	order.SyncedAt = now

	var orderDate interface{}
	if order.OrderDate != nil {
		orderDate = *order.OrderDate
	}

	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.AccountID, order.ProductID, order.Platform, order.OrderID, order.ParentOrderID,
		order.OrderStatus, orderDate, order.OrderItems, order.OrderAmount, order.ProductAmount,
		order.ShippingAmount, order.DiscountAmount, order.Currency, order.CommissionRate,
		order.CommissionAmount, order.CommissionStatus, order.CustomerID, order.TrafficSource,
		order.UTMSource, order.UTMMedium, order.UTMCampaign, order.SyncedAt, order.CreatedAt, order.UpdatedAt,
	)

	return err
}

// GetByID retrieves an order by ID
func (r *externalAffiliateOrderRepository) GetByID(ctx context.Context, id string) (*domain.ExternalAffiliateOrder, error) {
	query := `
		SELECT id, account_id, product_id, platform, order_id, parent_order_id, order_status, order_date,
			order_items, order_amount, product_amount, shipping_amount, discount_amount, currency,
			commission_rate, commission_amount, commission_status, commission_paid_at, customer_id,
			traffic_source, utm_source, utm_medium, utm_campaign, synced_at, created_at, updated_at
		FROM external_affiliate_orders WHERE id = $1
	`

	order := &domain.ExternalAffiliateOrder{}
	var orderDate, commissionPaidAt sql.NullTime
	var parentOrderID, productID, customerID, trafficSource, utmSource, utmMedium, utmCampaign sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.AccountID, &productID, &order.Platform, &order.OrderID, &parentOrderID,
		&order.OrderStatus, &orderDate, &order.OrderItems, &order.OrderAmount, &order.ProductAmount,
		&order.ShippingAmount, &order.DiscountAmount, &order.Currency, &order.CommissionRate,
		&order.CommissionAmount, &order.CommissionStatus, &commissionPaidAt, &customerID,
		&trafficSource, &utmSource, &utmMedium, &utmCampaign, &order.SyncedAt, &order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if orderDate.Valid {
		order.OrderDate = &orderDate.Time
	}
	if commissionPaidAt.Valid {
		order.CommissionPaidAt = &commissionPaidAt.Time
	}
	if parentOrderID.Valid {
		order.ParentOrderID = parentOrderID.String
	}
	if productID.Valid {
		order.ProductID = productID.String
	}
	if customerID.Valid {
		order.CustomerID = customerID.String
	}
	if trafficSource.Valid {
		order.TrafficSource = trafficSource.String
	}
	if utmSource.Valid {
		order.UTMSource = utmSource.String
	}
	if utmMedium.Valid {
		order.UTMMedium = utmMedium.String
	}
	if utmCampaign.Valid {
		order.UTMCampaign = utmCampaign.String
	}

	return order, nil
}

// GetByAccountAndOrderID retrieves an order by account and order ID
func (r *externalAffiliateOrderRepository) GetByAccountAndOrderID(ctx context.Context, accountID, orderID string) (*domain.ExternalAffiliateOrder, error) {
	query := `
		SELECT id, account_id, product_id, platform, order_id, parent_order_id, order_status, order_date,
			order_items, order_amount, product_amount, shipping_amount, discount_amount, currency,
			commission_rate, commission_amount, commission_status, commission_paid_at, customer_id,
			traffic_source, utm_source, utm_medium, utm_campaign, synced_at, created_at, updated_at
		FROM external_affiliate_orders WHERE account_id = $1 AND order_id = $2
	`

	order := &domain.ExternalAffiliateOrder{}
	var orderDate, commissionPaidAt sql.NullTime
	var parentOrderID, productID, customerID, trafficSource, utmSource, utmMedium, utmCampaign sql.NullString

	err := r.db.QueryRowContext(ctx, query, accountID, orderID).Scan(
		&order.ID, &order.AccountID, &productID, &order.Platform, &order.OrderID, &parentOrderID,
		&order.OrderStatus, &orderDate, &order.OrderItems, &order.OrderAmount, &order.ProductAmount,
		&order.ShippingAmount, &order.DiscountAmount, &order.Currency, &order.CommissionRate,
		&order.CommissionAmount, &order.CommissionStatus, &commissionPaidAt, &customerID,
		&trafficSource, &utmSource, &utmMedium, &utmCampaign, &order.SyncedAt, &order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if orderDate.Valid {
		order.OrderDate = &orderDate.Time
	}
	if commissionPaidAt.Valid {
		order.CommissionPaidAt = &commissionPaidAt.Time
	}
	if parentOrderID.Valid {
		order.ParentOrderID = parentOrderID.String
	}
	if productID.Valid {
		order.ProductID = productID.String
	}
	if customerID.Valid {
		order.CustomerID = customerID.String
	}
	if trafficSource.Valid {
		order.TrafficSource = trafficSource.String
	}
	if utmSource.Valid {
		order.UTMSource = utmSource.String
	}
	if utmMedium.Valid {
		order.UTMMedium = utmMedium.String
	}
	if utmCampaign.Valid {
		order.UTMCampaign = utmCampaign.String
	}

	return order, nil
}

// GetAllByAccountID retrieves all orders for an account
func (r *externalAffiliateOrderRepository) GetAllByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*domain.ExternalAffiliateOrder, error) {
	query := `
		SELECT id, account_id, product_id, platform, order_id, parent_order_id, order_status, order_date,
			order_items, order_amount, product_amount, shipping_amount, discount_amount, currency,
			commission_rate, commission_amount, commission_status, customer_id,
			traffic_source, utm_source, utm_medium, utm_campaign, synced_at, created_at, updated_at
		FROM external_affiliate_orders WHERE account_id = $1
		ORDER BY order_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanOrders(rows)
}

// GetRecentByUserID retrieves recent orders for a user
func (r *externalAffiliateOrderRepository) GetRecentByUserID(ctx context.Context, userID string, limit int) ([]*domain.ExternalAffiliateOrder, error) {
	query := `
		SELECT o.id, o.account_id, o.product_id, o.platform, o.order_id, o.parent_order_id, o.order_status, o.order_date,
			o.order_items, o.order_amount, o.product_amount, o.shipping_amount, o.discount_amount, o.currency,
			o.commission_rate, o.commission_amount, o.commission_status, o.customer_id,
			o.traffic_source, o.utm_source, o.utm_medium, o.utm_campaign, o.synced_at, o.created_at, o.updated_at
		FROM external_affiliate_orders o
		JOIN external_affiliate_accounts a ON o.account_id = a.id
		WHERE a.user_id = $1
		ORDER BY o.order_date DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanOrders(rows)
}

// GetByDateRange retrieves orders within a date range
func (r *externalAffiliateOrderRepository) GetByDateRange(ctx context.Context, accountID string, start, end time.Time) ([]*domain.ExternalAffiliateOrder, error) {
	query := `
		SELECT id, account_id, product_id, platform, order_id, parent_order_id, order_status, order_date,
			order_items, order_amount, product_amount, shipping_amount, discount_amount, currency,
			commission_rate, commission_amount, commission_status, customer_id,
			traffic_source, utm_source, utm_medium, utm_campaign, synced_at, created_at, updated_at
		FROM external_affiliate_orders WHERE account_id = $1 AND order_date >= $2 AND order_date <= $3
		ORDER BY order_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanOrders(rows)
}

// Update updates an order
func (r *externalAffiliateOrderRepository) Update(ctx context.Context, order *domain.ExternalAffiliateOrder) error {
	query := `
		UPDATE external_affiliate_orders SET
			order_status = $2, order_items = $3, order_amount = $4, commission_amount = $5,
			commission_status = $6, updated_at = $7
		WHERE id = $1
	`

	order.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.OrderStatus, order.OrderItems, order.OrderAmount, order.CommissionAmount,
		order.CommissionStatus, order.UpdatedAt,
	)

	return err
}

// Upsert creates or updates an order
func (r *externalAffiliateOrderRepository) Upsert(ctx context.Context, order *domain.ExternalAffiliateOrder) error {
	existing, err := r.GetByAccountAndOrderID(ctx, order.AccountID, order.OrderID)
	if err == nil {
		order.ID = existing.ID
		order.CreatedAt = existing.CreatedAt
		return r.Update(ctx, order)
	}

	return r.Create(ctx, order)
}

// Delete deletes an order
func (r *externalAffiliateOrderRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM external_affiliate_orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Helper function to scan orders
func (r *externalAffiliateOrderRepository) scanOrders(rows *sql.Rows) ([]*domain.ExternalAffiliateOrder, error) {
	orders := make([]*domain.ExternalAffiliateOrder, 0)
	for rows.Next() {
		order := &domain.ExternalAffiliateOrder{}
		var orderDate sql.NullTime
		var parentOrderID, productID, customerID, trafficSource, utmSource, utmMedium, utmCampaign sql.NullString

		err := rows.Scan(
			&order.ID, &order.AccountID, &productID, &order.Platform, &order.OrderID, &parentOrderID,
			&order.OrderStatus, &orderDate, &order.OrderItems, &order.OrderAmount, &order.ProductAmount,
			&order.ShippingAmount, &order.DiscountAmount, &order.Currency, &order.CommissionRate,
			&order.CommissionAmount, &order.CommissionStatus, &customerID,
			&trafficSource, &utmSource, &utmMedium, &utmCampaign, &order.SyncedAt, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if orderDate.Valid {
			order.OrderDate = &orderDate.Time
		}
		if parentOrderID.Valid {
			order.ParentOrderID = parentOrderID.String
		}
		if productID.Valid {
			order.ProductID = productID.String
		}
		if customerID.Valid {
			order.CustomerID = customerID.String
		}
		if trafficSource.Valid {
			order.TrafficSource = trafficSource.String
		}
		if utmSource.Valid {
			order.UTMSource = utmSource.String
		}
		if utmMedium.Valid {
			order.UTMMedium = utmMedium.String
		}
		if utmCampaign.Valid {
			order.UTMCampaign = utmCampaign.String
		}

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// externalAffiliateSyncLogRepository implements domain.ExternalAffiliateSyncLogRepository
type externalAffiliateSyncLogRepository struct {
	db *sql.DB
}

// NewExternalAffiliateSyncLogRepository creates a new sync log repository
func NewExternalAffiliateSyncLogRepository(db *sql.DB) domain.ExternalAffiliateSyncLogRepository {
	return &externalAffiliateSyncLogRepository{db: db}
}

// Create creates a new sync log
func (r *externalAffiliateSyncLogRepository) Create(ctx context.Context, log *domain.ExternalAffiliateSyncLog) error {
	query := `
		INSERT INTO external_affiliate_sync_logs (
			id, account_id, platform, sync_type, status, started_at, api_calls_count, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	log.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.AccountID, log.Platform, log.SyncType, log.Status, log.StartedAt, log.APICallsCount, log.CreatedAt,
	)

	return err
}

// GetByID retrieves a sync log by ID
func (r *externalAffiliateSyncLogRepository) GetByID(ctx context.Context, id string) (*domain.ExternalAffiliateSyncLog, error) {
	query := `
		SELECT id, account_id, platform, sync_type, status, started_at, completed_at, duration_ms,
			records_processed, records_created, records_updated, records_failed, error_message,
			api_calls_count, rate_limit_hit, created_at
		FROM external_affiliate_sync_logs WHERE id = $1
	`

	log := &domain.ExternalAffiliateSyncLog{}
	var completedAt sql.NullTime
	var durationMs sql.NullInt64
	var errorMessage sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID, &log.AccountID, &log.Platform, &log.SyncType, &log.Status, &log.StartedAt,
		&completedAt, &durationMs, &log.RecordsProcessed, &log.RecordsCreated, &log.RecordsUpdated,
		&log.RecordsFailed, &errorMessage, &log.APICallsCount, &log.RateLimitHit, &log.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		log.CompletedAt = &completedAt.Time
	}
	if durationMs.Valid {
		log.DurationMs = durationMs.Int64
	}
	if errorMessage.Valid {
		log.ErrorMessage = errorMessage.String
	}

	return log, nil
}

// GetRecentByAccountID retrieves recent sync logs for an account
func (r *externalAffiliateSyncLogRepository) GetRecentByAccountID(ctx context.Context, accountID string, limit int) ([]*domain.ExternalAffiliateSyncLog, error) {
	query := `
		SELECT id, account_id, platform, sync_type, status, started_at, completed_at, duration_ms,
			records_processed, records_created, records_updated, records_failed, error_message,
			api_calls_count, rate_limit_hit, created_at
		FROM external_affiliate_sync_logs WHERE account_id = $1
		ORDER BY started_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*domain.ExternalAffiliateSyncLog, 0)
	for rows.Next() {
		log := &domain.ExternalAffiliateSyncLog{}
		var completedAt sql.NullTime
		var durationMs sql.NullInt64
		var errorMessage sql.NullString

		err := rows.Scan(
			&log.ID, &log.AccountID, &log.Platform, &log.SyncType, &log.Status, &log.StartedAt,
			&completedAt, &durationMs, &log.RecordsProcessed, &log.RecordsCreated, &log.RecordsUpdated,
			&log.RecordsFailed, &errorMessage, &log.APICallsCount, &log.RateLimitHit, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			log.CompletedAt = &completedAt.Time
		}
		if durationMs.Valid {
			log.DurationMs = durationMs.Int64
		}
		if errorMessage.Valid {
			log.ErrorMessage = errorMessage.String
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// UpdateStatus updates the status of a sync log
func (r *externalAffiliateSyncLogRepository) UpdateStatus(ctx context.Context, id string, status string, completedAt time.Time, recordsCreated, recordsUpdated, recordsFailed int, errorMsg string) error {
	query := `
		UPDATE external_affiliate_sync_logs SET
			status = $2, completed_at = $3, duration_ms = $4, records_created = $5, records_updated = $6,
			records_failed = $7, error_message = $8
		WHERE id = $1
	`

	var errorMsgVal interface{}
	if errorMsg != "" {
		errorMsgVal = errorMsg
	}

	durationMs := completedAt.Sub(completedAt).Milliseconds()

	_, err := r.db.ExecContext(ctx, query,
		id, status, completedAt, durationMs, recordsCreated, recordsUpdated, recordsFailed, errorMsgVal,
	)

	return err
}