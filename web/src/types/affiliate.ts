// Affiliate Types

export interface Affiliate {
  id: string
  userId: string
  affiliateCode: string
  status: 'active' | 'pending' | 'suspended' | 'approved'
  commissionRate: number
  commissionType: 'percentage' | 'fixed'
  fixedCommissionAmount: number
  totalClicks: number
  totalSignups: number
  totalConversions: number
  totalEarnings: number
  totalPaid: number
  balance: number
  minimumPayout: number
  conversionRate: number
  signupRate: number
  canRequestPayout: boolean
  website?: string
  description?: string
  paypalEmail?: string
  createdAt: string
  updatedAt: string
  approvedAt?: string
}

export interface AffiliateLink {
  id: string
  affiliateId: string
  linkCode: string
  trackingUrl: string
  name: string
  destinationUrl: string
  campaignId?: string
  clickCount: number
  signupCount: number
  conversionCount: number
  earnings: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Referral {
  id: string
  affiliateId: string
  referralCode: string
  status: 'pending' | 'signup' | 'qualified' | 'converted' | 'cancelled'
  clickTimestamp: string
  signupTimestamp?: string
  conversionTimestamp?: string
  orderId?: string
  orderAmount: number
  commissionAmount: number
  commissionStatus: 'pending' | 'approved' | 'paid' | 'cancelled' | 'refunded'
  createdAt: string
}

export interface Commission {
  id: string
  affiliateId: string
  referralId?: string
  amount: number
  status: 'pending' | 'approved' | 'paid' | 'cancelled' | 'refunded'
  type: 'referral' | 'bonus' | 'adjustment' | 'referral_bonus'
  description?: string
  referenceId?: string
  createdAt: string
  updatedAt: string
  approvedAt?: string
  paidAt?: string
}

export interface AffiliateStats {
  totalAffiliates: number
  activeAffiliates: number
  totalClicks: number
  totalSignups: number
  totalConversions: number
  conversionRate: number
  totalEarnings: number
  totalPaid: number
  pendingPayouts: number
  topAffiliates: AffiliateLeaderboardEntry[]
}

export interface AffiliateLeaderboardEntry {
  affiliateId: string
  affiliateCode: string
  name: string
  conversions: number
  earnings: number
}

export interface AffiliateDashboard {
  affiliate: Affiliate
  recentReferrals: Referral[]
  recentCommissions: Commission[]
  topLinks: AffiliateLink[]
  earningsHistory: MonthlyEarnings[]
}

export interface MonthlyEarnings {
  month: string
  earnings: number
  conversions: number
}

export interface AffiliateSettings {
  id: string
  organizationId: string
  programName: string
  programDescription: string
  defaultCommissionRate: number
  defaultCommissionType: 'percentage' | 'fixed'
  defaultFixedCommission: number
  cookieDurationDays: number
  minimumPayoutAmount: number
  payoutSchedule: 'weekly' | 'biweekly' | 'monthly' | 'quarterly'
  requireApproval: boolean
  allowCustomCommissionRates: boolean
  allowAffiliateLinks: boolean
  termsAndConditions: string
  welcomeMessage: string
}

// Request Types

export interface RegisterAffiliateRequest {
  website?: string
  description?: string
}

export interface UpdateAffiliateRequest {
  commissionRate?: number
  commissionType?: 'percentage' | 'fixed'
  fixedCommissionAmount?: number
  paypalEmail?: string
  bankName?: string
  bankAccount?: string
  bankRouting?: string
  taxId?: string
  website?: string
  description?: string
}

export interface CreateAffiliateLinkRequest {
  name: string
  destinationUrl: string
  campaignId?: string
}

export interface TrackClickRequest {
  affiliateCode: string
  linkCode?: string
  ipAddress?: string
  userAgent?: string
  referrerUrl?: string
}

export interface TrackSignupRequest {
  referralCode: string
  referredUserId: string
}

export interface TrackConversionRequest {
  referralCode: string
  orderId: string
  orderAmount: number
}

export interface UpdateAffiliateSettingsRequest {
  programName?: string
  programDescription?: string
  defaultCommissionRate?: number
  defaultCommissionType?: 'percentage' | 'fixed'
  defaultFixedCommission?: number
  cookieDurationDays?: number
  minimumPayoutAmount?: number
  payoutSchedule?: 'weekly' | 'biweekly' | 'monthly' | 'quarterly'
  requireApproval?: boolean
  allowCustomCommissionRates?: boolean
  allowAffiliateLinks?: boolean
  termsAndConditions?: string
  welcomeMessage?: string
}