'use client'

import { useQuery } from '@tanstack/react-query'
import { Box, Card, CardContent, Typography, Grid, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Chip, Skeleton, LinearProgress, Button } from '@mui/material'
import { competitorApi } from '@/lib/api-client'
import { formatCurrency, formatNumber } from '@/lib/utils'

interface Competitor {
  id: string
  competitor_name: string
  domain: string
  type: string
  market_share: number
  seo_score: number
  sem_spend: number
  keywords: number
  backlinks: number
  social_followers: number
}

interface MarketOverview {
  total_market_size: number
  competitors_count: number
  avg_market_share: number
}

const defaultCompetitors: Competitor[] = [
  { id: '1', competitor_name: 'HubSpot', domain: 'hubspot.com', type: 'direct', market_share: 25, seo_score: 92, sem_spend: 50000, keywords: 15000, backlinks: 250000, social_followers: 500000 },
  { id: '2', competitor_name: 'Marketo', domain: 'marketo.com', type: 'direct', market_share: 18, seo_score: 85, sem_spend: 35000, keywords: 8000, backlinks: 180000, social_followers: 300000 },
  { id: '3', competitor_name: 'Mailchimp', domain: 'mailchimp.com', type: 'indirect', market_share: 15, seo_score: 88, sem_spend: 25000, keywords: 12000, backlinks: 200000, social_followers: 450000 },
  { id: '4', competitor_name: 'ActiveCampaign', domain: 'activecampaign.com', type: 'direct', market_share: 12, seo_score: 78, sem_spend: 20000, keywords: 5000, backlinks: 120000, social_followers: 150000 },
  { id: '5', competitor_name: 'ConvertKit', domain: 'convertkit.com', type: 'emerging', market_share: 8, seo_score: 72, sem_spend: 10000, keywords: 3000, backlinks: 80000, social_followers: 100000 },
]

export default function CompetitorsPage() {
  const { data: competitorsResponse, isLoading: competitorsLoading } = useQuery<{ competitors: Competitor[] }>({
    queryKey: ['competitors'],
    queryFn: async () => {
      return competitorApi.getAll()
    },
    initialData: { competitors: defaultCompetitors },
  })

  const { data: marketOverview, isLoading: overviewLoading } = useQuery<MarketOverview>({
    queryKey: ['market-overview'],
    queryFn: async () => {
      return competitorApi.marketOverview()
    },
    initialData: { total_market_size: 500000000, competitors_count: 5, avg_market_share: 16 },
  })

  const competitors = competitorsResponse?.competitors || defaultCompetitors

  const getTypeColor = (type: string): 'error' | 'warning' | 'success' => {
    switch (type) {
      case 'direct': return 'error'
      case 'indirect': return 'warning'
      case 'emerging': return 'success'
      default: return 'warning'
    }
  }

  return (
    <Box>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Competitor Analysis
        </Typography>
        <Button variant="contained" disabled title="Coming soon">Add Competitor</Button>
      </Box>

      {/* Market Overview */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              {overviewLoading ? (
                <Skeleton height={60} />
              ) : (
                <>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Total Market Size
                  </Typography>
                  <Typography variant="h4" fontWeight="bold" color="primary">
                    {formatCurrency(marketOverview?.total_market_size || 0)}
                  </Typography>
                </>
              )}
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              {overviewLoading ? (
                <Skeleton height={60} />
              ) : (
                <>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Competitors Tracked
                  </Typography>
                  <Typography variant="h4" fontWeight="bold" color="primary">
                    {marketOverview?.competitors_count || competitors?.length || 5}
                  </Typography>
                </>
              )}
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              {overviewLoading ? (
                <Skeleton height={60} />
              ) : (
                <>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Avg. Market Share
                  </Typography>
                  <Typography variant="h4" fontWeight="bold" color="primary">
                    {marketOverview?.avg_market_share?.toFixed(1) || '16'}%
                  </Typography>
                </>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Market Share Visualization */}
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Market Share Distribution
          </Typography>
          {competitorsLoading ? (
            <Skeleton height={100} />
          ) : (
            <Grid container spacing={2}>
              {competitors.slice(0, 4).map((comp, i) => (
                <Grid item xs={6} sm={3} key={i}>
                  <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'grey.50' }}>
                    <Typography variant="h5" fontWeight="bold" color="primary">
                      {comp.market_share}%
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {comp.competitor_name}
                    </Typography>
                    <LinearProgress
                      variant="determinate"
                      value={comp.market_share}
                      sx={{ mt: 1, height: 8, borderRadius: 4 }}
                    />
                  </Paper>
                </Grid>
              ))}
            </Grid>
          )}
        </CardContent>
      </Card>

      {/* Competitors Table */}
      <Card>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Competitor Details
          </Typography>
          {competitorsLoading ? (
            <Skeleton height={300} />
          ) : (
            <Box sx={{ overflowX: 'auto' }}>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Competitor</TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Market Share</TableCell>
                      <TableCell>SEO Score</TableCell>
                      <TableCell>Keywords</TableCell>
                      <TableCell>Backlinks</TableCell>
                      <TableCell>SEM Spend</TableCell>
                      <TableCell>Social</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {competitors.map((comp) => (
                      <TableRow key={comp.id}>
                        <TableCell>
                          <Typography fontWeight="medium">{comp.competitor_name}</Typography>
                          <Typography variant="caption" color="text.secondary">{comp.domain}</Typography>
                        </TableCell>
                        <TableCell>
                          <Chip label={comp.type} size="small" color={getTypeColor(comp.type)} />
                        </TableCell>
                        <TableCell>{comp.market_share}%</TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <LinearProgress
                              variant="determinate"
                              value={comp.seo_score}
                              sx={{ width: 60, height: 6, borderRadius: 3 }}
                              color={comp.seo_score >= 80 ? 'success' : 'warning'}
                            />
                            <Typography variant="body2">{comp.seo_score}</Typography>
                          </Box>
                        </TableCell>
                        <TableCell>{formatNumber(comp.keywords)}</TableCell>
                        <TableCell>{formatNumber(comp.backlinks)}</TableCell>
                        <TableCell>{formatCurrency(comp.sem_spend)}</TableCell>
                        <TableCell>{formatNumber(comp.social_followers)}</TableCell>
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
