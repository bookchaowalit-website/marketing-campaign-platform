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

// BudgetHandler handles HTTP requests for budget planning
type BudgetHandler struct {
	service *service.BudgetService
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(service *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// CreateBudgetPeriod handles POST /api/v1/budget/periods
func (h *BudgetHandler) CreateBudgetPeriod(c *gin.Context) {
	var req dto.CreateBudgetPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use YYYY-MM-DD"})
		return
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}

	period := &domain.BudgetPeriod{
		ID:             uuid.New().String(),
		OrganizationID: defaultOrganizationID,
		Name:           req.Name,
		PeriodType:     req.PeriodType,
		StartDate:      startDate,
		EndDate:        endDate,
		TotalBudget:    req.TotalBudget,
		Status:         status,
		Notes:          req.Notes,
	}

	if err := h.service.CreatePeriod(c.Request.Context(), period); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.BudgetPeriodResponse{
		ID:             period.ID,
		Name:           period.Name,
		PeriodType:     period.PeriodType,
		StartDate:      period.StartDate.Format("2006-01-02"),
		EndDate:        period.EndDate.Format("2006-01-02"),
		TotalBudget:    period.TotalBudget,
		TotalAllocated: period.TotalAllocated,
		TotalSpent:     period.TotalSpent,
		Status:         period.Status,
		Notes:          period.Notes,
		CreatedAt:      period.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      period.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// ListBudgetPeriods handles GET /api/v1/budget/periods
func (h *BudgetHandler) ListBudgetPeriods(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	periods, err := h.service.ListPeriods(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.BudgetPeriodResponse, len(periods))
	for i, p := range periods {
		response[i] = dto.BudgetPeriodResponse{
			ID:             p.ID,
			Name:           p.Name,
			PeriodType:     p.PeriodType,
			StartDate:      p.StartDate.Format("2006-01-02"),
			EndDate:        p.EndDate.Format("2006-01-02"),
			TotalBudget:    p.TotalBudget,
			TotalAllocated: p.TotalAllocated,
			TotalSpent:     p.TotalSpent,
			Status:         p.Status,
			Notes:          p.Notes,
			CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      p.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"periods": response, "total": len(periods)})
}

// GetBudgetPeriod handles GET /api/v1/budget/periods/:id
func (h *BudgetHandler) GetBudgetPeriod(c *gin.Context) {
	id := c.Param("id")

	period, err := h.service.GetPeriod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget period not found"})
		return
	}

	allocs, err := h.service.ListAllocations(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allocResponses := make([]dto.BudgetAllocationResponse, len(allocs))
	for i, a := range allocs {
		allocResponses[i] = dto.BudgetAllocationResponse{
			ID:              a.ID,
			PeriodID:        a.PeriodID,
			Channel:         a.Channel,
			ChannelCategory: a.ChannelCategory,
			AllocatedAmount: a.AllocatedAmount,
			SpentAmount:     a.SpentAmount,
			ForecastAmount:  a.ForecastAmount,
			Notes:           a.Notes,
			CreatedAt:       a.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       a.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.BudgetPeriodWithAllocationsResponse{
		ID:             period.ID,
		Name:           period.Name,
		PeriodType:     period.PeriodType,
		StartDate:      period.StartDate.Format("2006-01-02"),
		EndDate:        period.EndDate.Format("2006-01-02"),
		TotalBudget:    period.TotalBudget,
		TotalAllocated: period.TotalAllocated,
		TotalSpent:     period.TotalSpent,
		Status:         period.Status,
		Notes:          period.Notes,
		Allocations:    allocResponses,
		CreatedAt:      period.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      period.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// CreateBudgetAllocation handles POST /api/v1/budget/periods/:id/allocations
func (h *BudgetHandler) CreateBudgetAllocation(c *gin.Context) {
	periodID := c.Param("id")

	var req dto.CreateBudgetAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alloc := &domain.BudgetAllocation{
		ID:              uuid.New().String(),
		PeriodID:        periodID,
		Channel:         req.Channel,
		ChannelCategory: req.ChannelCategory,
		AllocatedAmount: req.AllocatedAmount,
		SpentAmount:     0,
		ForecastAmount:  req.ForecastAmount,
		Notes:           req.Notes,
	}

	if err := h.service.CreateAllocation(c.Request.Context(), alloc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.BudgetAllocationResponse{
		ID:              alloc.ID,
		PeriodID:        alloc.PeriodID,
		Channel:         alloc.Channel,
		ChannelCategory: alloc.ChannelCategory,
		AllocatedAmount: alloc.AllocatedAmount,
		SpentAmount:     alloc.SpentAmount,
		ForecastAmount:  alloc.ForecastAmount,
		Notes:           alloc.Notes,
		CreatedAt:       alloc.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       alloc.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// UpdateBudgetAllocation handles PUT /api/v1/budget/allocations/:id
func (h *BudgetHandler) UpdateBudgetAllocation(c *gin.Context) {
	id := c.Param("id")

	alloc, err := h.service.GetAllocation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Allocation not found"})
		return
	}

	var req dto.UpdateBudgetAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Channel != "" {
		alloc.Channel = req.Channel
	}
	if req.ChannelCategory != "" {
		alloc.ChannelCategory = req.ChannelCategory
	}
	if req.AllocatedAmount > 0 {
		alloc.AllocatedAmount = req.AllocatedAmount
	}
	if req.SpentAmount > 0 {
		alloc.SpentAmount = req.SpentAmount
	}
	if req.ForecastAmount != nil {
		alloc.ForecastAmount = req.ForecastAmount
	}
	if req.Notes != "" {
		alloc.Notes = req.Notes
	}

	if err := h.service.UpdateAllocation(c.Request.Context(), alloc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.BudgetAllocationResponse{
		ID:              alloc.ID,
		PeriodID:        alloc.PeriodID,
		Channel:         alloc.Channel,
		ChannelCategory: alloc.ChannelCategory,
		AllocatedAmount: alloc.AllocatedAmount,
		SpentAmount:     alloc.SpentAmount,
		ForecastAmount:  alloc.ForecastAmount,
		Notes:           alloc.Notes,
		CreatedAt:       alloc.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       alloc.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// RecordTransaction handles POST /api/v1/budget/transactions
func (h *BudgetHandler) RecordTransaction(c *gin.Context) {
	var req dto.RecordTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	txn := &domain.BudgetTransaction{
		ID:             uuid.New().String(),
		OrganizationID: defaultOrganizationID,
		CampaignID:     req.CampaignID,
		AllocationID:   req.AllocationID,
		Channel:        req.Channel,
		Amount:         req.Amount,
		Currency:       currency,
		TransactionDate: time.Now(),
		Description:    req.Description,
		InvoiceNumber:  req.InvoiceNumber,
		Vendor:         req.Vendor,
	}

	if err := h.service.RecordTransaction(c.Request.Context(), txn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.BudgetTransactionResponse{
		ID:             txn.ID,
		OrganizationID: txn.OrganizationID,
		CampaignID:     txn.CampaignID,
		AllocationID:   txn.AllocationID,
		Channel:        txn.Channel,
		Amount:         txn.Amount,
		Currency:       txn.Currency,
		TransactionDate: txn.TransactionDate.Format("2006-01-02 15:04:05"),
		Description:    txn.Description,
		InvoiceNumber:  txn.InvoiceNumber,
		Vendor:         txn.Vendor,
		CreatedAt:      txn.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// ListTransactions handles GET /api/v1/budget/transactions
func (h *BudgetHandler) ListTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	txns, err := h.service.ListTransactions(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.BudgetTransactionResponse, len(txns))
	for i, t := range txns {
		response[i] = dto.BudgetTransactionResponse{
			ID:             t.ID,
			OrganizationID: t.OrganizationID,
			CampaignID:     t.CampaignID,
			AllocationID:   t.AllocationID,
			Channel:        t.Channel,
			Amount:         t.Amount,
			Currency:       t.Currency,
			TransactionDate: t.TransactionDate.Format("2006-01-02 15:04:05"),
			Description:    t.Description,
			InvoiceNumber:  t.InvoiceNumber,
			Vendor:         t.Vendor,
			CreatedAt:      t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"transactions": response, "total": len(txns)})
}

// ListAlerts handles GET /api/v1/budget/alerts
func (h *BudgetHandler) ListAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Also check for new triggered alerts
	h.service.CheckAlerts(c.Request.Context(), defaultOrganizationID)

	alerts, err := h.service.ListAlerts(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.BudgetAlertResponse, len(alerts))
	for i, a := range alerts {
		response[i] = dto.BudgetAlertResponse{
			ID:               a.ID,
			OrganizationID:   a.OrganizationID,
			PeriodID:         a.PeriodID,
			CampaignID:       a.CampaignID,
			AlertType:        a.AlertType,
			ThresholdPercent: a.ThresholdPercent,
			Triggered:        a.Triggered,
			TriggeredAt:      a.TriggeredAt,
			NotifiedUsers:    a.NotifiedUsers,
			CreatedAt:        a.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"alerts": response, "total": len(alerts)})
}
