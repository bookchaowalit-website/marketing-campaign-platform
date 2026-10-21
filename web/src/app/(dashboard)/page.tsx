'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { alpha, Box, Card, CardContent, Typography, Grid, Paper, LinearProgress, Chip, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, IconButton, Skeleton, Button } from '@mui/material'
import PlayArrowIcon from '@mui/icons-material/PlayArrow'
import PauseIcon from '@mui/icons-material/Pause'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import AddIcon from '@mui/icons-material/Add'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import TrendingDownIcon from '@mui/icons-material/TrendingDown'
import { dashboardApi, campaignApi } from '@/lib/api-client'
import { CustomBarChart, CustomPieChart, CustomAreaChart } from '@/components/charts'
import { ConfirmDialog, useConfirmDialog } from '@/components/ConfirmDialog'
import { useRouter } from 'next/navigation'
import { formatCurrency, formatNumber } from '@/lib/utils'
import { paletteColor } from '@/lib/palette'

interface DashboardData {
  total_campaigns: number
  active_campaigns: number
  total_spend: number
  total_revenue: number
  roi: number
  conversion_rate: number
  top_channels: { name: string; value: number }[]
}

interface Campaign {
  id: string
  name: string
  type: string
  status: string
  budget: number
  spent: number
  impressions: number
  clicks: number
  conversions: number
  ctr: number
  conversion_rate: number
  budget_utilization: number
}

const defaultDashboard: DashboardData = {
  total_campaigns: 12,
  active_campaigns: 5,
  total_spend: 45280,
  total_revenue: 128450,
  roi: 183.6,
  conversion_rate: 3.24,
  top_channels: [
    { name: 'Google Ads', value: 45 },
    { name: 'Facebook', value: 30 },
    { name: 'Organic', value: 15 },
    { name: 'Email', value: 10 },
  ],
}

const defaultCampaigns: Campaign[] = [
  { id: '10000000-0000-0000-0000-000000000001', name: 'Summer Sale 2024', type: 'SEM', status: 'active', budget: 10000, spent: 6500, impressions: 125000, clicks: 4500, conversions: 180, ctr: 3.6, conversion_rate: 4.0, budget_utilization: 65 },
  { id: '10000000-0000-0000-0000-000000000002', name: 'SEO Organic Growth', type: 'SEO', status: 'active', budget: 5000, spent: 3200, impressions: 85000, clicks: 3200, conversions: 95, ctr: 3.8, conversion_rate: 3.0, budget_utilization: 64 },
  { id: '10000000-0000-0000-0000-000000000003', name: 'Facebook Brand Awareness', type: 'Social', status: 'paused', budget: 8000, spent: 4100, impressions: 250000, clicks: 8900, conversions: 220, ctr: 3.6, conversion_rate: 2.5, budget_utilization: 51 },
  { id: '10000000-0000-0000-0000-000000000004', name: 'Email Retargeting Q2', type: 'Email', status: 'completed', budget: 3000, spent: 2850, impressions: 45000, clicks: 2100, conversions: 310, ctr: 4.7, conversion_rate: 14.8, budget_utilization: 95 },
]

// Sample data for charts
const revenueData = [
  { name: 'Jan', revenue: 45000, spend: 12000 },
  { name: 'Feb', revenue: 52000, spend: 15000 },
  { name: 'Mar', revenue: 48000, spend: 14000 },
  { name: 'Apr', revenue: 61000, spend: 18000 },
  { name: 'May', revenue: 55000, spend: 16000 },
  { name: 'Jun', revenue: 72000, spend: 22000 },
]

const campaignPerformanceData = [
  { name: 'SEM', conversions: 450, spend: 25000 },
  { name: 'SEO', conversions: 280, spend: 12000 },
  { name: 'Social', conversions: 320, spend: 18000 },
  { name: 'Email', conversions: 520, spend: 8000 },
  { name: 'Display', conversions: 180, spend: 15000 },
]

export default function DashboardPage() {
  const { user } = useAuth()
  const toast = useToast()
  const router = useRouter()
  const queryClient = useQueryClient()
  const confirmDialog = useConfirmDialog()

  const { data: dashboard, isLoading: dashboardLoading } = useQuery<DashboardData>({
    queryKey: ['dashboard'],
    queryFn: () => dashboardApi.getOverview(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })

  const { data: campaignsResponse, isLoading: campaignsLoading } = useQuery({
    queryKey: ['campaigns-list'],
    queryFn: () => campaignApi.getAll(10, 0),
    staleTime: 5 * 60 * 1000,
  })

  const campaigns = campaignsResponse?.campaigns || []

  // Mutations for campaign actions
  const pauseMutation = useMutation({
    mutationFn: (id: string) => campaignApi.pause(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns-list'] })
      toast.success('Campaign paused')
    },
    onError: () => toast.error('Failed to pause campaign'),
  })

  const activateMutation = useMutation({
    mutationFn: (id: string) => campaignApi.activate(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns-list'] })
      toast.success('Campaign activated')
    },
    onError: () => toast.error('Failed to activate campaign'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => campaignApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns-list'] })
      toast.success('Campaign deleted')
    },
    onError: () => toast.error('Failed to delete campaign'),
  })

  const stats = [
    { label: 'Total Revenue', value: formatCurrency(dashboard?.total_revenue || 0), color: 'success.main', icon: TrendingUpIcon },
    { label: 'Total Spend', value: formatCurrency(dashboard?.total_spend || 0), color: 'error.main', icon: TrendingDownIcon },
    { label: 'ROI', value: `${dashboard?.roi?.toFixed(1) || '0'}%`, color: 'info.main', icon: TrendingUpIcon },
    { label: 'Conversion Rate', value: `${dashboard?.conversion_rate?.toFixed(2) || '0'}%`, color: 'secondary.main', icon: TrendingUpIcon },
  ]

  const getStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' => {
    switch (status) {
      case 'active': return 'success'
      case 'paused': return 'warning'
      case 'completed': return 'default'
      case 'draft': return 'info'
      default: return 'default'
    }
  }

  const getTypeColor = (type: string): 'primary' | 'success' | 'secondary' | 'warning' | 'info' => {
    switch (type) {
      case 'SEM': return 'primary'
      case 'SEO': return 'success'
      case 'Social': return 'secondary'
      case 'Email': return 'warning'
      case 'Display': return 'info'
      default: return 'primary'
    }
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Box>
          <Typography variant="h4" fontWeight="bold">
            Dashboard Overview
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Welcome back, {user?.full_name || user?.email?.split('@')[0] || 'User'}
          </Typography>
        </Box>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => router.push('/campaigns')} aria-label="Create new campaign">
          New Campaign
        </Button>
      </Box>

      {/* Stats Grid */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {stats.map((stat, i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <Card sx={{ height: '100%' }}>
              <CardContent>
                {dashboardLoading ? (
                  <Skeleton height={80} />
                ) : (
                  <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        {stat.label}
                      </Typography>
                      <Typography variant="h4" fontWeight="bold" sx={{ color: stat.color }}>
                        {stat.value}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                        {dashboard?.active_campaigns || 0} active campaigns
                      </Typography>
                    </Box>
                    <Box sx={{
                      p: 1,
                      borderRadius: 2,
                      bgcolor: (theme) => alpha(paletteColor(theme.palette, stat.color), 0.15),
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center'
                    }}>
                      <stat.icon sx={{ color: stat.color }} />
                    </Box>
                  </Box>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Charts Row */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
                Revenue vs Spend Trend
              </Typography>
              <CustomAreaChart
                data={revenueData}
                areas={[
                  { dataKey: 'revenue', color: '#10b981', name: 'Revenue' },
                  { dataKey: 'spend', color: '#ef4444', name: 'Spend' },
                ]}
                height={280}
              />
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
                Traffic by Channel
              </Typography>
              <CustomPieChart
                data={dashboard?.top_channels || defaultDashboard.top_channels}
                height={280}
                innerRadius={50}
                outerRadius={90}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Campaign Performance Chart */}
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Campaign Performance by Type
          </Typography>
          <CustomBarChart
            data={campaignPerformanceData}
            bars={[
              { dataKey: 'conversions', color: '#2563eb', name: 'Conversions' },
              { dataKey: 'spend', color: '#f59e0b', name: 'Spend ($)' },
            ]}
            height={250}
          />
        </CardContent>
      </Card>

      {/* Campaigns Table */}
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6" fontWeight="bold">
              Recent Campaigns
            </Typography>
            <Button size="small" variant="outlined">
              View All
            </Button>
          </Box>
          {campaignsLoading ? (
            <Skeleton height={200} />
          ) : (
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Campaign</TableCell>
                    <TableCell>Type</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Budget</TableCell>
                    <TableCell>Spent</TableCell>
                    <TableCell>Progress</TableCell>
                    <TableCell>CTR</TableCell>
                    <TableCell>Conversions</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {campaigns.slice(0, 5).map((campaign) => (
                    <TableRow key={campaign.id} hover>
                      <TableCell>
                        <Typography fontWeight="medium">{campaign.name}</Typography>
                      </TableCell>
                      <TableCell>
                        <Chip label={campaign.type} size="small" color={getTypeColor(campaign.type)} />
                      </TableCell>
                      <TableCell>
                        <Chip label={campaign.status} size="small" color={getStatusColor(campaign.status)} />
                      </TableCell>
                      <TableCell>{formatCurrency(campaign.budget)}</TableCell>
                      <TableCell>{formatCurrency(campaign.spent)}</TableCell>
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <LinearProgress
                            variant="determinate"
                            value={campaign.budget_utilization}
                            sx={{ width: 80, height: 8, borderRadius: 4 }}
                            color={campaign.budget_utilization > 90 ? 'error' : campaign.budget_utilization > 70 ? 'warning' : 'success'}
                          />
                          <Typography variant="body2">
                            {campaign.budget_utilization.toFixed(0)}%
                          </Typography>
                        </Box>
                      </TableCell>
                      <TableCell>{campaign.ctr.toFixed(2)}%</TableCell>
                      <TableCell>{formatNumber(campaign.conversions)}</TableCell>
                      <TableCell>
                        <IconButton size="small" onClick={() => router.push(`/campaigns?edit=${campaign.id}`)} aria-label={`Edit campaign ${campaign.name}`}>
                          <EditIcon fontSize="small" />
                        </IconButton>
                        <IconButton size="small" onClick={() => {
                          if (campaign.status === 'active') {
                            pauseMutation.mutate(campaign.id)
                          } else {
                            activateMutation.mutate(campaign.id)
                          }
                        }} aria-label={campaign.status === 'active' ? 'Pause campaign' : 'Activate campaign'}>
                          {campaign.status === 'active' ? <PauseIcon fontSize="small" /> : <PlayArrowIcon fontSize="small" />}
                        </IconButton>
                        <IconButton size="small" color="error" onClick={() => {
                          confirmDialog.showConfirm(
                            'Delete Campaign',
                            `Are you sure you want to delete "${campaign.name}"? This action cannot be undone.`,
                            () => deleteMutation.mutate(campaign.id)
                          )
                        }} aria-label={`Delete campaign ${campaign.name}`}>
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </CardContent>
      </Card>

      <ConfirmDialog
        open={confirmDialog.dialogProps.open}
        title={confirmDialog.dialogProps.title}
        message={confirmDialog.dialogProps.message}
        onConfirm={confirmDialog.dialogProps.onConfirm as () => void}
        onCancel={confirmDialog.dialogProps.onCancel}
        isLoading={confirmDialog.dialogProps.isLoading}
      />
    </Box>
  )
}
