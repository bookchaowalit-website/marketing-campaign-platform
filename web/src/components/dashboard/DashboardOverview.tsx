'use client'

import { useQuery } from '@tanstack/react-query'
import { Grid, Card, CardContent, Typography, Box, LinearProgress } from '@mui/material'
import {
  TrendingUp,
  TrendingDown,
  AttachMoney,
  Percent,
  People,
  ShoppingCart,
} from '@mui/icons-material'
import { dashboardApi } from '@/lib/api-client'

export function DashboardOverview() {
  const { data: dashboard, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: () => dashboardApi.getOverview(),
  })

  if (isLoading) {
    return <LinearProgress />
  }

  if (!dashboard) {
    return (
      <Card>
        <CardContent>
          <Typography>No dashboard data available. Create campaigns to see metrics.</Typography>
        </CardContent>
      </Card>
    )
  }

  const stats = [
    {
      title: 'Total Revenue',
      value: `$${dashboard.total_revenue?.toLocaleString() || '0'}`,
      icon: <AttachMoney />,
      color: 'success.main',
    },
    {
      title: 'Total Spend',
      value: `$${dashboard.total_spend?.toLocaleString() || '0'}`,
      icon: <AttachMoney />,
      color: 'warning.main',
    },
    {
      title: 'ROI',
      value: `${dashboard.roi?.toFixed(1) || '0'}%`,
      icon: <TrendingUp />,
      color: 'primary.main',
    },
    {
      title: 'Conversion Rate',
      value: `${dashboard.conversion_rate?.toFixed(2) || '0'}%`,
      icon: <Percent />,
      color: 'secondary.main',
    },
  ]

  return (
    <Box>
      {/* Stats Cards */}
      <Grid container spacing={3} mb={3}>
        {stats.map((stat, index) => (
          <Grid item xs={12} sm={6} md={3} key={index}>
            <Card sx={{ height: '100%' }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <Box>
                    <Typography color="text.secondary" variant="body2" gutterBottom>
                      {stat.title}
                    </Typography>
                    <Typography variant="h4" fontWeight="bold">
                      {stat.value}
                    </Typography>
                  </Box>
                  <Box
                    sx={{
                      p: 1,
                      borderRadius: 2,
                      bgcolor: `${stat.color}15`,
                      color: stat.color,
                    }}
                  >
                    {stat.icon}
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Campaigns Overview */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" gutterBottom>
                Campaign Performance
              </Typography>
              <Box sx={{ mt: 3 }}>
                {dashboard.top_channels && dashboard.top_channels.length > 0 ? (
                  dashboard.top_channels.map((channel: any, index: number) => (
                    <Box key={index} sx={{ mb: 2 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                        <Typography variant="body2">{channel.name}</Typography>
                        <Typography variant="body2" fontWeight="medium">
                          {channel.value}%
                        </Typography>
                      </Box>
                      <LinearProgress
                        variant="determinate"
                        value={channel.value}
                        sx={{
                          height: 8,
                          borderRadius: 4,
                          bgcolor: 'grey.200',
                        }}
                      />
                    </Box>
                  ))
                ) : (
                  <Typography color="text.secondary">No channel data available</Typography>
                )}
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" gutterBottom>
                Quick Stats
              </Typography>
              <Box sx={{ mt: 2 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  <People sx={{ mr: 2, color: 'primary.main' }} />
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Total Campaigns
                    </Typography>
                    <Typography variant="h6" fontWeight="bold">
                      {dashboard.total_campaigns || 0}
                    </Typography>
                  </Box>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  <Box sx={{ mr: 2, color: 'success.main' }}><ShoppingCart /></Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Active Campaigns
                    </Typography>
                    <Typography variant="h6" fontWeight="bold">
                      {dashboard.active_campaigns || 0}
                    </Typography>
                  </Box>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}
