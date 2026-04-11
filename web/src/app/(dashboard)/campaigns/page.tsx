'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useMemo } from 'react'
import { Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Chip, LinearProgress, Button, Dialog, DialogTitle, DialogContent, DialogActions, TextField, MenuItem, Skeleton, IconButton, InputAdornment, FormControl, InputLabel, Select, Grid, Pagination, TableSortLabel } from '@mui/material'
import PlayArrowIcon from '@mui/icons-material/PlayArrow'
import StopIcon from '@mui/icons-material/Stop'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import AddIcon from '@mui/icons-material/Add'
import SearchIcon from '@mui/icons-material/Search'
import FilterListIcon from '@mui/icons-material/FilterList'
import CampaignIcon from '@mui/icons-material/Campaign'
import { campaignApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { ConfirmDialog, useConfirmDialog } from '@/components/ConfirmDialog'
import { EmptyState } from '@/components/EmptyState'
import { formatCurrency, formatDate, formatNumber } from '@/lib/utils'

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
  start_date: string
}

const defaultCampaigns: Campaign[] = [
  { id: '10000000-0000-0000-0000-000000000001', name: 'Summer Sale 2024', type: 'SEM', status: 'active', budget: 10000, spent: 6500, impressions: 125000, clicks: 4500, conversions: 180, ctr: 3.6, conversion_rate: 4.0, budget_utilization: 65, start_date: '2024-06-01' },
  { id: '10000000-0000-0000-0000-000000000002', name: 'SEO Organic Growth', type: 'SEO', status: 'active', budget: 5000, spent: 3200, impressions: 85000, clicks: 3200, conversions: 95, ctr: 3.8, conversion_rate: 3.0, budget_utilization: 64, start_date: '2024-05-01' },
  { id: '10000000-0000-0000-0000-000000000003', name: 'Facebook Brand Awareness', type: 'Social', status: 'paused', budget: 8000, spent: 4100, impressions: 250000, clicks: 8900, conversions: 220, ctr: 3.6, conversion_rate: 2.5, budget_utilization: 51, start_date: '2024-04-15' },
  { id: '10000000-0000-0000-0000-000000000004', name: 'Email Retargeting Q2', type: 'Email', status: 'completed', budget: 3000, spent: 2850, impressions: 45000, clicks: 2100, conversions: 310, ctr: 4.7, conversion_rate: 14.8, budget_utilization: 95, start_date: '2024-04-01' },
  { id: '10000000-0000-0000-0000-000000000005', name: 'Google Display Ads', type: 'Display', status: 'active', budget: 6000, spent: 4200, impressions: 180000, clicks: 5400, conversions: 145, ctr: 3.0, conversion_rate: 2.7, budget_utilization: 70, start_date: '2024-05-15' },
]

export default function CampaignsPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const confirmDialog = useConfirmDialog()

  const [openDialog, setOpenDialog] = useState(false)
  const [editingCampaign, setEditingCampaign] = useState<Campaign | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    type: 'SEM',
    budget: 0,
    start_date: new Date().toISOString().split('T')[0],
  })

  // Search and filter state
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [typeFilter, setTypeFilter] = useState<string>('all')

  // Pagination state
  const [page, setPage] = useState(0)
  const pageSize = 20

  const { data: campaignsResponse, isLoading } = useQuery<{ campaigns: Campaign[] }>({
    queryKey: ['campaigns', page, pageSize],
    queryFn: async () => {
      return campaignApi.getAll(page * pageSize, pageSize)
    },
    initialData: { campaigns: defaultCampaigns },
  })

  const campaigns = campaignsResponse?.campaigns || defaultCampaigns

  // Filter campaigns based on search and filters
  const filteredCampaigns = useMemo(() => {
    return campaigns.filter((campaign) => {
      const matchesSearch = campaign.name.toLowerCase().includes(searchQuery.toLowerCase())
      const matchesStatus = statusFilter === 'all' || campaign.status === statusFilter
      const matchesType = typeFilter === 'all' || campaign.type === typeFilter
      return matchesSearch && matchesStatus && matchesType
    })
  }, [campaigns, searchQuery, statusFilter, typeFilter])

  // Sorting state
  const [orderBy, setOrderBy] = useState<'name' | 'budget' | 'spent' | 'conversions'>('name')
  const [order, setOrder] = useState<'asc' | 'desc'>('asc')

  const handleSort = (field: typeof orderBy) => {
    if (orderBy === field) {
      setOrder(order === 'asc' ? 'desc' : 'asc')
    } else {
      setOrderBy(field)
      setOrder('asc')
    }
  }

  const sortedCampaigns = useMemo(() => {
    return [...filteredCampaigns].sort((a, b) => {
      const mod = order === 'asc' ? 1 : -1
      if (orderBy === 'name') return a.name.localeCompare(b.name) * mod
      return ((a[orderBy] as number) - (b[orderBy] as number)) * mod
    })
  }, [filteredCampaigns, orderBy, order])

  const pauseMutation = useMutation({
    mutationFn: async (id: string) => {
      await campaignApi.pause(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
      toast.success('Campaign paused successfully')
    },
    onError: () => {
      toast.error('Failed to pause campaign')
    },
  })

  const activateMutation = useMutation({
    mutationFn: async (id: string) => {
      await campaignApi.activate(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
      toast.success('Campaign activated successfully')
    },
    onError: () => {
      toast.error('Failed to activate campaign')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await campaignApi.delete(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
      toast.success('Campaign deleted successfully')
    },
    onError: () => {
      toast.error('Failed to delete campaign')
    },
  })

  const createMutation = useMutation({
    mutationFn: async (data: typeof formData) => {
      await campaignApi.create({
        ...data,
        budget: Number(data.budget),
      } as any)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
      setOpenDialog(false)
      setFormData({ name: '', type: 'SEM', budget: 0, start_date: new Date().toISOString().split('T')[0] })
      setEditingCampaign(null)
      toast.success('Campaign created successfully')
    },
    onError: () => {
      toast.error('Failed to create campaign')
    },
  })

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: typeof formData }) => {
      await campaignApi.update(id, {
        ...data,
        budget: Number(data.budget),
      } as any)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
      setOpenDialog(false)
      setFormData({ name: '', type: 'SEM', budget: 0, start_date: new Date().toISOString().split('T')[0] })
      setEditingCampaign(null)
      toast.success('Campaign updated successfully')
    },
    onError: () => {
      toast.error('Failed to update campaign')
    },
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

  const handleOpenDialog = () => {
    setEditingCampaign(null)
    setFormData({ name: '', type: 'SEM', budget: 0, start_date: new Date().toISOString().split('T')[0] })
    setOpenDialog(true)
  }

  const handleSubmit = () => {
    if (editingCampaign) {
      updateMutation.mutate({ id: editingCampaign.id, data: formData })
    } else {
      createMutation.mutate(formData)
    }
  }

  const handleDelete = (campaign: Campaign) => {
    confirmDialog.showConfirm(
      'Delete Campaign',
      `Are you sure you want to delete "${campaign.name}"? This action cannot be undone.`,
      () => deleteMutation.mutate(campaign.id)
    )
  }

  const handleToggleStatus = (campaign: Campaign) => {
    if (campaign.status === 'active') {
      confirmDialog.showConfirm(
        'Pause Campaign',
        `Are you sure you want to pause "${campaign.name}"?`,
        () => pauseMutation.mutate(campaign.id)
      )
    } else {
      activateMutation.mutate(campaign.id)
    }
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Campaign Management
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleOpenDialog}
        >
          Create Campaign
        </Button>
      </Box>

      {/* Search and Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} md={4}>
              <TextField
                fullWidth
                placeholder="Search campaigns..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  ),
                }}
                size="small"
              />
            </Grid>
            <Grid item xs={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Status</InputLabel>
                <Select
                  value={statusFilter}
                  label="Status"
                  onChange={(e) => setStatusFilter(e.target.value)}
                >
                  <MenuItem value="all">All Status</MenuItem>
                  <MenuItem value="active">Active</MenuItem>
                  <MenuItem value="paused">Paused</MenuItem>
                  <MenuItem value="completed">Completed</MenuItem>
                  <MenuItem value="draft">Draft</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Type</InputLabel>
                <Select
                  value={typeFilter}
                  label="Type"
                  onChange={(e) => setTypeFilter(e.target.value)}
                >
                  <MenuItem value="all">All Types</MenuItem>
                  <MenuItem value="SEM">SEM</MenuItem>
                  <MenuItem value="SEO">SEO</MenuItem>
                  <MenuItem value="Social">Social</MenuItem>
                  <MenuItem value="Email">Email</MenuItem>
                  <MenuItem value="Display">Display</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={2}>
              <Typography variant="body2" color="text.secondary">
                {filteredCampaigns.length} of {campaigns.length} campaigns
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
          ) : (
            <Box sx={{ overflowX: 'auto' }}>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>
                        <TableSortLabel active={orderBy === 'name'} direction={orderBy === 'name' ? order : 'asc'} onClick={() => handleSort('name')}>
                          Campaign
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>
                        <TableSortLabel active={orderBy === 'budget'} direction={orderBy === 'budget' ? order : 'asc'} onClick={() => handleSort('budget')}>
                          Budget
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>
                        <TableSortLabel active={orderBy === 'spent'} direction={orderBy === 'spent' ? order : 'asc'} onClick={() => handleSort('spent')}>
                          Spent
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Impressions</TableCell>
                      <TableCell>Clicks</TableCell>
                      <TableCell>CTR</TableCell>
                      <TableCell>
                        <TableSortLabel active={orderBy === 'conversions'} direction={orderBy === 'conversions' ? order : 'asc'} onClick={() => handleSort('conversions')}>
                          Conversions
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Start Date</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredCampaigns.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={11} align="center">
                          <EmptyState icon={CampaignIcon} title="No campaigns found" description="Create your first campaign to start tracking marketing performance." actionLabel="Create Campaign" onAction={handleOpenDialog} />
                        </TableCell>
                      </TableRow>
                    ) : (
                      sortedCampaigns.map((campaign) => (
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
                          <TableCell>{formatNumber(campaign.impressions)}</TableCell>
                          <TableCell>{formatNumber(campaign.clicks)}</TableCell>
                          <TableCell>{campaign.ctr.toFixed(2)}%</TableCell>
                          <TableCell>{campaign.conversions}</TableCell>
                          <TableCell>{formatDate(campaign.start_date)}</TableCell>
                          <TableCell>
                            <IconButton size="small" onClick={() => {
                              setEditingCampaign(campaign)
                              setFormData({
                                name: campaign.name,
                                type: campaign.type,
                                budget: campaign.budget,
                                start_date: campaign.start_date,
                              })
                              setOpenDialog(true)
                            }} title="Edit" aria-label="Edit campaign">
                              <EditIcon fontSize="small" />
                            </IconButton>
                            <IconButton
                              size="small"
                              onClick={() => handleToggleStatus(campaign)}
                              title={campaign.status === 'active' ? 'Pause' : 'Activate'}
                              aria-label={campaign.status === 'active' ? 'Pause campaign' : 'Activate campaign'}
                            >
                              {campaign.status === 'active' ? <StopIcon fontSize="small" /> : <PlayArrowIcon fontSize="small" />}
                            </IconButton>
                            <IconButton
                              size="small"
                              color="error"
                              onClick={() => handleDelete(campaign)}
                              title="Delete"
                              aria-label="Delete campaign"
                            >
                              <DeleteIcon fontSize="small" />
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

          {/* Pagination */}
          {!isLoading && filteredCampaigns.length > 0 && (
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', pt: 2 }}>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                Showing {page * pageSize + 1}-{Math.min((page + 1) * pageSize, filteredCampaigns.length)} of {filteredCampaigns.length} results
              </Typography>
              <Pagination
                count={Math.ceil(filteredCampaigns.length / pageSize)}
                page={page + 1}
                onChange={(_, newPage) => setPage(newPage - 1)}
                color="primary"
              />
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
            <TextField
              label="Campaign Type"
              fullWidth
              select
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value })}
            >
              <MenuItem value="SEM">SEM</MenuItem>
              <MenuItem value="SEO">SEO</MenuItem>
              <MenuItem value="Social">Social</MenuItem>
              <MenuItem value="Email">Email</MenuItem>
              <MenuItem value="Display">Display</MenuItem>
            </TextField>
            <TextField
              label="Budget"
              fullWidth
              type="number"
              value={formData.budget}
              onChange={(e) => setFormData({ ...formData, budget: Number(e.target.value) })}
              required
            />
            <TextField
              label="Start Date"
              fullWidth
              type="date"
              value={formData.start_date}
              onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
              required
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

      {/* Confirmation Dialog */}
      <ConfirmDialog {...confirmDialog.dialogProps} />
    </Box>
  )
}
