package domain

import (
	"context"
	"time"
)

// ApprovalWorkflow represents an approval workflow template
type ApprovalWorkflow struct {
	ID                   string     `json:"id"`
	OrganizationID       string     `json:"organization_id"`
	Name                 string     `json:"name"`
	Description          string     `json:"description,omitempty"`
	TriggerType          string     `json:"trigger_type"` // campaign_create, campaign_publish, content_create, content_publish, budget_create, budget_approve, vendor_create, custom
	Channels             []string   `json:"channels,omitempty"`
	MinApprovers         int        `json:"min_approvers"`
	ApprovalOrder        string     `json:"approval_order"` // sequential, parallel, any
	AutoApproveAfterHours *int      `json:"auto_approve_after_hours,omitempty"`
	Status               string     `json:"status"` // active, inactive
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// ApprovalStep represents a step in an approval workflow
type ApprovalStep struct {
	ID            string     `json:"id"`
	WorkflowID    string     `json:"workflow_id"`
	StepOrder     int        `json:"step_order"`
	ApproverRole  string     `json:"approver_role,omitempty"`
	ApproverUser *string    `json:"approver_user,omitempty"`
	Required      bool       `json:"required"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ApprovalRequest represents an approval request instance
type ApprovalRequest struct {
	ID            string     `json:"id"`
	OrganizationID string    `json:"organization_id"`
	WorkflowID    string     `json:"workflow_id"`
	EntityType    string     `json:"entity_type"`
	EntityID      string     `json:"entity_id"`
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	SubmittedBy   *string    `json:"submitted_by,omitempty"`
	SubmittedAt   time.Time  `json:"submitted_at"`
	Status        string     `json:"status"` // pending, approved, rejected, cancelled, expired
	Deadline      *time.Time `json:"deadline,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	Comments      string     `json:"comments,omitempty"`
}

// ApprovalResponse represents a response to an approval request
type ApprovalResponse struct {
	ID           string    `json:"id"`
	RequestID    string    `json:"request_id"`
	StepID       string    `json:"step_id"`
	ApproverID   string    `json:"approver_id"`
	Action       string    `json:"action"` // approved, rejected, commented, requested_revision
	Comments     string    `json:"comments,omitempty"`
	RespondedAt  time.Time `json:"responded_at"`
}

// ApprovalWorkflowRepository defines the interface for approval workflow repository
type ApprovalWorkflowRepository interface {
	Create(ctx context.Context, workflow *ApprovalWorkflow) error
	GetByID(ctx context.Context, id string) (*ApprovalWorkflow, error)
	List(ctx context.Context, organizationID string, limit, offset int) ([]*ApprovalWorkflow, error)
	Update(ctx context.Context, workflow *ApprovalWorkflow) error
	Delete(ctx context.Context, id string) error
}

// ApprovalStepRepository defines the interface for approval step repository
type ApprovalStepRepository interface {
	Create(ctx context.Context, step *ApprovalStep) error
	GetByWorkflowID(ctx context.Context, workflowID string) ([]*ApprovalStep, error)
	GetByID(ctx context.Context, id string) (*ApprovalStep, error)
	DeleteByWorkflowID(ctx context.Context, workflowID string) error
}

// ApprovalRequestRepository defines the interface for approval request repository
type ApprovalRequestRepository interface {
	Create(ctx context.Context, request *ApprovalRequest) error
	GetByID(ctx context.Context, id string) (*ApprovalRequest, error)
	List(ctx context.Context, organizationID string, status, entityType string, limit, offset int) ([]*ApprovalRequest, error)
	UpdateStatus(ctx context.Context, id string, status string, completedAt *time.Time) error
}

// ApprovalResponseRepository defines the interface for approval response repository
type ApprovalResponseRepository interface {
	Create(ctx context.Context, response *ApprovalResponse) error
	ListByRequest(ctx context.Context, requestID string) ([]*ApprovalResponse, error)
	CountByRequestAndAction(ctx context.Context, requestID string, action string) (int, error)
	ExistsByRequestAndApprover(ctx context.Context, requestID string, approverID string) (bool, error)
}
