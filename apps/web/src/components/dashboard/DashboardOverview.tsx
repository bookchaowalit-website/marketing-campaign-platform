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
import { MarketingDashboard } from '@/types/marketing'

const mockDashboardData: MarketingDashboard = {
  totalCampaigns: 24,
  activeCampaigns: 12,
  totalSpend: 45280,
  totalRevenue: 128450,
  roi: 183.6,
  conversionRate: 3.24,
  topChannels: [
    { name: 'Google Ads', value: 45 },
    { name: 'Facebook', value: 30 },
    { name: 'Organic Search', value: 15 },
    { name: 'Email', value: 10 },
  ],
  recentCampaigns: [],
}

export function DashboardOverview() {
  const { data: dashboard, isLoading } = useQuery<MarketingDashboard>({
    queryKey: ['dashboard'],
    queryFn: async () => {
      // TODO: Replace with actual API call
      // const response = await apiGo.get('/api/v1/dashboard')
      // return response.data
      return mockDashboardData
    },
  })

  if (isLoading) {
    return <LinearProgress />
  }

  if (!dashboard) {
    return null
  }

  const stats = [
    {
      title: 'Total Revenue',
      value: `$${dashboard.totalRevenue.toLocaleString()}`,
      change: '+12.5%',
      isPositive: true,
      icon: <AttachMoney />,
      color: 'success.main',
    },
    {
      title: 'Total Spend',
      value: `$${dashboard.totalSpend.toLocaleString()}`,
      change: '-3.2%',
      isPositive: true,
      icon: <AttachMoney />,
      color: 'warning.main',
    },
    {
      title: 'ROI',
      value: `${dashboard.roi}%`,
      change: '+8.1%',
      isPositive: true,
      icon: <TrendingUp />,
      color: 'primary.main',
    },
    {
      title: 'Conversion Rate',
      value: `${dashboard.conversionRate}%`,
      change: '-0.5%',
      isPositive: false,
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
                <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
                  {stat.isPositive ? (
                    <TrendingUp sx={{ fontSize: 16, color: 'success.main', mr: 0.5 }} />
                  ) : (
                    <TrendingDown sx={{ fontSize: 16, color: 'error.main', mr: 0.5 }} />
                  )}
                  <Typography
                    variant="body2"
                    color={stat.isPositive ? 'success.main' : 'error.main'}
                    fontWeight="medium"
                  >
                    {stat.change}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" ml={0.5}>
                    vs last month
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Campaigns Overview */}
      <Grid container spacing={3}>
        <Grid size={{ xs: 12, md: 8 }}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" gutterBottom>
                Campaign Performance
              </Typography>
              <Box sx={{ mt: 3 }}>
                {dashboard.topChannels.map((channel, index) => (
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
                ))}
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid size={{ xs: 12, md: 4 }}>
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
                      Total Visitors
                    </Typography>
                    <Typography variant="h6" fontWeight="bold">
                      284,592
                    </Typography>
                  </Box>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  <Box sx={{ mr: 2, color: 'success.main' }}>
                    <ShoppingCart />
                  </Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Conversions
                    </Typography>
                    <Typography variant="h6" fontWeight="bold">
                      9,221
                    </Typography>
                  </Box>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                  <Box sx={{ mr: 2, color: 'secondary.main' }}>
                    <CampaignIcon />
                  </Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Active Campaigns
                    </Typography>
                    <Typography variant="h6" fontWeight="bold">
                      {dashboard.activeCampaigns} / {dashboard.totalCampaigns}
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

function CampaignIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor" {...props}>
      <path d="M21 3H3C1.89 3 1 3.89 1 5V19C1 20.11 1.89 21 3 21H21C22.11 21 23 20.11 23 19V5C23 3.89 22.11 3 21 3ZM21 19H3V5H21V19Z" />
      <path d="M17 10H15V14H17V10Z" />
      <path d="M13 10H11V14H13V10Z" />
      <path d="M9 10H7V14H9V10Z" />
    </svg>
  )
}
