package domain

import (
	"context"
	"time"
)

// BudgetPeriod represents a budget planning period
type BudgetPeriod struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Name           string     `json:"name"`
	PeriodType     string     `json:"period_type"` // monthly, quarterly, annual, custom
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	TotalBudget    float64    `json:"total_budget"`
	TotalAllocated float64    `json:"total_allocated"`
	TotalSpent     float64    `json:"total_spent"`
	Status         string     `json:"status"` // draft, approved, active, closed
	Notes          string     `json:"notes,omitempty"`
	CreatedBy      *string    `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// BudgetAllocation represents budget allocated to a channel within a period
type BudgetAllocation struct {
	ID              string    `json:"id"`
	PeriodID        string    `json:"period_id"`
	Channel         string    `json:"channel"`
	ChannelCategory string    `json:"channel_category"` // online, offline, hybrid
	AllocatedAmount float64   `json:"allocated_amount"`
	SpentAmount     float64   `json:"spent_amount"`
	ForecastAmount  *float64  `json:"forecast_amount,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BudgetTransaction represents an actual spend transaction
type BudgetTransaction struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	CampaignID     *string    `json:"campaign_id,omitempty"`
	AllocationID   *string    `json:"allocation_id,omitempty"`
	Channel        string     `json:"channel"`
	Amount         float64    `json:"amount"`
	Currency       string     `json:"currency"`
	TransactionDate time.Time `json:"transaction_date"`
	Description    string     `json:"description,omitempty"`
	InvoiceNumber  *string    `json:"invoice_number,omitempty"`
	Vendor         *string    `json:"vendor,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// BudgetAlert represents a budget alert configuration
type BudgetAlert struct {
	ID               string     `json:"id"`
	OrganizationID   string     `json:"organization_id"`
	PeriodID         string     `json:"period_id"`
	CampaignID       *string    `json:"campaign_id,omitempty"`
	AlertType        string     `json:"alert_type"` // threshold, overspend, forecast
	ThresholdPercent *float64   `json:"threshold_percent,omitempty"`
	Triggered        bool       `json:"triggered"`
	TriggeredAt      *time.Time `json:"triggered_at,omitempty"`
	NotifiedUsers    []string   `json:"notified_users,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// BudgetPeriodRepository defines the interface for budget period repository
type BudgetPeriodRepository interface {
	Create(ctx context.Context, period *BudgetPeriod) error
	GetByID(ctx context.Context, id string) (*BudgetPeriod, error)
	List(ctx context.Context, organizationID string, limit, offset int) ([]*BudgetPeriod, error)
	Update(ctx context.Context, period *BudgetPeriod) error
	Delete(ctx context.Context, id string) error
}

// BudgetAllocationRepository defines the interface for budget allocation repository
type BudgetAllocationRepository interface {
	Create(ctx context.Context, alloc *BudgetAllocation) error
	GetByPeriod(ctx context.Context, periodID string) ([]*BudgetAllocation, error)
	GetByID(ctx context.Context, id string) (*BudgetAllocation, error)
	Update(ctx context.Context, alloc *BudgetAllocation) error
	Delete(ctx context.Context, id string) error
}

// BudgetTransactionRepository defines the interface for budget transaction repository
type BudgetTransactionRepository interface {
	Create(ctx context.Context, txn *BudgetTransaction) error
	List(ctx context.Context, organizationID string, limit, offset int) ([]*BudgetTransaction, error)
}

// BudgetAlertRepository defines the interface for budget alert repository
type BudgetAlertRepository interface {
	Create(ctx context.Context, alert *BudgetAlert) error
	List(ctx context.Context, organizationID string, limit, offset int) ([]*BudgetAlert, error)
	MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) error
}
