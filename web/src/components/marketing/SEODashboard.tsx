'use client'

import { useQuery } from '@tanstack/react-query'
import {
  Card,
  CardContent,
  Typography,
  Box,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Button,
  TextField,
  InputAdornment,
  LinearProgress,
} from '@mui/material'
import { Search, TrendingUp, Link, Speed } from '@mui/icons-material'
import { seoApi } from '@/lib/api-client'

interface SEOMetricsData {
  domain: string
  organic_traffic?: number
  organicTraffic?: number
  keywords?: number
  avg_position?: number
  backlinks?: number
  domain_authority?: number
  domainAuthority?: number
  page_speed?: number
  pageSpeed?: number
  impressions?: number
  ctr?: number
  referring_domains?: number
}

interface KeywordData {
  keyword: string
  position: number
  volume?: number
  search_volume?: number
  difficulty?: number
  cpc?: number
  clicks?: number
  impressions?: number
  ctr?: number
  trend?: string
}

const DOMAIN = 'yoursite.com' // Can be made configurable

export function SEODashboard() {
  const { data: seoMetrics, isLoading: metricsLoading } = useQuery<SEOMetricsData>({
    queryKey: ['seoMetrics', DOMAIN],
    // The API returns { metrics: {...}, recommendations: [] }, not a flat object.
    queryFn: async () => (await seoApi.getMetrics(DOMAIN)).metrics,
  })

  const { data: keywords, isLoading: keywordsLoading } = useQuery<KeywordData[]>({
    queryKey: ['seoKeywords', DOMAIN],
    queryFn: () => seoApi.getKeywords(DOMAIN, 50),
  })

  // Normalize metrics data (handle both snake_case and camelCase from different sources)
  const metrics = seoMetrics ? {
    domain: seoMetrics.domain || DOMAIN,
    organicTraffic: seoMetrics.organic_traffic || seoMetrics.organicTraffic || 0,
    keywords: seoMetrics.keywords || 0,
    avgPosition: seoMetrics.avg_position || 0,
    backlinks: seoMetrics.backlinks || 0,
    domainAuthority: seoMetrics.domain_authority || seoMetrics.domainAuthority || 0,
    pageSpeed: seoMetrics.page_speed || seoMetrics.pageSpeed || 0,
  } : null

  const keywordsList = Array.isArray(keywords)
    ? keywords.map(k => ({
      keyword: k.keyword,
      position: k.position || 0,
      volume: k.search_volume || k.volume || 0,
      difficulty: k.difficulty || 0,
      trend: k.trend || 'stable',
    }))
    : []

  const getDifficultyColor = (difficulty: number) => {
    if (difficulty >= 70) return 'error'
    if (difficulty >= 50) return 'warning'
    return 'success'
  }

  const getPositionColor = (position: number) => {
    if (position <= 3) return 'success'
    if (position <= 10) return 'warning'
    return 'default'
  }

  if (metricsLoading || keywordsLoading) {
    return <LinearProgress />
  }

  if (!metrics) {
    return (
      <Card>
        <CardContent>
          <Typography>No SEO data available. Connect a data source in settings.</Typography>
        </CardContent>
      </Card>
    )
  }

  return (
    <Box>
      {/* Search Bar */}
      <Box mb={3}>
        <TextField
          fullWidth
          placeholder="Search keywords, domains, or URLs..."
          variant="outlined"
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search />
              </InputAdornment>
            ),
          }}
        />
      </Box>

      {/* Key Metrics */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Organic Traffic
              </Typography>
              <Typography variant="h4" fontWeight="bold">
                {metrics.organicTraffic.toLocaleString()}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Keywords Tracked
              </Typography>
              <Typography variant="h4" fontWeight="bold">
                {metrics.keywords.toLocaleString()}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Avg. Position
              </Typography>
              <Typography variant="h4" fontWeight="bold">
                {metrics.avgPosition.toFixed(1)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Domain Authority
              </Typography>
              <Typography variant="h4" fontWeight="bold">
                {metrics.domainAuthority}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Grid container spacing={3}>
        {/* Keywords Table */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" gutterBottom>
                Top Keywords
              </Typography>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Keyword</TableCell>
                      <TableCell>Position</TableCell>
                      <TableCell>Volume</TableCell>
                      <TableCell>Difficulty</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {keywordsList.slice(0, 10).map((item, index) => (
                      <TableRow key={index}>
                        <TableCell>
                          <Typography variant="body2" fontWeight="medium">
                            {item.keyword}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={item.position}
                            size="small"
                            color={getPositionColor(item.position)}
                          />
                        </TableCell>
                        <TableCell>{item.volume.toLocaleString()}</TableCell>
                        <TableCell>
                          <Chip
                            label={item.difficulty || '-'}
                            size="small"
                            color={item.difficulty ? getDifficultyColor(item.difficulty) : 'default'}
                            variant="outlined"
                          />
                        </TableCell>
                      </TableRow>
                    ))}
                    {keywordsList.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={4} align="center">
                          <Typography color="text.secondary">No keyword data available</Typography>
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Technical SEO */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" gutterBottom>
                Technical SEO
              </Typography>

              <Box sx={{ mt: 3 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                  <Speed sx={{ mr: 2, color: 'primary.main' }} />
                  <Box sx={{ flexGrow: 1 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                      <Typography variant="body2">Page Speed</Typography>
                      <Typography variant="body2" fontWeight="bold">
                        {metrics.pageSpeed}/100
                      </Typography>
                    </Box>
                    <LinearProgress
                      variant="determinate"
                      value={metrics.pageSpeed}
                      sx={{ height: 8, borderRadius: 4 }}
                      color={
                        metrics.pageSpeed >= 90
                          ? 'success'
                          : metrics.pageSpeed >= 70
                            ? 'warning'
                            : 'error'
                      }
                    />
                  </Box>
                </Box>

                <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                  <Link sx={{ mr: 2, color: 'secondary.main' }} />
                  <Box sx={{ flexGrow: 1 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                      <Typography variant="body2">Backlinks</Typography>
                      <Typography variant="body2" fontWeight="bold">
                        {metrics.backlinks.toLocaleString()}
                      </Typography>
                    </Box>
                    <LinearProgress
                      variant="determinate"
                      value={Math.min(metrics.backlinks / 10000 * 100, 100)}
                      sx={{ height: 8, borderRadius: 4 }}
                      color="secondary"
                    />
                  </Box>
                </Box>

                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                  <TrendingUp sx={{ mr: 2, color: 'success.main' }} />
                  <Box sx={{ flexGrow: 1 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                      <Typography variant="body2">Domain Authority</Typography>
                      <Typography variant="body2" fontWeight="bold">
                        {metrics.domainAuthority}/100
                      </Typography>
                    </Box>
                    <LinearProgress
                      variant="determinate"
                      value={metrics.domainAuthority}
                      sx={{ height: 8, borderRadius: 4 }}
                      color="success"
                    />
                  </Box>
                </Box>
              </Box>

              <Button fullWidth variant="outlined" sx={{ mt: 3 }}>
                Run SEO Audit
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}
