'use client'

import { useQuery } from '@tanstack/react-query'
import { Box, Card, CardContent, Typography, Grid, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Chip, Skeleton, LinearProgress } from '@mui/material'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import TrendingDownIcon from '@mui/icons-material/TrendingDown'
import TrendingFlatIcon from '@mui/icons-material/TrendingFlat'
import { seoApi } from '@/lib/api-client'
import { formatNumber } from '@/lib/utils'

interface SEOMetrics {
  organic_traffic: number
  keywords_tracked: number
  avg_position: number
  domain_authority: number
  page_speed_score: number
  backlinks: number
  referring_domains: number
}

interface KeywordRanking {
  keyword: string
  position: number
  volume: number
  difficulty: number
  trend: string
}

const defaultMetrics: SEOMetrics = {
  organic_traffic: 125000,
  keywords_tracked: 450,
  avg_position: 12.5,
  domain_authority: 65,
  page_speed_score: 78,
  backlinks: 12500,
  referring_domains: 850,
}

const defaultKeywords: KeywordRanking[] = [
  { keyword: 'marketing automation', position: 3, volume: 12000, difficulty: 65, trend: 'up' },
  { keyword: 'saas marketing tools', position: 5, volume: 8500, difficulty: 55, trend: 'up' },
  { keyword: 'email marketing software', position: 12, volume: 15000, difficulty: 70, trend: 'down' },
  { keyword: 'campaign management', position: 8, volume: 6000, difficulty: 45, trend: 'stable' },
  { keyword: 'marketing analytics', position: 2, volume: 9500, difficulty: 50, trend: 'up' },
]

export default function SEOPage() {
  const { data: metrics, isLoading: metricsLoading } = useQuery<SEOMetrics>({
    queryKey: ['seo-metrics'],
    queryFn: async () => {
      // The API returns { metrics: {...}, recommendations: [] } with
      // snake_case field names that don't all match this page's local
      // shape (real API: keywords, page_speed; no referring_domains).
      const { metrics: m } = await seoApi.getMetrics('example.com')
      return {
        organic_traffic: m.organic_traffic,
        keywords_tracked: m.keywords,
        avg_position: m.avg_position,
        domain_authority: m.domain_authority,
        page_speed_score: m.page_speed,
        backlinks: m.backlinks,
        referring_domains: 0,
      }
    },
    initialData: defaultMetrics,
  })

  const { data: keywords, isLoading: keywordsLoading } = useQuery<KeywordRanking[]>({
    queryKey: ['seo-keywords'],
    queryFn: async () => {
      // The API returns `search_volume`, not `volume`.
      const raw: Array<{ keyword: string; position: number; search_volume: number; difficulty: number; trend: string }> =
        await seoApi.getKeywords('example.com')
      return raw.map((kw) => ({
        keyword: kw.keyword,
        position: kw.position,
        volume: kw.search_volume,
        difficulty: kw.difficulty,
        trend: kw.trend,
      }))
    },
    initialData: defaultKeywords,
  })

  const stats = [
    { label: 'Organic Traffic', value: formatNumber(metrics?.organic_traffic || 0), color: 'success.main' },
    { label: 'Keywords Tracked', value: metrics?.keywords_tracked?.toString() || '0', color: 'primary.main' },
    { label: 'Avg. Position', value: metrics?.avg_position?.toFixed(1) || '0', color: 'secondary.main' },
    { label: 'Domain Authority', value: metrics?.domain_authority?.toString() || '0', color: 'warning.main' },
  ]

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return <TrendingUpIcon fontSize="small" color="success" />
      case 'down': return <TrendingDownIcon fontSize="small" color="error" />
      default: return <TrendingFlatIcon fontSize="small" color="action" />
    }
  }

  return (
    <Box>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          SEO & SEM Analytics
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Search engine optimization metrics
        </Typography>
      </Box>

      {/* Stats Grid */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {stats.map((stat, i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <Card>
              <CardContent>
                {metricsLoading ? (
                  <Skeleton height={60} />
                ) : (
                  <>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      {stat.label}
                    </Typography>
                    <Typography variant="h4" fontWeight="bold" sx={{ color: stat.color }}>
                      {stat.value}
                    </Typography>
                  </>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Additional Metrics */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" gutterBottom>
                Page Speed Score
              </Typography>
              {metricsLoading ? (
                <Skeleton height={100} />
              ) : (
                <Box sx={{ textAlign: 'center', py: 2 }}>
                  <Typography variant="h3" fontWeight="bold" color={(metrics?.page_speed_score ?? 0) >= 80 ? 'success.main' : 'warning.main'}>
                    {metrics?.page_speed_score ?? 0}
                  </Typography>
                  <LinearProgress
                    variant="determinate"
                    value={metrics?.page_speed_score ?? 0}
                    sx={{ mt: 2, height: 10, borderRadius: 5 }}
                    color={(metrics?.page_speed_score ?? 0) >= 80 ? 'success' : 'warning'}
                  />
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" gutterBottom>
                Backlink Profile
              </Typography>
              {metricsLoading ? (
                <Skeleton height={100} />
              ) : (
                <Grid container spacing={2}>
                  <Grid item xs={6}>
                    <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'grey.50' }}>
                      <Typography variant="h5" fontWeight="bold" color="primary">
                        {formatNumber(metrics?.backlinks || 0)}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Total Backlinks
                      </Typography>
                    </Paper>
                  </Grid>
                  <Grid item xs={6}>
                    <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'grey.50' }}>
                      <Typography variant="h5" fontWeight="bold" color="primary">
                        {formatNumber(metrics?.referring_domains || 0)}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Referring Domains
                      </Typography>
                    </Paper>
                  </Grid>
                </Grid>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Keyword Rankings */}
      <Card>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Keyword Rankings
          </Typography>
          {keywordsLoading ? (
            <Skeleton height={200} />
          ) : (
            <Box sx={{ overflowX: 'auto' }}>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Keyword</TableCell>
                      <TableCell>Position</TableCell>
                      <TableCell>Volume</TableCell>
                      <TableCell>Difficulty</TableCell>
                      <TableCell>Trend</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(keywords || defaultKeywords).map((kw, i) => (
                      <TableRow key={i}>
                        <TableCell>
                          <Typography fontWeight="medium">{kw.keyword}</Typography>
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={`#${kw.position}`}
                            size="small"
                            color={kw.position <= 10 ? 'success' : kw.position <= 20 ? 'warning' : 'default'}
                          />
                        </TableCell>
                        <TableCell>{formatNumber(kw.volume)}</TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <LinearProgress
                              variant="determinate"
                              value={kw.difficulty}
                              sx={{ width: 60, height: 6, borderRadius: 3 }}
                              color={kw.difficulty >= 70 ? 'error' : kw.difficulty >= 40 ? 'warning' : 'success'}
                            />
                            <Typography variant="body2">{kw.difficulty}</Typography>
                          </Box>
                        </TableCell>
                        <TableCell>{getTrendIcon(kw.trend)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          )}
        </CardContent>
      </Card>
    </Box>
  )
}
