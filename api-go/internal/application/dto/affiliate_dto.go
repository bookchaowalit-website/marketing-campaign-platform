package dto

// RegisterAffiliateRequest represents the request to register as an affiliate
type RegisterAffiliateRequest struct {
	Website     string `json:"website,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateAffiliateRequest represents the request to update an affiliate
type UpdateAffiliateRequest struct {
	CommissionRate       float64 `json:"commission_rate,omitempty"`
	CommissionType       string  `json:"commission_type,omitempty"` // percentage, fixed
	FixedCommissionAmount float64 `json:"fixed_commission_amount,omitempty"`
	PaypalEmail          string  `json:"paypal_email,omitempty"`
	BankName             string  `json:"bank_name,omitempty"`
	BankAccount          string  `json:"bank_account,omitempty"`
	BankRouting          string  `json:"bank_routing,omitempty"`
	TaxID                string  `json:"tax_id,omitempty"`
	Website              string  `json:"website,omitempty"`
	Description          string  `json:"description,omitempty"`
}

// AffiliateResponse represents the affiliate response
type AffiliateResponse struct {
	ID                    string  `json:"id"`
	UserID                string  `json:"user_id"`
	AffiliateCode         string  `json:"affiliate_code"`
	Status                string  `json:"status"`
	CommissionRate        float64 `json:"commission_rate"`
	CommissionType        string  `json:"commission_type"`
	FixedCommissionAmount float64 `json:"fixed_commission_amount"`
	TotalClicks           int64   `json:"total_clicks"`
	TotalSignups          int64   `json:"total_signups"`
	TotalConversions      int64   `json:"total_conversions"`
	TotalEarnings         float64 `json:"total_earnings"`
	TotalPaid             float64 `json:"total_paid"`
	Balance               float64 `json:"balance"`
	MinimumPayout         float64 `json:"minimum_payout"`
	ConversionRate        float64 `json:"conversion_rate"`
	SignupRate            float64 `json:"signup_rate"`
	CanRequestPayout      bool    `json:"can_request_payout"`
	Website               string  `json:"website,omitempty"`
	Description           string  `json:"description,omitempty"`
	PaypalEmail           string  `json:"paypal_email,omitempty"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
	ApprovedAt            string  `json:"approved_at,omitempty"`
}

// CreateAffiliateLinkRequest represents the request to create an affiliate link
type CreateAffiliateLinkRequest struct {
	Name           string `json:"name" binding:"required"`
	DestinationURL string `json:"destination_url" binding:"required"`
	CampaignID     string `json:"campaign_id,omitempty"`
}

// AffiliateLinkResponse represents the affiliate link response
type AffiliateLinkResponse struct {
	ID              string  `json:"id"`
	AffiliateID     string  `json:"affiliate_id"`
	LinkCode        string  `json:"link_code"`
	TrackingURL     string  `json:"tracking_url"`
	Name            string  `json:"name"`
	DestinationURL  string  `json:"destination_url"`
	CampaignID      string  `json:"campaign_id,omitempty"`
	ClickCount      int64   `json:"click_count"`
	SignupCount     int64   `json:"signup_count"`
	ConversionCount int64   `json:"conversion_count"`
	Earnings        float64 `json:"earnings"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// ReferralResponse represents the referral response
type ReferralResponse struct {
	ID                string  `json:"id"`
	AffiliateID       string  `json:"affiliate_id"`
	ReferralCode      string  `json:"referral_code"`
	Status            string  `json:"status"`
	ClickTimestamp    string  `json:"click_timestamp"`
	SignupTimestamp   string  `json:"signup_timestamp,omitempty"`
	ConversionTimestamp string `json:"conversion_timestamp,omitempty"`
	OrderID           string  `json:"order_id,omitempty"`
	OrderAmount       float64 `json:"order_amount"`
	CommissionAmount  float64 `json:"commission_amount"`
	CommissionStatus  string  `json:"commission_status"`
	CreatedAt         string  `json:"created_at"`
}

// CommissionResponse represents the commission response
type CommissionResponse struct {
	ID           string  `json:"id"`
	AffiliateID  string  `json:"affiliate_id"`
	ReferralID   string  `json:"referral_id,omitempty"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
	Type         string  `json:"type"`
	Description  string  `json:"description,omitempty"`
	ReferenceID  string  `json:"reference_id,omitempty"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	ApprovedAt   string  `json:"approved_at,omitempty"`
	PaidAt       string  `json:"paid_at,omitempty"`
}

// AffiliateStatsResponse represents the affiliate statistics response
type AffiliateStatsResponse struct {
	TotalAffiliates  int64                   `json:"total_affiliates"`
	ActiveAffiliates int64                   `json:"active_affiliates"`
	TotalClicks      int64                   `json:"total_clicks"`
	TotalSignups     int64                   `json:"total_signups"`
	TotalConversions int64                   `json:"total_conversions"`
	ConversionRate   float64                 `json:"conversion_rate"`
	TotalEarnings    float64                 `json:"total_earnings"`
	TotalPaid        float64                 `json:"total_paid"`
	PendingPayouts   float64                 `json:"pending_payouts"`
	TopAffiliates    []AffiliateLeaderboardResponse `json:"top_affiliates"`
}

// AffiliateLeaderboardResponse represents an affiliate on the leaderboard
type AffiliateLeaderboardResponse struct {
	AffiliateID   string  `json:"affiliate_id"`
	AffiliateCode string  `json:"affiliate_code"`
	Name          string  `json:"name"`
	Conversions   int64   `json:"conversions"`
	Earnings      float64 `json:"earnings"`
}

// AffiliateDashboardResponse represents the affiliate dashboard response
type AffiliateDashboardResponse struct {
	Affiliate         AffiliateResponse      `json:"affiliate"`
	RecentReferrals   []ReferralResponse     `json:"recent_referrals"`
	RecentCommissions []CommissionResponse   `json:"recent_commissions"`
	TopLinks          []AffiliateLinkResponse `json:"top_links"`
	EarningsHistory   []MonthlyEarningsResponse `json:"earnings_history"`
}

// MonthlyEarningsResponse represents monthly earnings
type MonthlyEarningsResponse struct {
	Month       string  `json:"month"`
	Earnings    float64 `json:"earnings"`
	Conversions int64   `json:"conversions"`
}

// AffiliateSettingsResponse represents affiliate program settings
type AffiliateSettingsResponse struct {
	ID                      string `json:"id"`
	OrganizationID          string `json:"organization_id"`
	ProgramName             string `json:"program_name"`
	ProgramDescription      string `json:"program_description"`
	DefaultCommissionRate   float64 `json:"default_commission_rate"`
	DefaultCommissionType   string `json:"default_commission_type"`
	DefaultFixedCommission  float64 `json:"default_fixed_commission"`
	CookieDurationDays      int    `json:"cookie_duration_days"`
	MinimumPayoutAmount     float64 `json:"minimum_payout_amount"`
	PayoutSchedule          string `json:"payout_schedule"`
	RequireApproval         bool   `json:"require_approval"`
	AllowCustomCommissionRates bool `json:"allow_custom_commission_rates"`
	AllowAffiliateLinks     bool   `json:"allow_affiliate_links"`
	TermsAndConditions      string `json:"terms_and_conditions"`
	WelcomeMessage          string `json:"welcome_message"`
}

// UpdateAffiliateSettingsRequest represents the request to update affiliate settings
type UpdateAffiliateSettingsRequest struct {
	ProgramName             string  `json:"program_name,omitempty"`
	ProgramDescription      string  `json:"program_description,omitempty"`
	DefaultCommissionRate   float64 `json:"default_commission_rate,omitempty"`
	DefaultCommissionType   string  `json:"default_commission_type,omitempty"`
	DefaultFixedCommission  float64 `json:"default_fixed_commission,omitempty"`
	CookieDurationDays      int     `json:"cookie_duration_days,omitempty"`
	MinimumPayoutAmount     float64 `json:"minimum_payout_amount,omitempty"`
	PayoutSchedule          string  `json:"payout_schedule,omitempty"`
	RequireApproval         bool    `json:"require_approval,omitempty"`
	AllowCustomCommissionRates bool  `json:"allow_custom_commission_rates,omitempty"`
	AllowAffiliateLinks     bool    `json:"allow_affiliate_links,omitempty"`
	TermsAndConditions      string  `json:"terms_and_conditions,omitempty"`
	WelcomeMessage          string  `json:"welcome_message,omitempty"`
}

// TrackClickRequest represents the request to track a click
type TrackClickRequest struct {
	LinkCode     string `json:"link_code,omitempty"`
	AffiliateCode string `json:"affiliate_code" binding:"required"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	ReferrerURL  string `json:"referrer_url,omitempty"`
}

// TrackSignupRequest represents the request to track a signup
type TrackSignupRequest struct {
	ReferralCode    string `json:"referral_code" binding:"required"`
	ReferredUserID  string `json:"referred_user_id" binding:"required"`
}

// TrackConversionRequest represents the request to track a conversion
type TrackConversionRequest struct {
	ReferralCode string  `json:"referral_code" binding:"required"`
	OrderID      string  `json:"order_id" binding:"required"`
	OrderAmount  float64 `json:"order_amount" binding:"required,min=0"`
}