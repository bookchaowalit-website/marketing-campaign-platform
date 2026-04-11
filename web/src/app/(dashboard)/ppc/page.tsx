'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, MenuItem, Tabs, Tab, Skeleton, IconButton, Grid, TableSortLabel
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import PlayArrowIcon from '@mui/icons-material/PlayArrow'
import PauseIcon from '@mui/icons-material/Pause'
import BarChartIcon from '@mui/icons-material/BarChart'
import CloudIcon from '@mui/icons-material/Cloud'
import { adsApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { TabPanel } from '@/components/TabPanel'
import { EmptyState } from '@/components/EmptyState'
import { formatCurrency, formatNumber } from '@/lib/utils'

const PLATFORM_OPTIONS = [
  { value: 'google_ads', label: 'Google Ads' },
  { value: 'facebook_ads', label: 'Facebook Ads' },
  { value: 'tiktok_ads', label: 'TikTok Ads' },
  { value: 'linkedin_ads', label: 'LinkedIn Ads' },
  { value: 'twitter_ads', label: 'Twitter Ads' },
]

const defaultPlatforms = [
  { id: '1', name: 'Google Ads', status: 'connected', account_id: 'google-123' },
  { id: '2', name: 'Facebook Ads', status: 'connected', account_id: 'fb-456' },
]

const defaultCampaigns = [
  { id: '1', name: 'Google Search - Brand', platform: 'Google Ads', budget: 5000, spent: 3200, impressions: 150000, clicks: 4500, conversions: 180, ctr: 3.0, status: 'active' },
  { id: '2', name: 'Facebook Retargeting', platform: 'Facebook Ads', budget: 3000, spent: 1800, impressions: 200000, clicks: 6000, conversions: 120, ctr: 3.0, status: 'active' },
  { id: '3', name: 'TikTok Awareness', platform: 'TikTok Ads', budget: 2000, spent: 500, impressions: 500000, clicks: 15000, conversions: 50, ctr: 3.0, status: 'paused' },
]

export default function PPCPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const [tabValue, setTabValue] = useState(0)
  const [orderBy, setOrderBy] = useState('name')
  const [order, setOrder] = useState<'asc' | 'desc'>('asc')

  // Dialog states
  const [openPlatformDialog, setOpenPlatformDialog] = useState(false)
  const [openCampaignDialog, setOpenCampaignDialog] = useState(false)
  const [platformForm, setPlatformForm] = useState({ platform: '', account_id: '', api_key: '' })
  const [campaignForm, setCampaignForm] = useState({
    name: '', platform: '', budget: '', targeting: '', start_date: '', end_date: '',
  })

  const { data: platformsData, isLoading: platformsLoading } = useQuery({
    queryKey: ['ads-platforms'],
    queryFn: async () => {
      try {
        return await adsApi.getPlatforms()
      } catch {
        return { platforms: defaultPlatforms }
      }
    },
  })

  const { data: campaignsData, isLoading: campaignsLoading } = useQuery({
    queryKey: ['ads-campaigns'],
    queryFn: async () => {
      try {
        return await adsApi.getCampaigns()
      } catch {
        return { campaigns: defaultCampaigns }
      }
    },
  })

  const platforms = platformsData?.platforms || defaultPlatforms
  const campaigns = campaignsData?.campaigns || defaultCampaigns

  const connectPlatformMutation = useMutation({
    mutationFn: async (data: typeof platformForm) => {
      return adsApi.connectPlatform(data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ads-platforms'] })
      setOpenPlatformDialog(false)
      setPlatformForm({ platform: '', account_id: '', api_key: '' })
      toast.success('Platform connected successfully')
    },
    onError: () => toast.error('Failed to connect platform'),
  })

  const createCampaignMutation = useMutation({
    mutationFn: async (data: typeof campaignForm) => {
      return adsApi.createCampaign({ ...data, budget: Number(data.budget) })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ads-campaigns'] })
      setOpenCampaignDialog(false)
      setCampaignForm({ name: '', platform: '', budget: '', targeting: '', start_date: '', end_date: '' })
      toast.success('Campaign created successfully')
    },
    onError: () => toast.error('Failed to create campaign'),
  })

  const pauseMutation = useMutation({
    mutationFn: async (id: string) => adsApi.pauseCampaign(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ads-campaigns'] })
      toast.success('Campaign paused')
    },
    onError: () => toast.error('Failed to pause campaign'),
  })

  const activateMutation = useMutation({
    mutationFn: async (id: string) => adsApi.activateCampaign(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ads-campaigns'] })
      toast.success('Campaign activated')
    },
    onError: () => toast.error('Failed to activate campaign'),
  })

  const getStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' => {
    switch (status) {
      case 'active': return 'success'
      case 'paused': return 'warning'
      case 'completed': return 'default'
      case 'draft': return 'info'
      default: return 'default'
    }
  }

  const handleSort = (field: string) => {
    if (orderBy === field) {
      setOrder(order === 'asc' ? 'desc' : 'asc')
    } else {
      setOrderBy(field)
      setOrder('asc')
    }
  }

  const sortedCampaigns = [...campaigns].sort((a, b) => {
    const aVal = a[orderBy as keyof typeof a]
    const bVal = b[orderBy as keyof typeof b]
    if (typeof aVal === 'string' && typeof bVal === 'string') {
      return order === 'asc' ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal)
    }
    if (typeof aVal === 'number' && typeof bVal === 'number') {
      return order === 'asc' ? aVal - bVal : bVal - aVal
    }
    return 0
  })

  const getPlatformColor = (platform: string): 'primary' | 'secondary' | 'success' | 'warning' | 'info' => {
    switch (platform) {
      case 'Google Ads': return 'primary'
      case 'Facebook Ads': return 'secondary'
      case 'TikTok Ads': return 'warning'
      case 'LinkedIn Ads': return 'info'
      case 'Twitter Ads': return 'success'
      default: return 'primary'
    }
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          PPC / Paid Ads
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setOpenCampaignDialog(true)}
        >
          New Campaign
        </Button>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={(e, v) => setTabValue(v)}>
          <Tab label="Platforms" />
          <Tab label="Campaigns" />
        </Tabs>
      </Box>

      {/* Platforms Tab */}
      <TabPanel value={tabValue} index={0}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">Connected Platforms</Typography>
              <Button startIcon={<AddIcon />} onClick={() => setOpenPlatformDialog(true)}>
                Connect Platform
              </Button>
            </Box>
            {platformsLoading ? <Skeleton height={200} /> : (
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Platform</TableCell>
                        <TableCell>Account ID</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {platforms.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={4} align="center">
                            <EmptyState icon={CloudIcon} title="No platforms connected" description="Connect your first ad platform to start managing PPC campaigns." actionLabel="Connect Platform" onAction={() => setOpenPlatformDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        platforms.map((platform: any) => (
                          <TableRow key={platform.id} hover>
                            <TableCell>
                              <Chip label={platform.name} size="small" color={getPlatformColor(platform.name)} />
                            </TableCell>
                            <TableCell>{platform.account_id || '-'}</TableCell>
                            <TableCell>
                              <Chip label={platform.status} size="small" color={platform.status === 'connected' ? 'success' : 'default'} />
                            </TableCell>
                            <TableCell>
                              <Button size="small" startIcon={<CloudIcon />} onClick={() => toast.info(`Syncing ${platform.name}...`)}>Sync</Button>
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Campaigns Tab */}
      <TabPanel value={tabValue} index={1}>
        <Card>
          <CardContent>
            {campaignsLoading ? <Skeleton height={300} /> : (
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell sx={{ minWidth: 150 }}>
                          <TableSortLabel active={orderBy === 'name'} direction={order} onClick={() => handleSort('name')}>Name</TableSortLabel>
                        </TableCell>
                        <TableCell>Platform</TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'budget'} direction={order} onClick={() => handleSort('budget')}>Budget</TableSortLabel>
                        </TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'spent'} direction={order} onClick={() => handleSort('spent')}>Spent</TableSortLabel>
                        </TableCell>
                        <TableCell>Impressions</TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'clicks'} direction={order} onClick={() => handleSort('clicks')}>Clicks</TableSortLabel>
                        </TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'conversions'} direction={order} onClick={() => handleSort('conversions')}>Conversions</TableSortLabel>
                        </TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'ctr'} direction={order} onClick={() => handleSort('ctr')}>CTR</TableSortLabel>
                        </TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {sortedCampaigns.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={10} align="center">
                            <EmptyState icon={BarChartIcon} title="No campaigns found" description="Create your first PPC campaign to start tracking ad performance." actionLabel="New Campaign" onAction={() => setOpenCampaignDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        sortedCampaigns.map((campaign: any) => (
                          <TableRow key={campaign.id} hover>
                            <TableCell>
                              <Typography fontWeight="medium">{campaign.name}</Typography>
                            </TableCell>
                            <TableCell>
                              <Chip label={campaign.platform} size="small" color={getPlatformColor(campaign.platform)} />
                            </TableCell>
                            <TableCell>{formatCurrency(campaign.budget)}</TableCell>
                            <TableCell>{formatCurrency(campaign.spent)}</TableCell>
                            <TableCell>{formatNumber(campaign.impressions)}</TableCell>
                            <TableCell>{formatNumber(campaign.clicks)}</TableCell>
                            <TableCell>{campaign.conversions}</TableCell>
                            <TableCell>{campaign.ctr?.toFixed(2)}%</TableCell>
                            <TableCell>
                              <Chip label={campaign.status} size="small" color={getStatusColor(campaign.status)} />
                            </TableCell>
                            <TableCell>
                              {campaign.status === 'active' ? (
                                <IconButton size="small" onClick={() => pauseMutation.mutate(campaign.id)} title="Pause" aria-label="Pause campaign">
                                  <PauseIcon fontSize="small" />
                                </IconButton>
                              ) : (
                                <IconButton size="small" onClick={() => activateMutation.mutate(campaign.id)} title="Activate" aria-label="Activate campaign">
                                  <PlayArrowIcon fontSize="small" />
                                </IconButton>
                              )}
                              <IconButton size="small" title="View Performance" aria-label="View performance">
                                <BarChartIcon fontSize="small" />
                              </IconButton>
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Connect Platform Dialog */}
      <Dialog open={openPlatformDialog} onClose={() => setOpenPlatformDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Connect Ad Platform</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Platform"
              fullWidth
              select
              value={platformForm.platform}
              onChange={(e) => setPlatformForm({ ...platformForm, platform: e.target.value })}
            >
              {PLATFORM_OPTIONS.map((opt) => (
                <MenuItem key={opt.value} value={opt.value}>{opt.label}</MenuItem>
              ))}
            </TextField>
            <TextField
              label="Account ID"
              fullWidth
              value={platformForm.account_id}
              onChange={(e) => setPlatformForm({ ...platformForm, account_id: e.target.value })}
            />
            <TextField
              label="API Key / Token"
              fullWidth
              type="password"
              value={platformForm.api_key}
              onChange={(e) => setPlatformForm({ ...platformForm, api_key: e.target.value })}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenPlatformDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => connectPlatformMutation.mutate(platformForm)}
            disabled={connectPlatformMutation.isPending}>
            {connectPlatformMutation.isPending ? 'Connecting...' : 'Connect'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* New Campaign Dialog */}
      <Dialog open={openCampaignDialog} onClose={() => setOpenCampaignDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>New Ad Campaign</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Campaign Name"
              fullWidth
              required
              value={campaignForm.name}
              onChange={(e) => setCampaignForm({ ...campaignForm, name: e.target.value })}
            />
            <TextField
              label="Platform"
              fullWidth
              select
              value={campaignForm.platform}
              onChange={(e) => setCampaignForm({ ...campaignForm, platform: e.target.value })}
            >
              {PLATFORM_OPTIONS.map((opt) => (
                <MenuItem key={opt.value} value={opt.value}>{opt.label}</MenuItem>
              ))}
            </TextField>
            <TextField
              label="Budget"
              fullWidth
              type="number"
              inputProps={{ min: 0 }}
              value={campaignForm.budget}
              onChange={(e) => setCampaignForm({ ...campaignForm, budget: e.target.value })}
            />
            <TextField
              label="Targeting (JSON or description)"
              fullWidth
              multiline
              rows={2}
              value={campaignForm.targeting}
              onChange={(e) => setCampaignForm({ ...campaignForm, targeting: e.target.value })}
            />
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField
                  label="Start Date"
                  fullWidth
                  type="date"
                  InputLabelProps={{ shrink: true }}
                  value={campaignForm.start_date}
                  onChange={(e) => setCampaignForm({ ...campaignForm, start_date: e.target.value })}
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  label="End Date"
                  fullWidth
                  type="date"
                  InputLabelProps={{ shrink: true }}
                  value={campaignForm.end_date}
                  onChange={(e) => setCampaignForm({ ...campaignForm, end_date: e.target.value })}
                />
              </Grid>
            </Grid>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenCampaignDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => createCampaignMutation.mutate(campaignForm)}
            disabled={createCampaignMutation.isPending || !campaignForm.name.trim() || !campaignForm.budget || Number(campaignForm.budget) <= 0}>
            {createCampaignMutation.isPending ? 'Creating...' : 'Create Campaign'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
