package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type budgetPeriodRepository struct {
	db *sql.DB
}

func NewBudgetPeriodRepository(db *sql.DB) domain.BudgetPeriodRepository {
	return &budgetPeriodRepository{db: db}
}

func (r *budgetPeriodRepository) Create(ctx context.Context, period *domain.BudgetPeriod) error {
	query := `
		INSERT INTO budget_periods (id, organization_id, name, period_type, start_date, end_date, total_budget, total_allocated, total_spent, status, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	now := time.Now()
	period.CreatedAt = now
	period.UpdatedAt = now

	var createdBy interface{}
	if period.CreatedBy != nil {
		createdBy = *period.CreatedBy
	}

	_, err := r.db.ExecContext(ctx, query,
		period.ID,
		period.OrganizationID,
		period.Name,
		period.PeriodType,
		period.StartDate,
		period.EndDate,
		period.TotalBudget,
		period.TotalAllocated,
		period.TotalSpent,
		period.Status,
		period.Notes,
		createdBy,
		period.CreatedAt,
		period.UpdatedAt,
	)
	return err
}

func (r *budgetPeriodRepository) GetByID(ctx context.Context, id string) (*domain.BudgetPeriod, error) {
	query := `
		SELECT id, organization_id, name, period_type, start_date, end_date, total_budget, total_allocated, total_spent, status, notes, created_by, created_at, updated_at
		FROM budget_periods
		WHERE id = $1
	`
	period := &domain.BudgetPeriod{}
	var createdBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&period.ID,
		&period.OrganizationID,
		&period.Name,
		&period.PeriodType,
		&period.StartDate,
		&period.EndDate,
		&period.TotalBudget,
		&period.TotalAllocated,
		&period.TotalSpent,
		&period.Status,
		&period.Notes,
		&createdBy,
		&period.CreatedAt,
		&period.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if createdBy.Valid {
		period.CreatedBy = &createdBy.String
	}
	return period, nil
}

func (r *budgetPeriodRepository) List(ctx context.Context, organizationID string, limit, offset int) ([]*domain.BudgetPeriod, error) {
	query := `
		SELECT id, organization_id, name, period_type, start_date, end_date, total_budget, total_allocated, total_spent, status, notes, created_by, created_at, updated_at
		FROM budget_periods
		WHERE organization_id = $1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	periods := make([]*domain.BudgetPeriod, 0)
	for rows.Next() {
		period := &domain.BudgetPeriod{}
		var createdBy sql.NullString

		err := rows.Scan(
			&period.ID,
			&period.OrganizationID,
			&period.Name,
			&period.PeriodType,
			&period.StartDate,
			&period.EndDate,
			&period.TotalBudget,
			&period.TotalAllocated,
			&period.TotalSpent,
			&period.Status,
			&period.Notes,
			&createdBy,
			&period.CreatedAt,
			&period.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if createdBy.Valid {
			period.CreatedBy = &createdBy.String
		}
		periods = append(periods, period)
	}
	return periods, rows.Err()
}

func (r *budgetPeriodRepository) Update(ctx context.Context, period *domain.BudgetPeriod) error {
	query := `
		UPDATE budget_periods
		SET name = $2, period_type = $3, start_date = $4, end_date = $5, total_budget = $6,
		    total_allocated = $7, total_spent = $8, status = $9, notes = $10, updated_at = $11
		WHERE id = $1
	`
	period.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		period.ID,
		period.Name,
		period.PeriodType,
		period.StartDate,
		period.EndDate,
		period.TotalBudget,
		period.TotalAllocated,
		period.TotalSpent,
		period.Status,
		period.Notes,
		period.UpdatedAt,
	)
	return err
}

func (r *budgetPeriodRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM budget_periods WHERE id = $1`, id)
	return err
}

// ===========================
// Budget Allocation Repository
// ===========================

type budgetAllocationRepository struct {
	db *sql.DB
}

func NewBudgetAllocationRepository(db *sql.DB) domain.BudgetAllocationRepository {
	return &budgetAllocationRepository{db: db}
}

func (r *budgetAllocationRepository) Create(ctx context.Context, alloc *domain.BudgetAllocation) error {
	query := `
		INSERT INTO budget_allocations (id, period_id, channel, channel_category, allocated_amount, spent_amount, forecast_amount, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	now := time.Now()
	alloc.CreatedAt = now
	alloc.UpdatedAt = now

	var forecastAmt interface{}
	if alloc.ForecastAmount != nil {
		forecastAmt = *alloc.ForecastAmount
	}

	_, err := r.db.ExecContext(ctx, query,
		alloc.ID,
		alloc.PeriodID,
		alloc.Channel,
		alloc.ChannelCategory,
		alloc.AllocatedAmount,
		alloc.SpentAmount,
		forecastAmt,
		alloc.Notes,
		alloc.CreatedAt,
		alloc.UpdatedAt,
	)
	return err
}

func (r *budgetAllocationRepository) GetByPeriod(ctx context.Context, periodID string) ([]*domain.BudgetAllocation, error) {
	query := `
		SELECT id, period_id, channel, channel_category, allocated_amount, spent_amount, forecast_amount, notes, created_at, updated_at
		FROM budget_allocations
		WHERE period_id = $1
		ORDER BY channel
	`
	rows, err := r.db.QueryContext(ctx, query, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allocs := make([]*domain.BudgetAllocation, 0)
	for rows.Next() {
		a := &domain.BudgetAllocation{}
		var forecastAmt sql.NullFloat64

		err := rows.Scan(
			&a.ID,
			&a.PeriodID,
			&a.Channel,
			&a.ChannelCategory,
			&a.AllocatedAmount,
			&a.SpentAmount,
			&forecastAmt,
			&a.Notes,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if forecastAmt.Valid {
			a.ForecastAmount = &forecastAmt.Float64
		}
		allocs = append(allocs, a)
	}
	return allocs, rows.Err()
}

func (r *budgetAllocationRepository) GetByID(ctx context.Context, id string) (*domain.BudgetAllocation, error) {
	query := `
		SELECT id, period_id, channel, channel_category, allocated_amount, spent_amount, forecast_amount, notes, created_at, updated_at
		FROM budget_allocations
		WHERE id = $1
	`
	a := &domain.BudgetAllocation{}
	var forecastAmt sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.PeriodID,
		&a.Channel,
		&a.ChannelCategory,
		&a.AllocatedAmount,
		&a.SpentAmount,
		&forecastAmt,
		&a.Notes,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if forecastAmt.Valid {
		a.ForecastAmount = &forecastAmt.Float64
	}
	return a, nil
}

func (r *budgetAllocationRepository) Update(ctx context.Context, alloc *domain.BudgetAllocation) error {
	query := `
		UPDATE budget_allocations
		SET channel = $2, channel_category = $3, allocated_amount = $4, spent_amount = $5,
		    forecast_amount = $6, notes = $7, updated_at = $8
		WHERE id = $1
	`
	alloc.UpdatedAt = time.Now()

	var forecastAmt interface{}
	if alloc.ForecastAmount != nil {
		forecastAmt = *alloc.ForecastAmount
	}

	_, err := r.db.ExecContext(ctx, query,
		alloc.ID,
		alloc.Channel,
		alloc.ChannelCategory,
		alloc.AllocatedAmount,
		alloc.SpentAmount,
		forecastAmt,
		alloc.Notes,
		alloc.UpdatedAt,
	)
	return err
}

func (r *budgetAllocationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM budget_allocations WHERE id = $1`, id)
	return err
}

// ===========================
// Budget Transaction Repository
// ===========================

type budgetTransactionRepository struct {
	db *sql.DB
}

func NewBudgetTransactionRepository(db *sql.DB) domain.BudgetTransactionRepository {
	return &budgetTransactionRepository{db: db}
}

func (r *budgetTransactionRepository) Create(ctx context.Context, txn *domain.BudgetTransaction) error {
	query := `
		INSERT INTO budget_transactions (id, organization_id, campaign_id, allocation_id, channel, amount, currency, transaction_date, description, invoice_number, vendor, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	now := time.Now()
	txn.CreatedAt = now

	var campaignID, allocationID, invoiceNumber, vendor interface{}
	if txn.CampaignID != nil {
		campaignID = *txn.CampaignID
	}
	if txn.AllocationID != nil {
		allocationID = *txn.AllocationID
	}
	if txn.InvoiceNumber != nil {
		invoiceNumber = *txn.InvoiceNumber
	}
	if txn.Vendor != nil {
		vendor = *txn.Vendor
	}

	_, err := r.db.ExecContext(ctx, query,
		txn.ID,
		txn.OrganizationID,
		campaignID,
		allocationID,
		txn.Channel,
		txn.Amount,
		txn.Currency,
		txn.TransactionDate,
		txn.Description,
		invoiceNumber,
		vendor,
		txn.CreatedAt,
	)
	return err
}

func (r *budgetTransactionRepository) List(ctx context.Context, organizationID string, limit, offset int) ([]*domain.BudgetTransaction, error) {
	query := `
		SELECT id, organization_id, campaign_id, allocation_id, channel, amount, currency, transaction_date, description, invoice_number, vendor, created_at
		FROM budget_transactions
		WHERE organization_id = $1
		ORDER BY transaction_date DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txns := make([]*domain.BudgetTransaction, 0)
	for rows.Next() {
		t := &domain.BudgetTransaction{}
		var campaignID, allocationID, invoiceNumber, vendor sql.NullString

		err := rows.Scan(
			&t.ID,
			&t.OrganizationID,
			&campaignID,
			&allocationID,
			&t.Channel,
			&t.Amount,
			&t.Currency,
			&t.TransactionDate,
			&t.Description,
			&invoiceNumber,
			&vendor,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if campaignID.Valid {
			t.CampaignID = &campaignID.String
		}
		if allocationID.Valid {
			t.AllocationID = &allocationID.String
		}
		if invoiceNumber.Valid {
			t.InvoiceNumber = &invoiceNumber.String
		}
		if vendor.Valid {
			t.Vendor = &vendor.String
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

// ===========================
// Budget Alert Repository
// ===========================

type budgetAlertRepository struct {
	db *sql.DB
}

func NewBudgetAlertRepository(db *sql.DB) domain.BudgetAlertRepository {
	return &budgetAlertRepository{db: db}
}

func (r *budgetAlertRepository) Create(ctx context.Context, alert *domain.BudgetAlert) error {
	query := `
		INSERT INTO budget_alerts (id, organization_id, period_id, campaign_id, alert_type, threshold_percent, triggered, triggered_at, notified_users, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	now := time.Now()
	alert.CreatedAt = now

	var campaignID, triggeredAt interface{}
	if alert.CampaignID != nil {
		campaignID = *alert.CampaignID
	}
	if alert.TriggeredAt != nil {
		triggeredAt = *alert.TriggeredAt
	}

	_, err := r.db.ExecContext(ctx, query,
		alert.ID,
		alert.OrganizationID,
		alert.PeriodID,
		campaignID,
		alert.AlertType,
		alert.ThresholdPercent,
		alert.Triggered,
		triggeredAt,
		alert.NotifiedUsers,
		alert.CreatedAt,
	)
	return err
}

func (r *budgetAlertRepository) List(ctx context.Context, organizationID string, limit, offset int) ([]*domain.BudgetAlert, error) {
	query := `
		SELECT id, organization_id, period_id, campaign_id, alert_type, threshold_percent, triggered, triggered_at, notified_users, created_at
		FROM budget_alerts
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alerts := make([]*domain.BudgetAlert, 0)
	for rows.Next() {
		a := &domain.BudgetAlert{}
		var campaignID sql.NullString
		var triggeredAt sql.NullTime
		var notifiedUsers []string

		err := rows.Scan(
			&a.ID,
			&a.OrganizationID,
			&a.PeriodID,
			&campaignID,
			&a.AlertType,
			&a.ThresholdPercent,
			&a.Triggered,
			&triggeredAt,
			&notifiedUsers,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if campaignID.Valid {
			a.CampaignID = &campaignID.String
		}
		if triggeredAt.Valid {
			a.TriggeredAt = &triggeredAt.Time
		}
		if notifiedUsers != nil {
			a.NotifiedUsers = notifiedUsers
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (r *budgetAlertRepository) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) error {
	query := `UPDATE budget_alerts SET triggered = true, triggered_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, triggeredAt)
	return err
}
