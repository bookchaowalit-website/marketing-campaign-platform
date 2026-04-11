package http

import (
	"net/http"
	"strconv"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/dto"
	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
)

// ExternalAffiliateHandler handles HTTP requests for external affiliate integrations
type ExternalAffiliateHandler struct {
	service *service.ExternalAffiliateService
}

// NewExternalAffiliateHandler creates a new external affiliate handler
func NewExternalAffiliateHandler(service *service.ExternalAffiliateService) *ExternalAffiliateHandler {
	return &ExternalAffiliateHandler{service: service}
}

// ConnectPlatform handles POST /api/v1/integrations/connect
func (h *ExternalAffiliateHandler) ConnectPlatform(c *gin.Context) {
	var req dto.ConnectPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	account, err := h.service.ConnectPlatform(c.Request.Context(), userID.(string), defaultOrganizationID, &domain.PlatformConnectionRequest{
		Platform:          domain.Platform(req.Platform),
		APIKey:            req.APIKey,
		APISecret:         req.APISecret,
		AccessToken:       req.AccessToken,
		RefreshToken:      req.RefreshToken,
		PlatformAccountID: req.PlatformAccountID,
		PlatformUsername:  req.PlatformUsername,
		AuthorizationCode: req.AuthorizationCode,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.toAccountResponse(account))
}

// DisconnectPlatform handles POST /api/v1/integrations/:id/disconnect
func (h *ExternalAffiliateHandler) DisconnectPlatform(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DisconnectPlatform(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account disconnected successfully"})
}

// GetConnectedAccounts handles GET /api/v1/integrations/accounts
func (h *ExternalAffiliateHandler) GetConnectedAccounts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	accounts, err := h.service.GetConnectedAccounts(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ExternalAccountResponse, len(accounts))
	for i, account := range accounts {
		response[i] = h.toAccountResponse(account)
	}

	c.JSON(http.StatusOK, gin.H{"accounts": response, "total": len(accounts)})
}

// GetAccount handles GET /api/v1/integrations/accounts/:id
func (h *ExternalAffiliateHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")

	account, err := h.service.GetAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, h.toAccountResponse(account))
}

// SyncAccount handles POST /api/v1/integrations/accounts/:id/sync
func (h *ExternalAffiliateHandler) SyncAccount(c *gin.Context) {
	id := c.Param("id")

	var req dto.SyncAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.SyncType = "full"
	}

	if err := h.service.SyncAccount(c.Request.Context(), id, req.SyncType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync completed successfully"})
}

// GetUnifiedDashboard handles GET /api/v1/integrations/dashboard
func (h *ExternalAffiliateHandler) GetUnifiedDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dashboard, err := h.service.GetUnifiedDashboard(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.UnifiedDashboardResponse{
		TotalEarnings:         dashboard.TotalEarnings,
		TotalPendingCommission: dashboard.TotalPendingCommission,
		TotalOrders:           dashboard.TotalOrders,
		TotalClicks:           dashboard.TotalClicks,
		TotalProducts:         dashboard.TotalProducts,
		PlatformStats:         make([]dto.PlatformStatsResponse, len(dashboard.PlatformStats)),
		RecentOrders:          make([]dto.ExternalOrderResponse, len(dashboard.RecentOrders)),
		TopProducts:           make([]dto.ExternalProductResponse, len(dashboard.TopProducts)),
		EarningsTrend:         make([]dto.DailyEarningsResponse, len(dashboard.EarningsTrend)),
		ConnectedAccounts:     make([]dto.ExternalAccountResponse, len(dashboard.ConnectedAccounts)),
	}

	for i, ps := range dashboard.PlatformStats {
		response.PlatformStats[i] = dto.PlatformStatsResponse{
			Platform:          string(ps.Platform),
			AccountCount:      ps.AccountCount,
			TotalEarnings:     ps.TotalEarnings,
			TotalOrders:       ps.TotalOrders,
			TotalClicks:       ps.TotalClicks,
			PendingCommission: ps.PendingCommission,
			AvgConversionRate: ps.AvgConversionRate,
		}
		if ps.LastSyncAt != nil {
			response.PlatformStats[i].LastSyncAt = ps.LastSyncAt.Format("2006-01-02 15:04:05")
		}
	}

	for i, order := range dashboard.RecentOrders {
		response.RecentOrders[i] = h.toOrderResponse(&order)
	}

	for i, product := range dashboard.TopProducts {
		response.TopProducts[i] = h.toProductResponse(&product)
	}

	for i, e := range dashboard.EarningsTrend {
		response.EarningsTrend[i] = dto.DailyEarningsResponse{
			Date:           e.Date.Format("2006-01-02"),
			TotalEarnings:  e.TotalEarnings,
			OrdersCount:    e.OrdersCount,
			Clicks:         e.Clicks,
			ConversionRate: e.ConversionRate,
		}
	}

	for i, account := range dashboard.ConnectedAccounts {
		response.ConnectedAccounts[i] = h.toAccountResponse(&account)
	}

	c.JSON(http.StatusOK, response)
}

// GetProducts handles GET /api/v1/integrations/accounts/:id/products
func (h *ExternalAffiliateHandler) GetProducts(c *gin.Context) {
	id := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.service.GetProducts(c.Request.Context(), id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ExternalProductResponse, len(products))
	for i, product := range products {
		response[i] = h.toProductResponse(product)
	}

	c.JSON(http.StatusOK, gin.H{"products": response, "total": len(products)})
}

// GetOrders handles GET /api/v1/integrations/accounts/:id/orders
func (h *ExternalAffiliateHandler) GetOrders(c *gin.Context) {
	id := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	orders, err := h.service.GetOrders(c.Request.Context(), id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ExternalOrderResponse, len(orders))
	for i, order := range orders {
		response[i] = h.toOrderResponse(order)
	}

	c.JSON(http.StatusOK, gin.H{"orders": response, "total": len(orders)})
}

// GetSyncLogs handles GET /api/v1/integrations/accounts/:id/logs
func (h *ExternalAffiliateHandler) GetSyncLogs(c *gin.Context) {
	id := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	logs, err := h.service.GetSyncLogs(c.Request.Context(), id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.SyncLogResponse, len(logs))
	for i, log := range logs {
		response[i] = dto.SyncLogResponse{
			ID:               log.ID,
			SyncType:         log.SyncType,
			Status:           log.Status,
			StartedAt:        log.StartedAt.Format("2006-01-02 15:04:05"),
			RecordsCreated:   log.RecordsCreated,
			RecordsUpdated:   log.RecordsUpdated,
			RecordsFailed:    log.RecordsFailed,
			ErrorMessage:     log.ErrorMessage,
		}
		if log.CompletedAt != nil {
			response[i].CompletedAt = log.CompletedAt.Format("2006-01-02 15:04:05")
		}
	}

	c.JSON(http.StatusOK, gin.H{"logs": response, "total": len(logs)})
}

// GenerateAffiliateLink handles POST /api/v1/integrations/accounts/:id/link
func (h *ExternalAffiliateHandler) GenerateAffiliateLink(c *gin.Context) {
	id := c.Param("id")

	var req dto.GenerateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.service.GenerateAffiliateLink(c.Request.Context(), id, req.ProductURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"affiliate_link": link})
}

// GetOAuthURL handles GET /api/v1/integrations/oauth/:platform
func (h *ExternalAffiliateHandler) GetOAuthURL(c *gin.Context) {
	platform := c.Param("platform")
	state := c.Query("state")
	if state == "" {
		state = "random-state"
	}

	url, err := h.service.GetOAuthURL(c.Request.Context(), domain.Platform(platform), state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"oauth_url": url})
}

// HandleOAuthCallback handles GET /api/v1/integrations/oauth/:platform/callback
func (h *ExternalAffiliateHandler) HandleOAuthCallback(c *gin.Context) {
	platform := c.Param("platform")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	account, err := h.service.HandleOAuthCallback(c.Request.Context(), userID.(string), defaultOrganizationID, domain.Platform(platform), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// In production, redirect back to frontend with success
	c.JSON(http.StatusOK, gin.H{
		"message":     "Account connected successfully",
		"account_id":  account.ID,
		"platform":    account.Platform,
		"state":       state,
	})
}

// Helper conversion functions
func (h *ExternalAffiliateHandler) toAccountResponse(account *domain.ExternalAffiliateAccount) dto.ExternalAccountResponse {
	response := dto.ExternalAccountResponse{
		ID:                  account.ID,
		Platform:            string(account.Platform),
		PlatformAccountID:   account.PlatformAccountID,
		PlatformAccountName: account.PlatformAccountName,
		PlatformUsername:    account.PlatformUsername,
		Status:              account.Status,
		AutoSync:            account.AutoSync,
		SyncIntervalMins:    account.SyncIntervalMins,
		LastSyncStatus:      account.LastSyncStatus,
		LastSyncError:       account.LastSyncError,
		TotalEarnings:       account.TotalEarnings,
		TotalOrders:         account.TotalOrders,
		TotalClicks:         account.TotalClicks,
		TotalProducts:       account.TotalProducts,
		PendingCommission:   account.PendingCommission,
		CreatedAt:           account.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           account.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if account.LastSyncAt != nil {
		response.LastSyncAt = account.LastSyncAt.Format("2006-01-02 15:04:05")
	}

	return response
}

func (h *ExternalAffiliateHandler) toProductResponse(product *domain.ExternalAffiliateProduct) dto.ExternalProductResponse {
	return dto.ExternalProductResponse{
		ID:               product.ID,
		Platform:         string(product.Platform),
		ProductID:        product.ProductID,
		ProductName:      product.ProductName,
		ProductURL:       product.ProductURL,
		ProductImageURL:  product.ProductImageURL,
		Category:         product.Category,
		Brand:            product.Brand,
		OriginalPrice:    product.OriginalPrice,
		SalePrice:        product.SalePrice,
		Currency:         product.Currency,
		CommissionRate:   product.CommissionRate,
		CommissionType:   product.CommissionType,
		EstimatedEarnings: product.EstimatedEarnings,
		TotalClicks:      product.TotalClicks,
		TotalOrders:      product.TotalOrders,
		TotalEarnings:    product.TotalEarnings,
		ProductStatus:    product.ProductStatus,
		AffiliateLink:    product.AffiliateLink,
	}
}

func (h *ExternalAffiliateHandler) toOrderResponse(order *domain.ExternalAffiliateOrder) dto.ExternalOrderResponse {
	response := dto.ExternalOrderResponse{
		ID:               order.ID,
		Platform:         string(order.Platform),
		OrderID:          order.OrderID,
		OrderStatus:      order.OrderStatus,
		OrderItems:       order.OrderItems,
		OrderAmount:      order.OrderAmount,
		ProductAmount:    order.ProductAmount,
		ShippingAmount:   order.ShippingAmount,
		DiscountAmount:   order.DiscountAmount,
		Currency:         order.Currency,
		CommissionRate:   order.CommissionRate,
		CommissionAmount: order.CommissionAmount,
		CommissionStatus: order.CommissionStatus,
		TrafficSource:    order.TrafficSource,
	}

	if order.OrderDate != nil {
		response.OrderDate = order.OrderDate.Format("2006-01-02 15:04:05")
	}

	return response
}