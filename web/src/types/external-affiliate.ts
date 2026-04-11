// External Affiliate Integration Types

export type Platform = 'shopee' | 'tiktok' | 'lazada' | 'amazon' | 'custom'

export interface ExternalAccount {
  id: string
  platform: Platform
  platformAccountId: string
  platformAccountName: string
  platformUsername: string
  status: 'active' | 'pending' | 'disconnected' | 'expired' | 'error'
  autoSync: boolean
  syncIntervalMinutes: number
  lastSyncAt?: string
  lastSyncStatus: string
  lastSyncError?: string
  totalEarnings: number
  totalOrders: number
  totalClicks: number
  totalProducts: number
  pendingCommission: number
  createdAt: string
  updatedAt: string
}

export interface ExternalProduct {
  id: string
  platform: Platform
  productId: string
  productName: string
  productUrl: string
  productImageUrl: string
  category: string
  brand: string
  originalPrice: number
  salePrice: number
  currency: string
  commissionRate: number
  commissionType: 'percentage' | 'fixed'
  estimatedEarnings: number
  totalClicks: number
  totalOrders: number
  totalEarnings: number
  productStatus: string
  affiliateLink: string
}

export interface ExternalOrder {
  id: string
  platform: Platform
  orderId: string
  orderStatus: string
  orderDate?: string
  orderItems: number
  orderAmount: number
  productAmount: number
  shippingAmount: number
  discountAmount: number
  currency: string
  commissionRate: number
  commissionAmount: number
  commissionStatus: string
  trafficSource?: string
}

export interface PlatformStats {
  platform: Platform
  accountCount: number
  totalEarnings: number
  totalOrders: number
  totalClicks: number
  pendingCommission: number
  avgConversionRate: number
  lastSyncAt?: string
}

export interface DailyEarnings {
  date: string
  totalEarnings: number
  ordersCount: number
  clicks: number
  conversionRate: number
}

export interface SyncLog {
  id: string
  syncType: string
  status: string
  startedAt: string
  completedAt?: string
  recordsCreated: number
  recordsUpdated: number
  recordsFailed: number
  errorMessage?: string
}

export interface UnifiedDashboard {
  totalEarnings: number
  totalPendingCommission: number
  totalOrders: number
  totalClicks: number
  totalProducts: number
  platformStats: PlatformStats[]
  recentOrders: ExternalOrder[]
  topProducts: ExternalProduct[]
  earningsTrend: DailyEarnings[]
  connectedAccounts: ExternalAccount[]
}

// Request Types

export interface ConnectPlatformRequest {
  platform: Platform
  apiKey?: string
  apiSecret?: string
  accessToken?: string
  refreshToken?: string
  platformAccountId?: string
  platformUsername?: string
  authorizationCode?: string
}

export interface SyncAccountRequest {
  syncType?: 'full' | 'incremental' | 'products' | 'orders'
}

export interface GenerateLinkRequest {
  productUrl: string
  customName?: string
}