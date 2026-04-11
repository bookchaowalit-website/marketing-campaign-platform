package domain

import "time"

// Affiliate represents an affiliate partner account
type Affiliate struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	OrganizationID      string     `json:"organization_id"`
	AffiliateCode       string     `json:"affiliate_code"`
	Status              string     `json:"status"` // active, pending, suspended, approved
	CommissionRate      float64    `json:"commission_rate"`
	CommissionType      string     `json:"commission_type"` // percentage, fixed
	FixedCommissionAmount float64  `json:"fixed_commission_amount"`
	TotalClicks         int64      `json:"total_clicks"`
	TotalSignups        int64      `json:"total_signups"`
	TotalConversions    int64      `json:"total_conversions"`
	TotalEarnings       float64    `json:"total_earnings"`
	TotalPaid           float64    `json:"total_paid"`
	Balance             float64    `json:"balance"`
	MinimumPayout       float64    `json:"minimum_payout"`
	PaypalEmail         string     `json:"paypal_email,omitempty"`
	BankName            string     `json:"bank_name,omitempty"`
	BankAccount         string     `json:"bank_account,omitempty"`
	BankRouting         string     `json:"bank_routing,omitempty"`
	TaxID               string     `json:"tax_id,omitempty"`
	Website             string     `json:"website,omitempty"`
	Description         string     `json:"description,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	ApprovedAt          *time.Time `json:"approved_at,omitempty"`
	ApprovedBy          string     `json:"approved_by,omitempty"`
}

// AffiliateLink represents a tracking link for an affiliate
type AffiliateLink struct {
	ID               string     `json:"id"`
	AffiliateID      string     `json:"affiliate_id"`
	CampaignID       string     `json:"campaign_id,omitempty"`
	LinkCode         string     `json:"link_code"`
	Name             string     `json:"name"`
	DestinationURL   string     `json:"destination_url"`
	ClickCount       int64      `json:"click_count"`
	SignupCount      int64      `json:"signup_count"`
	ConversionCount  int64      `json:"conversion_count"`
	Earnings         float64    `json:"earnings"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Referral represents an individual referral tracked through an affiliate link
type Referral struct {
	ID                string     `json:"id"`
	AffiliateID       string     `json:"affiliate_id"`
	AffiliateLinkID   string     `json:"affiliate_link_id,omitempty"`
	ReferredUserID    string     `json:"referred_user_id,omitempty"`
	ReferralCode      string     `json:"referral_code"`
	Status            string     `json:"status"` // pending, signup, qualified, converted, cancelled
	ClickTimestamp    time.Time  `json:"click_timestamp"`
	SignupTimestamp   *time.Time `json:"signup_timestamp,omitempty"`
	ConversionTimestamp *time.Time `json:"conversion_timestamp,omitempty"`
	OrderID           string     `json:"order_id,omitempty"`
	OrderAmount       float64    `json:"order_amount"`
	CommissionAmount  float64    `json:"commission_amount"`
	CommissionStatus  string     `json:"commission_status"` // pending, approved, paid, cancelled, refunded
	IPAddress         string     `json:"ip_address,omitempty"`
	UserAgent         string     `json:"user_agent,omitempty"`
	ReferrerURL       string     `json:"referrer_url,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Commission represents a commission record for an affiliate
type Commission struct {
	ID           string     `json:"id"`
	AffiliateID  string     `json:"affiliate_id"`
	ReferralID   string     `json:"referral_id,omitempty"`
	Amount       float64    `json:"amount"`
	Status       string     `json:"status"` // pending, approved, paid, cancelled, refunded
	Type         string     `json:"type"`   // referral, bonus, adjustment, referral_bonus
	Description  string     `json:"description,omitempty"`
	ReferenceID  string     `json:"reference_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ApprovedBy   string     `json:"approved_by,omitempty"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
	PayoutID     string     `json:"payout_id,omitempty"`
}

// AffiliatePayout represents a payment to an affiliate
type AffiliatePayout struct {
	ID              string     `json:"id"`
	AffiliateID     string     `json:"affiliate_id"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"` // pending, processing, completed, failed, cancelled
	PaymentMethod   string     `json:"payment_method"` // paypal, bank_transfer, stripe, manual
	PaymentReference string    `json:"payment_reference,omitempty"`
	PaymentNotes    string     `json:"payment_notes,omitempty"`
	CommissionIDs   []string   `json:"commission_ids,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	ProcessedBy     string     `json:"processed_by,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

// AffiliateSettings represents affiliate program configuration for an organization
type AffiliateSettings struct {
	ID                      string    `json:"id"`
	OrganizationID          string    `json:"organization_id"`
	ProgramName             string    `json:"program_name"`
	ProgramDescription      string    `json:"program_description"`
	DefaultCommissionRate   float64   `json:"default_commission_rate"`
	DefaultCommissionType   string    `json:"default_commission_type"` // percentage, fixed
	DefaultFixedCommission  float64   `json:"default_fixed_commission"`
	CookieDurationDays      int       `json:"cookie_duration_days"`
	MinimumPayoutAmount     float64   `json:"minimum_payout_amount"`
	PayoutSchedule          string    `json:"payout_schedule"` // weekly, biweekly, monthly, quarterly
	RequireApproval         bool      `json:"require_approval"`
	AllowCustomCommissionRates bool   `json:"allow_custom_commission_rates"`
	AllowAffiliateLinks     bool      `json:"allow_affiliate_links"`
	TermsAndConditions      string    `json:"terms_and_conditions"`
	WelcomeMessage          string    `json:"welcome_message"`
	EmailNotifications      bool      `json:"email_notifications"`
	NotifyOnSignup          bool      `json:"notify_on_signup"`
	NotifyOnConversion      bool      `json:"notify_on_conversion"`
	NotifyOnPayout          bool      `json:"notify_on_payout"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// AffiliateStats represents statistics for an affiliate
type AffiliateStats struct {
	TotalAffiliates      int64   `json:"total_affiliates"`
	ActiveAffiliates     int64   `json:"active_affiliates"`
	TotalClicks          int64   `json:"total_clicks"`
	TotalSignups         int64   `json:"total_signups"`
	TotalConversions     int64   `json:"total_conversions"`
	ConversionRate       float64 `json:"conversion_rate"`
	TotalEarnings        float64 `json:"total_earnings"`
	TotalPaid            float64 `json:"total_paid"`
	PendingPayouts       float64 `json:"pending_payouts"`
	TopAffiliates        []AffiliateLeaderboard `json:"top_affiliates"`
}

// AffiliateLeaderboard represents an affiliate on the leaderboard
type AffiliateLeaderboard struct {
	AffiliateID   string  `json:"affiliate_id"`
	AffiliateCode string  `json:"affiliate_code"`
	Name          string  `json:"name"`
	Conversions   int64   `json:"conversions"`
	Earnings      float64 `json:"earnings"`
}

// AffiliateDashboard represents dashboard data for an affiliate
type AffiliateDashboard struct {
	Affiliate         Affiliate          `json:"affiliate"`
	RecentReferrals   []Referral         `json:"recent_referrals"`
	RecentCommissions []Commission       `json:"recent_commissions"`
	TopLinks          []AffiliateLink    `json:"top_links"`
	EarningsHistory   []MonthlyEarnings  `json:"earnings_history"`
}

// MonthlyEarnings represents earnings for a specific month
type MonthlyEarnings struct {
	Month     string  `json:"month"`
	Earnings  float64 `json:"earnings"`
	Conversions int64 `json:"conversions"`
}

// CalculateEarnings calculates the commission amount based on order amount
func (a *Affiliate) CalculateEarnings(orderAmount float64) float64 {
	if a.CommissionType == "fixed" {
		return a.FixedCommissionAmount
	}
	return orderAmount * a.CommissionRate
}

// ConversionRate calculates the affiliate's conversion rate
func (a *Affiliate) ConversionRate() float64 {
	if a.TotalClicks == 0 {
		return 0
	}
	return float64(a.TotalConversions) / float64(a.TotalClicks) * 100
}

// SignupRate calculates the signup rate from clicks
func (a *Affiliate) SignupRate() float64 {
	if a.TotalClicks == 0 {
		return 0
	}
	return float64(a.TotalSignups) / float64(a.TotalClicks) * 100
}

// CanRequestPayout checks if affiliate can request a payout
func (a *Affiliate) CanRequestPayout() bool {
	return a.Balance >= a.MinimumPayout && a.Status == "active"
}