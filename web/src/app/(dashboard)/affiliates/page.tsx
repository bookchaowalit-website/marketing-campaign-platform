'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { affiliateApi } from '@/lib/affiliate-api'
import { Box, Card, CardContent, Typography, Grid, Paper, Chip, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, IconButton, Skeleton, Button, Tabs, Tab, LinearProgress, Avatar, Dialog, DialogTitle, DialogContent, DialogActions, TextField, Alert, TableSortLabel, alpha } from '@mui/material'
import { paletteColor } from '@/lib/palette'
import PersonAddIcon from '@mui/icons-material/PersonAdd'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import BlockIcon from '@mui/icons-material/Block'
import EditIcon from '@mui/icons-material/Edit'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import AttachMoneyIcon from '@mui/icons-material/AttachMoney'
import LinkIcon from '@mui/icons-material/Link'
import PeopleIcon from '@mui/icons-material/People'
import { useState } from 'react'
import * as React from 'react'
import { Affiliate, AffiliateStats, RegisterAffiliateRequest } from '@/types/affiliate'
import { TabPanel } from '@/components/TabPanel'
import { EmptyState } from '@/components/EmptyState'
import { formatCurrency } from '@/lib/utils'

export default function AffiliatesPage() {
  const { user } = useAuth()
  const toast = useToast()
  const queryClient = useQueryClient()
  const [tabValue, setTabValue] = useState(0)
  const [registerDialogOpen, setRegisterDialogOpen] = useState(false)
  const [registerForm, setRegisterForm] = useState<RegisterAffiliateRequest>({
    website: '',
    description: '',
  })
  const [affOrderBy, setAffOrderBy] = useState('name')
  const [affOrder, setAffOrder] = useState<'asc' | 'desc'>('asc')

  const handleAffSort = (field: string) => {
    if (affOrderBy === field) {
      setAffOrder(affOrder === 'asc' ? 'desc' : 'asc')
    } else {
      setAffOrderBy(field)
      setAffOrder('asc')
    }
  }

  // Check if user has affiliate account
  const { data: myAffiliate, isLoading: myAffiliateLoading, error: myAffiliateError } = useQuery<Affiliate>({
    queryKey: ['myAffiliate'],
    queryFn: async () => {
      try {
        return await affiliateApi.getMyAffiliate()
      } catch (err: unknown) {
        const error = err as { response?: { status?: number } }
        if (error.response?.status === 404) {
          return null as unknown as Affiliate
        }
        throw err
      }
    },
  })

  // Get affiliate stats
  const { data: stats, isLoading: statsLoading } = useQuery<AffiliateStats>({
    queryKey: ['affiliateStats'],
    queryFn: async () => {
      try {
        return await affiliateApi.getStats()
      } catch {
        return {
          totalAffiliates: 0,
          activeAffiliates: 0,
          totalClicks: 0,
          totalSignups: 0,
          totalConversions: 0,
          conversionRate: 0,
          totalEarnings: 0,
          totalPaid: 0,
          pendingPayouts: 0,
          topAffiliates: [],
        }
      }
    },
  })

  // Get all affiliates
  const { data: affiliatesResponse, isLoading: affiliatesLoading } = useQuery<{ affiliates: Affiliate[]; total: number }>({
    queryKey: ['affiliates'],
    queryFn: async () => {
      try {
        return await affiliateApi.getAllAffiliates(20, 0)
      } catch {
        return { affiliates: [], total: 0 }
      }
    },
  })

  const sortedAffiliates = React.useMemo(() => {
    const affiliates = affiliatesResponse?.affiliates || []
    return [...affiliates].sort((a, b) => {
      const mod = affOrder === 'asc' ? 1 : -1
      if (affOrderBy === 'name') return (a.affiliateCode || '').localeCompare(b.affiliateCode || '') * mod
      if (affOrderBy === 'commission') return ((a.commissionRate || 0) - (b.commissionRate || 0)) * mod
      if (affOrderBy === 'earnings') return ((a.totalEarnings || 0) - (b.totalEarnings || 0)) * mod
      if (affOrderBy === 'status') return (a.status || '').localeCompare(b.status || '') * mod
      return 0
    })
  }, [affiliatesResponse, affOrderBy, affOrder])

  // Register affiliate mutation
  const registerMutation = useMutation({
    mutationFn: (data: RegisterAffiliateRequest) => affiliateApi.register(data),
    onSuccess: () => {
      toast.success('Affiliate account registered successfully!')
      queryClient.invalidateQueries({ queryKey: ['myAffiliate'] })
      setRegisterDialogOpen(false)
      setRegisterForm({ website: '', description: '' })
    },
    onError: () => {
      toast.error('Failed to register affiliate account')
    },
  })

  // Approve affiliate mutation
  const approveMutation = useMutation({
    mutationFn: (id: string) => affiliateApi.approveAffiliate(id),
    onSuccess: () => {
      toast.success('Affiliate approved successfully!')
      queryClient.invalidateQueries({ queryKey: ['affiliates'] })
    },
    onError: () => {
      toast.error('Failed to approve affiliate')
    },
  })

  // Suspend affiliate mutation
  const suspendMutation = useMutation({
    mutationFn: (id: string) => affiliateApi.suspendAffiliate(id),
    onSuccess: () => {
      toast.success('Affiliate suspended successfully!')
      queryClient.invalidateQueries({ queryKey: ['affiliates'] })
    },
    onError: () => {
      toast.error('Failed to suspend affiliate')
    },
  })

  const handleRegister = () => {
    registerMutation.mutate(registerForm)
  }

  const getStatusColor = (status: string): 'success' | 'warning' | 'error' | 'default' => {
    switch (status) {
      case 'active': return 'success'
      case 'pending': return 'warning'
      case 'suspended': return 'error'
      default: return 'default'
    }
  }

  const statsCards = [
    { label: 'Total Affiliates', value: stats?.totalAffiliates || 0, color: 'primary.main', icon: PeopleIcon },
    { label: 'Active Affiliates', value: stats?.activeAffiliates || 0, color: 'success.main', icon: CheckCircleIcon },
    { label: 'Total Earnings', value: formatCurrency(stats?.totalEarnings || 0), color: 'secondary.main', icon: AttachMoneyIcon },
    { label: 'Pending Payouts', value: formatCurrency(stats?.pendingPayouts || 0), color: 'warning.main', icon: TrendingUpIcon },
  ]

  const hasAffiliateAccount = myAffiliate && myAffiliate.id

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Box>
          <Typography variant="h4" fontWeight="bold">
            Affiliate Program
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Manage your affiliate partners and track performance
          </Typography>
        </Box>
        {!hasAffiliateAccount && (
          <Button variant="contained" startIcon={<PersonAddIcon />} onClick={() => setRegisterDialogOpen(true)}>
            Join Program
          </Button>
        )}
      </Box>

      {/* Stats Grid */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {statsCards.map((stat, i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <Card sx={{ height: '100%' }}>
              <CardContent>
                {statsLoading ? (
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
                    </Box>
                    <Box sx={{
                      p: 1,
                      borderRadius: 2,
                      bgcolor: (theme) => alpha(paletteColor(theme.palette, stat.color), 0.125),
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

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={(e, v) => setTabValue(v)}>
          <Tab label="All Affiliates" />
          <Tab label="My Affiliate" />
          <Tab label="Settings" />
        </Tabs>
      </Box>

      {/* All Affiliates Tab */}
      <TabPanel value={tabValue} index={0}>
        <Card>
          <CardContent>
            {affiliatesLoading ? (
              <Skeleton height={300} />
            ) : (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>
                        <TableSortLabel active={affOrderBy === 'name'} direction={affOrderBy === 'name' ? affOrder : 'asc'} onClick={() => handleAffSort('name')}>
                          Affiliate
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Code</TableCell>
                      <TableCell>
                        <TableSortLabel active={affOrderBy === 'status'} direction={affOrderBy === 'status' ? affOrder : 'asc'} onClick={() => handleAffSort('status')}>
                          Status
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>
                        <TableSortLabel active={affOrderBy === 'commission'} direction={affOrderBy === 'commission' ? affOrder : 'asc'} onClick={() => handleAffSort('commission')}>
                          Commission
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Clicks</TableCell>
                      <TableCell>Conversions</TableCell>
                      <TableCell>
                        <TableSortLabel active={affOrderBy === 'earnings'} direction={affOrderBy === 'earnings' ? affOrder : 'asc'} onClick={() => handleAffSort('earnings')}>
                          Earnings
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Balance</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {sortedAffiliates?.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={9} align="center">
                          <EmptyState icon={PeopleIcon} title="No affiliates yet" description="Affiliates will appear here once they join your program." />
                        </TableCell>
                      </TableRow>
                    ) : (
                      sortedAffiliates?.map((affiliate) => (
                        <TableRow key={affiliate.id} hover>
                          <TableCell>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                              <Avatar sx={{ bgcolor: 'primary.main' }}>
                                {affiliate.affiliateCode.charAt(0)}
                              </Avatar>
                              <Box>
                                <Typography fontWeight="medium">{affiliate.affiliateCode}</Typography>
                                <Typography variant="body2" color="text.secondary">
                                  {affiliate.website || 'No website'}
                                </Typography>
                              </Box>
                            </Box>
                          </TableCell>
                          <TableCell>
                            <Chip label={affiliate.affiliateCode} size="small" color="primary" />
                          </TableCell>
                          <TableCell>
                            <Chip label={affiliate.status} size="small" color={getStatusColor(affiliate.status)} />
                          </TableCell>
                          <TableCell>
                            {affiliate.commissionType === 'percentage'
                              ? `${(affiliate.commissionRate * 100).toFixed(1)}%`
                              : `$${affiliate.fixedCommissionAmount}`}
                          </TableCell>
                          <TableCell>{affiliate.totalClicks}</TableCell>
                          <TableCell>{affiliate.totalConversions}</TableCell>
                          <TableCell>{formatCurrency(affiliate.totalEarnings)}</TableCell>
                          <TableCell>{formatCurrency(affiliate.balance)}</TableCell>
                          <TableCell>
                            {affiliate.status === 'pending' && (
                              <IconButton size="small" title="Approve" color="success" onClick={() => approveMutation.mutate(affiliate.id)}>
                                <CheckCircleIcon fontSize="small" />
                              </IconButton>
                            )}
                            {affiliate.status === 'active' && (
                              <IconButton size="small" title="Suspend" color="error" onClick={() => suspendMutation.mutate(affiliate.id)}>
                                <BlockIcon fontSize="small" />
                              </IconButton>
                            )}
                            <IconButton size="small" title="Edit" disabled aria-label="Edit affiliate (coming soon)">
                              <EditIcon fontSize="small" />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* My Affiliate Tab */}
      <TabPanel value={tabValue} index={1}>
        {myAffiliateLoading ? (
          <Skeleton height={300} />
        ) : !hasAffiliateAccount ? (
          <Card>
            <CardContent sx={{ textAlign: 'center', py: 6 }}>
              <PersonAddIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" gutterBottom>
                You haven&apos;t joined the affiliate program yet
              </Typography>
              <Typography color="text.secondary" sx={{ mb: 3 }}>
                Join our affiliate program and earn commissions by referring new customers
              </Typography>
              <Button variant="contained" startIcon={<PersonAddIcon />} onClick={() => setRegisterDialogOpen(true)}>
                Join Affiliate Program
              </Button>
            </CardContent>
          </Card>
        ) : (
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Card>
                <CardContent>
                  <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
                    Your Affiliate Account
                  </Typography>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="body2" color="text.secondary">Affiliate Code</Typography>
                    <Typography variant="h5" fontWeight="bold" color="primary">
                      {myAffiliate?.affiliateCode}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="body2" color="text.secondary">Status</Typography>
                    <Chip label={myAffiliate?.status} color={getStatusColor(myAffiliate?.status || '')} />
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="body2" color="text.secondary">Commission Rate</Typography>
                    <Typography variant="h6">
                      {myAffiliate?.commissionType === 'percentage'
                        ? `${(myAffiliate?.commissionRate * 100).toFixed(1)}%`
                        : `$${myAffiliate?.fixedCommissionAmount}`}
                    </Typography>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={6}>
              <Card>
                <CardContent>
                  <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
                    Performance Stats
                  </Typography>
                  <Grid container spacing={2}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Total Clicks</Typography>
                      <Typography variant="h5">{myAffiliate?.totalClicks || 0}</Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Signups</Typography>
                      <Typography variant="h5">{myAffiliate?.totalSignups || 0}</Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Conversions</Typography>
                      <Typography variant="h5">{myAffiliate?.totalConversions || 0}</Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Earnings</Typography>
                      <Typography variant="h5" color="success.main">{formatCurrency(myAffiliate?.totalEarnings || 0)}</Typography>
                    </Grid>
                    <Grid item xs={12}>
                      <Typography variant="body2" color="text.secondary">Conversion Rate</Typography>
                      <LinearProgress
                        variant="determinate"
                        value={myAffiliate?.conversionRate || 0}
                        sx={{ height: 10, borderRadius: 5, mt: 1 }}
                        color="success"
                      />
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        {(myAffiliate?.conversionRate || 0).toFixed(2)}%
                      </Typography>
                    </Grid>
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                    <Typography variant="h6" fontWeight="bold">
                      Your Affiliate Links
                    </Typography>
                    <Button variant="outlined" startIcon={<LinkIcon />} disabled title="Coming soon">
                      Create Link
                    </Button>
                  </Box>
                  <Alert severity="info" sx={{ mb: 2 }}>
                    Share your affiliate link: <strong>https://bookmarketing.com/ref/{myAffiliate?.affiliateCode}</strong>
                  </Alert>
                  <Typography color="text.secondary">
                    Create custom tracking links for different campaigns and platforms
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        )}
      </TabPanel>

      {/* Settings Tab */}
      <TabPanel value={tabValue} index={2}>
        <Card>
          <CardContent>
            <Typography variant="h6" fontWeight="bold" sx={{ mb: 3 }}>
              Affiliate Program Settings
            </Typography>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <TextField
                  label="Program Name"
                  fullWidth
                  defaultValue="BookMarketing Affiliate Program"
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  label="Default Commission Rate (%)"
                  fullWidth
                  type="number"
                  defaultValue="10"
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  label="Minimum Payout Amount ($)"
                  fullWidth
                  type="number"
                  defaultValue="50"
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  label="Cookie Duration (days)"
                  fullWidth
                  type="number"
                  defaultValue="30"
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  label="Program Description"
                  fullWidth
                  multiline
                  rows={4}
                  defaultValue="Join our affiliate program and earn commissions by referring new customers to our marketing SaaS platform."
                />
              </Grid>
              <Grid item xs={12}>
                <Button variant="contained" disabled title="Coming soon">
                  Save Settings
                </Button>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      </TabPanel>

      {/* Register Dialog */}
      <Dialog open={registerDialogOpen} onClose={() => setRegisterDialogOpen(false)}>
        <DialogTitle>Join Affiliate Program</DialogTitle>
        <DialogContent>
          <Typography color="text.secondary" sx={{ mb: 2 }}>
            Fill in your details to join our affiliate program
          </Typography>
          <TextField
            label="Website (optional)"
            fullWidth
            value={registerForm.website}
            onChange={(e) => setRegisterForm({ ...registerForm, website: e.target.value })}
            sx={{ mb: 2 }}
          />
          <TextField
            label="Description (optional)"
            fullWidth
            multiline
            rows={3}
            value={registerForm.description}
            onChange={(e) => setRegisterForm({ ...registerForm, description: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRegisterDialogOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleRegister} disabled={registerMutation.isPending}>
            {registerMutation.isPending ? 'Registering...' : 'Register'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
