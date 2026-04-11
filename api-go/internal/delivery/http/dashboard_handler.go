package http

import (
	"net/http"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles HTTP requests for dashboard
type DashboardHandler struct {
	service *service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(service *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetOverview handles GET /api/v1/dashboard
func (h *DashboardHandler) GetOverview(c *gin.Context) {
	// Use default organization for demo
	organizationID := "00000000-0000-0000-0000-000000000001"

	dashboard, err := h.service.GetOverview(c.Request.Context(), organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_campaigns":   dashboard.TotalCampaigns,
		"active_campaigns":  dashboard.ActiveCampaigns,
		"total_spend":       dashboard.TotalSpend,
		"total_revenue":     dashboard.TotalRevenue,
		"roi":               dashboard.ROI,
		"conversion_rate":   dashboard.ConversionRate,
		"top_channels":      dashboard.TopChannels,
	})
}
