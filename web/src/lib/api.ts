// This file is deprecated. Use @/lib/api-client instead.
// The new client uses Next.js API proxy to avoid CORS issues.
//
// Migration guide:
// - Replace `import { apiGo } from '@/lib/api'` with `import { campaignApi, authApi } from '@/lib/api-client'`
// - Replace `import { apiPython } from '@/lib/api'` with `import { competitorApi, seoApi, analyticsApi } from '@/lib/api-client'`
//
// Example:
// Old: const response = await apiGo.get('/api/v1/campaigns')
// New: const response = await campaignApi.getAll()

// Re-export API client for convenience
export * from './api-client'
