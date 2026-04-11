import { apiGo } from './api'
import {
  EmailList,
  EmailSubscriber,
  EmailCampaign,
  SocialPost,
  LandingPage,
  AutomationWorkflow,
  Influencer,
  MarketingOverview,
} from '@/types/marketing-features'

// Marketing Features API

export const marketingApi = {
  // Email Marketing
  getEmailLists: async (): Promise<{ lists: EmailList[]; total: number }> => {
    const response = await apiGo.get('/api/v1/email/lists')
    return response.data
  },

  createEmailList: async (data: { name: string; description?: string }): Promise<EmailList> => {
    const response = await apiGo.post('/api/v1/email/lists', data)
    return response.data
  },

  getSubscribers: async (limit = 20, offset = 0): Promise<{ subscribers: EmailSubscriber[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/email/subscribers?limit=${limit}&offset=${offset}`)
    return response.data
  },

  addSubscriber: async (data: { email: string; first_name?: string; last_name?: string; source?: string; list_ids?: string[] }): Promise<EmailSubscriber> => {
    const response = await apiGo.post('/api/v1/email/subscribers', data)
    return response.data
  },

  getEmailCampaigns: async (limit = 20, offset = 0): Promise<{ campaigns: EmailCampaign[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/email/campaigns?limit=${limit}&offset=${offset}`)
    return response.data
  },

  createEmailCampaign: async (data: { name: string; subject: string; from_name?: string; from_email?: string; html_content?: string; list_ids?: string[] }): Promise<EmailCampaign> => {
    const response = await apiGo.post('/api/v1/email/campaigns', data)
    return response.data
  },

  // Social Media
  getSocialPosts: async (limit = 20, offset = 0): Promise<{ posts: SocialPost[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/social/posts?limit=${limit}&offset=${offset}`)
    return response.data
  },

  createSocialPost: async (data: { content: string; media_urls?: string[]; hashtags?: string[] }): Promise<SocialPost> => {
    const response = await apiGo.post('/api/v1/social/posts', data)
    return response.data
  },

  scheduleSocialPost: async (postId: string, data: { account_id: string; scheduled_time: string }): Promise<{ message: string }> => {
    const response = await apiGo.post(`/api/v1/social/posts/${postId}/schedule`, data)
    return response.data
  },

  // Landing Pages
  getLandingPages: async (limit = 20, offset = 0): Promise<{ pages: LandingPage[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/landing-pages?limit=${limit}&offset=${offset}`)
    return response.data
  },

  createLandingPage: async (data: { name: string; slug?: string; title?: string; html_content?: string }): Promise<LandingPage> => {
    const response = await apiGo.post('/api/v1/landing-pages', data)
    return response.data
  },

  createLead: async (pageId: string, data: { email: string; first_name?: string; last_name?: string; phone?: string; custom_data?: Record<string, unknown> }): Promise<unknown> => {
    const response = await apiGo.post(`/api/v1/landing-pages/${pageId}/leads`, data)
    return response.data
  },

  // Automation
  getWorkflows: async (): Promise<{ workflows: AutomationWorkflow[]; total: number }> => {
    const response = await apiGo.get('/api/v1/automation/workflows')
    return response.data
  },

  createWorkflow: async (data: { name: string; trigger_type: string; trigger_config?: Record<string, unknown> }): Promise<AutomationWorkflow> => {
    const response = await apiGo.post('/api/v1/automation/workflows', data)
    return response.data
  },

  // Influencers
  getInfluencers: async (limit = 20, offset = 0): Promise<{ influencers: Influencer[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/influencers?limit=${limit}&offset=${offset}`)
    return response.data
  },

  createInfluencer: async (data: { name: string; email?: string; category?: string; niches?: string[] }): Promise<Influencer> => {
    const response = await apiGo.post('/api/v1/influencers', data)
    return response.data
  },

  createInfluencerCampaign: async (data: { name: string; campaign_type?: string; budget?: number }): Promise<unknown> => {
    const response = await apiGo.post('/api/v1/influencers/campaigns', data)
    return response.data
  },

  // Reporting
  getMarketingOverview: async (): Promise<MarketingOverview> => {
    const response = await apiGo.get('/api/v1/reporting/overview')
    return response.data
  },

  trackEvent: async (data: { event_type: string; event_category?: string; event_action?: string; event_data?: Record<string, unknown> }): Promise<{ message: string }> => {
    const response = await apiGo.post('/api/v1/reporting/events', data)
    return response.data
  },
}