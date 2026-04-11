package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/gin-gonic/gin"
)

// MarketingFeaturesHandler handles HTTP requests for marketing features
type MarketingFeaturesHandler struct {
	emailService      *service.EmailMarketingService
	socialService     *service.SocialMediaService
	landingService    *service.LandingPageService
	automationService *service.AutomationService
	influencerService *service.InfluencerService
	reportingService  *service.ReportingService
}

// NewMarketingFeaturesHandler creates a new marketing features handler
func NewMarketingFeaturesHandler(
	emailService *service.EmailMarketingService,
	socialService *service.SocialMediaService,
	landingService *service.LandingPageService,
	automationService *service.AutomationService,
	influencerService *service.InfluencerService,
	reportingService *service.ReportingService,
) *MarketingFeaturesHandler {
	return &MarketingFeaturesHandler{
		emailService:      emailService,
		socialService:     socialService,
		landingService:    landingService,
		automationService: automationService,
		influencerService: influencerService,
		reportingService:  reportingService,
	}
}

// =====================================================
// EMAIL MARKETING HANDLERS
// =====================================================

// CreateEmailList handles POST /api/v1/email/lists
func (h *MarketingFeaturesHandler) CreateEmailList(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list, err := h.emailService.CreateList(c.Request.Context(), defaultOrganizationID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, list)
}

// GetEmailLists handles GET /api/v1/email/lists
func (h *MarketingFeaturesHandler) GetEmailLists(c *gin.Context) {
	lists, err := h.emailService.GetLists(c.Request.Context(), defaultOrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lists": lists, "total": len(lists)})
}

// AddSubscriber handles POST /api/v1/email/subscribers
func (h *MarketingFeaturesHandler) AddSubscriber(c *gin.Context) {
	var req struct {
		Email     string   `json:"email" binding:"required"`
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Source    string   `json:"source"`
		ListIDs   []string `json:"list_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscriber, err := h.emailService.AddSubscriber(c.Request.Context(), defaultOrganizationID, req.Email, req.FirstName, req.LastName, req.Source, req.ListIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscriber)
}

// GetSubscribers handles GET /api/v1/email/subscribers
func (h *MarketingFeaturesHandler) GetSubscribers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	subscribers, total, err := h.emailService.GetSubscribers(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscribers": subscribers, "total": total})
}

// CreateEmailCampaign handles POST /api/v1/email/campaigns
func (h *MarketingFeaturesHandler) CreateEmailCampaign(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Subject     string   `json:"subject" binding:"required"`
		FromName    string   `json:"from_name"`
		FromEmail   string   `json:"from_email"`
		HTMLContent string   `json:"html_content"`
		ListIDs     []string `json:"list_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := h.emailService.CreateCampaign(c.Request.Context(), defaultOrganizationID, req.Name, req.Subject, req.FromName, req.FromEmail, req.HTMLContent, req.ListIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// GetEmailCampaigns handles GET /api/v1/email/campaigns
func (h *MarketingFeaturesHandler) GetEmailCampaigns(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	campaigns, total, err := h.emailService.GetCampaigns(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns, "total": total})
}

// =====================================================
// SOCIAL MEDIA HANDLERS
// =====================================================

// CreateSocialPost handles POST /api/v1/social/posts
func (h *MarketingFeaturesHandler) CreateSocialPost(c *gin.Context) {
	var req struct {
		Content   string   `json:"content" binding:"required"`
		MediaURLs []string `json:"media_urls"`
		Hashtags  []string `json:"hashtags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.socialService.CreatePost(c.Request.Context(), defaultOrganizationID, req.Content, req.MediaURLs, req.Hashtags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// ScheduleSocialPost handles POST /api/v1/social/posts/:id/schedule
func (h *MarketingFeaturesHandler) ScheduleSocialPost(c *gin.Context) {
	postID := c.Param("id")

	var req struct {
		AccountID     string `json:"account_id" binding:"required"`
		ScheduledTime string `json:"scheduled_time" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled_time format"})
		return
	}

	if err := h.socialService.SchedulePost(c.Request.Context(), postID, req.AccountID, scheduledTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post scheduled successfully"})
}

// GetSocialPosts handles GET /api/v1/social/posts
func (h *MarketingFeaturesHandler) GetSocialPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	posts, err := h.socialService.GetPosts(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts, "total": len(posts)})
}

// =====================================================
// LANDING PAGE HANDLERS
// =====================================================

// CreateLandingPage handles POST /api/v1/landing-pages
func (h *MarketingFeaturesHandler) CreateLandingPage(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		HTMLContent string `json:"html_content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := req.Slug
	if slug == "" {
		slug = req.Name // Should generate a slug
	}

	page, err := h.landingService.CreatePage(c.Request.Context(), defaultOrganizationID, req.Name, slug, req.Title, req.HTMLContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, page)
}

// GetLandingPages handles GET /api/v1/landing-pages
func (h *MarketingFeaturesHandler) GetLandingPages(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	pages, err := h.landingService.GetPages(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages, "total": len(pages)})
}

// CreateLead handles POST /api/v1/landing-pages/:id/leads
func (h *MarketingFeaturesHandler) CreateLead(c *gin.Context) {
	pageID := c.Param("id")

	var req struct {
		Email      string                 `json:"email" binding:"required"`
		FirstName  string                 `json:"first_name"`
		LastName   string                 `json:"last_name"`
		Phone      string                 `json:"phone"`
		CustomData map[string]interface{} `json:"custom_data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lead, err := h.landingService.CreateLead(c.Request.Context(), defaultOrganizationID, pageID, req.Email, req.FirstName, req.LastName, req.Phone, req.CustomData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lead)
}

// =====================================================
// AUTOMATION HANDLERS
// =====================================================

// CreateWorkflow handles POST /api/v1/automation/workflows
func (h *MarketingFeaturesHandler) CreateWorkflow(c *gin.Context) {
	var req struct {
		Name          string                 `json:"name" binding:"required"`
		TriggerType   string                 `json:"trigger_type" binding:"required"`
		TriggerConfig map[string]interface{} `json:"trigger_config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.automationService.CreateWorkflow(c.Request.Context(), defaultOrganizationID, req.Name, req.TriggerType, req.TriggerConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

// GetWorkflows handles GET /api/v1/automation/workflows
func (h *MarketingFeaturesHandler) GetWorkflows(c *gin.Context) {
	workflows, err := h.automationService.GetWorkflows(c.Request.Context(), defaultOrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflows": workflows, "total": len(workflows)})
}

// =====================================================
// INFLUENCER HANDLERS
// =====================================================

// CreateInfluencer handles POST /api/v1/influencers
func (h *MarketingFeaturesHandler) CreateInfluencer(c *gin.Context) {
	var req struct {
		Name     string   `json:"name" binding:"required"`
		Email    string   `json:"email"`
		Category string   `json:"category"`
		Niches   []string `json:"niches"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	influencer, err := h.influencerService.CreateInfluencer(c.Request.Context(), defaultOrganizationID, req.Name, req.Email, req.Category, req.Niches)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, influencer)
}

// GetInfluencers handles GET /api/v1/influencers
func (h *MarketingFeaturesHandler) GetInfluencers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	influencers, err := h.influencerService.GetInfluencers(c.Request.Context(), defaultOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"influencers": influencers, "total": len(influencers)})
}

// CreateInfluencerCampaign handles POST /api/v1/influencers/campaigns
func (h *MarketingFeaturesHandler) CreateInfluencerCampaign(c *gin.Context) {
	var req struct {
		Name         string  `json:"name" binding:"required"`
		CampaignType string  `json:"campaign_type"`
		Budget       float64 `json:"budget"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := h.influencerService.CreateCampaign(c.Request.Context(), defaultOrganizationID, req.Name, req.CampaignType, req.Budget)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// =====================================================
// REPORTING HANDLERS
// =====================================================

// GetMarketingOverview handles GET /api/v1/reporting/overview
func (h *MarketingFeaturesHandler) GetMarketingOverview(c *gin.Context) {
	overview, err := h.reportingService.GetMarketingOverview(c.Request.Context(), defaultOrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

// TrackEvent handles POST /api/v1/reporting/events
func (h *MarketingFeaturesHandler) TrackEvent(c *gin.Context) {
	var req struct {
		EventType     string                 `json:"event_type" binding:"required"`
		EventCategory string                 `json:"event_category"`
		EventAction   string                 `json:"event_action"`
		EventData     map[string]interface{} `json:"event_data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reportingService.TrackEvent(c.Request.Context(), defaultOrganizationID, req.EventType, req.EventCategory, req.EventAction, req.EventData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event tracked successfully"})
}

// Helper function to convert to JSON
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}