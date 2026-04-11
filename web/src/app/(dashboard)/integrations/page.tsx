'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { externalAffiliateApi } from '@/lib/external-affiliate-api'
import { Box, Card, CardContent, Typography, Grid, Chip, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, IconButton, Skeleton, Button, Tabs, Tab, LinearProgress, Avatar, Dialog, DialogTitle, DialogContent, DialogActions, TextField, Alert, Menu, MenuItem, ListItemIcon, ListItemText, Divider, TableSortLabel } from '@mui/material'
import RefreshIcon from '@mui/icons-material/Refresh'
import LinkIcon from '@mui/icons-material/Link'
import AddIcon from '@mui/icons-material/Add'
import DeleteIcon from '@mui/icons-material/Delete'
import SyncIcon from '@mui/icons-material/Sync'
import StorefrontIcon from '@mui/icons-material/Storefront'
import MusicNoteIcon from '@mui/icons-material/MusicNote'
import ShoppingBagIcon from '@mui/icons-material/ShoppingBag'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import AttachMoneyIcon from '@mui/icons-material/AttachMoney'
import MoreVertIcon from '@mui/icons-material/MoreVert'
import ContentCopyIcon from '@mui/icons-material/ContentCopy'
import OpenInNewIcon from '@mui/icons-material/OpenInNew'
import { useState } from 'react'
import * as React from 'react'
import { ExternalAccount, Platform, ConnectPlatformRequest } from '@/types/external-affiliate'
import { TabPanel } from '@/components/TabPanel'
import { EmptyState } from '@/components/EmptyState'
import { formatCurrency, formatNumber } from '@/lib/utils'

const platformIcons: Record<Platform, typeof StorefrontIcon> = {
  shopee: StorefrontIcon,
  tiktok: MusicNoteIcon,
  lazada: ShoppingBagIcon,
  amazon: ShoppingBagIcon,
  custom: LinkIcon,
}

const platformColors: Record<Platform, string> = {
  shopee: 'error.main',
  tiktok: 'text.primary',
  lazada: 'primary.dark',
  amazon: 'warning.main',
  custom: 'text.secondary',
}

export default function IntegrationsPage() {
  const { user } = useAuth()
  const toast = useToast()
  const queryClient = useQueryClient()
  const [tabValue, setTabValue] = useState(0)
  const [connectDialogOpen, setConnectDialogOpen] = useState(false)
  const [selectedPlatform, setSelectedPlatform] = useState<Platform>('shopee')
  const [connectForm, setConnectForm] = useState<ConnectPlatformRequest>({
    platform: 'shopee',
  })
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [selectedAccount, setSelectedAccount] = useState<ExternalAccount | null>(null)
  const [linkDialogOpen, setLinkDialogOpen] = useState(false)
  const [productUrl, setProductUrl] = useState('')
  const [generatedLink, setGeneratedLink] = useState('')

  // Get connected accounts
  const { data: accountsData, isLoading: accountsLoading } = useQuery({
    queryKey: ['externalAccounts'],
    queryFn: async () => {
      try {
        return await externalAffiliateApi.getAccounts()
      } catch {
        return { accounts: [], total: 0 }
      }
    },
  })

  // Get unified dashboard
  const { data: dashboard, isLoading: dashboardLoading } = useQuery({
    queryKey: ['externalDashboard'],
    queryFn: async () => {
      try {
        return await externalAffiliateApi.getDashboard()
      } catch {
        return null
      }
    },
  })

  // Connect platform mutation
  const connectMutation = useMutation({
    mutationFn: (data: ConnectPlatformRequest) => externalAffiliateApi.connect(data),
    onSuccess: () => {
      toast.success('Platform connected successfully!')
      queryClient.invalidateQueries({ queryKey: ['externalAccounts'] })
      queryClient.invalidateQueries({ queryKey: ['externalDashboard'] })
      setConnectDialogOpen(false)
      setConnectForm({ platform: selectedPlatform })
    },
    onError: () => {
      toast.error('Failed to connect platform')
    },
  })

  // Disconnect mutation
  const disconnectMutation = useMutation({
    mutationFn: (id: string) => externalAffiliateApi.disconnect(id),
    onSuccess: () => {
      toast.success('Account disconnected')
      queryClient.invalidateQueries({ queryKey: ['externalAccounts'] })
      queryClient.invalidateQueries({ queryKey: ['externalDashboard'] })
    },
    onError: () => {
      toast.error('Failed to disconnect account')
    },
  })

  // Sync mutation
  const syncMutation = useMutation({
    mutationFn: (id: string) => externalAffiliateApi.sync(id),
    onSuccess: () => {
      toast.success('Sync started successfully!')
      queryClient.invalidateQueries({ queryKey: ['externalAccounts'] })
    },
    onError: () => {
      toast.error('Failed to sync account')
    },
  })

  // Generate link mutation
  const generateLinkMutation = useMutation({
    mutationFn: ({ accountId, url }: { accountId: string; url: string }) =>
      externalAffiliateApi.generateLink(accountId, { productUrl: url }),
    onSuccess: (data) => {
      setGeneratedLink(data.affiliate_link)
      toast.success('Affiliate link generated!')
    },
    onError: () => {
      toast.error('Failed to generate link')
    },
  })

  const handleConnect = () => {
    connectMutation.mutate(connectForm)
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, account: ExternalAccount) => {
    setAnchorEl(event.currentTarget)
    setSelectedAccount(account)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedAccount(null)
  }

  const handleSync = () => {
    if (selectedAccount) {
      syncMutation.mutate(selectedAccount.id)
    }
    handleMenuClose()
  }

  const handleDisconnect = () => {
    if (selectedAccount) {
      disconnectMutation.mutate(selectedAccount.id)
    }
    handleMenuClose()
  }

  const handleGenerateLink = () => {
    if (selectedAccount && productUrl) {
      generateLinkMutation.mutate({ accountId: selectedAccount.id, url: productUrl })
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success('Copied to clipboard!')
  }

  const accounts = accountsData?.accounts || []

  // Orders table sorting
  const [orderOrderBy, setOrderOrderBy] = useState('platform')
  const [orderOrder, setOrderOrder] = useState<'asc' | 'desc'>('asc')

  const handleOrderSort = (field: string) => {
    if (orderOrderBy === field) {
      setOrderOrder(orderOrder === 'asc' ? 'desc' : 'asc')
    } else {
      setOrderOrderBy(field)
      setOrderOrder('asc')
    }
  }

  const sortedOrders = React.useMemo(() => {
    const orders = dashboard?.recentOrders || []
    return [...orders].sort((a, b) => {
      const mod = orderOrder === 'asc' ? 1 : -1
      if (orderOrderBy === 'platform') return (a.platform || '').localeCompare(b.platform || '') * mod
      if (orderOrderBy === 'status') return (a.orderStatus || '').localeCompare(b.orderStatus || '') * mod
      if (orderOrderBy === 'amount') return ((a.orderAmount || 0) - (b.orderAmount || 0)) * mod
      if (orderOrderBy === 'date') return new Date(a.orderDate || 0).getTime() - new Date(b.orderDate || 0).getTime() * mod
      return 0
    })
  }, [dashboard, orderOrderBy, orderOrder])

  const statsCards = [
    { label: 'Total Earnings', value: formatCurrency(dashboard?.totalEarnings || 0), color: 'success.main', icon: AttachMoneyIcon },
    { label: 'Pending Commission', value: formatCurrency(dashboard?.totalPendingCommission || 0), color: 'warning.main', icon: TrendingUpIcon },
    { label: 'Total Orders', value: dashboard?.totalOrders || 0, color: 'primary.main', icon: ShoppingBagIcon },
    { label: 'Connected Platforms', value: accounts.length, color: 'secondary.main', icon: LinkIcon },
  ]

  const getStatusColor = (status: string): 'success' | 'warning' | 'error' | 'default' => {
    switch (status) {
      case 'active': return 'success'
      case 'pending': return 'warning'
      case 'error': return 'error'
      default: return 'default'
    }
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Box>
          <Typography variant="h4" fontWeight="bold">
            Affiliate Integrations
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Connect and manage your external affiliate accounts
          </Typography>
        </Box>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => setConnectDialogOpen(true)}>
          Connect Platform
        </Button>
      </Box>

      {/* Stats Grid */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {statsCards.map((stat, i) => (
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
                    </Box>
                    <Box sx={{
                      p: 1,
                      borderRadius: 2,
                      bgcolor: `${stat.color}20`,
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
          <Tab label="Connected Accounts" />
          <Tab label="All Orders" />
          <Tab label="Top Products" />
        </Tabs>
      </Box>

      {/* Connected Accounts Tab */}
      <TabPanel value={tabValue} index={0}>
        {accountsLoading ? (
          <Skeleton height={300} />
        ) : accounts.length === 0 ? (
          <Card>
            <CardContent sx={{ textAlign: 'center', py: 6 }}>
              <LinkIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" gutterBottom>
                No platforms connected yet
              </Typography>
              <Typography color="text.secondary" sx={{ mb: 3 }}>
                Connect your Shopee, TikTok, or other affiliate accounts to get started
              </Typography>
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => setConnectDialogOpen(true)}>
                Connect Platform
              </Button>
            </CardContent>
          </Card>
        ) : (
          <Grid container spacing={3}>
            {accounts.map((account) => {
              const PlatformIcon = platformIcons[account.platform]
              return (
                <Grid item xs={12} md={6} lg={4} key={account.id}>
                  <Card>
                    <CardContent>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Avatar sx={{ bgcolor: platformColors[account.platform] }}>
                            <PlatformIcon />
                          </Avatar>
                          <Box>
                            <Typography fontWeight="medium" textTransform="capitalize">
                              {account.platform}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              {account.platformUsername || account.platformAccountName || 'Connected'}
                            </Typography>
                          </Box>
                        </Box>
                        <IconButton size="small" onClick={(e) => handleMenuOpen(e, account)}>
                          <MoreVertIcon />
                        </IconButton>
                      </Box>

                      <Chip
                        label={account.status}
                        size="small"
                        color={getStatusColor(account.status)}
                        sx={{ mb: 2 }}
                      />

                      <Grid container spacing={2}>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="text.secondary">Earnings</Typography>
                          <Typography variant="h6" color="success.main">
                            {formatCurrency(account.totalEarnings)}
                          </Typography>
                        </Grid>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="text.secondary">Orders</Typography>
                          <Typography variant="h6">{account.totalOrders}</Typography>
                        </Grid>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="text.secondary">Clicks</Typography>
                          <Typography variant="h6">{formatNumber(account.totalClicks)}</Typography>
                        </Grid>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="text.secondary">Pending</Typography>
                          <Typography variant="h6" color="warning.main">
                            {formatCurrency(account.pendingCommission)}
                          </Typography>
                        </Grid>
                      </Grid>

                      {account.lastSyncAt && (
                        <Typography variant="caption" color="text.secondary" sx={{ mt: 2, display: 'block' }}>
                          Last synced: {new Date(account.lastSyncAt).toLocaleString()}
                        </Typography>
                      )}
                    </CardContent>
                  </Card>
                </Grid>
              )
            })}
          </Grid>
        )}
      </TabPanel>

      {/* All Orders Tab */}
      <TabPanel value={tabValue} index={1}>
        <Card>
          <CardContent>
            <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
              Recent Orders
            </Typography>
            {dashboard?.recentOrders && dashboard.recentOrders.length > 0 ? (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>
                        <TableSortLabel active={orderOrderBy === 'platform'} direction={orderOrderBy === 'platform' ? orderOrder : 'asc'} onClick={() => handleOrderSort('platform')}>
                          Platform
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Order ID</TableCell>
                      <TableCell>
                        <TableSortLabel active={orderOrderBy === 'status'} direction={orderOrderBy === 'status' ? orderOrder : 'asc'} onClick={() => handleOrderSort('status')}>
                          Status
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>
                        <TableSortLabel active={orderOrderBy === 'amount'} direction={orderOrderBy === 'amount' ? orderOrder : 'asc'} onClick={() => handleOrderSort('amount')}>
                          Amount
                        </TableSortLabel>
                      </TableCell>
                      <TableCell>Commission</TableCell>
                      <TableCell>
                        <TableSortLabel active={orderOrderBy === 'date'} direction={orderOrderBy === 'date' ? orderOrder : 'asc'} onClick={() => handleOrderSort('date')}>
                          Date
                        </TableSortLabel>
                      </TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {sortedOrders.map((order) => (
                      <TableRow key={order.id} hover>
                        <TableCell>
                          <Chip label={order.platform} size="small" />
                        </TableCell>
                        <TableCell>{order.orderId}</TableCell>
                        <TableCell>
                          <Chip label={order.orderStatus} size="small" color={order.orderStatus === 'delivered' ? 'success' : 'default'} />
                        </TableCell>
                        <TableCell>{formatCurrency(order.orderAmount)}</TableCell>
                        <TableCell>{formatCurrency(order.commissionAmount)}</TableCell>
                        <TableCell>{order.orderDate ? new Date(order.orderDate).toLocaleDateString() : '-'}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            ) : (
              <EmptyState icon={ShoppingBagIcon} title="No orders yet" description="Orders from connected platforms will appear here." />
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Top Products Tab */}
      <TabPanel value={tabValue} index={2}>
        <Card>
          <CardContent>
            <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
              Top Performing Products
            </Typography>
            {dashboard?.topProducts && dashboard.topProducts.length > 0 ? (
              <Grid container spacing={2}>
                {dashboard.topProducts.map((product) => (
                  <Grid item xs={12} sm={6} md={4} key={product.id}>
                    <Card variant="outlined">
                      <CardContent>
                        <Typography variant="body2" color="text.secondary">{product.platform}</Typography>
                        <Typography fontWeight="medium" noWrap>{product.productName}</Typography>
                        <Box sx={{ mt: 1 }}>
                          <Typography variant="body2">Earnings: {formatCurrency(product.totalEarnings)}</Typography>
                          <Typography variant="body2">Orders: {product.totalOrders}</Typography>
                        </Box>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            ) : (
              <EmptyState icon={TrendingUpIcon} title="No products yet" description="Top products from your affiliate networks will appear here." />
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Account Menu */}
      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
        <MenuItem onClick={handleSync}>
          <ListItemIcon><SyncIcon fontSize="small" /></ListItemIcon>
          <ListItemText>Sync Now</ListItemText>
        </MenuItem>
        <MenuItem onClick={() => { setLinkDialogOpen(true); handleMenuClose() }}>
          <ListItemIcon><LinkIcon fontSize="small" /></ListItemIcon>
          <ListItemText>Generate Link</ListItemText>
        </MenuItem>
        <Divider />
        <MenuItem onClick={handleDisconnect} sx={{ color: 'error.main' }}>
          <ListItemIcon><DeleteIcon fontSize="small" color="error" /></ListItemIcon>
          <ListItemText>Disconnect</ListItemText>
        </MenuItem>
      </Menu>

      {/* Connect Dialog */}
      <Dialog open={connectDialogOpen} onClose={() => setConnectDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Connect Platform</DialogTitle>
        <DialogContent>
          <Typography color="text.secondary" sx={{ mb: 2 }}>
            Select a platform and enter your credentials
          </Typography>

          <Typography variant="subtitle2" sx={{ mb: 1 }}>Platform</Typography>
          <Grid container spacing={1} sx={{ mb: 3 }}>
            {(['shopee', 'tiktok'] as Platform[]).map((platform) => {
              const Icon = platformIcons[platform]
              return (
                <Grid item key={platform}>
                  <Button
                    variant={selectedPlatform === platform ? 'contained' : 'outlined'}
                    startIcon={<Icon />}
                    onClick={() => {
                      setSelectedPlatform(platform)
                      setConnectForm({ ...connectForm, platform })
                    }}
                    sx={{ textTransform: 'capitalize' }}
                  >
                    {platform}
                  </Button>
                </Grid>
              )
            })}
          </Grid>

          <Alert severity="info" sx={{ mb: 2 }}>
            You&apos;ll need API credentials from your {selectedPlatform} affiliate dashboard
          </Alert>

          <TextField
            label="API Key / Partner ID"
            fullWidth
            value={connectForm.apiKey || ''}
            onChange={(e) => setConnectForm({ ...connectForm, apiKey: e.target.value })}
            sx={{ mb: 2 }}
          />
          <TextField
            label="API Secret / Partner Key"
            fullWidth
            type="password"
            value={connectForm.apiSecret || ''}
            onChange={(e) => setConnectForm({ ...connectForm, apiSecret: e.target.value })}
            sx={{ mb: 2 }}
          />
          <TextField
            label="Account ID (optional)"
            fullWidth
            value={connectForm.platformAccountId || ''}
            onChange={(e) => setConnectForm({ ...connectForm, platformAccountId: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConnectDialogOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleConnect} disabled={connectMutation.isPending}>
            {connectMutation.isPending ? 'Connecting...' : 'Connect'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Generate Link Dialog */}
      <Dialog open={linkDialogOpen} onClose={() => setLinkDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Generate Affiliate Link</DialogTitle>
        <DialogContent>
          <Typography color="text.secondary" sx={{ mb: 2 }}>
            Enter a product URL to generate your affiliate link
          </Typography>
          <TextField
            label="Product URL"
            fullWidth
            value={productUrl}
            onChange={(e) => setProductUrl(e.target.value)}
            placeholder="https://shopee.co.th/product/..."
            sx={{ mb: 2 }}
          />
          {generatedLink && (
            <Alert severity="success" sx={{ mb: 2 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Typography variant="body2" sx={{ wordBreak: 'break-all', flex: 1 }}>
                  {generatedLink}
                </Typography>
                <IconButton size="small" onClick={() => copyToClipboard(generatedLink)}>
                  <ContentCopyIcon fontSize="small" />
                </IconButton>
                <IconButton size="small" onClick={() => window.open(generatedLink, '_blank')}>
                  <OpenInNewIcon fontSize="small" />
                </IconButton>
              </Box>
            </Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => { setLinkDialogOpen(false); setGeneratedLink(''); setProductUrl('') }}>
            Close
          </Button>
          <Button
            variant="contained"
            onClick={handleGenerateLink}
            disabled={!productUrl || generateLinkMutation.isPending}
          >
            Generate
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
