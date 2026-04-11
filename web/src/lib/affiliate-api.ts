import { apiGo } from './api'
import {
  Affiliate,
  AffiliateLink,
  AffiliateStats,
  AffiliateDashboard,
  AffiliateSettings,
  Commission,
  RegisterAffiliateRequest,
  UpdateAffiliateRequest,
  CreateAffiliateLinkRequest,
  TrackClickRequest,
  TrackSignupRequest,
  TrackConversionRequest,
  UpdateAffiliateSettingsRequest,
} from '@/types/affiliate'

// Affiliate API

export const affiliateApi = {
  // Register as an affiliate
  register: async (data: RegisterAffiliateRequest): Promise<Affiliate> => {
    const response = await apiGo.post('/api/v1/affiliates/register', data)
    return response.data
  },

  // Get current user's affiliate account
  getMyAffiliate: async (): Promise<Affiliate> => {
    const response = await apiGo.get('/api/v1/affiliates/me')
    return response.data
  },

  // Get affiliate by ID
  getAffiliate: async (id: string): Promise<Affiliate> => {
    const response = await apiGo.get(`/api/v1/affiliates/${id}`)
    return response.data
  },

  // Get all affiliates (admin)
  getAllAffiliates: async (limit = 20, offset = 0, status?: string): Promise<{ affiliates: Affiliate[]; total: number }> => {
    const params = new URLSearchParams({ limit: String(limit), offset: String(offset) })
    if (status) params.append('status', status)
    const response = await apiGo.get(`/api/v1/affiliates?${params.toString()}`)
    return response.data
  },

  // Update affiliate
  updateAffiliate: async (id: string, data: UpdateAffiliateRequest): Promise<Affiliate> => {
    const response = await apiGo.put(`/api/v1/affiliates/${id}`, data)
    return response.data
  },

  // Approve affiliate (admin)
  approveAffiliate: async (id: string): Promise<{ message: string }> => {
    const response = await apiGo.post(`/api/v1/affiliates/${id}/approve`)
    return response.data
  },

  // Suspend affiliate (admin)
  suspendAffiliate: async (id: string): Promise<{ message: string }> => {
    const response = await apiGo.post(`/api/v1/affiliates/${id}/suspend`)
    return response.data
  },

  // Get affiliate stats (admin)
  getStats: async (): Promise<AffiliateStats> => {
    const response = await apiGo.get('/api/v1/affiliates/stats')
    return response.data
  },

  // Get affiliate dashboard
  getDashboard: async (): Promise<AffiliateDashboard> => {
    const response = await apiGo.get('/api/v1/affiliates/dashboard')
    return response.data
  },

  // Affiliate Links

  // Create affiliate link
  createLink: async (data: CreateAffiliateLinkRequest): Promise<AffiliateLink> => {
    const response = await apiGo.post('/api/v1/affiliates/links', data)
    return response.data
  },

  // Get affiliate links
  getLinks: async (limit = 20, offset = 0): Promise<{ links: AffiliateLink[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/affiliates/links?limit=${limit}&offset=${offset}`)
    return response.data
  },

  // Delete affiliate link
  deleteLink: async (id: string): Promise<{ message: string }> => {
    const response = await apiGo.delete(`/api/v1/affiliates/links/${id}`)
    return response.data
  },

  // Commissions

  // Get commissions
  getCommissions: async (limit = 20, offset = 0, status?: string): Promise<{ commissions: Commission[]; total: number }> => {
    const params = new URLSearchParams({ limit: String(limit), offset: String(offset) })
    if (status) params.append('status', status)
    const response = await apiGo.get(`/api/v1/affiliates/commissions?${params.toString()}`)
    return response.data
  },

  // Affiliate Settings (admin)

  // Get affiliate settings
  getSettings: async (): Promise<AffiliateSettings> => {
    const response = await apiGo.get('/api/v1/affiliates/settings')
    return response.data
  },

  // Update affiliate settings
  updateSettings: async (data: UpdateAffiliateSettingsRequest): Promise<AffiliateSettings> => {
    const response = await apiGo.put('/api/v1/affiliates/settings', data)
    return response.data
  },

  // Public Tracking Endpoints

  // Track click (public)
  trackClick: async (data: TrackClickRequest): Promise<{ referralCode: string; message: string }> => {
    const response = await apiGo.post('/api/v1/affiliates/track/click', data)
    return response.data
  },

  // Track signup (internal)
  trackSignup: async (data: TrackSignupRequest): Promise<{ message: string }> => {
    const response = await apiGo.post('/api/v1/affiliates/track/signup', data)
    return response.data
  },

  // Track conversion (internal)
  trackConversion: async (data: TrackConversionRequest): Promise<{ commissionId: string; commissionAmount: number; message: string }> => {
    const response = await apiGo.post('/api/v1/affiliates/track/conversion', data)
    return response.data
  },
}