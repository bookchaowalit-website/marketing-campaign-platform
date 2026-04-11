package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ABTestHandler struct {
	service *service.ABTestService
}

func NewABTestHandler(service *service.ABTestService) *ABTestHandler {
	return &ABTestHandler{service: service}
}

type CreateABTestRequest struct {
	Name       string  `json:"name" binding:"required"`
	Variant    string  `json:"variant" binding:"required,oneof=A B"`
	CampaignID string  `json:"campaign_id"`
	StartDate  string  `json:"start_date" binding:"required"`
}

type UpdateABTestRequest struct {
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	Visitors    int64   `json:"visitors"`
	Conversions int64   `json:"conversions"`
}

// CreateABTest handles POST /api/v1/ab-tests
func (h *ABTestHandler) CreateABTest(c *gin.Context) {
	var req CreateABTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM-DD"})
		return
	}

	test := &domain.ABTest{
		ID:          uuid.New().String(),
		UserID:      userID.(string),
		CampaignID:  req.CampaignID,
		Name:        req.Name,
		Variant:     req.Variant,
		Status:      "draft",
		StartDate:   startDate,
		ConversionRate: 0,
		Confidence:  0,
	}

	if err := h.service.CreateABTest(c.Request.Context(), test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, test)
}

// GetABTest handles GET /api/v1/ab-tests/:id
func (h *ABTestHandler) GetABTest(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	test, err := h.service.GetABTest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "A/B test not found"})
		return
	}

	// Validate ownership
	if test.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, test)
}

// GetAllABTests handles GET /api/v1/ab-tests
func (h *ABTestHandler) GetAllABTests(c *gin.Context) {
	userID, _ := c.Get("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")

	var tests []*domain.ABTest
	var err error

	if status != "" {
		tests, err = h.service.GetByStatus(c.Request.Context(), userID.(string), status)
	} else {
		tests, err = h.service.GetAllABTests(c.Request.Context(), userID.(string), limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tests": tests, "total": len(tests)})
}

// UpdateABTest handles PUT /api/v1/ab-tests/:id
func (h *ABTestHandler) UpdateABTest(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	test, err := h.service.GetABTest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "A/B test not found"})
		return
	}

	// Validate ownership
	if test.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req UpdateABTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Name != "" {
		test.Name = req.Name
	}
	if req.Status != "" {
		test.Status = req.Status
	}
	test.Visitors = req.Visitors
	test.Conversions = req.Conversions

	// Calculate conversion rate
	if test.Visitors > 0 {
		test.ConversionRate = float64(test.Conversions) / float64(test.Visitors) * 100
	}

	if err := h.service.UpdateABTest(c.Request.Context(), test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, test)
}

// DeleteABTest handles DELETE /api/v1/ab-tests/:id
func (h *ABTestHandler) DeleteABTest(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	test, err := h.service.GetABTest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "A/B test not found"})
		return
	}

	// Validate ownership
	if test.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.service.DeleteABTest(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "A/B test deleted successfully"})
}

// StartABTest handles POST /api/v1/ab-tests/:id/start
func (h *ABTestHandler) StartABTest(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	test, err := h.service.GetABTest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "A/B test not found"})
		return
	}

	// Validate ownership
	if test.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	test.Status = "running"
	test.StartDate = time.Now()

	if err := h.service.UpdateABTest(c.Request.Context(), test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, test)
}

// StopABTest handles POST /api/v1/ab-tests/:id/stop
func (h *ABTestHandler) StopABTest(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	test, err := h.service.GetABTest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "A/B test not found"})
		return
	}

	// Validate ownership
	if test.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	test.Status = "completed"
	now := time.Now()
	test.EndDate = &now

	if err := h.service.UpdateABTest(c.Request.Context(), test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, test)
}