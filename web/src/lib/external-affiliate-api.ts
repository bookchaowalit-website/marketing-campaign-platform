import { apiGo } from './api'
import {
  ExternalAccount,
  ExternalProduct,
  ExternalOrder,
  UnifiedDashboard,
  SyncLog,
  ConnectPlatformRequest,
  SyncAccountRequest,
  GenerateLinkRequest,
} from '@/types/external-affiliate'

// External Affiliate API

export const externalAffiliateApi = {
  // Platform Connections

  // Connect a platform account
  connect: async (data: ConnectPlatformRequest): Promise<ExternalAccount> => {
    const response = await apiGo.post('/api/v1/integrations/connect', data)
    return response.data
  },

  // Disconnect a platform account
  disconnect: async (id: string): Promise<{ message: string }> => {
    const response = await apiGo.post(`/api/v1/integrations/accounts/${id}/disconnect`)
    return response.data
  },

  // Get all connected accounts
  getAccounts: async (): Promise<{ accounts: ExternalAccount[]; total: number }> => {
    const response = await apiGo.get('/api/v1/integrations/accounts')
    return response.data
  },

  // Get account by ID
  getAccount: async (id: string): Promise<ExternalAccount> => {
    const response = await apiGo.get(`/api/v1/integrations/accounts/${id}`)
    return response.data
  },

  // Sync account data
  sync: async (id: string, data?: SyncAccountRequest): Promise<{ message: string }> => {
    const response = await apiGo.post(`/api/v1/integrations/accounts/${id}/sync`, data || {})
    return response.data
  },

  // Unified Dashboard

  // Get unified dashboard
  getDashboard: async (): Promise<UnifiedDashboard> => {
    const response = await apiGo.get('/api/v1/integrations/dashboard')
    return response.data
  },

  // Products

  // Get products for an account
  getProducts: async (accountId: string, limit = 20, offset = 0): Promise<{ products: ExternalProduct[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/integrations/accounts/${accountId}/products?limit=${limit}&offset=${offset}`)
    return response.data
  },

  // Orders

  // Get orders for an account
  getOrders: async (accountId: string, limit = 20, offset = 0): Promise<{ orders: ExternalOrder[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/integrations/accounts/${accountId}/orders?limit=${limit}&offset=${offset}`)
    return response.data
  },

  // Sync Logs

  // Get sync logs for an account
  getSyncLogs: async (accountId: string, limit = 10): Promise<{ logs: SyncLog[]; total: number }> => {
    const response = await apiGo.get(`/api/v1/integrations/accounts/${accountId}/logs?limit=${limit}`)
    return response.data
  },

  // Affiliate Links

  // Generate affiliate link
  generateLink: async (accountId: string, data: GenerateLinkRequest): Promise<{ affiliate_link: string }> => {
    const response = await apiGo.post(`/api/v1/integrations/accounts/${accountId}/link`, data)
    return response.data
  },

  // OAuth

  // Get OAuth URL for a platform
  getOAuthUrl: async (platform: string, state?: string): Promise<{ oauth_url: string }> => {
    const stateParam = state ? `?state=${state}` : ''
    const response = await apiGo.get(`/api/v1/integrations/oauth/${platform}${stateParam}`)
    return response.data
  },
}