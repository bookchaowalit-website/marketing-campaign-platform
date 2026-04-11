'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, Tabs, Tab, Skeleton, IconButton, Grid, Select, MenuItem, FormControl,
  InputLabel, LinearProgress, TableSortLabel
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import SendIcon from '@mui/icons-material/Send'
import BarChartIcon from '@mui/icons-material/BarChart'
import ScienceIcon from '@mui/icons-material/Science'
import EmailIcon from '@mui/icons-material/Email'
import { campaignApi, emailCampaignApi } from '@/lib/api-client'
import { marketingApi } from '@/lib/marketing-features-api'
import { useToast } from '@/hooks/useToast'
import { TabPanel } from '@/components/TabPanel'
import { EmptyState } from '@/components/EmptyState'
import { formatNumber } from '@/lib/utils'

const defaultCampaigns = [
  { id: '1', name: 'Welcome Series', subject: 'Welcome to our community!', list_name: 'New Subscribers', status: 'sent', sent: 500, opened: 250, clicked: 100, bounced: 10, open_rate: 50.0, click_rate: 20.0 },
  { id: '2', name: 'Monthly Newsletter', subject: 'April Newsletter - What\'s New', list_name: 'All Subscribers', status: 'sent', sent: 2000, opened: 800, clicked: 200, bounced: 30, open_rate: 40.0, click_rate: 10.0 },
  { id: '3', name: 'Product Launch', subject: 'Introducing our newest feature', list_name: 'Active Users', status: 'draft', sent: 0, opened: 0, clicked: 0, bounced: 0, open_rate: 0, click_rate: 0 },
]

const defaultABTests = [
  { id: '1', campaign_a: 'Welcome A', campaign_b: 'Welcome B', status: 'completed', metric: 'open_rate', winner: 'Welcome A', winner_rate: 55.0 },
  { id: '2', name: 'Subject Line Test', campaign_a: 'Newsletter A', campaign_b: 'Newsletter B', status: 'running', metric: 'click_rate', winner: null, winner_rate: 0 },
]

const defaultLists = [
  { id: '1', name: 'New Subscribers', subscriber_count: 500 },
  { id: '2', name: 'All Subscribers', subscriber_count: 2000 },
  { id: '3', name: 'Active Users', subscriber_count: 800 },
]

export default function EmailCampaignsPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const [tabValue, setTabValue] = useState(0)
  const [orderBy, setOrderBy] = useState('name')
  const [order, setOrder] = useState<'asc' | 'desc'>('asc')

  // Dialog states
  const [openCampaignDialog, setOpenCampaignDialog] = useState(false)
  const [openABTestDialog, setOpenABTestDialog] = useState(false)
  const [openAnalyticsModal, setOpenAnalyticsModal] = useState(false)
  const [analyticsData, setAnalyticsData] = useState<any>(null)

  const [campaignForm, setCampaignForm] = useState({
    name: '', subject: '', preheader: '', content: '', list_ids: [] as string[],
  })
  const [abTestForm, setABTestForm] = useState({
    name: '', campaign_a_id: '', campaign_b_id: '', metric: 'open_rate',
  })

  const { data: campaignsData, isLoading: campaignsLoading } = useQuery({
    queryKey: ['email-campaigns'],
    queryFn: async () => {
      try {
        return await marketingApi.getEmailCampaigns()
      } catch {
        return { campaigns: defaultCampaigns }
      }
    },
  })

  const { data: listsData } = useQuery({
    queryKey: ['email-lists'],
    queryFn: async () => {
      try {
        return await marketingApi.getEmailLists()
      } catch {
        return { lists: defaultLists }
      }
    },
  })

  const { data: abTestsData, isLoading: abTestsLoading } = useQuery({
    queryKey: ['email-ab-tests'],
    queryFn: async () => {
      try {
        const response = await fetch('/api/go/v1/email/ab-tests', {
          headers: { Authorization: `Bearer ${typeof window !== 'undefined' ? localStorage.getItem('token') : ''}` },
        })
        if (!response.ok) throw new Error('Failed')
        return response.json()
      } catch {
        return { tests: defaultABTests }
      }
    },
  })

  const campaigns = campaignsData?.campaigns || defaultCampaigns
  const lists = listsData?.lists || defaultLists
  const abTests = abTestsData?.tests || defaultABTests

  const createCampaignMutation = useMutation({
    mutationFn: async (data: typeof campaignForm) => {
      return marketingApi.createEmailCampaign({
        name: data.name,
        subject: data.subject,
        html_content: data.content,
        list_ids: data.list_ids,
      } as any)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['email-campaigns'] })
      setOpenCampaignDialog(false)
      setCampaignForm({ name: '', subject: '', preheader: '', content: '', list_ids: [] })
      toast.success('Campaign created successfully')
    },
    onError: () => toast.error('Failed to create campaign'),
  })

  const sendCampaignMutation = useMutation({
    mutationFn: async (id: string) => emailCampaignApi.sendCampaign(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['email-campaigns'] })
      toast.success('Campaign sent successfully')
    },
    onError: () => toast.error('Failed to send campaign'),
  })

  const createABTestMutation = useMutation({
    mutationFn: async (data: typeof abTestForm) => emailCampaignApi.createABTest(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['email-ab-tests'] })
      setOpenABTestDialog(false)
      setABTestForm({ name: '', campaign_a_id: '', campaign_b_id: '', metric: 'open_rate' })
      toast.success('A/B test created successfully')
    },
    onError: () => toast.error('Failed to create A/B test'),
  })

  const handleViewAnalytics = async (campaign: any) => {
    try {
      const data = await emailCampaignApi.getAnalytics(campaign.id)
      setAnalyticsData(data)
    } catch {
      setAnalyticsData({
        sent: campaign.sent,
        delivered: campaign.sent - campaign.bounced,
        opened: campaign.opened,
        clicked: campaign.clicked,
        bounced: campaign.bounced,
      })
    }
    setOpenAnalyticsModal(true)
  }

  const getStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' | 'error' => {
    switch (status) {
      case 'sent': case 'delivered': case 'completed': return 'success'
      case 'draft': case 'pending': case 'scheduled': return 'warning'
      case 'running': case 'active': return 'info'
      case 'failed': case 'error': return 'error'
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

  const barData = analyticsData
    ? [
      { label: 'Sent', value: analyticsData.sent || 0, color: '#2196f3' },
      { label: 'Delivered', value: analyticsData.delivered || 0, color: '#4caf50' },
      { label: 'Opened', value: analyticsData.opened || 0, color: '#ff9800' },
      { label: 'Clicked', value: analyticsData.clicked || 0, color: '#9c27b0' },
      { label: 'Bounced', value: analyticsData.bounced || 0, color: '#f44336' },
    ]
    : []

  const maxValue = Math.max(...barData.map((b) => b.value), 1)

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Email Campaigns
        </Typography>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            variant="outlined"
            startIcon={<ScienceIcon />}
            onClick={() => setOpenABTestDialog(true)}
          >
            New A/B Test
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => setOpenCampaignDialog(true)}
          >
            New Campaign
          </Button>
        </Box>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={(e, v) => setTabValue(v)}>
          <Tab label="Campaigns" />
          <Tab label="A/B Tests" />
        </Tabs>
      </Box>

      {/* Campaigns Tab */}
      <TabPanel value={tabValue} index={0}>
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
                        <TableCell>Subject</TableCell>
                        <TableCell>List</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'sent'} direction={order} onClick={() => handleSort('sent')}>Sent</TableSortLabel>
                        </TableCell>
                        <TableCell>Opened</TableCell>
                        <TableCell>Clicked</TableCell>
                        <TableCell>Bounced</TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'open_rate'} direction={order} onClick={() => handleSort('open_rate')}>Open Rate</TableSortLabel>
                        </TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'click_rate'} direction={order} onClick={() => handleSort('click_rate')}>Click Rate</TableSortLabel>
                        </TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {campaigns.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={11} align="center">
                            <EmptyState icon={EmailIcon} title="No campaigns found" description="Create your first email campaign to start reaching subscribers." actionLabel="New Campaign" onAction={() => setOpenCampaignDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        sortedCampaigns.map((campaign: any) => (
                          <TableRow key={campaign.id} hover>
                            <TableCell><Typography fontWeight="medium">{campaign.name}</Typography></TableCell>
                            <TableCell>{campaign.subject}</TableCell>
                            <TableCell>{campaign.list_name || '-'}</TableCell>
                            <TableCell>
                              <Chip label={campaign.status} size="small" color={getStatusColor(campaign.status)} />
                            </TableCell>
                            <TableCell>{campaign.sent}</TableCell>
                            <TableCell>{campaign.opened}</TableCell>
                            <TableCell>{campaign.clicked}</TableCell>
                            <TableCell>{campaign.bounced}</TableCell>
                            <TableCell>{campaign.open_rate?.toFixed(1)}%</TableCell>
                            <TableCell>{campaign.click_rate?.toFixed(1)}%</TableCell>
                            <TableCell>
                              {campaign.status === 'draft' && (
                                <IconButton size="small" onClick={() => sendCampaignMutation.mutate(campaign.id)}
                                  title="Send" aria-label="Send campaign">
                                  <SendIcon fontSize="small" />
                                </IconButton>
                              )}
                              {campaign.status === 'sent' && (
                                <IconButton size="small" onClick={() => handleViewAnalytics(campaign)}
                                  title="Analytics" aria-label="View analytics">
                                  <BarChartIcon fontSize="small" />
                                </IconButton>
                              )}
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

      {/* A/B Tests Tab */}
      <TabPanel value={tabValue} index={1}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">A/B Tests</Typography>
            </Box>
            {abTestsLoading ? <Skeleton height={200} /> : (
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell>Campaign A</TableCell>
                        <TableCell>Campaign B</TableCell>
                        <TableCell>Test Metric</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Winner</TableCell>
                        <TableCell>Winner Rate</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {abTests.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={7} align="center">
                            <EmptyState icon={ScienceIcon} title="No A/B tests found" description="Create an A/B test to optimize your email performance." actionLabel="New A/B Test" onAction={() => setOpenABTestDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        abTests.map((test: any) => (
                          <TableRow key={test.id} hover>
                            <TableCell><Typography fontWeight="medium">{test.name || `${test.campaign_a} vs ${test.campaign_b}`}</Typography></TableCell>
                            <TableCell>{test.campaign_a}</TableCell>
                            <TableCell>{test.campaign_b}</TableCell>
                            <TableCell>
                              <Chip label={test.metric} size="small" variant="outlined" />
                            </TableCell>
                            <TableCell>
                              <Chip label={test.status} size="small" color={getStatusColor(test.status)} />
                            </TableCell>
                            <TableCell>
                              {test.winner ? (
                                <Chip label={test.winner} size="small" color="success" />
                              ) : (
                                <Typography variant="body2" color="text.secondary">-</Typography>
                              )}
                            </TableCell>
                            <TableCell>{test.winner_rate ? `${test.winner_rate}%` : '-'}</TableCell>
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

      {/* New Campaign Dialog */}
      <Dialog open={openCampaignDialog} onClose={() => setOpenCampaignDialog(false)} maxWidth="md" fullWidth>
        <DialogTitle>New Email Campaign</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Campaign Name"
              fullWidth
              value={campaignForm.name}
              onChange={(e) => setCampaignForm({ ...campaignForm, name: e.target.value })}
            />
            <TextField
              label="Subject Line"
              fullWidth
              value={campaignForm.subject}
              onChange={(e) => setCampaignForm({ ...campaignForm, subject: e.target.value })}
            />
            <TextField
              label="Preheader"
              fullWidth
              value={campaignForm.preheader}
              onChange={(e) => setCampaignForm({ ...campaignForm, preheader: e.target.value })}
              helperText="Preview text shown after subject line"
            />
            <FormControl fullWidth>
              <InputLabel>Select Lists</InputLabel>
              <Select
                multiple
                value={campaignForm.list_ids}
                label="Select Lists"
                onChange={(e) => setCampaignForm({ ...campaignForm, list_ids: e.target.value as string[] })}
                renderValue={(selected) =>
                  (selected as string[]).map((id) => lists.find((l: any) => l.id === id)?.name).filter(Boolean).join(', ')
                }
              >
                {lists.map((list: any) => (
                  <MenuItem key={list.id} value={list.id}>
                    {list.name} ({list.subscriber_count} subscribers)
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <TextField
              label="Content (HTML)"
              fullWidth
              multiline
              rows={8}
              value={campaignForm.content}
              onChange={(e) => setCampaignForm({ ...campaignForm, content: e.target.value })}
              placeholder="<html>Your email content here...</html>"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenCampaignDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => createCampaignMutation.mutate(campaignForm)}
            disabled={createCampaignMutation.isPending}>
            {createCampaignMutation.isPending ? 'Creating...' : 'Create Campaign'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* New A/B Test Dialog */}
      <Dialog open={openABTestDialog} onClose={() => setOpenABTestDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>New A/B Test</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Test Name"
              fullWidth
              value={abTestForm.name}
              onChange={(e) => setABTestForm({ ...abTestForm, name: e.target.value })}
            />
            <TextField
              label="Campaign A"
              fullWidth
              select
              SelectProps={{ native: true }}
              value={abTestForm.campaign_a_id}
              onChange={(e) => setABTestForm({ ...abTestForm, campaign_a_id: e.target.value })}
            >
              <option value="">-- Select Campaign A --</option>
              {campaigns.map((c: any) => (
                <option key={`a-${c.id}`} value={c.id}>{c.name}</option>
              ))}
            </TextField>
            <TextField
              label="Campaign B"
              fullWidth
              select
              SelectProps={{ native: true }}
              value={abTestForm.campaign_b_id}
              onChange={(e) => setABTestForm({ ...abTestForm, campaign_b_id: e.target.value })}
            >
              <option value="">-- Select Campaign B --</option>
              {campaigns.map((c: any) => (
                <option key={`b-${c.id}`} value={c.id}>{c.name}</option>
              ))}
            </TextField>
            <TextField
              label="Test Metric"
              fullWidth
              select
              value={abTestForm.metric}
              onChange={(e) => setABTestForm({ ...abTestForm, metric: e.target.value })}
            >
              <MenuItem value="open_rate">Open Rate</MenuItem>
              <MenuItem value="click_rate">Click Rate</MenuItem>
            </TextField>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenABTestDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => createABTestMutation.mutate(abTestForm)}
            disabled={createABTestMutation.isPending}>
            {createABTestMutation.isPending ? 'Creating...' : 'Create A/B Test'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Analytics Dialog */}
      <Dialog open={openAnalyticsModal} onClose={() => setOpenAnalyticsModal(false)} maxWidth="md" fullWidth>
        <DialogTitle>Campaign Analytics</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {barData.map((bar) => (
              <Box key={bar.label}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                  <Typography variant="body2">{bar.label}</Typography>
                  <Typography variant="body2" fontWeight="medium">{formatNumber(bar.value)}</Typography>
                </Box>
                <LinearProgress
                  variant="determinate"
                  value={(bar.value / maxValue) * 100}
                  sx={{ height: 8, borderRadius: 4, bgcolor: 'grey.200', '& .MuiLinearProgress-bar': { bgcolor: bar.color } }}
                />
              </Box>
            ))}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAnalyticsModal(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
