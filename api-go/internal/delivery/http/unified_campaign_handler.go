package http

import (
	"net/http"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type UnifiedCampaignHandler struct {
	svc *service.UnifiedCampaignService
}

func NewUnifiedCampaignHandler(svc *service.UnifiedCampaignService) *UnifiedCampaignHandler {
	return &UnifiedCampaignHandler{svc: svc}
}

type CreateCampaignRequest struct {
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	ChannelCategory string   `json:"channel_category" binding:"required,oneof=online offline hybrid"`
	Channel         string   `json:"channel" binding:"required"`
	Budget          float64  `json:"budget"`
	TargetAudience  string   `json:"target_audience"`
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	Timezone        string   `json:"timezone"`
	TrackingURL     string   `json:"tracking_url"`
	CreativeBrief   string   `json:"creative_brief"`
	GroupID         string   `json:"group_id"`
}

// CreateCampaign handles POST /api/v1/unified-campaigns
func (h *UnifiedCampaignHandler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	userID := c.GetString("user_id")

	campaign := &domain.UnifiedCampaign{
		OrganizationID:  orgID,
		Name:            req.Name,
		Description:     req.Description,
		ChannelCategory: req.ChannelCategory,
		Channel:         req.Channel,
		Budget:          req.Budget,
		TargetAudience:  req.TargetAudience,
		Status:          "draft",
		Timezone:        req.Timezone,
		TrackingURL:     req.TrackingURL,
		CreativeBrief:   req.CreativeBrief,
		OwnerID:         userID,
		CreatedBy:       userID,
	}

	if req.GroupID != "" {
		campaign.GroupID = &req.GroupID
	}

	if err := h.svc.CreateCampaign(c.Request.Context(), campaign); err != nil {
		if err == domain.ErrInvalidChannel || err == domain.ErrInvalidChannelCategory {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// GetCampaign handles GET /api/v1/unified-campaigns/:id
func (h *UnifiedCampaignHandler) GetCampaign(c *gin.Context) {
	id := c.Param("id")
	campaign, err := h.svc.GetCampaign(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}
	c.JSON(http.StatusOK, campaign)
}

// ListCampaigns handles GET /api/v1/unified-campaigns
func (h *UnifiedCampaignHandler) ListCampaigns(c *gin.Context) {
	orgID := c.GetString("organization_id")
	channel := c.Query("channel")
	status := c.Query("status")
	limit := 50
	offset := 0

	campaigns, total, err := h.svc.ListCampaigns(c.Request.Context(), orgID, channel, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"campaigns": campaigns,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateCampaign handles PUT /api/v1/unified-campaigns/:id
func (h *UnifiedCampaignHandler) UpdateCampaign(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.svc.GetCampaign(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Budget = req.Budget
	existing.TargetAudience = req.TargetAudience
	existing.TrackingURL = req.TrackingURL
	existing.CreativeBrief = req.CreativeBrief

	if err := h.svc.UpdateCampaign(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteCampaign handles DELETE /api/v1/unified-campaigns/:id
func (h *UnifiedCampaignHandler) DeleteCampaign(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted"})
}

// SubmitForApproval handles POST /api/v1/unified-campaigns/:id/submit
func (h *UnifiedCampaignHandler) SubmitForApproval(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.SubmitForApproval(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Submitted for approval"})
}

// ApproveCampaign handles POST /api/v1/unified-campaigns/:id/approve
func (h *UnifiedCampaignHandler) ApproveCampaign(c *gin.Context) {
	id := c.Param("id")
	approverID := c.GetString("user_id")
	if err := h.svc.ApproveCampaign(c.Request.Context(), id, approverID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Campaign approved"})
}

// RejectCampaign handles POST /api/v1/unified-campaigns/:id/reject
func (h *UnifiedCampaignHandler) RejectCampaign(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)
	if err := h.svc.RejectCampaign(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Campaign rejected"})
}
