// Shared API client using Next.js API proxy (avoids CORS)
import axios, { AxiosInstance } from 'axios';
import type {
  AuthResponse,
  RefreshTokenResponse,
  User,
  Campaign,
  ABTest,
  Competitor,
  DashboardOverview,
  SEOMetricsResponse,
  MarketingOverview,
  AffiliateDashboard,
  EmailList,
  EmailCampaign,
  SocialPost,
  LandingPage,
  ExternalAccount,
  APIError,
} from '@bookmarketing/shared-types';

// Use Next.js API proxy instead of direct backend URLs
const apiGo: AxiosInstance = axios.create({
  baseURL: '/api/go', // Proxied to Go API
  headers: {
    'Content-Type': 'application/json',
  },
});

const apiPython: AxiosInstance = axios.create({
  baseURL: '/api/python', // Proxied to Python API
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth interceptor
apiGo.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

apiPython.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Add response interceptor for error handling
const handleApiError = (error: any) => {
  if (error.response?.status === 401) {
    // Token expired or invalid - redirect to login
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('refresh_token');
      window.location.href = '/login';
    }
  }
  return Promise.reject(error);
};

apiGo.interceptors.response.use((response) => response, handleApiError);
apiPython.interceptors.response.use((response) => response, handleApiError);

// Auth API
export const authApi = {
  login: async (email: string, password: string): Promise<AuthResponse> => {
    const response = await apiGo.post('/v1/auth/login', { email, password });
    return response.data;
  },

  register: async (email: string, password: string, full_name: string): Promise<AuthResponse> => {
    const response = await apiGo.post('/v1/auth/register', { email, password, full_name });
    return response.data;
  },

  refreshToken: async (refresh_token: string): Promise<RefreshTokenResponse> => {
    const response = await apiGo.post('/v1/auth/refresh', { refresh_token });
    return response.data;
  },

  logout: async (token: string, refresh_token: string): Promise<void> => {
    await apiGo.post('/v1/auth/logout', { token, refresh_token });
  },

  getMe: async (): Promise<User> => {
    const response = await apiGo.get('/v1/auth/me');
    return response.data;
  },

  forgotPassword: async (email: string): Promise<void> => {
    await apiGo.post('/v1/auth/forgot-password', { email });
  },

  resetPassword: async (token: string, new_password: string): Promise<void> => {
    await apiGo.post('/v1/auth/reset-password', { token, new_password });
  },

  verifyEmail: async (token: string): Promise<void> => {
    await apiGo.post('/v1/auth/verify-email', { token });
  },

  resendVerification: async (): Promise<void> => {
    await apiGo.post('/v1/auth/resend-verification');
  },
};

// Campaign API
export const campaignApi = {
  getAll: async (limit = 20, offset = 0): Promise<{ campaigns: Campaign[]; total: number }> => {
    const response = await apiGo.get(`/v1/campaigns?limit=${limit}&offset=${offset}`);
    return response.data;
  },

  getById: async (id: string): Promise<Campaign> => {
    const response = await apiGo.get(`/v1/campaigns/${id}`);
    return response.data;
  },

  create: async (data: Partial<Campaign>): Promise<Campaign> => {
    const response = await apiGo.post('/v1/campaigns', data);
    return response.data;
  },

  update: async (id: string, data: Partial<Campaign>): Promise<Campaign> => {
    const response = await apiGo.put(`/v1/campaigns/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiGo.delete(`/v1/campaigns/${id}`);
  },

  pause: async (id: string): Promise<void> => {
    await apiGo.post(`/v1/campaigns/${id}/pause`);
  },

  activate: async (id: string): Promise<void> => {
    await apiGo.post(`/v1/campaigns/${id}/activate`);
  },
};

// A/B Test API
export const abTestApi = {
  getAll: async (limit = 20, offset = 0, status?: string): Promise<{ tests: ABTest[]; total: number }> => {
    const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
    if (status) params.append('status', status);
    const response = await apiGo.get(`/v1/ab-tests?${params}`);
    return response.data;
  },

  getById: async (id: string): Promise<ABTest> => {
    const response = await apiGo.get(`/v1/ab-tests/${id}`);
    return response.data;
  },

  create: async (data: Partial<ABTest>): Promise<ABTest> => {
    const response = await apiGo.post('/v1/ab-tests', data);
    return response.data;
  },

  update: async (id: string, data: Partial<ABTest>): Promise<ABTest> => {
    const response = await apiGo.put(`/v1/ab-tests/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiGo.delete(`/v1/ab-tests/${id}`);
  },

  start: async (id: string): Promise<void> => {
    await apiGo.post(`/v1/ab-tests/${id}/start`);
  },

  stop: async (id: string): Promise<void> => {
    await apiGo.post(`/v1/ab-tests/${id}/stop`);
  },
};

// Dashboard API
export const dashboardApi = {
  getOverview: async (): Promise<DashboardOverview> => {
    const response = await apiGo.get('/v1/dashboard');
    return response.data;
  },
};

// Competitor API (Python)
export const competitorApi = {
  getAll: async (skip = 0, limit = 20): Promise<{ competitors: Competitor[]; total: number }> => {
    const response = await apiPython.get(`/v1/competitors?skip=${skip}&limit=${limit}`);
    return response.data;
  },

  getById: async (id: string): Promise<Competitor> => {
    const response = await apiPython.get(`/v1/competitors/${id}`);
    return response.data;
  },

  create: async (data: Partial<Competitor>): Promise<Competitor> => {
    const response = await apiPython.post('/v1/competitors', data);
    return response.data;
  },

  update: async (id: string, data: Partial<Competitor>): Promise<Competitor> => {
    const response = await apiPython.put(`/v1/competitors/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiPython.delete(`/v1/competitors/${id}`);
  },

  analyze: async (id: string): Promise<any> => {
    const response = await apiPython.get(`/v1/competitors/${id}/analyze`);
    return response.data;
  },

  marketOverview: async (): Promise<any> => {
    const response = await apiPython.get('/v1/competitors/market/overview');
    return response.data;
  },
};

// SEO API (Python)
export const seoApi = {
  getMetrics: async (domain: string): Promise<SEOMetricsResponse> => {
    const response = await apiPython.get(`/v1/seo/metrics?domain=${domain}`);
    return response.data;
  },

  getKeywords: async (domain: string, limit = 50): Promise<any> => {
    const response = await apiPython.get(`/v1/seo/keywords?domain=${domain}&limit=${limit}`);
    return response.data;
  },

  audit: async (domain: string): Promise<any> => {
    const response = await apiPython.get(`/v1/seo/audit?domain=${domain}`);
    return response.data;
  },

  backlinks: async (domain: string, limit = 100): Promise<any> => {
    const response = await apiPython.get(`/v1/seo/backlinks?domain=${domain}&limit=${limit}`);
    return response.data;
  },
};

// Analytics API (Python)
export const analyticsApi = {
  getSummary: async (days = 30): Promise<any> => {
    const response = await apiPython.get(`/v1/analytics/summary?days=${days}`);
    return response.data;
  },

  getChannels: async (): Promise<any> => {
    const response = await apiPython.get('/v1/analytics/channels');
    return response.data;
  },

  getConversionFunnel: async (): Promise<any> => {
    const response = await apiPython.get('/v1/analytics/conversion/funnel');
    return response.data;
  },

  getTrends: async (metric: string, period: string): Promise<any> => {
    const response = await apiPython.get(`/v1/analytics/trends?metric=${metric}&period=${period}`);
    return response.data;
  },

  predict: async (days = 30): Promise<any> => {
    const response = await apiPython.post(`/v1/analytics/predict?days=${days}`);
    return response.data;
  },
};

// Marketing Features API (Go)
export const marketingApi = {
  // Email
  getEmailLists: async (): Promise<{ lists: EmailList[] }> => {
    const response = await apiGo.get('/v1/email/lists');
    return response.data;
  },

  createEmailList: async (data: Partial<EmailList>): Promise<EmailList> => {
    const response = await apiGo.post('/v1/email/lists', data);
    return response.data;
  },

  getEmailCampaigns: async (): Promise<{ campaigns: EmailCampaign[] }> => {
    const response = await apiGo.get('/v1/email/campaigns');
    return response.data;
  },

  // Social
  getSocialPosts: async (): Promise<{ posts: SocialPost[] }> => {
    const response = await apiGo.get('/v1/social/posts');
    return response.data;
  },

  createSocialPost: async (data: Partial<SocialPost>): Promise<SocialPost> => {
    const response = await apiGo.post('/v1/social/posts', data);
    return response.data;
  },

  // Landing Pages
  getLandingPages: async (): Promise<{ pages: LandingPage[] }> => {
    const response = await apiGo.get('/v1/landing-pages');
    return response.data;
  },

  createLandingPage: async (data: Partial<LandingPage>): Promise<LandingPage> => {
    const response = await apiGo.post('/v1/landing-pages', data);
    return response.data;
  },

  // Reporting
  getMarketingOverview: async (): Promise<MarketingOverview> => {
    const response = await apiGo.get('/v1/reporting/overview');
    return response.data;
  },
};

// Affiliate API (Go)
export const affiliateApi = {
  register: async (data: any): Promise<any> => {
    const response = await apiGo.post('/v1/affiliates/register', data);
    return response.data;
  },

  getMyAffiliate: async (): Promise<any> => {
    const response = await apiGo.get('/v1/affiliates/me');
    return response.data;
  },

  getDashboard: async (): Promise<AffiliateDashboard> => {
    const response = await apiGo.get('/v1/affiliates/dashboard');
    return response.data;
  },

  getCommissions: async (): Promise<any> => {
    const response = await apiGo.get('/v1/affiliates/commissions');
    return response.data;
  },
};

// External Affiliate API (Go)
export const externalAffiliateApi = {
  connectPlatform: async (data: any): Promise<any> => {
    const response = await apiGo.post('/v1/integrations/connect', data);
    return response.data;
  },

  getConnectedAccounts: async (): Promise<{ accounts: ExternalAccount[] }> => {
    const response = await apiGo.get('/v1/integrations/accounts');
    return response.data;
  },

  syncAccount: async (id: string): Promise<any> => {
    const response = await apiGo.post(`/v1/integrations/accounts/${id}/sync`);
    return response.data;
  },

  getUnifiedDashboard: async (): Promise<any> => {
    const response = await apiGo.get('/v1/integrations/dashboard');
    return response.data;
  },
};

// Unified Campaigns API (Go) - Phase 1: supports ALL online + offline channels
export const unifiedCampaignApi = {
  getAll: async (channel?: string, status?: string, limit = 50, offset = 0): Promise<{ campaigns: any[]; total: number }> => {
    const params = new URLSearchParams();
    if (channel) params.set('channel', channel);
    if (status) params.set('status', status);
    params.set('limit', String(limit));
    params.set('offset', String(offset));
    const response = await apiGo.get(`/v1/unified-campaigns?${params}`);
    return response.data;
  },

  getById: async (id: string): Promise<any> => {
    const response = await apiGo.get(`/v1/unified-campaigns/${id}`);
    return response.data;
  },

  create: async (data: any): Promise<any> => {
    const response = await apiGo.post('/v1/unified-campaigns', data);
    return response.data;
  },

  update: async (id: string, data: any): Promise<any> => {
    const response = await apiGo.put(`/v1/unified-campaigns/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiGo.delete(`/v1/unified-campaigns/${id}`);
  },

  submitForApproval: async (id: string): Promise<void> => {
    await apiGo.post(`/v1/unified-campaigns/${id}/submit`);
  },

  approve: async (id: string): Promise<void> => {
    await apiGo.post(`/v1/unified-campaigns/${id}/approve`);
  },

  reject: async (id: string, reason?: string): Promise<void> => {
    await apiGo.post(`/v1/unified-campaigns/${id}/reject`, { reason });
  },
};

// Budget API
export const budgetApi = {
  getPeriods: async () => { const r = await apiGo.get('/v1/budget/periods'); return r.data; },
  getPeriod: async (id: string) => { const r = await apiGo.get(`/v1/budget/periods/${id}`); return r.data; },
  createPeriod: async (data: any) => { const r = await apiGo.post('/v1/budget/periods', data); return r.data; },
  createAllocation: async (periodId: string, data: any) => { const r = await apiGo.post(`/v1/budget/periods/${periodId}/allocations`, data); return r.data; },
  recordTransaction: async (data: any) => { const r = await apiGo.post('/v1/budget/transactions', data); return r.data; },
  getTransactions: async () => { const r = await apiGo.get('/v1/budget/transactions'); return r.data; },
  getAlerts: async () => { const r = await apiGo.get('/v1/budget/alerts'); return r.data; },
};

// Calendar API
export const calendarApi = {
  getEvents: async (params: any) => { const r = await apiGo.get('/v1/calendar/events', { params }); return r.data; },
  createEvent: async (data: any) => { const r = await apiGo.post('/v1/calendar/events', data); return r.data; },
  updateEvent: async (id: string, data: any) => { const r = await apiGo.put(`/v1/calendar/events/${id}`, data); return r.data; },
  deleteEvent: async (id: string) => { await apiGo.delete(`/v1/calendar/events/${id}`); },
  getViews: async () => { const r = await apiGo.get('/v1/calendar/views'); return r.data; },
};

// Approval API
export const approvalApi = {
  getWorkflows: async () => { const r = await apiGo.get('/v1/approval/workflows'); return r.data; },
  getRequests: async (status?: string) => { const r = await apiGo.get('/v1/approval/requests', { params: status ? { status } : {} }); return r.data; },
  getRequest: async (id: string) => { const r = await apiGo.get(`/v1/approval/requests/${id}`); return r.data; },
  createRequest: async (data: any) => { const r = await apiGo.post('/v1/approval/requests', data); return r.data; },
  respond: async (id: string, data: { action: string; comments: string }) => { const r = await apiGo.post(`/v1/approval/requests/${id}/respond`, data); return r.data; },
};

// Email Campaign API (Complete)
export const emailCampaignApi = {
  sendCampaign: async (id: string) => apiGo.post(`/v1/email/campaigns/${id}/send`),
  getAnalytics: async (id: string) => { const r = await apiGo.get(`/v1/email/campaigns/${id}/analytics`); return r.data; },
  createABTest: async (data: any) => { const r = await apiGo.post('/v1/email/ab-tests', data); return r.data; },
  completeABTest: async (id: string) => apiGo.post(`/v1/email/ab-tests/${id}/complete`),
  trackEvent: async (data: any) => apiGo.post('/v1/email/track', data),
};

// PPC Ads API
export const adsApi = {
  connectPlatform: async (data: any) => { const r = await apiGo.post('/v1/ads/platforms/connect', data); return r.data; },
  getPlatforms: async () => { const r = await apiGo.get('/v1/ads/platforms'); return r.data; },
  createCampaign: async (data: any) => { const r = await apiGo.post('/v1/ads/campaigns', data); return r.data; },
  getCampaigns: async () => { const r = await apiGo.get('/v1/ads/campaigns'); return r.data; },
  pauseCampaign: async (id: string) => apiGo.post(`/v1/ads/campaigns/${id}/pause`),
  activateCampaign: async (id: string) => apiGo.post(`/v1/ads/campaigns/${id}/activate`),
};

// SMS API
export const smsApi = {
  getLists: async () => { const r = await apiGo.get('/v1/sms/lists'); return r.data; },
  createList: async (data: any) => { const r = await apiGo.post('/v1/sms/lists', data); return r.data; },
  getSubscribers: async (listId?: string) => { const r = await apiGo.get('/v1/sms/subscribers', { params: listId ? { list_id: listId } : {} }); return r.data; },
  addSubscriber: async (data: any) => { const r = await apiGo.post('/v1/sms/subscribers', data); return r.data; },
  getCampaigns: async () => { const r = await apiGo.get('/v1/sms/campaigns'); return r.data; },
  createCampaign: async (data: any) => { const r = await apiGo.post('/v1/sms/campaigns', data); return r.data; },
  sendCampaign: async (id: string) => { const r = await apiGo.post(`/v1/sms/campaigns/${id}/send`); return r.data; },
};
