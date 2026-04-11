package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// ApprovalService handles business logic for approval workflows
type ApprovalService struct {
	workflowRepo domain.ApprovalWorkflowRepository
	stepRepo     domain.ApprovalStepRepository
	requestRepo  domain.ApprovalRequestRepository
	responseRepo domain.ApprovalResponseRepository
}

// NewApprovalService creates a new approval service
func NewApprovalService(
	workflowRepo domain.ApprovalWorkflowRepository,
	stepRepo domain.ApprovalStepRepository,
	requestRepo domain.ApprovalRequestRepository,
	responseRepo domain.ApprovalResponseRepository,
) *ApprovalService {
	return &ApprovalService{
		workflowRepo: workflowRepo,
		stepRepo:     stepRepo,
		requestRepo:  requestRepo,
		responseRepo: responseRepo,
	}
}

// CreateWorkflow creates a new approval workflow with optional steps
func (s *ApprovalService) CreateWorkflow(ctx context.Context, workflow *domain.ApprovalWorkflow, steps []*domain.ApprovalStep) error {
	if err := s.workflowRepo.Create(ctx, workflow); err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}

	for _, step := range steps {
		step.WorkflowID = workflow.ID
		if err := s.stepRepo.Create(ctx, step); err != nil {
			return fmt.Errorf("failed to create step %d: %w", step.StepOrder, err)
		}
	}

	return nil
}

// GetWorkflow retrieves an approval workflow by ID with its steps
func (s *ApprovalService) GetWorkflow(ctx context.Context, id string) (*domain.ApprovalWorkflow, []*domain.ApprovalStep, error) {
	workflow, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	steps, err := s.stepRepo.GetByWorkflowID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get steps: %w", err)
	}

	return workflow, steps, nil
}

// ListWorkflows retrieves all approval workflows for an organization
func (s *ApprovalService) ListWorkflows(ctx context.Context, organizationID string, limit, offset int) ([]*domain.ApprovalWorkflow, error) {
	return s.workflowRepo.List(ctx, organizationID, limit, offset)
}

// UpdateWorkflow updates an approval workflow
func (s *ApprovalService) UpdateWorkflow(ctx context.Context, workflow *domain.ApprovalWorkflow) error {
	return s.workflowRepo.Update(ctx, workflow)
}

// DeleteWorkflow deletes an approval workflow and its steps
func (s *ApprovalService) DeleteWorkflow(ctx context.Context, id string) error {
	if err := s.stepRepo.DeleteByWorkflowID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete steps: %w", err)
	}
	return s.workflowRepo.Delete(ctx, id)
}

// CreateRequest creates a new approval request
func (s *ApprovalService) CreateRequest(ctx context.Context, request *domain.ApprovalRequest) error {
	// Verify the workflow exists
	_, err := s.workflowRepo.GetByID(ctx, request.WorkflowID)
	if err != nil {
		return fmt.Errorf("workflow not found: %w", err)
	}

	request.Status = "pending"

	if err := s.requestRepo.Create(ctx, request); err != nil {
		return fmt.Errorf("failed to create approval request: %w", err)
	}

	return nil
}

// GetRequest retrieves an approval request by ID with its responses
func (s *ApprovalService) GetRequest(ctx context.Context, id string) (*domain.ApprovalRequest, []*domain.ApprovalResponse, error) {
	request, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get request: %w", err)
	}

	responses, err := s.responseRepo.ListByRequest(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get responses: %w", err)
	}

	return request, responses, nil
}

// ListRequests retrieves approval requests with optional filters
func (s *ApprovalService) ListRequests(ctx context.Context, organizationID string, status, entityType string, limit, offset int) ([]*domain.ApprovalRequest, error) {
	return s.requestRepo.List(ctx, organizationID, status, entityType, limit, offset)
}

// Respond submits a response to an approval request and checks if all required responses are collected
func (s *ApprovalService) Respond(ctx context.Context, requestID string, approverID string, action string, comments string) (*domain.ApprovalRequest, error) {
	// Get the request
	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}

	// Check if request is still pending
	if request.Status != "pending" {
		return nil, fmt.Errorf("request is no longer pending (current status: %s)", request.Status)
	}

	// Check if approver has already responded
	alreadyResponded, err := s.responseRepo.ExistsByRequestAndApprover(ctx, requestID, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing response: %w", err)
	}
	if alreadyResponded {
		return nil, fmt.Errorf("approver has already responded to this request")
	}

	// Create the response
	response := &domain.ApprovalResponse{
		ID:         uuid.New().String(),
		RequestID:  requestID,
		StepID:     "", // Will be determined by workflow steps
		ApproverID: approverID,
		Action:     action,
		Comments:   comments,
	}

	if err := s.responseRepo.Create(ctx, response); err != nil {
		return nil, fmt.Errorf("failed to create response: %w", err)
	}

	// Get the workflow to determine approval requirements
	workflow, err := s.workflowRepo.GetByID(ctx, request.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Check if we have enough responses to finalize the request
	approvedCount, err := s.responseRepo.CountByRequestAndAction(ctx, requestID, "approved")
	if err != nil {
		return nil, fmt.Errorf("failed to count approved responses: %w", err)
	}

	rejectedCount, err := s.responseRepo.CountByRequestAndAction(ctx, requestID, "rejected")
	if err != nil {
		return nil, fmt.Errorf("failed to count rejected responses: %w", err)
	}

	minApprovers := workflow.MinApprovers
	if minApprovers < 1 {
		minApprovers = 1
	}

	now := time.Now()
	var newStatus string

	// Determine if request can be finalized based on approval order
	switch workflow.ApprovalOrder {
	case "any":
		// Any single approval is enough
		if approvedCount >= 1 {
			newStatus = "approved"
		} else if rejectedCount >= 1 {
			newStatus = "rejected"
		}
	case "sequential":
		// Need sequential approvals - check against min_approvers
		if approvedCount >= minApprovers {
			newStatus = "approved"
		} else if rejectedCount >= 1 {
			// Any rejection in sequential flow rejects the request
			newStatus = "rejected"
		}
	case "parallel", "":
		// Need all required approvers (default to parallel)
		if approvedCount >= minApprovers {
			newStatus = "approved"
		} else if rejectedCount >= minApprovers {
			newStatus = "rejected"
		}
	default:
		// Default behavior: need min_approvers approvals
		if approvedCount >= minApprovers {
			newStatus = "approved"
		} else if rejectedCount >= 1 {
			newStatus = "rejected"
		}
	}

	if newStatus != "" {
		if err := s.requestRepo.UpdateStatus(ctx, requestID, newStatus, &now); err != nil {
			return nil, fmt.Errorf("failed to update request status: %w", err)
		}
		request.Status = newStatus
		request.CompletedAt = &now
	}

	return request, nil
}
