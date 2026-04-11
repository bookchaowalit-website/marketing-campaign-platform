package http

import (
	"net/http"
	"strconv"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/dto"
	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AffiliateHandler handles HTTP requests for affiliate operations
type AffiliateHandler struct {
	service *service.AffiliateService
}

// NewAffiliateHandler creates a new affiliate handler
func NewAffiliateHandler(service *service.AffiliateService) *AffiliateHandler {
	return &AffiliateHandler{service: service}
}

// RegisterAffiliate handles POST /api/v1/affiliates/register
func (h *AffiliateHandler) RegisterAffiliate(c *gin.Context) {
	var req dto.RegisterAffiliateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	affiliate, err := h.service.RegisterAffiliate(
		c.Request.Context(),
		userID.(string),
		defaultOrganizationID,
		req.Website,
		req.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.toAffiliateResponse(affiliate)
	c.JSON(http.StatusCreated, response)
}

// GetMyAffiliate handles GET /api/v1/affiliates/me
func (h *AffiliateHandler) GetMyAffiliate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	affiliate, err := h.service.GetAffiliateByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate account not found"})
		return
	}

	response := h.toAffiliateResponse(affiliate)
	c.JSON(http.StatusOK, response)
}

// GetAffiliate handles GET /api/v1/affiliates/:id
func (h *AffiliateHandler) GetAffiliate(c *gin.Context) {
	id := c.Param("id")

	affiliate, err := h.service.GetAffiliate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate not found"})
		return
	}

	response := h.toAffiliateResponse(affiliate)
	c.JSON(http.StatusOK, response)
}

// GetAllAffiliates handles GET /api/v1/affiliates
func (h *AffiliateHandler) GetAllAffiliates(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")

	var affiliates []*domain.Affiliate
	var err error

	if status != "" {
		affiliates, err = h.service.affiliateRepo.GetByStatus(c.Request.Context(), defaultOrganizationID, status, limit, offset)
	} else {
		affiliates, err = h.service.GetAllAffiliates(c.Request.Context(), defaultOrganizationID, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.AffiliateResponse, len(affiliates))
	for i, affiliate := range affiliates {
		response[i] = h.toAffiliateResponse(affiliate)
	}

	c.JSON(http.StatusOK, gin.H{"affiliates": response, "total": len(affiliates)})
}

// ApproveAffiliate handles POST /api/v1/affiliates/:id/approve
func (h *AffiliateHandler) ApproveAffiliate(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.ApproveAffiliate(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Affiliate approved successfully"})
}

// SuspendAffiliate handles POST /api/v1/affiliates/:id/suspend
func (h *AffiliateHandler) SuspendAffiliate(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.SuspendAffiliate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Affiliate suspended successfully"})
}

// UpdateAffiliate handles PUT /api/v1/affiliates/:id
func (h *AffiliateHandler) UpdateAffiliate(c *gin.Context) {
	id := c.Param("id")

	affiliate, err := h.service.GetAffiliate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate not found"})
		return
	}

	var req dto.UpdateAffiliateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.CommissionRate > 0 {
		affiliate.CommissionRate = req.CommissionRate
	}
	if req.CommissionType != "" {
		affiliate.CommissionType = req.CommissionType
	}
	if req.FixedCommissionAmount > 0 {
		affiliate.FixedCommissionAmount = req.FixedCommissionAmount
	}
	if req.PaypalEmail != "" {
		affiliate.PaypalEmail = req.PaypalEmail
	}
	if req.BankName != "" {
		affiliate.BankName = req.BankName
	}
	if req.BankAccount != "" {
		affiliate.BankAccount = req.BankAccount
	}
	if req.BankRouting != "" {
		affiliate.BankRouting = req.BankRouting
	}
	if req.TaxID != "" {
		affiliate.TaxID = req.TaxID
	}
	if req.Website != "" {
		affiliate.Website = req.Website
	}
	if req.Description != "" {
		affiliate.Description = req.Description
	}

	if err := h.service.affiliateRepo.Update(c.Request.Context(), affiliate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.toAffiliateResponse(affiliate)
	c.JSON(http.StatusOK, response)
}

// GetAffiliateStats handles GET /api/v1/affiliates/stats
func (h *AffiliateHandler) GetAffiliateStats(c *gin.Context) {
	stats, err := h.service.GetAffiliateStats(c.Request.Context(), defaultOrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.AffiliateStatsResponse{
		TotalAffiliates:  stats.TotalAffiliates,
		ActiveAffiliates: stats.ActiveAffiliates,
		TotalClicks:      stats.TotalClicks,
		TotalSignups:     stats.TotalSignups,
		TotalConversions: stats.TotalConversions,
		ConversionRate:   stats.ConversionRate,
		TotalEarnings:    stats.TotalEarnings,
		TotalPaid:        stats.TotalPaid,
		PendingPayouts:   stats.PendingPayouts,
		TopAffiliates:    make([]dto.AffiliateLeaderboardResponse, len(stats.TopAffiliates)),
	}

	for i, entry := range stats.TopAffiliates {
		response.TopAffiliates[i] = dto.AffiliateLeaderboardResponse{
			AffiliateID:   entry.AffiliateID,
			AffiliateCode: entry.AffiliateCode,
			Name:          entry.Name,
			Conversions:   entry.Conversions,
			Earnings:      entry.Earnings,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetAffiliateDashboard handles GET /api/v1/affiliates/dashboard
func (h *AffiliateHandler) GetAffiliateDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	affiliate, err := h.service.GetAffiliateByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate account not found"})
		return
	}

	dashboard, err := h.service.GetAffiliateDashboard(c.Request.Context(), affiliate.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.AffiliateDashboardResponse{
		Affiliate:         h.toAffiliateResponse(&dashboard.Affiliate),
		RecentReferrals:   make([]dto.ReferralResponse, len(dashboard.RecentReferrals)),
		RecentCommissions: make([]dto.CommissionResponse, len(dashboard.RecentCommissions)),
		TopLinks:          make([]dto.AffiliateLinkResponse, len(dashboard.TopLinks)),
		EarningsHistory:   make([]dto.MonthlyEarningsResponse, len(dashboard.EarningsHistory)),
	}

	for i, ref := range dashboard.RecentReferrals {
		response.RecentReferrals[i] = h.toReferralResponse(&ref)
	}

	for i, comm := range dashboard.RecentCommissions {
		response.RecentCommissions[i] = h.toCommissionResponse(&comm)
	}

	for i, link := range dashboard.TopLinks {
		response.TopLinks[i] = h.toAffiliateLinkResponse(&link)
	}

	for i, earning := range dashboard.EarningsHistory {
		response.EarningsHistory[i] = dto.MonthlyEarningsResponse{
			Month:       earning.Month,
			Earnings:    earning.Earnings,
			Conversions: earning.Conversions,
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateAffiliateLink handles POST /api/v1/affiliates/links
func (h *AffiliateHandler) CreateAffiliateLink(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	affiliate, err := h.service.GetAffiliateByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate account not found"})
		return
	}

	var req dto.CreateAffiliateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.service.CreateAffiliateLink(
		c.Request.Context(),
		affiliate.ID,
		req.Name,
		req.DestinationURL,
		req.CampaignID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.toAffiliateLinkResponse(link)
	c.JSON(http.StatusCreated, response)
}

// GetAffiliateLinks handles GET /api/v1/affiliates/links
func (h *AffiliateHandler) GetAffiliateLinks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	affiliate, err := h.service.GetAffiliateByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate account not found"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	links, err := h.service.GetAffiliateLinks(c.Request.Context(), affiliate.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.AffiliateLinkResponse, len(links))
	for i, link := range links {
		response[i] = h.toAffiliateLinkResponse(link)
	}

	c.JSON(http.StatusOK, gin.H{"links": response, "total": len(links)})
}

// DeleteAffiliateLink handles DELETE /api/v1/affiliates/links/:id
func (h *AffiliateHandler) DeleteAffiliateLink(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.linkRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}

// TrackClick handles POST /api/v1/affiliates/track/click (public endpoint)
func (h *AffiliateHandler) TrackClick(c *gin.Context) {
	var req dto.TrackClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referrerURL := c.GetHeader("Referer")

	referral, err := h.service.TrackClick(
		c.Request.Context(),
		req.LinkCode,
		req.AffiliateCode,
		ip,
		userAgent,
		referrerURL,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"referral_code": referral.ReferralCode,
		"message":       "Click tracked successfully",
	})
}

// TrackSignup handles POST /api/v1/affiliates/track/signup (internal)
func (h *AffiliateHandler) TrackSignup(c *gin.Context) {
	var req dto.TrackSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.TrackSignup(c.Request.Context(), req.ReferralCode, req.ReferredUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup tracked successfully"})
}

// TrackConversion handles POST /api/v1/affiliates/track/conversion (internal)
func (h *AffiliateHandler) TrackConversion(c *gin.Context) {
	var req dto.TrackConversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commission, err := h.service.TrackConversion(
		c.Request.Context(),
		req.ReferralCode,
		req.OrderID,
		req.OrderAmount,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"commission_id":     commission.ID,
		"commission_amount": commission.Amount,
		"message":           "Conversion tracked successfully",
	})
}

// GetCommissions handles GET /api/v1/affiliates/commissions
func (h *AffiliateHandler) GetCommissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	affiliate, err := h.service.GetAffiliateByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate account not found"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")

	var commissions []*domain.Commission
	if status != "" {
		commissions, err = h.service.commissionRepo.GetByStatus(c.Request.Context(), affiliate.ID, status, limit, offset)
	} else {
		commissions, err = h.service.GetCommissions(c.Request.Context(), affiliate.ID, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.CommissionResponse, len(commissions))
	for i, comm := range commissions {
		response[i] = h.toCommissionResponse(comm)
	}

	c.JSON(http.StatusOK, gin.H{"commissions": response, "total": len(commissions)})
}

// GetAffiliateSettings handles GET /api/v1/affiliates/settings
func (h *AffiliateHandler) GetAffiliateSettings(c *gin.Context) {
	settings, err := h.service.GetAffiliateSettings(c.Request.Context(), defaultOrganizationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate settings not found"})
		return
	}

	response := dto.AffiliateSettingsResponse{
		ID:                      settings.ID,
		OrganizationID:          settings.OrganizationID,
		ProgramName:             settings.ProgramName,
		ProgramDescription:      settings.ProgramDescription,
		DefaultCommissionRate:   settings.DefaultCommissionRate,
		DefaultCommissionType:   settings.DefaultCommissionType,
		DefaultFixedCommission:  settings.DefaultFixedCommission,
		CookieDurationDays:      settings.CookieDurationDays,
		MinimumPayoutAmount:     settings.MinimumPayoutAmount,
		PayoutSchedule:          settings.PayoutSchedule,
		RequireApproval:         settings.RequireApproval,
		AllowCustomCommissionRates: settings.AllowCustomCommissionRates,
		AllowAffiliateLinks:     settings.AllowAffiliateLinks,
		TermsAndConditions:      settings.TermsAndConditions,
		WelcomeMessage:          settings.WelcomeMessage,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAffiliateSettings handles PUT /api/v1/affiliates/settings
func (h *AffiliateHandler) UpdateAffiliateSettings(c *gin.Context) {
	settings, err := h.service.GetAffiliateSettings(c.Request.Context(), defaultOrganizationID)
	if err != nil {
		// Create default settings
		settings = &domain.AffiliateSettings{
			ID:                    uuid.New().String(),
			OrganizationID:        defaultOrganizationID,
			ProgramName:           "Affiliate Program",
			DefaultCommissionRate: 0.10,
			DefaultCommissionType: "percentage",
			MinimumPayoutAmount:   50.00,
			RequireApproval:       true,
			AllowCustomCommissionRates: true,
			AllowAffiliateLinks:   true,
		}
	}

	var req dto.UpdateAffiliateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.ProgramName != "" {
		settings.ProgramName = req.ProgramName
	}
	if req.ProgramDescription != "" {
		settings.ProgramDescription = req.ProgramDescription
	}
	if req.DefaultCommissionRate > 0 {
		settings.DefaultCommissionRate = req.DefaultCommissionRate
	}
	if req.DefaultCommissionType != "" {
		settings.DefaultCommissionType = req.DefaultCommissionType
	}
	if req.DefaultFixedCommission > 0 {
		settings.DefaultFixedCommission = req.DefaultFixedCommission
	}
	if req.CookieDurationDays > 0 {
		settings.CookieDurationDays = req.CookieDurationDays
	}
	if req.MinimumPayoutAmount > 0 {
		settings.MinimumPayoutAmount = req.MinimumPayoutAmount
	}
	if req.PayoutSchedule != "" {
		settings.PayoutSchedule = req.PayoutSchedule
	}

	settings.RequireApproval = req.RequireApproval
	settings.AllowCustomCommissionRates = req.AllowCustomCommissionRates
	settings.AllowAffiliateLinks = req.AllowAffiliateLinks

	if req.TermsAndConditions != "" {
		settings.TermsAndConditions = req.TermsAndConditions
	}
	if req.WelcomeMessage != "" {
		settings.WelcomeMessage = req.WelcomeMessage
	}

	if err := h.service.UpdateAffiliateSettings(c.Request.Context(), settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.AffiliateSettingsResponse{
		ID:                      settings.ID,
		OrganizationID:          settings.OrganizationID,
		ProgramName:             settings.ProgramName,
		ProgramDescription:      settings.ProgramDescription,
		DefaultCommissionRate:   settings.DefaultCommissionRate,
		DefaultCommissionType:   settings.DefaultCommissionType,
		DefaultFixedCommission:  settings.DefaultFixedCommission,
		CookieDurationDays:      settings.CookieDurationDays,
		MinimumPayoutAmount:     settings.MinimumPayoutAmount,
		PayoutSchedule:          settings.PayoutSchedule,
		RequireApproval:         settings.RequireApproval,
		AllowCustomCommissionRates: settings.AllowCustomCommissionRates,
		AllowAffiliateLinks:     settings.AllowAffiliateLinks,
		TermsAndConditions:      settings.TermsAndConditions,
		WelcomeMessage:          settings.WelcomeMessage,
	}

	c.JSON(http.StatusOK, response)
}

// Helper conversion functions
func (h *AffiliateHandler) toAffiliateResponse(affiliate *domain.Affiliate) dto.AffiliateResponse {
	response := dto.AffiliateResponse{
		ID:                    affiliate.ID,
		UserID:                affiliate.UserID,
		AffiliateCode:         affiliate.AffiliateCode,
		Status:                affiliate.Status,
		CommissionRate:        affiliate.CommissionRate,
		CommissionType:        affiliate.CommissionType,
		FixedCommissionAmount: affiliate.FixedCommissionAmount,
		TotalClicks:           affiliate.TotalClicks,
		TotalSignups:          affiliate.TotalSignups,
		TotalConversions:      affiliate.TotalConversions,
		TotalEarnings:         affiliate.TotalEarnings,
		TotalPaid:             affiliate.TotalPaid,
		Balance:               affiliate.Balance,
		MinimumPayout:         affiliate.MinimumPayout,
		ConversionRate:        affiliate.ConversionRate(),
		SignupRate:            affiliate.SignupRate(),
		CanRequestPayout:      affiliate.CanRequestPayout(),
		Website:               affiliate.Website,
		Description:           affiliate.Description,
		PaypalEmail:           affiliate.PaypalEmail,
		CreatedAt:             affiliate.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             affiliate.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if affiliate.ApprovedAt != nil {
		response.ApprovedAt = affiliate.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	return response
}

func (h *AffiliateHandler) toAffiliateLinkResponse(link *domain.AffiliateLink) dto.AffiliateLinkResponse {
	return dto.AffiliateLinkResponse{
		ID:              link.ID,
		AffiliateID:     link.AffiliateID,
		LinkCode:        link.LinkCode,
		TrackingURL:     "https://bookmarketing.com/go/" + link.LinkCode,
		Name:            link.Name,
		DestinationURL:  link.DestinationURL,
		CampaignID:      link.CampaignID,
		ClickCount:      link.ClickCount,
		SignupCount:     link.SignupCount,
		ConversionCount: link.ConversionCount,
		Earnings:        link.Earnings,
		IsActive:        link.IsActive,
		CreatedAt:       link.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       link.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (h *AffiliateHandler) toReferralResponse(referral *domain.Referral) dto.ReferralResponse {
	response := dto.ReferralResponse{
		ID:               referral.ID,
		AffiliateID:      referral.AffiliateID,
		ReferralCode:     referral.ReferralCode,
		Status:           referral.Status,
		ClickTimestamp:   referral.ClickTimestamp.Format("2006-01-02 15:04:05"),
		OrderAmount:      referral.OrderAmount,
		CommissionAmount: referral.CommissionAmount,
		CommissionStatus: referral.CommissionStatus,
		CreatedAt:        referral.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if referral.SignupTimestamp != nil {
		response.SignupTimestamp = referral.SignupTimestamp.Format("2006-01-02 15:04:05")
	}
	if referral.ConversionTimestamp != nil {
		response.ConversionTimestamp = referral.ConversionTimestamp.Format("2006-01-02 15:04:05")
	}
	if referral.OrderID != "" {
		response.OrderID = referral.OrderID
	}

	return response
}

func (h *AffiliateHandler) toCommissionResponse(commission *domain.Commission) dto.CommissionResponse {
	response := dto.CommissionResponse{
		ID:           commission.ID,
		AffiliateID:  commission.AffiliateID,
		ReferralID:   commission.ReferralID,
		Amount:       commission.Amount,
		Status:       commission.Status,
		Type:         commission.Type,
		Description:  commission.Description,
		ReferenceID:  commission.ReferenceID,
		CreatedAt:    commission.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    commission.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if commission.ApprovedAt != nil {
		response.ApprovedAt = commission.ApprovedAt.Format("2006-01-02 15:04:05")
	}
	if commission.PaidAt != nil {
		response.PaidAt = commission.PaidAt.Format("2006-01-02 15:04:05")
	}

	return response
}