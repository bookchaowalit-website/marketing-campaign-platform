package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/dto"
	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ApprovalHandler handles HTTP requests for approval workflows
type ApprovalHandler struct {
	service *service.ApprovalService
}

// NewApprovalHandler creates a new approval handler
func NewApprovalHandler(svc *service.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{service: svc}
}

// CreateWorkflow handles POST /api/v1/approval/workflows
func (h *ApprovalHandler) CreateWorkflow(c *gin.Context) {
	var req dto.CreateApprovalWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	approvalOrder := req.ApprovalOrder
	if approvalOrder == "" {
		approvalOrder = "parallel"
	}

	minApprovers := req.MinApprovers
	if minApprovers < 1 {
		minApprovers = 1
	}

	workflow := &domain.ApprovalWorkflow{
		ID:                    uuid.New().String(),
		OrganizationID:        defaultOrganizationID,
		Name:                  req.Name,
		Description:           req.Description,
		TriggerType:           req.TriggerType,
		Channels:              req.Channels,
		MinApprovers:          minApprovers,
		ApprovalOrder:         approvalOrder,
		AutoApproveAfterHours: req.AutoApproveAfterHours,
		Status:                status,
	}

	if err := h.service.CreateWorkflow(c.Request.Context(), workflow, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ApprovalWorkflowResponse{
		ID:                    workflow.ID,
		OrganizationID:        workflow.OrganizationID,
		Name:                  workflow.Name,
		Description:           workflow.Description,
		TriggerType:           workflow.TriggerType,
		Channels:              workflow.Channels,
		MinApprovers:          workflow.MinApprovers,
		ApprovalOrder:         workflow.ApprovalOrder,
		AutoApproveAfterHours: workflow.AutoApproveAfterHours,
		Status:                workflow.Status,
		CreatedAt:             workflow.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             workflow.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// ListWorkflows handles GET /api/v1/approval/workflows
func (h *ApprovalHandler) ListWorkflows(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	workflows, err := h.service.ListWorkflows(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ApprovalWorkflowResponse, len(workflows))
	for i, w := range workflows {
		response[i] = dto.ApprovalWorkflowResponse{
			ID:                    w.ID,
			OrganizationID:        w.OrganizationID,
			Name:                  w.Name,
			Description:           w.Description,
			TriggerType:           w.TriggerType,
			Channels:              w.Channels,
			MinApprovers:          w.MinApprovers,
			ApprovalOrder:         w.ApprovalOrder,
			AutoApproveAfterHours: w.AutoApproveAfterHours,
			Status:                w.Status,
			CreatedAt:             w.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:             w.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"workflows": response, "total": len(workflows)})
}

// CreateRequest handles POST /api/v1/approval/requests
func (h *ApprovalHandler) CreateRequest(c *gin.Context) {
	var req dto.CreateApprovalRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request := &domain.ApprovalRequest{
		ID:             uuid.New().String(),
		OrganizationID: defaultOrganizationID,
		WorkflowID:     req.WorkflowID,
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		Title:          req.Title,
		Description:    req.Description,
		Status:         "pending",
	}

	// Parse deadline if provided
	if req.Deadline != "" {
		deadline, err := time.Parse("2006-01-02T15:04:05Z", req.Deadline)
		if err != nil {
			deadline, err = time.Parse("2006-01-02 15:04:05", req.Deadline)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format, use ISO 8601 or YYYY-MM-DD HH:MM:SS"})
				return
			}
		}
		request.Deadline = &deadline
	}

	// Get user ID from context (set by auth middleware)
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(string)
		request.SubmittedBy = &uid
	}

	if err := h.service.CreateRequest(c.Request.Context(), request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ApprovalRequestResponse{
		ID:             request.ID,
		OrganizationID: request.OrganizationID,
		WorkflowID:     request.WorkflowID,
		EntityType:     request.EntityType,
		EntityID:       request.EntityID,
		Title:          request.Title,
		Description:    request.Description,
		SubmittedBy:    request.SubmittedBy,
		SubmittedAt:    request.SubmittedAt.Format("2006-01-02 15:04:05"),
		Status:         request.Status,
	}

	if request.Deadline != nil {
		deadlineStr := request.Deadline.Format("2006-01-02 15:04:05")
		resp.Deadline = &deadlineStr
	}

	c.JSON(http.StatusCreated, resp)
}

// ListRequests handles GET /api/v1/approval/requests
func (h *ApprovalHandler) ListRequests(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")
	entityType := c.Query("entity_type")

	requests, err := h.service.ListRequests(c.Request.Context(), defaultOrganizationID, status, entityType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ApprovalRequestResponse, len(requests))
	for i, req := range requests {
		r := dto.ApprovalRequestResponse{
			ID:             req.ID,
			OrganizationID: req.OrganizationID,
			WorkflowID:     req.WorkflowID,
			EntityType:     req.EntityType,
			EntityID:       req.EntityID,
			Title:          req.Title,
			Description:    req.Description,
			SubmittedBy:    req.SubmittedBy,
			SubmittedAt:    req.SubmittedAt.Format("2006-01-02 15:04:05"),
			Status:         req.Status,
			Comments:       req.Comments,
		}

		if req.Deadline != nil {
			deadlineStr := req.Deadline.Format("2006-01-02 15:04:05")
			r.Deadline = &deadlineStr
		}
		if req.CompletedAt != nil {
			completedStr := req.CompletedAt.Format("2006-01-02 15:04:05")
			r.CompletedAt = &completedStr
		}

		response[i] = r
	}

	c.JSON(http.StatusOK, gin.H{"requests": response, "total": len(requests)})
}

// GetRequest handles GET /api/v1/approval/requests/:id
func (h *ApprovalHandler) GetRequest(c *gin.Context) {
	id := c.Param("id")

	request, responses, err := h.service.GetRequest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approval request not found"})
		return
	}

	resp := dto.ApprovalRequestResponse{
		ID:             request.ID,
		OrganizationID: request.OrganizationID,
		WorkflowID:     request.WorkflowID,
		EntityType:     request.EntityType,
		EntityID:       request.EntityID,
		Title:          request.Title,
		Description:    request.Description,
		SubmittedBy:    request.SubmittedBy,
		SubmittedAt:    request.SubmittedAt.Format("2006-01-02 15:04:05"),
		Status:         request.Status,
		Comments:       request.Comments,
	}

	if request.Deadline != nil {
		deadlineStr := request.Deadline.Format("2006-01-02 15:04:05")
		resp.Deadline = &deadlineStr
	}
	if request.CompletedAt != nil {
		completedStr := request.CompletedAt.Format("2006-01-02 15:04:05")
		resp.CompletedAt = &completedStr
	}

	// Include responses
	if len(responses) > 0 {
		resp.Responses = make([]dto.ApprovalResponseResponse, len(responses))
		for i, r := range responses {
			resp.Responses[i] = dto.ApprovalResponseResponse{
				ID:          r.ID,
				RequestID:   r.RequestID,
				StepID:      r.StepID,
				ApproverID:  r.ApproverID,
				Action:      r.Action,
				Comments:    r.Comments,
				RespondedAt: r.RespondedAt.Format("2006-01-02 15:04:05"),
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

// Respond handles POST /api/v1/approval/requests/:id/respond
func (h *ApprovalHandler) Respond(c *gin.Context) {
	id := c.Param("id")

	var req dto.RespondToRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get approver ID from context (set by auth middleware)
	approverID := ""
	if uid, exists := c.Get("user_id"); exists {
		approverID = uid.(string)
	}

	if approverID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to determine approver identity"})
		return
	}

	updatedRequest, err := h.service.Respond(c.Request.Context(), id, approverID, req.Action, req.Comments)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ApprovalRequestResponse{
		ID:             updatedRequest.ID,
		OrganizationID: updatedRequest.OrganizationID,
		WorkflowID:     updatedRequest.WorkflowID,
		EntityType:     updatedRequest.EntityType,
		EntityID:       updatedRequest.EntityID,
		Title:          updatedRequest.Title,
		Description:    updatedRequest.Description,
		SubmittedBy:    updatedRequest.SubmittedBy,
		SubmittedAt:    updatedRequest.SubmittedAt.Format("2006-01-02 15:04:05"),
		Status:         updatedRequest.Status,
		Comments:       updatedRequest.Comments,
	}

	if updatedRequest.Deadline != nil {
		deadlineStr := updatedRequest.Deadline.Format("2006-01-02 15:04:05")
		resp.Deadline = &deadlineStr
	}
	if updatedRequest.CompletedAt != nil {
		completedStr := updatedRequest.CompletedAt.Format("2006-01-02 15:04:05")
		resp.CompletedAt = &completedStr
	}

	c.JSON(http.StatusOK, resp)
}
