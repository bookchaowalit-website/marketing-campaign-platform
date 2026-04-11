package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// BudgetService handles business logic for budget planning
type BudgetService struct {
	periodRepo domain.BudgetPeriodRepository
	allocRepo  domain.BudgetAllocationRepository
	txnRepo    domain.BudgetTransactionRepository
	alertRepo  domain.BudgetAlertRepository
}

// NewBudgetService creates a new budget service
func NewBudgetService(
	periodRepo domain.BudgetPeriodRepository,
	allocRepo domain.BudgetAllocationRepository,
	txnRepo domain.BudgetTransactionRepository,
	alertRepo domain.BudgetAlertRepository,
) *BudgetService {
	return &BudgetService{
		periodRepo: periodRepo,
		allocRepo:  allocRepo,
		txnRepo:    txnRepo,
		alertRepo:  alertRepo,
	}
}

// CreatePeriod creates a new budget period
func (s *BudgetService) CreatePeriod(ctx context.Context, period *domain.BudgetPeriod) error {
	return s.periodRepo.Create(ctx, period)
}

// GetPeriod retrieves a budget period by ID
func (s *BudgetService) GetPeriod(ctx context.Context, id string) (*domain.BudgetPeriod, error) {
	return s.periodRepo.GetByID(ctx, id)
}

// ListPeriods retrieves all budget periods for an organization
func (s *BudgetService) ListPeriods(ctx context.Context, organizationID string, limit, offset int) ([]*domain.BudgetPeriod, error) {
	return s.periodRepo.List(ctx, organizationID, limit, offset)
}

// UpdatePeriod updates a budget period
func (s *BudgetService) UpdatePeriod(ctx context.Context, period *domain.BudgetPeriod) error {
	return s.periodRepo.Update(ctx, period)
}

// DeletePeriod deletes a budget period
func (s *BudgetService) DeletePeriod(ctx context.Context, id string) error {
	return s.periodRepo.Delete(ctx, id)
}

// CreateAllocation adds a budget allocation to a period
func (s *BudgetService) CreateAllocation(ctx context.Context, alloc *domain.BudgetAllocation) error {
	return s.allocRepo.Create(ctx, alloc)
}

// ListAllocations retrieves all allocations for a period
func (s *BudgetService) ListAllocations(ctx context.Context, periodID string) ([]*domain.BudgetAllocation, error) {
	return s.allocRepo.GetByPeriod(ctx, periodID)
}

// GetAllocation retrieves a single allocation by ID
func (s *BudgetService) GetAllocation(ctx context.Context, id string) (*domain.BudgetAllocation, error) {
	return s.allocRepo.GetByID(ctx, id)
}

// UpdateAllocation updates a budget allocation
func (s *BudgetService) UpdateAllocation(ctx context.Context, alloc *domain.BudgetAllocation) error {
	return s.allocRepo.Update(ctx, alloc)
}

// DeleteAllocation deletes a budget allocation
func (s *BudgetService) DeleteAllocation(ctx context.Context, id string) error {
	return s.allocRepo.Delete(ctx, id)
}

// RecordTransaction records a spend transaction and updates allocation spent amount
func (s *BudgetService) RecordTransaction(ctx context.Context, txn *domain.BudgetTransaction) error {
	// Create the transaction
	if err := s.txnRepo.Create(ctx, txn); err != nil {
		return err
	}

	// If the transaction is linked to an allocation, update the allocation's spent_amount
	if txn.AllocationID != nil {
		alloc, err := s.allocRepo.GetByID(ctx, *txn.AllocationID)
		if err != nil {
			return fmt.Errorf("failed to get allocation: %w", err)
		}

		alloc.SpentAmount += txn.Amount
		if err := s.allocRepo.Update(ctx, alloc); err != nil {
			return fmt.Errorf("failed to update allocation spent amount: %w", err)
		}
	}

	return nil
}

// CheckAlerts checks all alerts for a period and triggers those whose threshold is exceeded
func (s *BudgetService) CheckAlerts(ctx context.Context, organizationID string) ([]*domain.BudgetAlert, error) {
	alerts, err := s.alertRepo.List(ctx, organizationID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}

	triggered := make([]*domain.BudgetAlert, 0)

	for _, alert := range alerts {
		if alert.Triggered {
			continue
		}

		// Get the period to check spending
		period, err := s.periodRepo.GetByID(ctx, alert.PeriodID)
		if err != nil {
			continue
		}

		shouldTrigger := false

		if alert.ThresholdPercent != nil && period.TotalBudget > 0 {
			spendPercent := (period.TotalSpent / period.TotalBudget) * 100
			if spendPercent >= *alert.ThresholdPercent {
				shouldTrigger = true
			}
		}

		if shouldTrigger {
			now := time.Now()
			if err := s.alertRepo.MarkTriggered(ctx, alert.ID, now); err != nil {
				continue
			}
			alert.Triggered = true
			alert.TriggeredAt = &now
			triggered = append(triggered, alert)
		}
	}

	return triggered, nil
}

// ListTransactions retrieves all transactions for an organization
func (s *BudgetService) ListTransactions(ctx context.Context, organizationID string, limit, offset int) ([]*domain.BudgetTransaction, error) {
	return s.txnRepo.List(ctx, organizationID, limit, offset)
}

// ListAlerts retrieves all alerts for an organization
func (s *BudgetService) ListAlerts(ctx context.Context, organizationID string, limit, offset int) ([]*domain.BudgetAlert, error) {
	return s.alertRepo.List(ctx, organizationID, limit, offset)
}

// CreateAlert creates a new budget alert
func (s *BudgetService) CreateAlert(ctx context.Context, alert *domain.BudgetAlert) error {
	return s.alertRepo.Create(ctx, alert)
}
