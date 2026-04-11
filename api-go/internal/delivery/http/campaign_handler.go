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

// CampaignHandler handles HTTP requests for campaigns
type CampaignHandler struct {
	service *service.CampaignService
}

// NewCampaignHandler creates a new campaign handler
func NewCampaignHandler(service *service.CampaignService) *CampaignHandler {
	return &CampaignHandler{service: service}
}

// Default organization ID for demo purposes
const defaultOrganizationID = "00000000-0000-0000-0000-000000000001"

// CreateCampaign handles POST /api/v1/campaigns
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req dto.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM-DD"})
		return
	}

	campaign := &domain.Campaign{
		ID:             uuid.New().String(),
		OrganizationID: defaultOrganizationID,
		Name:           req.Name,
		Type:           req.Type,
		Status:         "draft",
		Budget:         req.Budget,
		Spent:          0,
		Impressions:    0,
		Clicks:         0,
		Conversions:    0,
		StartDate:      startDate,
	}

	if err := h.service.CreateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// GetCampaign handles GET /api/v1/campaigns/:id
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	id := c.Param("id")

	campaign, err := h.service.GetCampaign(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	response := dto.CampaignResponse{
		ID:                campaign.ID,
		Name:              campaign.Name,
		Type:              campaign.Type,
		Status:            campaign.Status,
		Budget:            campaign.Budget,
		Spent:             campaign.Spent,
		Impressions:       campaign.Impressions,
		Clicks:            campaign.Clicks,
		Conversions:       campaign.Conversions,
		CTR:               campaign.CTR(),
		ConversionRate:    campaign.ConversionRate(),
		BudgetUtilization: campaign.BudgetUtilization(),
		StartDate:         campaign.StartDate.Format("2006-01-02"),
		CreatedAt:         campaign.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         campaign.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if campaign.EndDate != nil {
		endDate := campaign.EndDate.Format("2006-01-02")
		response.EndDate = &endDate
	}

	c.JSON(http.StatusOK, response)
}

// GetAllCampaigns handles GET /api/v1/campaigns
func (h *CampaignHandler) GetAllCampaigns(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	campaigns, err := h.service.GetAllCampaigns(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.CampaignResponse, len(campaigns))
	for i, campaign := range campaigns {
		response[i] = dto.CampaignResponse{
			ID:                campaign.ID,
			Name:              campaign.Name,
			Type:              campaign.Type,
			Status:            campaign.Status,
			Budget:            campaign.Budget,
			Spent:             campaign.Spent,
			Impressions:       campaign.Impressions,
			Clicks:            campaign.Clicks,
			Conversions:       campaign.Conversions,
			CTR:               campaign.CTR(),
			ConversionRate:    campaign.ConversionRate(),
			BudgetUtilization: campaign.BudgetUtilization(),
			StartDate:         campaign.StartDate.Format("2006-01-02"),
			CreatedAt:         campaign.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         campaign.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": response, "total": len(campaigns)})
}

// UpdateCampaign handles PUT /api/v1/campaigns/:id
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	id := c.Param("id")

	campaign, err := h.service.GetCampaign(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	var req dto.UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Name != "" {
		campaign.Name = req.Name
	}
	if req.Type != "" {
		campaign.Type = req.Type
	}
	if req.Status != "" {
		campaign.Status = req.Status
	}
	if req.Budget > 0 {
		campaign.Budget = req.Budget
	}

	campaign.Spent = req.Spent
	campaign.Impressions = req.Impressions
	campaign.Clicks = req.Clicks
	campaign.Conversions = req.Conversions

	if err := h.service.UpdateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// DeleteCampaign handles DELETE /api/v1/campaigns/:id
func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}

// PauseCampaign handles POST /api/v1/campaigns/:id/pause
func (h *CampaignHandler) PauseCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.PauseCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign paused successfully"})
}

// ActivateCampaign handles POST /api/v1/campaigns/:id/activate
func (h *CampaignHandler) ActivateCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.ActivateCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign activated successfully"})
}
