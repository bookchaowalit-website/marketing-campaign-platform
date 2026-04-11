package http

import (
	"net/http"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type SMSHandler struct {
	smsService *service.SMSService
}

func NewSMSHandler(smsService *service.SMSService) *SMSHandler {
	return &SMSHandler{smsService: smsService}
}

// ===========================
// REQUEST STRUCTS
// ===========================

type CreateSMSListRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AddSubscriberRequest struct {
	Phone        string                 `json:"phone" binding:"required"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	ListID       *string                `json:"list_id"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type CreateSMSCampaignRequest struct {
	Name            string  `json:"name" binding:"required"`
	ListID          *string `json:"list_id"`
	MessageTemplate string  `json:"message_template" binding:"required"`
	ScheduledAt     *string `json:"scheduled_at"`
}

// ===========================
// SMS LIST HANDLERS
// ===========================

// ListSMSLists handles GET /api/v1/sms/lists
func (h *SMSHandler) ListSMSLists(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	lists, err := h.smsService.ListLists(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lists": lists, "total": len(lists)})
}

// CreateSMSList handles POST /api/v1/sms/lists
func (h *SMSHandler) CreateSMSList(c *gin.Context) {
	var req CreateSMSListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	list := &domain.SMSList{
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
	}

	if err := h.smsService.CreateList(c.Request.Context(), list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, list)
}

// ===========================
// SMS SUBSCRIBER HANDLERS
// ===========================

// ListSubscribers handles GET /api/v1/sms/subscribers
func (h *SMSHandler) ListSubscribers(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	status := c.Query("status")
	limit := 50
	offset := 0

	subscribers, total, err := h.smsService.ListSubscribers(c.Request.Context(), orgID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscribers": subscribers,
		"total":       total,
		"limit":       limit,
		"offset":      offset,
	})
}

// AddSubscriber handles POST /api/v1/sms/subscribers
func (h *SMSHandler) AddSubscriber(c *gin.Context) {
	var req AddSubscriberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	subscriber := &domain.SMSSubscriber{
		OrganizationID: orgID,
		Phone:          req.Phone,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		ListID:         req.ListID,
		CustomFields:   req.CustomFields,
		Status:         "active",
	}

	if err := h.smsService.AddSubscriber(c.Request.Context(), subscriber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscriber)
}

// ===========================
// SMS CAMPAIGN HANDLERS
// ===========================

// ListCampaigns handles GET /api/v1/sms/campaigns
func (h *SMSHandler) ListCampaigns(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	status := c.Query("status")
	limit := 50
	offset := 0

	campaigns, total, err := h.smsService.ListCampaigns(c.Request.Context(), orgID, status, limit, offset)
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

// CreateCampaign handles POST /api/v1/sms/campaigns
func (h *SMSHandler) CreateCampaign(c *gin.Context) {
	var req CreateSMSCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	campaign := &domain.SMSCampaign{
		OrganizationID:  orgID,
		Name:            req.Name,
		ListID:          req.ListID,
		MessageTemplate: req.MessageTemplate,
		Status:          "draft",
	}

	if req.ScheduledAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err == nil {
			campaign.ScheduledAt = &t
			campaign.Status = "scheduled"
		}
	}

	if err := h.smsService.CreateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// SendCampaign handles POST /api/v1/sms/campaigns/:id/send
func (h *SMSHandler) SendCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.smsService.SendCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign sending initiated"})
}
