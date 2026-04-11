package http

import (
	"net/http"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	socialService *service.SocialPublishService
}

func NewSocialHandler(socialService *service.SocialPublishService) *SocialHandler {
	return &SocialHandler{socialService: socialService}
}

type ConnectAccountRequest struct {
	Platform string `json:"platform" binding:"required,oneof=facebook instagram"`
	Code     string `json:"code" binding:"required"`
}

type CreatePostRequest struct {
	AccountID   string     `json:"account_id" binding:"required"`
	Content     string     `json:"content" binding:"required"`
	MediaURLs   []string   `json:"media_urls"`
	Hashtags    []string   `json:"hashtags"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

type PublishRequest struct {
	AccountID string   `json:"account_id" binding:"required"`
	Content   string   `json:"content" binding:"required"`
	MediaURLs []string `json:"media_urls"`
}

// ConnectAccount handles POST /api/v1/social/connect
func (h *SocialHandler) ConnectAccount(c *gin.Context) {
	var req ConnectAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001" // Default for development
	}

	account, err := h.socialService.ConnectAccount(c.Request.Context(), orgID, req.Platform, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetOAuthURL handles GET /api/v1/social/oauth/:platform
func (h *SocialHandler) GetOAuthURL(c *gin.Context) {
	platform := c.Param("platform")
	state := c.Query("state")
	
	var url string
	switch platform {
	case "facebook":
		url = h.socialService.GetFacebookOAuthURL(state)
	case "instagram":
		url = h.socialService.GetInstagramOAuthURL(state)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported platform"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"oauth_url": url})
}

// GetAccounts handles GET /api/v1/social/accounts
func (h *SocialHandler) GetAccounts(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	accounts, err := h.socialService.GetAccounts(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// DisconnectAccount handles POST /api/v1/social/accounts/:id/disconnect
func (h *SocialHandler) DisconnectAccount(c *gin.Context) {
	accountID := c.Param("id")
	
	err := h.socialService.DisconnectAccount(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account disconnected"})
}

// GetPosts handles GET /api/v1/social/posts
func (h *SocialHandler) GetPosts(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	limit := 20
	offset := 0
	posts, total, err := h.socialService.GetPosts(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
	})
}

// CreatePost handles POST /api/v1/social/posts
func (h *SocialHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	post, err := h.socialService.CreatePost(c.Request.Context(), orgID, req.AccountID, req.Content, req.MediaURLs, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// PublishPost handles POST /api/v1/social/posts/:id/publish
func (h *SocialHandler) PublishPost(c *gin.Context) {
	postID := c.Param("id")
	
	err := h.socialService.PublishPost(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post published successfully"})
}

// DeletePost handles DELETE /api/v1/social/posts/:id
func (h *SocialHandler) DeletePost(c *gin.Context) {
	postID := c.Param("id")
	
	err := h.socialService.DeletePost(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// GetAnalytics handles GET /api/v1/social/posts/:id/analytics
func (h *SocialHandler) GetAnalytics(c *gin.Context) {
	postID := c.Param("id")
	
	analytics, err := h.socialService.GetPostAnalytics(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetScheduledPosts handles GET /api/v1/social/posts/scheduled
func (h *SocialHandler) GetScheduledPosts(c *gin.Context) {
	orgID := c.GetString("organization_id")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}

	posts, err := h.socialService.GetScheduledPosts(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
