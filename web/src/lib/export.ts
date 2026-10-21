import { toCSV } from './csv'
import { MarketingCampaign, CompetitorAnalysis, ABTest, SEOMetrics } from '@/types/marketing'

export function exportToCSV(
  data: Record<string, unknown>[],
  filename: string,
  columns?: string[],
  onProgress?: (progress: number) => void
): void {
  const csvContent = toCSV(data, columns)
  if (onProgress) for (let i = 1; i <= data.length; i++) onProgress(i)

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `${filename}_${new Date().toISOString().split('T')[0]}.csv`
  link.click()
  URL.revokeObjectURL(link.href)
}

export function exportCampaigns(campaigns: MarketingCampaign[]): void {
  const data = campaigns.map(c => ({
    name: c.name,
    type: c.type,
    status: c.status,
    budget: c.budget,
    spent: c.spent,
    impressions: c.impressions,
    clicks: c.clicks,
    conversions: c.conversions,
    ctr: ((c.clicks / c.impressions) * 100).toFixed(2),
    conversion_rate: ((c.conversions / c.clicks) * 100).toFixed(2),
    start_date: c.startDate,
    end_date: c.endDate || 'N/A',
  }))
  exportToCSV(data, 'campaigns')
}

export function exportCompetitors(competitors: CompetitorAnalysis[]): void {
  const data = competitors.map(c => ({
    name: c.competitorName,
    domain: c.domain,
    market_share: c.marketShare,
    seo_score: c.seoScore,
    keywords: c.keywords,
    backlinks: c.backlinks,
    sem_spend: c.semSpend,
    social_followers: c.socialFollowers,
  }))
  exportToCSV(data, 'competitors')
}

export function exportABTests(tests: ABTest[]): void {
  const data = tests.map(t => ({
    name: t.name,
    variant: t.variant,
    status: t.status,
    visitors: t.visitors,
    conversions: t.conversions,
    conversion_rate: t.conversionRate,
    confidence: t.confidence,
    start_date: t.startDate,
    end_date: t.endDate || 'Ongoing',
  }))
  exportToCSV(data, 'ab_tests')
}

export function exportSEOMetrics(metrics: SEOMetrics): void {
  const data = [{
    domain: metrics.domain,
    organic_traffic: metrics.organicTraffic,
    keywords_tracked: metrics.keywords,
    avg_position: metrics.avgPosition,
    domain_authority: metrics.domainAuthority,
    page_speed_score: metrics.pageSpeed,
    backlinks: metrics.backlinks,
  }]
  exportToCSV(data, 'seo_metrics')
}
