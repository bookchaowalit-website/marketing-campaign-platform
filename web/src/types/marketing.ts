export interface MarketingCampaign {
  id: string
  name: string
  type: 'SEO' | 'SEM' | 'Social' | 'Email' | 'Display'
  status: 'active' | 'paused' | 'completed' | 'draft'
  budget: number
  spent: number
  impressions: number
  clicks: number
  conversions: number
  startDate: string
  endDate?: string
  createdAt: string
  updatedAt: string
}

export interface CompetitorAnalysis {
  id: string
  competitorName: string
  domain: string
  marketShare: number
  seoScore: number
  semSpend: number
  keywords: number
  backlinks: number
  socialFollowers: number
  lastAnalyzed: string
}

export interface ABTest {
  id: string
  name: string
  variant: 'A' | 'B'
  status: 'running' | 'completed' | 'draft'
  visitors: number
  conversions: number
  conversionRate: number
  confidence: number
  startDate: string
  endDate?: string
}

export interface SEOMetrics {
  domain: string
  organicTraffic: number
  keywords: number
  avgPosition: number
  backlinks: number
  domainAuthority: number
  pageSpeed: number
  lastUpdated: string
}

export interface MarketingDashboard {
  totalCampaigns: number
  activeCampaigns: number
  totalSpend: number
  totalRevenue: number
  roi: number
  conversionRate: number
  topChannels: Array<{ name: string; value: number }>
  recentCampaigns: MarketingCampaign[]
}
