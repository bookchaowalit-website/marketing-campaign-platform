'use client'

import { useQuery } from '@tanstack/react-query'
import { Box, Card, CardContent, Typography, Grid, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Skeleton, LinearProgress, Button } from '@mui/material'
import DownloadIcon from '@mui/icons-material/Download'
import { analyticsApi } from '@/lib/api-client'
import { exportToCSV } from '@/lib/export'
import { formatNumber } from '@/lib/utils'

interface AnalyticsSummary {
  total_visitors: number
  total_sessions: number
  bounce_rate: number
  avg_session_duration: string
}

interface ChannelData {
  name: string
  visitors: number
  conversions: number
  rate: number
}

const defaultChannels: ChannelData[] = [
  { name: 'Organic Search', visitors: 85000, conversions: 2550, rate: 3.0 },
  { name: 'Paid Search', visitors: 62000, conversions: 2480, rate: 4.0 },
  { name: 'Social Media', visitors: 45000, conversions: 900, rate: 2.0 },
  { name: 'Email', visitors: 28000, conversions: 1400, rate: 5.0 },
  { name: 'Direct', visitors: 38000, conversions: 760, rate: 2.0 },
  { name: 'Referral', visitors: 26592, conversions: 532, rate: 2.0 },
]

const defaultFunnel = [
  { stage: 'Visitors', count: 284592, percent: 100 },
  { stage: 'Product Views', count: 142296, percent: 50 },
  { stage: 'Add to Cart', count: 42689, percent: 15 },
  { stage: 'Checkout', count: 21344, percent: 7.5 },
  { stage: 'Purchase', count: 8538, percent: 3 },
]

export default function AnalyticsPage() {
  const { data: summary, isLoading: summaryLoading } = useQuery<AnalyticsSummary>({
    queryKey: ['analytics-summary'],
    queryFn: async () => {
      return analyticsApi.getSummary()
    },
    initialData: { total_visitors: 284592, total_sessions: 425890, bounce_rate: 42.5, avg_session_duration: '4m 5s' },
  })

  const { data: channels, isLoading: channelsLoading } = useQuery<ChannelData[]>({
    queryKey: ['analytics-channels'],
    queryFn: async () => {
      return analyticsApi.getChannels()
    },
    initialData: defaultChannels,
  })

  const { data: funnel, isLoading: funnelLoading } = useQuery({
    queryKey: ['analytics-funnel'],
    queryFn: async () => {
      return analyticsApi.getConversionFunnel()
    },
    initialData: defaultFunnel,
  })

  const stats = [
    { label: 'Total Visitors', value: formatNumber(summary?.total_visitors || 0), color: 'primary.main' },
    { label: 'Total Sessions', value: formatNumber(summary?.total_sessions || 0), color: 'secondary.main' },
    { label: 'Bounce Rate', value: `${summary?.bounce_rate || 0}%`, color: 'error.main' },
    { label: 'Avg. Session', value: summary?.avg_session_duration || '0m 0s', color: 'success.main' },
  ]

  const totalVisitors = channels?.reduce((sum, c) => sum + c.visitors, 0) || 284592

  const handleExportChannels = () => {
    const data = (channels || defaultChannels).map(c => ({
      channel: c.name,
      visitors: c.visitors,
      conversions: c.conversions,
      conversion_rate: `${c.rate}%`,
      share: `${((c.visitors / totalVisitors) * 100).toFixed(1)}%`,
    }))
    exportToCSV(data, 'analytics_channels')
  }

  const handleExportFunnel = () => {
    const data = (funnel || defaultFunnel).map((f: { stage: string; count: number; percent: number }) => ({
      stage: f.stage,
      count: f.count,
      percentage: `${f.percent}%`,
    }))
    exportToCSV(data, 'conversion_funnel')
  }

  return (
    <Box>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Analytics Dashboard
        </Typography>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            variant="outlined"
            startIcon={<DownloadIcon />}
            onClick={handleExportChannels}
          >
            Export Channels
          </Button>
          <Button
            variant="outlined"
            startIcon={<DownloadIcon />}
            onClick={handleExportFunnel}
          >
            Export Funnel
          </Button>
        </Box>
      </Box>

      {/* Stats Grid */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {stats.map((stat, i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <Card>
              <CardContent>
                {summaryLoading ? (
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

      {/* Traffic by Channel */}
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Traffic by Channel
          </Typography>
          {channelsLoading ? (
            <Skeleton height={200} />
          ) : (
            <Box sx={{ overflowX: 'auto' }}>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Channel</TableCell>
                      <TableCell>Visitors</TableCell>
                      <TableCell>Conversions</TableCell>
                      <TableCell>Conv. Rate</TableCell>
                      <TableCell>Share</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(channels || defaultChannels).map((channel, i) => (
                      <TableRow key={i}>
                        <TableCell>
                          <Typography fontWeight="medium">{channel.name}</Typography>
                        </TableCell>
                        <TableCell>{formatNumber(channel.visitors)}</TableCell>
                        <TableCell>{formatNumber(channel.conversions)}</TableCell>
                        <TableCell>{channel.rate}%</TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <LinearProgress
                              variant="determinate"
                              value={(channel.visitors / totalVisitors) * 100}
                              sx={{ width: 80, height: 8, borderRadius: 4 }}
                            />
                            <Typography variant="body2">
                              {((channel.visitors / totalVisitors) * 100).toFixed(1)}%
                            </Typography>
                          </Box>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          )}
        </CardContent>
      </Card>

      {/* Conversion Funnel */}
      <Card>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Conversion Funnel
          </Typography>
          {funnelLoading ? (
            <Skeleton height={200} />
          ) : (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              {(funnel || defaultFunnel).map((step: { stage: string; count: number; percent: number }, i: number) => (
                <Box key={i} sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  <Typography sx={{ width: 120, fontWeight: 'medium' }}>{step.stage}</Typography>
                  <Box sx={{ flex: 1 }}>
                    <LinearProgress
                      variant="determinate"
                      value={step.percent}
                      sx={{ height: 24, borderRadius: 2 }}
                    />
                  </Box>
                  <Box sx={{ width: 120, textAlign: 'right' }}>
                    <Typography fontWeight="bold">{formatNumber(step.count)}</Typography>
                    <Typography variant="body2" color="text.secondary">({step.percent}%)</Typography>
                  </Box>
                </Box>
              ))}
            </Box>
          )}
        </CardContent>
      </Card>
    </Box>
  )
}
