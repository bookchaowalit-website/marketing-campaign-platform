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
  LinearProgress,
  IconButton,
} from '@mui/material'
import { Add, TrendingUp, TrendingDown, Refresh } from '@mui/icons-material'
import { competitorApi } from '@/lib/api-client'
import type { Competitor } from '@bookmarketing/shared-types'

export function CompetitorAnalysis() {
  const { data: competitorsResponse, isLoading, refetch } = useQuery({
    queryKey: ['competitors'],
    queryFn: () => competitorApi.getAll(0, 50),
  })

  const competitors = competitorsResponse?.competitors || []

  const getMarketShareColor = (share: number) => {
    if (share > 25) return 'error'
    if (share > 15) return 'warning'
    return 'success'
  }

  const getSEOScoreColor = (score: number) => {
    if (score >= 80) return 'success'
    if (score >= 60) return 'warning'
    return 'error'
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h6" fontWeight="bold">
          Competitor Analysis
        </Typography>
        <Button
          variant="outlined"
          startIcon={<Refresh />}
          onClick={() => refetch()}
        >
          Refresh Data
        </Button>
      </Box>

      {/* Market Share Overview */}
      <Grid container spacing={3} mb={3}>
        {competitors.length > 0 ? competitors.map((competitor: Competitor) => (
          <Grid item xs={12} sm={6} md={4} key={competitor.id}>
            <Card>
              <CardContent>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  {competitor.competitor_name}
                </Typography>
                <Typography variant="body2" color="text.secondary" mb={1}>
                  {competitor.domain}
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h4" fontWeight="bold" mr={1}>
                    {competitor.market_share}%
                  </Typography>
                  <Chip
                    label="Market Share"
                    size="small"
                    color={getMarketShareColor(competitor.market_share)}
                  />
                </Box>
                <LinearProgress
                  variant="determinate"
                  value={competitor.market_share * 2}
                  sx={{ height: 8, borderRadius: 4 }}
                  color={getMarketShareColor(competitor.market_share)}
                />
              </CardContent>
            </Card>
          </Grid>
        )) : (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" align="center">
                  No competitors added yet. Add your first competitor to start tracking!
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>

      {/* Detailed Comparison Table */}
      <Card>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" gutterBottom>
            Detailed Comparison
          </Typography>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Competitor</TableCell>
                  <TableCell>Domain</TableCell>
                  <TableCell>SEO Score</TableCell>
                  <TableCell>Keywords</TableCell>
                  <TableCell>Backlinks</TableCell>
                  <TableCell>SEM Spend</TableCell>
                  <TableCell>Social Followers</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {competitors.length > 0 ? competitors.map((competitor: Competitor) => (
                  <TableRow key={competitor.id}>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {competitor.competitor_name}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" color="primary">
                        {competitor.domain}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={competitor.seo_score}
                        size="small"
                        color={getSEOScoreColor(competitor.seo_score)}
                      />
                    </TableCell>
                    <TableCell>{competitor.keywords?.toLocaleString() || '-'}</TableCell>
                    <TableCell>{competitor.backlinks?.toLocaleString() || '-'}</TableCell>
                    <TableCell>${competitor.sem_spend?.toLocaleString() || '-'}</TableCell>
                    <TableCell>{competitor.social_followers?.toLocaleString() || '-'}</TableCell>
                  </TableRow>
                )) : (
                  <TableRow>
                    <TableCell colSpan={7} align="center">
                      <Typography color="text.secondary">No competitor data available</Typography>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>

      {/* Add Competitor Button */}
      <Box sx={{ mt: 3, display: 'flex', justifyContent: 'center' }}>
        <Button variant="outlined" startIcon={<Add />}>
          Add New Competitor
        </Button>
      </Box>
    </Box>
  )
}
