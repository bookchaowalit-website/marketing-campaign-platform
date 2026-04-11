// Shared TypeScript types for the Marketing SaaS Platform
// These types are used across frontend and backend services

export interface User {
  id: string;
  email: string;
  full_name: string;
  role: string;
  is_active: boolean;
  email_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

export interface RefreshTokenResponse {
  token: string;
  refresh_token: string;
  expires_in: number;
}

export interface Campaign {
  id: string;
  name: string;
  type: 'SEO' | 'SEM' | 'Social' | 'Email' | 'Display';
  status: 'active' | 'paused' | 'completed' | 'draft';
  budget: number;
  spent: number;
  impressions: number;
  clicks: number;
  conversions: number;
  ctr: number;
  conversion_rate: number;
  budget_utilization: number;
  start_date: string;
  end_date?: string;
  created_at: string;
  updated_at: string;
}

export interface ABTest {
  id: string;
  name: string;
  variant: 'A' | 'B';
  status: 'running' | 'completed' | 'draft';
  visitors: number;
  conversions: number;
  conversion_rate: number;
  confidence: number;
  start_date: string;
  end_date?: string;
  created_at: string;
  updated_at: string;
}

export interface Competitor {
  id: string;
  competitor_name: string;
  domain: string;
  type: 'direct' | 'indirect' | 'emerging';
  market_share: number;
  seo_score: number;
  sem_spend: number;
  keywords: number;
  backlinks: number;
  social_followers: number;
  last_analyzed?: string;
  created_at: string;
  updated_at: string;
}

export interface DashboardOverview {
  total_campaigns: number;
  active_campaigns: number;
  total_spend: number;
  total_revenue: number;
  roi: number;
  conversion_rate: number;
  top_channels: ChannelPerformance[];
}

export interface ChannelPerformance {
  name: string;
  value: number;
}

export interface SEOMetrics {
  domain: string;
  organic_traffic: number;
  keywords: number;
  avg_position: number;
  domain_authority: number;
  page_speed: number;
  backlinks: number;
  referring_domains: number;
}

export interface SEOMetricsResponse {
  metrics: SEOMetrics;
  recommendations: string[];
}

export interface Affiliate {
  id: string;
  user_id: string;
  affiliate_code: string;
  status: 'pending' | 'approved' | 'suspended';
  commission_rate: number;
  commission_type: 'percentage' | 'fixed';
  total_clicks: number;
  total_signups: number;
  total_conversions: number;
  total_earnings: number;
  created_at: string;
  updated_at: string;
}

export interface AffiliateDashboard {
  affiliate: Affiliate;
  total_earnings: number;
  pending_commissions: number;
  recent_referrals: any[];
  monthly_earnings: MonthlyEarnings[];
}

export interface MonthlyEarnings {
  month: string;
  earnings: number;
}

export interface EmailList {
  id: string;
  name: string;
  subscriber_count: number;
  status: 'active' | 'inactive';
  created_at: string;
}

export interface EmailCampaign {
  id: string;
  name: string;
  subject: string;
  status: 'draft' | 'scheduled' | 'sending' | 'sent' | 'failed';
  sent_count: number;
  open_count: number;
  click_count: number;
  open_rate: number;
  click_rate: number;
  created_at: string;
}

export interface SocialPost {
  id: string;
  content: string;
  status: 'draft' | 'scheduled' | 'published' | 'failed';
  scheduled_at?: string;
  published_at?: string;
  created_at: string;
}

export interface LandingPage {
  id: string;
  name: string;
  slug: string;
  status: 'draft' | 'published' | 'archived';
  view_count: number;
  lead_count: number;
  conversion_rate: number;
  created_at: string;
}

export interface APIError {
  error: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  has_next: boolean;
  has_prev: boolean;
}

export type Platform = 'shopee' | 'tiktok' | 'lazada' | 'amazon' | 'custom';

export interface ExternalAccount {
  id: string;
  platform: Platform;
  status: 'connected' | 'disconnected' | 'error';
  last_sync?: string;
  created_at: string;
}

export interface MarketingOverview {
  email: {
    total_campaigns: number;
    total_sent: number;
    total_opens: number;
    total_clicks: number;
    avg_open_rate: number;
    avg_click_rate: number;
    total_subscribers: number;
  };
  social: {
    total_posts: number;
    scheduled_posts: number;
    published_posts: number;
    total_impressions: number;
    total_engagement: number;
  };
  landing_pages: {
    total_pages: number;
    active_pages: number;
    total_views: number;
    total_leads: number;
    avg_conversion_rate: number;
  };
  influencers: {
    total_influencers: number;
    active_campaigns: number;
    total_collaborations: number;
    total_reach: number;
    total_spend: number;
  };
}
