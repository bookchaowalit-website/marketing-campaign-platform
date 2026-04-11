'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useMemo } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent,
  DialogActions, TextField, MenuItem, Skeleton, IconButton, Grid,
  FormControl, InputLabel, Select,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import SendIcon from '@mui/icons-material/Send'
import CheckIcon from '@mui/icons-material/Check'
import CloseIcon from '@mui/icons-material/Close'
import CampaignIcon from '@mui/icons-material/Campaign'
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline'
import { unifiedCampaignApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { ConfirmDialog, useConfirmDialog } from '@/components/ConfirmDialog'
import { EmptyState } from '@/components/EmptyState'
import { formatCurrency } from '@/lib/utils'

const CHANNELS = [
  'Google Ads', 'Facebook', 'Instagram', 'TikTok', 'Twitter/X', 'LinkedIn', 'YouTube',
  'Pinterest', 'Snapchat', 'Reddit', 'Google Display', 'Programmatic', 'Native Ads',
  'Email', 'SMS', 'Push Notifications', 'WhatsApp', 'Telegram',
  'SEO', 'Content Marketing', 'Blog', 'Podcast', 'Webinar',
  'TV', 'Radio', 'Print', 'OOH (Billboard)', 'Transit Ads', 'Newspaper', 'Magazine',
  'Events', 'Trade Shows', 'Sponsorships', 'PR', 'Influencer', 'Affiliate',
  'Partnerships', 'Referral', 'Direct Mail', 'Catalog', 'Telemarketing',
]

const CHANNEL_CATEGORIES = ['all', 'online', 'offline', 'hybrid']

const ONLINE_CHANNELS = new Set([
  'Google Ads', 'Facebook', 'Instagram', 'TikTok', 'Twitter/X', 'LinkedIn', 'YouTube',
  'Pinterest', 'Snapchat', 'Reddit', 'Google Display', 'Programmatic', 'Native Ads',
  'Email', 'SMS', 'Push Notifications', 'WhatsApp', 'Telegram',
  'SEO', 'Content Marketing', 'Blog', 'Podcast', 'Webinar',
])

const OFFLINE_CHANNELS = new Set([
  'TV', 'Radio', 'Print', 'OOH (Billboard)', 'Transit Ads', 'Newspaper', 'Magazine',
  'Events', 'Trade Shows', 'Sponsorships', 'PR', 'Direct Mail', 'Catalog', 'Telemarketing',
])

const defaultCampaigns: any[] = []

export default function UnifiedCampaignsPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const confirmDialog = useConfirmDialog()

  const [openDialog, setOpenDialog] = useState(false)
  const [editingCampaign, setEditingCampaign] = useState<any>(null)
  const [formData, setFormData] = useState({
    name: '',
    channel: 'Google Ads',
    category: 'online',
    budget: 0,
    start_date: new Date().toISOString().split('T')[0],
    end_date: '',
    description: '',
  })

  const [categoryFilter, setCategoryFilter] = useState<string>('all')
  const [statusFilter, setStatusFilter] = useState<string>('all')

  const { data: campaignsData, isLoading, error, refetch } = useQuery<{ campaigns: any[]; total: number }>({
    queryKey: ['unified-campaigns', categoryFilter, statusFilter],
    queryFn: async () => {
      const channel = categoryFilter === 'all' ? undefined : undefined
      const status = statusFilter === 'all' ? undefined : statusFilter
      return unifiedCampaignApi.getAll(channel, status)
    },
    initialData: { campaigns: defaultCampaigns, total: 0 },
  })

  const campaigns = campaignsData?.campaigns || defaultCampaigns

  const filteredCampaigns = useMemo(() => {
    return campaigns.filter((c: any) => {
      const cat = (c.category || '').toLowerCase()
      const matchesCategory = categoryFilter === 'all' || cat === categoryFilter
      return matchesCategory
    })
  }, [campaigns, categoryFilter])

  const getChannelCategory = (channel: string): string => {
    if (ONLINE_CHANNELS.has(channel)) return 'online'
    if (OFFLINE_CHANNELS.has(channel)) return 'offline'
    return 'hybrid'
  }

  const createMutation = useMutation({
    mutationFn: async (data: typeof formData) => {
      await unifiedCampaignApi.create({
        ...data,
        budget: Number(data.budget),
        category: getChannelCategory(data.channel),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['unified-campaigns'] })
      setOpenDialog(false)
      setFormData({ name: '', channel: 'Google Ads', category: 'online', budget: 0, start_date: new Date().toISOString().split('T')[0], end_date: '', description: '' })
      setEditingCampaign(null)
      toast.success('Campaign created successfully')
    },
    onError: () => {
      toast.error('Failed to create campaign')
    },
  })

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: typeof formData }) => {
      await unifiedCampaignApi.update(id, {
        ...data,
        budget: Number(data.budget),
        category: getChannelCategory(data.channel),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['unified-campaigns'] })
      setOpenDialog(false)
      setEditingCampaign(null)
      toast.success('Campaign updated successfully')
    },
    onError: () => {
      toast.error('Failed to update campaign')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await unifiedCampaignApi.delete(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['unified-campaigns'] })
      toast.success('Campaign deleted successfully')
    },
    onError: () => {
      toast.error('Failed to delete campaign')
    },
  })

  const submitMutation = useMutation({
    mutationFn: async (id: string) => {
      await unifiedCampaignApi.submitForApproval(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['unified-campaigns'] })
      toast.success('Campaign submitted for approval')
    },
    onError: () => {
      toast.error('Failed to submit for approval')
    },
  })

  const approveMutation = useMutation({
    mutationFn: async (id: string) => {
      await unifiedCampaignApi.approve(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['unified-campaigns'] })
      toast.success('Campaign approved')
    },
    onError: () => {
      toast.error('Failed to approve campaign')
    },
  })

  const rejectMutation = useMutation({
    mutationFn: async (id: string) => {
      await unifiedCampaignApi.reject(id, 'Rejected by reviewer')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['unified-campaigns'] })
      toast.success('Campaign rejected')
    },
    onError: () => {
      toast.error('Failed to reject campaign')
    },
  })

  const getStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' | 'error' => {
    switch (status) {
      case 'active': return 'success'
      case 'paused': return 'warning'
      case 'completed': return 'default'
      case 'draft': return 'info'
      default: return 'default'
    }
  }

  const getApprovalStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' | 'error' => {
    switch (status) {
      case 'approved': return 'success'
      case 'rejected': return 'error'
      case 'pending': return 'warning'
      case 'submitted': return 'info'
      default: return 'default'
    }
  }

  const handleOpenDialog = () => {
    setEditingCampaign(null)
    setFormData({ name: '', channel: 'Google Ads', category: 'online', budget: 0, start_date: new Date().toISOString().split('T')[0], end_date: '', description: '' })
    setOpenDialog(true)
  }

  const handleSubmit = () => {
    if (editingCampaign) {
      updateMutation.mutate({ id: editingCampaign.id, data: formData })
    } else {
      createMutation.mutate(formData)
    }
  }

  const handleDelete = (campaign: any) => {
    confirmDialog.showConfirm(
      'Delete Campaign',
      `Are you sure you want to delete "${campaign.name}"? This action cannot be undone.`,
      () => deleteMutation.mutate(campaign.id)
    )
  }

  const handleChannelChange = (channel: string) => {
    setFormData({ ...formData, channel, category: getChannelCategory(channel) })
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Unified Campaigns
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleOpenDialog}
        >
          New Campaign
        </Button>
      </Box>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={6} md={4}>
              <FormControl fullWidth size="small">
                <InputLabel>Channel Category</InputLabel>
                <Select
                  value={categoryFilter}
                  label="Channel Category"
                  onChange={(e) => setCategoryFilter(e.target.value)}
                >
                  <MenuItem value="all">All Categories</MenuItem>
                  <MenuItem value="online">Online</MenuItem>
                  <MenuItem value="offline">Offline</MenuItem>
                  <MenuItem value="hybrid">Hybrid</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={6} md={4}>
              <FormControl fullWidth size="small">
                <InputLabel>Status</InputLabel>
                <Select
                  value={statusFilter}
                  label="Status"
                  onChange={(e) => setStatusFilter(e.target.value)}
                >
                  <MenuItem value="all">All Status</MenuItem>
                  <MenuItem value="draft">Draft</MenuItem>
                  <MenuItem value="active">Active</MenuItem>
                  <MenuItem value="paused">Paused</MenuItem>
                  <MenuItem value="completed">Completed</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <Typography variant="body2" color="text.secondary">
                {filteredCampaigns.length} campaign(s) found
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Campaigns Table */}
      <Card>
        <CardContent>
          {isLoading ? (
            <Skeleton height={300} />
          ) : error ? (
            <Box sx={{ py: 8, textAlign: 'center' }}>
              <ErrorOutlineIcon sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
              <Typography variant="h6" color="error" gutterBottom>Failed to load campaigns</Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                {error instanceof Error ? error.message : 'An unexpected error occurred'}
              </Typography>
              <Button variant="contained" onClick={() => refetch()}>Retry</Button>
            </Box>
          ) : (
            <Box sx={{ overflowX: 'auto' }}>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Channel</TableCell>
                      <TableCell>Category</TableCell>
                      <TableCell>Budget</TableCell>
                      <TableCell>Spent</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Approval</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredCampaigns.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={8} align="center">
                          <EmptyState
                            icon={CampaignIcon}
                            title="No campaigns yet"
                            description="Create your first campaign to get started."
                            actionLabel="New Campaign"
                            onAction={() => setOpenDialog(true)}
                          />
                        </TableCell>
                      </TableRow>
                    ) : (
                      filteredCampaigns.map((campaign: any) => (
                        <TableRow key={campaign.id} hover>
                          <TableCell>
                            <Typography fontWeight="medium">{campaign.name}</Typography>
                          </TableCell>
                          <TableCell>
                            <Chip label={campaign.channel} size="small" color="primary" variant="outlined" />
                          </TableCell>
                          <TableCell>
                            <Chip
                              label={campaign.category || getChannelCategory(campaign.channel)}
                              size="small"
                              color={campaign.category === 'online' ? 'success' : campaign.category === 'offline' ? 'warning' : 'info'}
                            />
                          </TableCell>
                          <TableCell>{formatCurrency(campaign.budget || 0)}</TableCell>
                          <TableCell>{formatCurrency(campaign.spent || 0)}</TableCell>
                          <TableCell>
                            <Chip label={campaign.status || 'draft'} size="small" color={getStatusColor(campaign.status || 'draft')} />
                          </TableCell>
                          <TableCell>
                            <Chip
                              label={campaign.approvalStatus || 'N/A'}
                              size="small"
                              color={getApprovalStatusColor(campaign.approvalStatus || 'draft')}
                            />
                          </TableCell>
                          <TableCell>
                            <Box sx={{ display: 'flex', gap: 0.5 }}>
                              {(campaign.status === 'draft' || !campaign.approvalStatus) && (
                                <IconButton size="small" onClick={() => submitMutation.mutate(campaign.id)} title="Submit for Approval" aria-label="Submit for approval">
                                  <SendIcon fontSize="small" />
                                </IconButton>
                              )}
                              {campaign.approvalStatus === 'pending' && (
                                <>
                                  <IconButton size="small" color="success" onClick={() => approveMutation.mutate(campaign.id)} title="Approve" aria-label="Approve campaign">
                                    <CheckIcon fontSize="small" />
                                  </IconButton>
                                  <IconButton size="small" color="error" onClick={() => rejectMutation.mutate(campaign.id)} title="Reject" aria-label="Reject campaign">
                                    <CloseIcon fontSize="small" />
                                  </IconButton>
                                </>
                              )}
                              <IconButton size="small" onClick={() => {
                                setEditingCampaign(campaign)
                                setFormData({
                                  name: campaign.name || '',
                                  channel: campaign.channel || 'Google Ads',
                                  category: campaign.category || 'online',
                                  budget: campaign.budget || 0,
                                  start_date: campaign.start_date || new Date().toISOString().split('T')[0],
                                  end_date: campaign.end_date || '',
                                  description: campaign.description || '',
                                })
                                setOpenDialog(true)
                              }} title="Edit" aria-label="Edit campaign">
                                <EditIcon fontSize="small" />
                              </IconButton>
                              <IconButton size="small" color="error" onClick={() => handleDelete(campaign)} title="Delete" aria-label="Delete campaign">
                                <DeleteIcon fontSize="small" />
                              </IconButton>
                            </Box>
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

      {/* Create/Edit Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingCampaign ? 'Edit Campaign' : 'Create New Campaign'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Campaign Name"
              fullWidth
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
            <FormControl fullWidth>
              <InputLabel>Channel</InputLabel>
              <Select
                value={formData.channel}
                label="Channel"
                onChange={(e) => handleChannelChange(e.target.value)}
              >
                {CHANNELS.map((ch) => (
                  <MenuItem key={ch} value={ch}>{ch}</MenuItem>
                ))}
              </Select>
            </FormControl>
            <TextField
              label="Category"
              fullWidth
              value={formData.category}
              InputProps={{ readOnly: true }}
              helperText="Auto-detected from channel selection"
            />
            <TextField
              label="Budget"
              fullWidth
              type="number"
              value={formData.budget}
              onChange={(e) => setFormData({ ...formData, budget: Number(e.target.value) })}
              required
            />
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField
                  label="Start Date"
                  fullWidth
                  type="date"
                  value={formData.start_date}
                  onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                  InputLabelProps={{ shrink: true }}
                  required
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  label="End Date"
                  fullWidth
                  type="date"
                  value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
            </Grid>
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={3}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handleSubmit}
            disabled={createMutation.isPending || updateMutation.isPending}
          >
            {createMutation.isPending || updateMutation.isPending
              ? (editingCampaign ? 'Updating...' : 'Creating...')
              : (editingCampaign ? 'Update' : 'Create')}
          </Button>
        </DialogActions>
      </Dialog>

      <ConfirmDialog {...confirmDialog.dialogProps} />
    </Box>
  )
}
