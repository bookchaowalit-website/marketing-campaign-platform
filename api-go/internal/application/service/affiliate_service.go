package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrAffiliateNotFound      = errors.New("affiliate not found")
	ErrAffiliateAlreadyExists = errors.New("affiliate already exists for this user")
	ErrAffiliateNotActive     = errors.New("affiliate is not active")
	ErrInvalidAffiliateCode   = errors.New("invalid affiliate code")
	ErrLinkNotFound           = errors.New("affiliate link not found")
	ErrReferralNotFound       = errors.New("referral not found")
	ErrCommissionNotFound     = errors.New("commission not found")
	ErrInsufficientBalance    = errors.New("insufficient balance for payout")
	ErrSettingsNotFound       = errors.New("affiliate settings not found")
)

// AffiliateService handles business logic for affiliate operations
type AffiliateService struct {
	affiliateRepo  domain.AffiliateRepository
	linkRepo       domain.AffiliateLinkRepository
	referralRepo   domain.ReferralRepository
	commissionRepo domain.CommissionRepository
	settingsRepo   domain.AffiliateSettingsRepository
}

// NewAffiliateService creates a new affiliate service
func NewAffiliateService(
	affiliateRepo domain.AffiliateRepository,
	linkRepo domain.AffiliateLinkRepository,
	referralRepo domain.ReferralRepository,
	commissionRepo domain.CommissionRepository,
	settingsRepo domain.AffiliateSettingsRepository,
) *AffiliateService {
	return &AffiliateService{
		affiliateRepo:  affiliateRepo,
		linkRepo:       linkRepo,
		referralRepo:   referralRepo,
		commissionRepo: commissionRepo,
		settingsRepo:   settingsRepo,
	}
}

// RegisterAffiliate registers a new affiliate
func (s *AffiliateService) RegisterAffiliate(ctx context.Context, userID, organizationID string, website, description string) (*domain.Affiliate, error) {
	// Check if user already has an affiliate account
	existing, err := s.affiliateRepo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrAffiliateAlreadyExists
	}

	// Get affiliate settings
	settings, err := s.settingsRepo.GetByOrganizationID(ctx, organizationID)
	if err != nil {
		// Create default settings if not exists
		settings = &domain.AffiliateSettings{
			ID:                    uuid.New().String(),
			OrganizationID:        organizationID,
			ProgramName:           "Affiliate Program",
			DefaultCommissionRate: 0.10,
			DefaultCommissionType: "percentage",
			MinimumPayoutAmount:   50.00,
			RequireApproval:       true,
		}
		if err := s.settingsRepo.Create(ctx, settings); err != nil {
			return nil, err
		}
	}

	// Generate unique affiliate code
	code, err := generateAffiliateCode()
	if err != nil {
		return nil, err
	}

	// Determine initial status
	status := "active"
	if settings.RequireApproval {
		status = "pending"
	}

	affiliate := &domain.Affiliate{
		ID:                  uuid.New().String(),
		UserID:              userID,
		OrganizationID:      organizationID,
		AffiliateCode:       code,
		Status:              status,
		CommissionRate:      settings.DefaultCommissionRate,
		CommissionType:      settings.DefaultCommissionType,
		FixedCommissionAmount: settings.DefaultFixedCommission,
		MinimumPayout:       settings.MinimumPayoutAmount,
		Website:             website,
		Description:         description,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.affiliateRepo.Create(ctx, affiliate); err != nil {
		return nil, err
	}

	return affiliate, nil
}

// GetAffiliate retrieves an affiliate by ID
func (s *AffiliateService) GetAffiliate(ctx context.Context, id string) (*domain.Affiliate, error) {
	affiliate, err := s.affiliateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrAffiliateNotFound
	}
	return affiliate, nil
}

// GetAffiliateByUserID retrieves an affiliate by user ID
func (s *AffiliateService) GetAffiliateByUserID(ctx context.Context, userID string) (*domain.Affiliate, error) {
	affiliate, err := s.affiliateRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrAffiliateNotFound
	}
	return affiliate, nil
}

// GetAffiliateByCode retrieves an affiliate by affiliate code
func (s *AffiliateService) GetAffiliateByCode(ctx context.Context, code string) (*domain.Affiliate, error) {
	affiliate, err := s.affiliateRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, ErrInvalidAffiliateCode
	}
	return affiliate, nil
}

// GetAllAffiliates retrieves all affiliates for an organization
func (s *AffiliateService) GetAllAffiliates(ctx context.Context, organizationID string, limit, offset int) ([]*domain.Affiliate, error) {
	return s.affiliateRepo.GetAllByOrganizationID(ctx, organizationID, limit, offset)
}

// GetAffiliatesByStatus returns affiliates filtered by status
func (s *AffiliateService) GetAffiliatesByStatus(ctx context.Context, organizationID, status string, limit, offset int) ([]*domain.Affiliate, error) {
	return s.affiliateRepo.GetByStatus(ctx, organizationID, status, limit, offset)
}

// UpdateAffiliate persists changes to an affiliate
func (s *AffiliateService) UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	return s.affiliateRepo.Update(ctx, affiliate)
}

// DeleteAffiliateLink deletes an affiliate link
func (s *AffiliateService) DeleteAffiliateLink(ctx context.Context, id string) error {
	return s.linkRepo.Delete(ctx, id)
}

// GetCommissionsByStatus returns commissions filtered by status
func (s *AffiliateService) GetCommissionsByStatus(ctx context.Context, affiliateID, status string, limit, offset int) ([]*domain.Commission, error) {
	return s.commissionRepo.GetByStatus(ctx, affiliateID, status, limit, offset)
}

// ApproveAffiliate approves a pending affiliate
func (s *AffiliateService) ApproveAffiliate(ctx context.Context, id, approvedBy string) error {
	affiliate, err := s.affiliateRepo.GetByID(ctx, id)
	if err != nil {
		return ErrAffiliateNotFound
	}

	now := time.Now()
	affiliate.Status = "active"
	affiliate.ApprovedAt = &now
	affiliate.ApprovedBy = approvedBy
	affiliate.UpdatedAt = now

	return s.affiliateRepo.Update(ctx, affiliate)
}

// SuspendAffiliate suspends an affiliate
func (s *AffiliateService) SuspendAffiliate(ctx context.Context, id string) error {
	affiliate, err := s.affiliateRepo.GetByID(ctx, id)
	if err != nil {
		return ErrAffiliateNotFound
	}

	affiliate.Status = "suspended"
	affiliate.UpdatedAt = time.Now()

	return s.affiliateRepo.Update(ctx, affiliate)
}

// UpdateAffiliateCommission updates an affiliate's commission settings
func (s *AffiliateService) UpdateAffiliateCommission(ctx context.Context, id string, rate float64, commissionType string, fixedAmount float64) error {
	affiliate, err := s.affiliateRepo.GetByID(ctx, id)
	if err != nil {
		return ErrAffiliateNotFound
	}

	// Check if custom rates are allowed
	settings, err := s.settingsRepo.GetByOrganizationID(ctx, affiliate.OrganizationID)
	if err == nil && !settings.AllowCustomCommissionRates {
		return errors.New("custom commission rates are not allowed")
	}

	affiliate.CommissionRate = rate
	affiliate.CommissionType = commissionType
	affiliate.FixedCommissionAmount = fixedAmount
	affiliate.UpdatedAt = time.Now()

	return s.affiliateRepo.Update(ctx, affiliate)
}

// GetAffiliateStats retrieves affiliate statistics for an organization
func (s *AffiliateService) GetAffiliateStats(ctx context.Context, organizationID string) (*domain.AffiliateStats, error) {
	return s.affiliateRepo.GetStats(ctx, organizationID)
}

// GetAffiliateDashboard retrieves dashboard data for an affiliate
func (s *AffiliateService) GetAffiliateDashboard(ctx context.Context, affiliateID string) (*domain.AffiliateDashboard, error) {
	return s.affiliateRepo.GetDashboard(ctx, affiliateID)
}

// CreateAffiliateLink creates a new affiliate link
func (s *AffiliateService) CreateAffiliateLink(ctx context.Context, affiliateID, name, destinationURL, campaignID string) (*domain.AffiliateLink, error) {
	// Verify affiliate is active
	affiliate, err := s.affiliateRepo.GetByID(ctx, affiliateID)
	if err != nil {
		return nil, ErrAffiliateNotFound
	}
	if affiliate.Status != "active" {
		return nil, ErrAffiliateNotActive
	}

	// Check if affiliate links are allowed
	settings, err := s.settingsRepo.GetByOrganizationID(ctx, affiliate.OrganizationID)
	if err == nil && !settings.AllowAffiliateLinks {
		return nil, errors.New("affiliate links are not allowed")
	}

	// Generate unique link code
	linkCode, err := generateLinkCode()
	if err != nil {
		return nil, err
	}

	link := &domain.AffiliateLink{
		ID:             uuid.New().String(),
		AffiliateID:    affiliateID,
		CampaignID:     campaignID,
		LinkCode:       linkCode,
		Name:           name,
		DestinationURL: destinationURL,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

// GetAffiliateLinks retrieves all links for an affiliate
func (s *AffiliateService) GetAffiliateLinks(ctx context.Context, affiliateID string, limit, offset int) ([]*domain.AffiliateLink, error) {
	return s.linkRepo.GetAllByAffiliateID(ctx, affiliateID, limit, offset)
}

// GetAffiliateLinkByCode retrieves an affiliate link by code
func (s *AffiliateService) GetAffiliateLinkByCode(ctx context.Context, code string) (*domain.AffiliateLink, error) {
	link, err := s.linkRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, ErrLinkNotFound
	}
	return link, nil
}

// TrackClick tracks a click on an affiliate link
func (s *AffiliateService) TrackClick(ctx context.Context, linkCode, affiliateCode, ip, userAgent, referrerURL string) (*domain.Referral, error) {
	// Get affiliate
	affiliate, err := s.affiliateRepo.GetByCode(ctx, affiliateCode)
	if err != nil {
		return nil, ErrInvalidAffiliateCode
	}
	if affiliate.Status != "active" {
		return nil, ErrAffiliateNotActive
	}

	// Get link if provided
	var linkID string
	if linkCode != "" {
		link, err := s.linkRepo.GetByCode(ctx, linkCode)
		if err == nil {
			linkID = link.ID
			// Increment link click count
			if err := s.linkRepo.IncrementClick(ctx, link.ID); err != nil {
				// Log error but don't fail
			}
		}
	}

	// Generate referral code
	referralCode, err := generateReferralCode()
	if err != nil {
		return nil, err
	}

	// Create referral
	referral := &domain.Referral{
		ID:              uuid.New().String(),
		AffiliateID:     affiliate.ID,
		AffiliateLinkID: linkID,
		ReferralCode:    referralCode,
		Status:          "pending",
		ClickTimestamp:  time.Now(),
		IPAddress:       ip,
		UserAgent:       userAgent,
		ReferrerURL:     referrerURL,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.referralRepo.Create(ctx, referral); err != nil {
		return nil, err
	}

	// Update affiliate total clicks
	affiliate.TotalClicks++
	affiliate.UpdatedAt = time.Now()
	if err := s.affiliateRepo.Update(ctx, affiliate); err != nil {
		// Log error but don't fail
	}

	return referral, nil
}

// TrackSignup tracks a signup from a referral
func (s *AffiliateService) TrackSignup(ctx context.Context, referralCode, referredUserID string) error {
	referral, err := s.referralRepo.GetByReferralCode(ctx, referralCode)
	if err != nil {
		return ErrReferralNotFound
	}

	now := time.Now()
	referral.Status = "signup"
	referral.ReferredUserID = referredUserID
	referral.SignupTimestamp = &now
	referral.UpdatedAt = now

	if err := s.referralRepo.Update(ctx, referral); err != nil {
		return err
	}

	// Update affiliate stats
	affiliate, err := s.affiliateRepo.GetByID(ctx, referral.AffiliateID)
	if err != nil {
		return err
	}
	affiliate.TotalSignups++
	affiliate.UpdatedAt = time.Now()
	if err := s.affiliateRepo.Update(ctx, affiliate); err != nil {
		return err
	}

	// Increment link signup count if link exists
	if referral.AffiliateLinkID != "" {
		if err := s.linkRepo.IncrementSignup(ctx, referral.AffiliateLinkID); err != nil {
			// Log error but don't fail
		}
	}

	return nil
}

// TrackConversion tracks a conversion from a referral
func (s *AffiliateService) TrackConversion(ctx context.Context, referralCode, orderID string, orderAmount float64) (*domain.Commission, error) {
	referral, err := s.referralRepo.GetByReferralCode(ctx, referralCode)
	if err != nil {
		return nil, ErrReferralNotFound
	}

	// Get affiliate to calculate commission
	affiliate, err := s.affiliateRepo.GetByID(ctx, referral.AffiliateID)
	if err != nil {
		return nil, err
	}
	if affiliate.Status != "active" {
		return nil, ErrAffiliateNotActive
	}

	// Calculate commission amount
	commissionAmount := affiliate.CalculateEarnings(orderAmount)

	// Update referral
	now := time.Now()
	referral.Status = "converted"
	referral.ConversionTimestamp = &now
	referral.OrderID = orderID
	referral.OrderAmount = orderAmount
	referral.CommissionAmount = commissionAmount
	referral.CommissionStatus = "pending"
	referral.UpdatedAt = now

	if err := s.referralRepo.Update(ctx, referral); err != nil {
		return nil, err
	}

	// Create commission record
	commission := &domain.Commission{
		ID:           uuid.New().String(),
		AffiliateID:  affiliate.ID,
		ReferralID:   referral.ID,
		Amount:       commissionAmount,
		Status:       "pending",
		Type:         "referral",
		Description:  "Commission from referral conversion",
		ReferenceID:  orderID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.commissionRepo.Create(ctx, commission); err != nil {
		return nil, err
	}

	// Update affiliate stats
	affiliate.TotalConversions++
	affiliate.TotalEarnings += commissionAmount
	affiliate.Balance += commissionAmount
	affiliate.UpdatedAt = time.Now()
	if err := s.affiliateRepo.Update(ctx, affiliate); err != nil {
		return nil, err
	}

	// Increment link conversion count if link exists
	if referral.AffiliateLinkID != "" {
		if err := s.linkRepo.IncrementConversion(ctx, referral.AffiliateLinkID, commissionAmount); err != nil {
			// Log error but don't fail
		}
	}

	return commission, nil
}

// GetCommissions retrieves commissions for an affiliate
func (s *AffiliateService) GetCommissions(ctx context.Context, affiliateID string, limit, offset int) ([]*domain.Commission, error) {
	return s.commissionRepo.GetAllByAffiliateID(ctx, affiliateID, limit, offset)
}

// GetPendingCommissions retrieves pending commissions for an affiliate
func (s *AffiliateService) GetPendingCommissions(ctx context.Context, affiliateID string) ([]*domain.Commission, error) {
	return s.commissionRepo.GetPendingByAffiliateID(ctx, affiliateID)
}

// ApproveCommission approves a pending commission
func (s *AffiliateService) ApproveCommission(ctx context.Context, id, approvedBy string) error {
	return s.commissionRepo.Approve(ctx, id, approvedBy)
}

// GetAffiliateSettings retrieves affiliate settings for an organization
func (s *AffiliateService) GetAffiliateSettings(ctx context.Context, organizationID string) (*domain.AffiliateSettings, error) {
	settings, err := s.settingsRepo.GetByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, ErrSettingsNotFound
	}
	return settings, nil
}

// UpdateAffiliateSettings updates affiliate settings for an organization
func (s *AffiliateService) UpdateAffiliateSettings(ctx context.Context, settings *domain.AffiliateSettings) error {
	return s.settingsRepo.Update(ctx, settings)
}

// Helper functions to generate codes
func generateAffiliateCode() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "AFF" + hex.EncodeToString(bytes), nil
}

func generateLinkCode() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateReferralCode() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}