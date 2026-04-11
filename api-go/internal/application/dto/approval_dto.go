package dto

// CreateApprovalWorkflowRequest represents the request to create an approval workflow
type CreateApprovalWorkflowRequest struct {
	Name                  string   `json:"name" binding:"required"`
	Description           string   `json:"description"`
	TriggerType           string   `json:"trigger_type" binding:"required"`
	Channels              []string `json:"channels"`
	MinApprovers          int      `json:"min_approvers"`
	ApprovalOrder         string   `json:"approval_order"`
	AutoApproveAfterHours *int     `json:"auto_approve_after_hours"`
	Status                string   `json:"status"`
}

// ApprovalWorkflowResponse represents the approval workflow response
type ApprovalWorkflowResponse struct {
	ID                    string   `json:"id"`
	OrganizationID        string   `json:"organization_id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description,omitempty"`
	TriggerType           string   `json:"trigger_type"`
	Channels              []string `json:"channels,omitempty"`
	MinApprovers          int      `json:"min_approvers"`
	ApprovalOrder         string   `json:"approval_order"`
	AutoApproveAfterHours *int     `json:"auto_approve_after_hours,omitempty"`
	Status                string   `json:"status"`
	CreatedAt             string   `json:"created_at"`
	UpdatedAt             string   `json:"updated_at"`
}

// CreateApprovalRequestRequest represents the request to create an approval request
type CreateApprovalRequestRequest struct {
	WorkflowID  string `json:"workflow_id" binding:"required"`
	EntityType  string `json:"entity_type" binding:"required"`
	EntityID    string `json:"entity_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

// RespondToRequest represents the request to respond to an approval request
type RespondToRequest struct {
	Action   string `json:"action" binding:"required,oneof=approved rejected commented requested_revision"`
	Comments string `json:"comments"`
}

// ApprovalResponseResponse represents an approval response in the API
type ApprovalResponseResponse struct {
	ID          string `json:"id"`
	RequestID   string `json:"request_id"`
	StepID      string `json:"step_id"`
	ApproverID  string `json:"approver_id"`
	Action      string `json:"action"`
	Comments    string `json:"comments,omitempty"`
	RespondedAt string `json:"responded_at"`
}

// ApprovalRequestResponse represents the approval request response
type ApprovalRequestResponse struct {
	ID           string                     `json:"id"`
	OrganizationID string                   `json:"organization_id"`
	WorkflowID   string                     `json:"workflow_id"`
	EntityType   string                     `json:"entity_type"`
	EntityID     string                     `json:"entity_id"`
	Title        string                     `json:"title"`
	Description  string                     `json:"description,omitempty"`
	SubmittedBy  *string                    `json:"submitted_by,omitempty"`
	SubmittedAt  string                     `json:"submitted_at"`
	Status       string                     `json:"status"`
	Deadline     *string                    `json:"deadline,omitempty"`
	CompletedAt  *string                    `json:"completed_at,omitempty"`
	Comments     string                     `json:"comments,omitempty"`
	Responses    []ApprovalResponseResponse `json:"responses,omitempty"`
}
