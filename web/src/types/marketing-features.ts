// Marketing Features Types

export interface EmailList {
  id: string
  name: string
  description: string
  status: string
  subscriber_count: number
  created_at: string
}

export interface EmailSubscriber {
  id: string
  email: string
  first_name: string
  last_name: string
  status: string
  source: string
  created_at: string
}

export interface EmailCampaign {
  id: string
  name: string
  subject: string
  status: string
  recipient_count: number
  sent_count: number
  opened_count: number
  clicked_count: number
  open_rate: number
  click_rate: number
  created_at: string
}

export interface SocialPost {
  id: string
  content: string
  media_urls: string[]
  hashtags: string[]
  status: string
  schedules: SocialSchedule[]
  created_at: string
}

export interface SocialSchedule {
  account_id: string
  scheduled_time: string
  status: string
}

export interface LandingPage {
  id: string
  name: string
  slug: string
  title: string
  status: string
  view_count: number
  conversion_count: number
  created_at: string
}

export interface AutomationWorkflow {
  id: string
  name: string
  trigger_type: string
  status: string
  total_contacts: number
  active_contacts: number
  created_at: string
}

export interface Influencer {
  id: string
  name: string
  email: string
  category: string
  niches: string[]
  status: string
  rating: number
  total_collaborations: number
  created_at: string
}

export interface MarketingOverview {
  email: {
    total_campaigns: number
    total_sent: number
    total_opens: number
    total_clicks: number
    avg_open_rate: number
    avg_click_rate: number
    total_subscribers: number
  }
  social: {
    total_posts: number
    scheduled_posts: number
    published_posts: number
    total_impressions: number
    total_engagement: number
  }
  landing_pages: {
    total_pages: number
    active_pages: number
    total_views: number
    total_leads: number
  }
  influencers: {
    total_influencers: number
    active_campaigns: number
    total_collaborations: number
    total_reach: number
    total_spend: number
  }
}