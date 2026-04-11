package dto

import "time"

// CreateBudgetPeriodRequest represents the request to create a budget period
type CreateBudgetPeriodRequest struct {
	Name        string  `json:"name" binding:"required"`
	PeriodType  string  `json:"period_type" binding:"required"` // monthly, quarterly, annual, custom
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	TotalBudget float64 `json:"total_budget" binding:"required,min=0"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}

// UpdateBudgetPeriodRequest represents the request to update a budget period
type UpdateBudgetPeriodRequest struct {
	Name        string  `json:"name"`
	PeriodType  string  `json:"period_type"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	TotalBudget float64 `json:"total_budget"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}

// CreateBudgetAllocationRequest represents the request to add an allocation
type CreateBudgetAllocationRequest struct {
	Channel         string  `json:"channel" binding:"required"`
	ChannelCategory string  `json:"channel_category" binding:"required"` // online, offline, hybrid
	AllocatedAmount float64 `json:"allocated_amount" binding:"required,min=0"`
	ForecastAmount  *float64 `json:"forecast_amount,omitempty"`
	Notes           string  `json:"notes"`
}

// UpdateBudgetAllocationRequest represents the request to update an allocation
type UpdateBudgetAllocationRequest struct {
	Channel         string  `json:"channel"`
	ChannelCategory string  `json:"channel_category"`
	AllocatedAmount float64 `json:"allocated_amount"`
	SpentAmount     float64 `json:"spent_amount"`
	ForecastAmount  *float64 `json:"forecast_amount,omitempty"`
	Notes           string  `json:"notes"`
}

// RecordTransactionRequest represents the request to record a spend transaction
type RecordTransactionRequest struct {
	CampaignID    *string   `json:"campaign_id,omitempty"`
	AllocationID  *string   `json:"allocation_id,omitempty"`
	Channel       string    `json:"channel" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,min=0"`
	Currency      string    `json:"currency"`
	Description   string    `json:"description"`
	InvoiceNumber *string   `json:"invoice_number,omitempty"`
	Vendor        *string   `json:"vendor,omitempty"`
}

// BudgetPeriodResponse represents the budget period response
type BudgetPeriodResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	PeriodType     string    `json:"period_type"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	TotalBudget    float64   `json:"total_budget"`
	TotalAllocated float64   `json:"total_allocated"`
	TotalSpent     float64   `json:"total_spent"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

// BudgetPeriodWithAllocationsResponse includes allocations in the period response
type BudgetPeriodWithAllocationsResponse struct {
	ID             string                      `json:"id"`
	Name           string                      `json:"name"`
	PeriodType     string                      `json:"period_type"`
	StartDate      string                      `json:"start_date"`
	EndDate        string                      `json:"end_date"`
	TotalBudget    float64                     `json:"total_budget"`
	TotalAllocated float64                     `json:"total_allocated"`
	TotalSpent     float64                     `json:"total_spent"`
	Status         string                      `json:"status"`
	Notes          string                      `json:"notes,omitempty"`
	Allocations    []BudgetAllocationResponse  `json:"allocations"`
	CreatedAt      string                      `json:"created_at"`
	UpdatedAt      string                      `json:"updated_at"`
}

// BudgetAllocationResponse represents the budget allocation response
type BudgetAllocationResponse struct {
	ID              string    `json:"id"`
	PeriodID        string    `json:"period_id"`
	Channel         string    `json:"channel"`
	ChannelCategory string    `json:"channel_category"`
	AllocatedAmount float64   `json:"allocated_amount"`
	SpentAmount     float64   `json:"spent_amount"`
	ForecastAmount  *float64  `json:"forecast_amount,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

// BudgetTransactionResponse represents the budget transaction response
type BudgetTransactionResponse struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	CampaignID     *string `json:"campaign_id,omitempty"`
	AllocationID   *string `json:"allocation_id,omitempty"`
	Channel        string  `json:"channel"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	TransactionDate string `json:"transaction_date"`
	Description    string  `json:"description,omitempty"`
	InvoiceNumber  *string `json:"invoice_number,omitempty"`
	Vendor         *string `json:"vendor,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

// BudgetAlertResponse represents the budget alert response
type BudgetAlertResponse struct {
	ID               string     `json:"id"`
	OrganizationID   string     `json:"organization_id"`
	PeriodID         string     `json:"period_id"`
	CampaignID       *string    `json:"campaign_id,omitempty"`
	AlertType        string     `json:"alert_type"`
	ThresholdPercent *float64   `json:"threshold_percent,omitempty"`
	Triggered        bool       `json:"triggered"`
	TriggeredAt      *time.Time `json:"triggered_at,omitempty"`
	NotifiedUsers    []string   `json:"notified_users,omitempty"`
	CreatedAt        string     `json:"created_at"`
}
