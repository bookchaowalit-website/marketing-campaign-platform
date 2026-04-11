package http

import (
	"net/http"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type PPCAdsHandler struct {
	adsService *service.PPCAdsService
}

func NewPPCAdsHandler(adsService *service.PPCAdsService) *PPCAdsHandler {
	return &PPCAdsHandler{adsService: adsService}
}

// ===========================
// REQUEST STRUCTS
// ===========================

type ConnectPlatformRequest struct {
	Platform  string `json:"platform" binding:"required,oneof=google_ads facebook_ads tiktok_ads linkedin_ads twitter_ads"`
	AccountID string `json:"account_id" binding:"required"`
}

type CreateAdCampaignRequest struct {
	Name              string                 `json:"name" binding:"required"`
	PlatformID        string                 `json:"platform_id" binding:"required"`
	UnifiedCampaignID *string                `json:"unified_campaign_id,omitempty"`
	Budget            float64                `json:"budget"`
	StartDate         *string                `json:"start_date"`
	EndDate           *string                `json:"end_date"`
	Targeting         map[string]interface{} `json:"targeting,omitempty"`
}

// ===========================
// PLATFORM HANDLERS
// ===========================

// ConnectPlatform handles POST /api/v1/ads/platforms/connect
func (h *PPCAdsHandler) ConnectPlatform(c *gin.Context) {
	var req ConnectPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	platform := &domain.AdPlatform{
		OrganizationID: orgID,
		Platform:       req.Platform,
		AccountID:      req.AccountID,
	}

	if err := h.adsService.ConnectPlatform(c.Request.Context(), platform); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, platform)
}

// ListPlatforms handles GET /api/v1/ads/platforms
func (h *PPCAdsHandler) ListPlatforms(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	platforms, err := h.adsService.ListPlatforms(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"platforms": platforms, "total": len(platforms)})
}

// ===========================
// CAMPAIGN HANDLERS
// ===========================

// CreateCampaign handles POST /api/v1/ads/campaigns
func (h *PPCAdsHandler) CreateCampaign(c *gin.Context) {
	var req CreateAdCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	campaign := &domain.AdCampaign{
		OrganizationID:    orgID,
		Name:              req.Name,
		PlatformID:        req.PlatformID,
		UnifiedCampaignID: req.UnifiedCampaignID,
		Budget:            req.Budget,
		Status:            "draft",
		Targeting:         req.Targeting,
	}

	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			campaign.StartDate = &t
		}
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			campaign.EndDate = &t
		}
	}

	if err := h.adsService.CreateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// ListCampaigns handles GET /api/v1/ads/campaigns
func (h *PPCAdsHandler) ListCampaigns(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	status := c.Query("status")
	limit := 50
	offset := 0

	campaigns, total, err := h.adsService.ListCampaigns(c.Request.Context(), orgID, status, limit, offset)
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

// GetCampaign handles GET /api/v1/ads/campaigns/:id
func (h *PPCAdsHandler) GetCampaign(c *gin.Context) {
	id := c.Param("id")

	campaign, err := h.adsService.GetCampaign(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	// Get performance data
	startDate := time.Now().AddDate(0, 0, -30) // Last 30 days
	endDate := time.Now()
	performance, _ := h.adsService.GetCampaignPerformance(c.Request.Context(), id, startDate, endDate)

	result := gin.H{
		"campaign":    campaign,
		"performance": performance,
	}

	c.JSON(http.StatusOK, result)
}

// SyncCampaign handles POST /api/v1/ads/campaigns/:id/sync
func (h *PPCAdsHandler) SyncCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.adsService.SyncCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign sync initiated"})
}

// PauseCampaign handles POST /api/v1/ads/campaigns/:id/pause
func (h *PPCAdsHandler) PauseCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.adsService.PauseCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign paused"})
}

// ActivateCampaign handles POST /api/v1/ads/campaigns/:id/activate
func (h *PPCAdsHandler) ActivateCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.adsService.ActivateCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign activated"})
}
