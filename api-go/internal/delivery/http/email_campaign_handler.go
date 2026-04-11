package http

import (
	"net/http"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/gin-gonic/gin"
)

// EmailCampaignHandler handles email campaign endpoints
type EmailCampaignHandler struct {
	emailCampaignService *service.EmailCampaignService
}

// NewEmailCampaignHandler creates a new email campaign handler
func NewEmailCampaignHandler(svc *service.EmailCampaignService) *EmailCampaignHandler {
	return &EmailCampaignHandler{emailCampaignService: svc}
}

// ===========================
// REQUEST STRUCTS
// ===========================

type CreateABTestRequest struct {
	CampaignIDA       string `json:"campaign_id_a" binding:"required"`
	CampaignIDB       string `json:"campaign_id_b" binding:"required"`
	TestMetric        string `json:"test_metric"`
	SampleSizePercent int    `json:"sample_size_percent"`
}

type TrackEventRequest struct {
	SendID     string            `json:"send_id" binding:"required"`
	EventType  string            `json:"event_type" binding:"required"`
	Metadata   map[string]string `json:"metadata"`
}

// ===========================
// CAMPAIGN SEND
// ===========================

// SendCampaignEmail handles POST /api/v1/email/campaigns/:id/send
func (h *EmailCampaignHandler) SendCampaignEmail(c *gin.Context) {
	id := c.Param("id")

	if err := h.emailCampaignService.SendEmailCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign sending initiated", "campaign_id": id})
}

// ===========================
// CAMPAIGN ANALYTICS
// ===========================

// GetCampaignAnalytics handles GET /api/v1/email/campaigns/:id/analytics
func (h *EmailCampaignHandler) GetCampaignAnalytics(c *gin.Context) {
	id := c.Param("id")

	analytics, err := h.emailCampaignService.GetCampaignAnalytics(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// ===========================
// A/B TESTING
// ===========================

// CreateABTest handles POST /api/v1/email/ab-tests
func (h *EmailCampaignHandler) CreateABTest(c *gin.Context) {
	var req CreateABTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	abTest, err := h.emailCampaignService.CreateABTest(c.Request.Context(), service.CreateABTestRequest{
		CampaignIDA:       req.CampaignIDA,
		CampaignIDB:       req.CampaignIDB,
		TestMetric:        req.TestMetric,
		SampleSizePercent: req.SampleSizePercent,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, abTest)
}

// ListABTests handles GET /api/v1/email/ab-tests
func (h *EmailCampaignHandler) ListABTests(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	tests, err := h.emailCampaignService.ListABTests(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ab_tests": tests, "total": len(tests)})
}

// CompleteABTest handles POST /api/v1/email/ab-tests/:id/complete
func (h *EmailCampaignHandler) CompleteABTest(c *gin.Context) {
	id := c.Param("id")

	abTest, err := h.emailCampaignService.CompleteABTest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "A/B test completed",
		"winner":  abTest.Winner,
		"test":    abTest,
	})
}

// ===========================
// EMAIL TRACKING (Pixel Webhook)
// ===========================

// TrackEvent handles POST /api/v1/email/track
func (h *EmailCampaignHandler) TrackEvent(c *gin.Context) {
	var req TrackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailCampaignService.TrackEmailEvent(c.Request.Context(), req.SendID, req.EventType, req.Metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event tracked"})
}

// TrackOpen handles GET /api/v1/email/track/open (for tracking pixel)
func (h *EmailCampaignHandler) TrackOpen(c *gin.Context) {
	sendID := c.Query("send_id")
	if sendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "send_id is required"})
		return
	}

	metadata := map[string]string{
		"ip_address":  c.ClientIP(),
		"user_agent":  c.GetHeader("User-Agent"),
	}

	if err := h.emailCampaignService.TrackEmailEvent(c.Request.Context(), sendID, "opened", metadata); err != nil {
		// Don't fail the request even if tracking fails (broken pixel)
	}

	// Return 1x1 transparent GIF
	gifData := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00,
		0x80, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x21,
		0xF9, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2C, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
		0x01, 0x00, 0x3B,
	}
	c.Data(http.StatusOK, "image/gif", gifData)
}

// TrackClick handles GET /api/v1/email/track/click (for link tracking)
func (h *EmailCampaignHandler) TrackClick(c *gin.Context) {
	sendID := c.Query("send_id")
	redirectURL := c.Query("url")

	if sendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "send_id is required"})
		return
	}

	metadata := map[string]string{
		"url": redirectURL,
	}

	if err := h.emailCampaignService.TrackEmailEvent(c.Request.Context(), sendID, "clicked", metadata); err != nil {
		// Don't fail redirect even if tracking fails
	}

	if redirectURL != "" {
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
}
