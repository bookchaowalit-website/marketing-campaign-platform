package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// ===========================
// Approval Workflow Repository
// ===========================

type approvalWorkflowRepository struct {
	db *sql.DB
}

func NewApprovalWorkflowRepository(db *sql.DB) domain.ApprovalWorkflowRepository {
	return &approvalWorkflowRepository{db: db}
}

func (r *approvalWorkflowRepository) Create(ctx context.Context, workflow *domain.ApprovalWorkflow) error {
	query := `
		INSERT INTO approval_workflows (id, organization_id, name, description, trigger_type, channels, min_approvers, approval_order, auto_approve_after_hours, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	now := time.Now()
	workflow.CreatedAt = now
	workflow.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		workflow.ID,
		workflow.OrganizationID,
		workflow.Name,
		workflow.Description,
		workflow.TriggerType,
		workflow.Channels,
		workflow.MinApprovers,
		workflow.ApprovalOrder,
		workflow.AutoApproveAfterHours,
		workflow.Status,
		workflow.CreatedAt,
		workflow.UpdatedAt,
	)
	return err
}

func (r *approvalWorkflowRepository) GetByID(ctx context.Context, id string) (*domain.ApprovalWorkflow, error) {
	query := `
		SELECT id, organization_id, name, description, trigger_type, channels, min_approvers, approval_order, auto_approve_after_hours, status, created_at, updated_at
		FROM approval_workflows
		WHERE id = $1
	`
	w := &domain.ApprovalWorkflow{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID,
		&w.OrganizationID,
		&w.Name,
		&w.Description,
		&w.TriggerType,
		&w.Channels,
		&w.MinApprovers,
		&w.ApprovalOrder,
		&w.AutoApproveAfterHours,
		&w.Status,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *approvalWorkflowRepository) List(ctx context.Context, organizationID string, limit, offset int) ([]*domain.ApprovalWorkflow, error) {
	query := `
		SELECT id, organization_id, name, description, trigger_type, channels, min_approvers, approval_order, auto_approve_after_hours, status, created_at, updated_at
		FROM approval_workflows
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := make([]*domain.ApprovalWorkflow, 0)
	for rows.Next() {
		w := &domain.ApprovalWorkflow{}
		err := rows.Scan(
			&w.ID,
			&w.OrganizationID,
			&w.Name,
			&w.Description,
			&w.TriggerType,
			&w.Channels,
			&w.MinApprovers,
			&w.ApprovalOrder,
			&w.AutoApproveAfterHours,
			&w.Status,
			&w.CreatedAt,
			&w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, w)
	}
	return workflows, rows.Err()
}

func (r *approvalWorkflowRepository) Update(ctx context.Context, workflow *domain.ApprovalWorkflow) error {
	query := `
		UPDATE approval_workflows
		SET name = $2, description = $3, trigger_type = $4, channels = $5, min_approvers = $6,
		    approval_order = $7, auto_approve_after_hours = $8, status = $9, updated_at = $10
		WHERE id = $1
	`
	workflow.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		workflow.ID,
		workflow.Name,
		workflow.Description,
		workflow.TriggerType,
		workflow.Channels,
		workflow.MinApprovers,
		workflow.ApprovalOrder,
		workflow.AutoApproveAfterHours,
		workflow.Status,
		workflow.UpdatedAt,
	)
	return err
}

func (r *approvalWorkflowRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM approval_workflows WHERE id = $1`, id)
	return err
}

// ===========================
// Approval Step Repository
// ===========================

type approvalStepRepository struct {
	db *sql.DB
}

func NewApprovalStepRepository(db *sql.DB) domain.ApprovalStepRepository {
	return &approvalStepRepository{db: db}
}

func (r *approvalStepRepository) Create(ctx context.Context, step *domain.ApprovalStep) error {
	query := `
		INSERT INTO approval_steps (id, workflow_id, step_order, approver_role, approver_user, required, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	step.CreatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		step.ID,
		step.WorkflowID,
		step.StepOrder,
		step.ApproverRole,
		step.ApproverUser,
		step.Required,
		step.CreatedAt,
	)
	return err
}

func (r *approvalStepRepository) GetByWorkflowID(ctx context.Context, workflowID string) ([]*domain.ApprovalStep, error) {
	query := `
		SELECT id, workflow_id, step_order, approver_role, approver_user, required, created_at
		FROM approval_steps
		WHERE workflow_id = $1
		ORDER BY step_order ASC
	`
	rows, err := r.db.QueryContext(ctx, query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := make([]*domain.ApprovalStep, 0)
	for rows.Next() {
		s := &domain.ApprovalStep{}
		err := rows.Scan(
			&s.ID,
			&s.WorkflowID,
			&s.StepOrder,
			&s.ApproverRole,
			&s.ApproverUser,
			&s.Required,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}
	return steps, rows.Err()
}

func (r *approvalStepRepository) GetByID(ctx context.Context, id string) (*domain.ApprovalStep, error) {
	query := `
		SELECT id, workflow_id, step_order, approver_role, approver_user, required, created_at
		FROM approval_steps
		WHERE id = $1
	`
	s := &domain.ApprovalStep{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.WorkflowID,
		&s.StepOrder,
		&s.ApproverRole,
		&s.ApproverUser,
		&s.Required,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *approvalStepRepository) DeleteByWorkflowID(ctx context.Context, workflowID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM approval_steps WHERE workflow_id = $1`, workflowID)
	return err
}

// ===========================
// Approval Request Repository
// ===========================

type approvalRequestRepository struct {
	db *sql.DB
}

func NewApprovalRequestRepository(db *sql.DB) domain.ApprovalRequestRepository {
	return &approvalRequestRepository{db: db}
}

func (r *approvalRequestRepository) Create(ctx context.Context, request *domain.ApprovalRequest) error {
	query := `
		INSERT INTO approval_requests (id, organization_id, workflow_id, entity_type, entity_id, title, description, submitted_by, submitted_at, status, deadline, completed_at, comments)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	now := time.Now()
	request.SubmittedAt = now

	var submittedBy, deadline, completedAt interface{}
	if request.SubmittedBy != nil {
		submittedBy = *request.SubmittedBy
	}
	if request.Deadline != nil {
		deadline = *request.Deadline
	}
	if request.CompletedAt != nil {
		completedAt = *request.CompletedAt
	}

	_, err := r.db.ExecContext(ctx, query,
		request.ID,
		request.OrganizationID,
		request.WorkflowID,
		request.EntityType,
		request.EntityID,
		request.Title,
		request.Description,
		submittedBy,
		request.SubmittedAt,
		request.Status,
		deadline,
		completedAt,
		request.Comments,
	)
	return err
}

func (r *approvalRequestRepository) GetByID(ctx context.Context, id string) (*domain.ApprovalRequest, error) {
	query := `
		SELECT id, organization_id, workflow_id, entity_type, entity_id, title, description, submitted_by, submitted_at, status, deadline, completed_at, comments
		FROM approval_requests
		WHERE id = $1
	`
	req := &domain.ApprovalRequest{}
	var submittedBy sql.NullString
	var deadline, completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID,
		&req.OrganizationID,
		&req.WorkflowID,
		&req.EntityType,
		&req.EntityID,
		&req.Title,
		&req.Description,
		&submittedBy,
		&req.SubmittedAt,
		&req.Status,
		&deadline,
		&completedAt,
		&req.Comments,
	)
	if err != nil {
		return nil, err
	}
	if submittedBy.Valid {
		req.SubmittedBy = &submittedBy.String
	}
	if deadline.Valid {
		req.Deadline = &deadline.Time
	}
	if completedAt.Valid {
		req.CompletedAt = &completedAt.Time
	}
	return req, nil
}

func (r *approvalRequestRepository) List(ctx context.Context, organizationID string, status, entityType string, limit, offset int) ([]*domain.ApprovalRequest, error) {
	query := `
		SELECT id, organization_id, workflow_id, entity_type, entity_id, title, description, submitted_by, submitted_at, status, deadline, completed_at, comments
		FROM approval_requests
		WHERE organization_id = $1
	`
	args := []interface{}{organizationID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if entityType != "" {
		query += fmt.Sprintf(" AND entity_type = $%d", argIdx)
		args = append(args, entityType)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY submitted_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]*domain.ApprovalRequest, 0)
	for rows.Next() {
		req := &domain.ApprovalRequest{}
		var submittedBy sql.NullString
		var deadline, completedAt sql.NullTime

		err := rows.Scan(
			&req.ID,
			&req.OrganizationID,
			&req.WorkflowID,
			&req.EntityType,
			&req.EntityID,
			&req.Title,
			&req.Description,
			&submittedBy,
			&req.SubmittedAt,
			&req.Status,
			&deadline,
			&completedAt,
			&req.Comments,
		)
		if err != nil {
			return nil, err
		}
		if submittedBy.Valid {
			req.SubmittedBy = &submittedBy.String
		}
		if deadline.Valid {
			req.Deadline = &deadline.Time
		}
		if completedAt.Valid {
			req.CompletedAt = &completedAt.Time
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (r *approvalRequestRepository) UpdateStatus(ctx context.Context, id string, status string, completedAt *time.Time) error {
	query := `UPDATE approval_requests SET status = $2, completed_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status, completedAt)
	return err
}

// ===========================
// Approval Response Repository
// ===========================

type approvalResponseRepository struct {
	db *sql.DB
}

func NewApprovalResponseRepository(db *sql.DB) domain.ApprovalResponseRepository {
	return &approvalResponseRepository{db: db}
}

func (r *approvalResponseRepository) Create(ctx context.Context, response *domain.ApprovalResponse) error {
	query := `
		INSERT INTO approval_responses (id, request_id, step_id, approver_id, action, comments, responded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	response.RespondedAt = now

	_, err := r.db.ExecContext(ctx, query,
		response.ID,
		response.RequestID,
		response.StepID,
		response.ApproverID,
		response.Action,
		response.Comments,
		response.RespondedAt,
	)
	return err
}

func (r *approvalResponseRepository) ListByRequest(ctx context.Context, requestID string) ([]*domain.ApprovalResponse, error) {
	query := `
		SELECT id, request_id, step_id, approver_id, action, comments, responded_at
		FROM approval_responses
		WHERE request_id = $1
		ORDER BY responded_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	responses := make([]*domain.ApprovalResponse, 0)
	for rows.Next() {
		resp := &domain.ApprovalResponse{}
		err := rows.Scan(
			&resp.ID,
			&resp.RequestID,
			&resp.StepID,
			&resp.ApproverID,
			&resp.Action,
			&resp.Comments,
			&resp.RespondedAt,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, rows.Err()
}

func (r *approvalResponseRepository) CountByRequestAndAction(ctx context.Context, requestID string, action string) (int, error) {
	query := `SELECT COUNT(*) FROM approval_responses WHERE request_id = $1 AND action = $2`
	var count int
	err := r.db.QueryRowContext(ctx, query, requestID, action).Scan(&count)
	return count, err
}

func (r *approvalResponseRepository) ExistsByRequestAndApprover(ctx context.Context, requestID string, approverID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM approval_responses WHERE request_id = $1 AND approver_id = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, requestID, approverID).Scan(&exists)
	return exists, err
}
