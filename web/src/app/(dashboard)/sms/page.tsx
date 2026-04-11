'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, Tabs, Tab, Skeleton, IconButton, Pagination, TableSortLabel
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import SendIcon from '@mui/icons-material/Send'
import DeleteIcon from '@mui/icons-material/Delete'
import SmsIcon from '@mui/icons-material/Sms'
import GroupIcon from '@mui/icons-material/Group'
import { smsApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { TabPanel } from '@/components/TabPanel'
import { EmptyState } from '@/components/EmptyState'
import { formatCurrency } from '@/lib/utils'

const defaultLists = [
  { id: '1', name: 'VIP Customers', subscriber_count: 250, status: 'active' },
  { id: '2', name: 'Newsletter Subscribers', subscriber_count: 1200, status: 'active' },
  { id: '3', name: 'Promotional', subscriber_count: 800, status: 'active' },
]

const defaultSubscribers = [
  { id: '1', phone: '+66812345678', first_name: 'Somchai', last_name: 'Jaidee', status: 'active', subscribed_at: '2024-03-15' },
  { id: '2', phone: '+66898765432', first_name: 'Suda', last_name: 'Rakdee', status: 'active', subscribed_at: '2024-03-10' },
  { id: '3', phone: '+66811223344', first_name: 'Prasert', last_name: 'Manket', status: 'unsubscribed', subscribed_at: '2024-02-20' },
]

const defaultCampaigns = [
  { id: '1', name: 'Flash Sale March', list_name: 'VIP Customers', message: 'Hi {{first_name}}! Flash sale 50% off today only!', status: 'sent', sent: 250, delivered: 245, failed: 5, cost: 12.50 },
  { id: '2', name: 'Welcome Message', list_name: 'Newsletter Subscribers', message: 'Welcome {{first_name}}! Thanks for subscribing.', status: 'sent', sent: 1200, delivered: 1180, failed: 20, cost: 60.00 },
  { id: '3', name: 'Promo April', list_name: 'Promotional', message: 'Special offer for you, {{first_name}}!', status: 'draft', sent: 0, delivered: 0, failed: 0, cost: 0 },
]

export default function SMSPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const [tabValue, setTabValue] = useState(0)
  const [subPage, setSubPage] = useState(0)
  const subPageSize = 10
  const [orderBy, setOrderBy] = useState('name')
  const [order, setOrder] = useState<'asc' | 'desc'>('asc')

  // Dialog states
  const [openListDialog, setOpenListDialog] = useState(false)
  const [openSubscriberDialog, setOpenSubscriberDialog] = useState(false)
  const [openCampaignDialog, setOpenCampaignDialog] = useState(false)

  const [listForm, setListForm] = useState({ name: '', description: '' })
  const [subscriberForm, setSubscriberForm] = useState({ phone: '', first_name: '', last_name: '', list_id: '' })
  const [campaignForm, setCampaignForm] = useState({ name: '', list_id: '', message: '' })

  const { data: listsData, isLoading: listsLoading } = useQuery({
    queryKey: ['sms-lists'],
    queryFn: async () => {
      try { return await smsApi.getLists() } catch { return { lists: defaultLists } }
    },
  })

  const { data: subscribersData, isLoading: subscribersLoading } = useQuery({
    queryKey: ['sms-subscribers'],
    queryFn: async () => {
      try { return await smsApi.getSubscribers() } catch { return { subscribers: defaultSubscribers } }
    },
  })

  const { data: campaignsData, isLoading: campaignsLoading } = useQuery({
    queryKey: ['sms-campaigns'],
    queryFn: async () => {
      try { return await smsApi.getCampaigns() } catch { return { campaigns: defaultCampaigns } }
    },
  })

  const lists = listsData?.lists || defaultLists
  const subscribers = subscribersData?.subscribers || defaultSubscribers
  const campaigns = campaignsData?.campaigns || defaultCampaigns

  const createListMutation = useMutation({
    mutationFn: async (data: typeof listForm) => smsApi.createList(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sms-lists'] })
      setOpenListDialog(false)
      setListForm({ name: '', description: '' })
      toast.success('List created successfully')
    },
    onError: () => toast.error('Failed to create list'),
  })

  const addSubscriberMutation = useMutation({
    mutationFn: async (data: typeof subscriberForm) => smsApi.addSubscriber(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sms-subscribers'] })
      queryClient.invalidateQueries({ queryKey: ['sms-lists'] })
      setOpenSubscriberDialog(false)
      setSubscriberForm({ phone: '', first_name: '', last_name: '', list_id: '' })
      toast.success('Subscriber added successfully')
    },
    onError: () => toast.error('Failed to add subscriber'),
  })

  const createCampaignMutation = useMutation({
    mutationFn: async (data: typeof campaignForm) => smsApi.createCampaign(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sms-campaigns'] })
      setOpenCampaignDialog(false)
      setCampaignForm({ name: '', list_id: '', message: '' })
      toast.success('Campaign created successfully')
    },
    onError: () => toast.error('Failed to create campaign'),
  })

  const sendCampaignMutation = useMutation({
    mutationFn: async (id: string) => smsApi.sendCampaign(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sms-campaigns'] })
      toast.success('Campaign sent successfully')
    },
    onError: () => toast.error('Failed to send campaign'),
  })

  const getStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' | 'error' => {
    switch (status) {
      case 'active': case 'sent': case 'delivered': return 'success'
      case 'draft': case 'pending': return 'warning'
      case 'unsubscribed': return 'default'
      case 'failed': return 'error'
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

  // Validation for Add Subscriber
  const isSubscriberFormValid = subscriberForm.phone.trim() !== ''

  // Validation for New Campaign
  const isCampaignFormValid =
    campaignForm.name.trim() !== '' &&
    campaignForm.message.trim() !== ''

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          SMS Marketing
        </Typography>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={(e, v) => setTabValue(v)}>
          <Tab label="Lists" />
          <Tab label="Subscribers" />
          <Tab label="Campaigns" />
        </Tabs>
      </Box>

      {/* Lists Tab */}
      <TabPanel value={tabValue} index={0}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">Contact Lists</Typography>
              <Button startIcon={<AddIcon />} onClick={() => setOpenListDialog(true)}>
                New List
              </Button>
            </Box>
            {listsLoading ? <Skeleton height={200} /> : (
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell>Subscriber Count</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {lists.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={4} align="center">
                            <EmptyState icon={SmsIcon} title="No lists found" description="Create your first contact list to start organizing SMS subscribers." actionLabel="New List" onAction={() => setOpenListDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        lists.map((list: any) => (
                          <TableRow key={list.id} hover>
                            <TableCell><Typography fontWeight="medium">{list.name}</Typography></TableCell>
                            <TableCell>{list.subscriber_count}</TableCell>
                            <TableCell>
                              <Chip label={list.status} size="small" color={getStatusColor(list.status)} />
                            </TableCell>
                            <TableCell>
                              <Button size="small" onClick={() => {
                                setSubscriberForm({ ...subscriberForm, list_id: list.id })
                                setOpenSubscriberDialog(true)
                              }}>Add Subscribers</Button>
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

      {/* Subscribers Tab */}
      <TabPanel value={tabValue} index={1}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">Subscribers</Typography>
              <Button startIcon={<AddIcon />} onClick={() => setOpenSubscriberDialog(true)}>
                Add Subscriber
              </Button>
            </Box>
            {subscribersLoading ? <Skeleton height={200} /> : (
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Phone</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Subscribed Date</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {subscribers.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={4} align="center">
                            <EmptyState icon={GroupIcon} title="No subscribers found" description="Add contacts to your lists to start sending SMS campaigns." actionLabel="Add Subscriber" onAction={() => setOpenSubscriberDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        subscribers.slice(subPage * subPageSize, (subPage + 1) * subPageSize).map((sub: any) => (
                          <TableRow key={sub.id} hover>
                            <TableCell>{sub.phone}</TableCell>
                            <TableCell>{sub.first_name} {sub.last_name}</TableCell>
                            <TableCell>
                              <Chip label={sub.status} size="small" color={getStatusColor(sub.status)} />
                            </TableCell>
                            <TableCell>{new Date(sub.subscribed_at).toLocaleDateString()}</TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
            )}
            {subscribers.length > subPageSize && (
              <Box sx={{ display: 'flex', justifyContent: 'center', pt: 2 }}>
                <Pagination
                  count={Math.ceil(subscribers.length / subPageSize)}
                  page={subPage + 1}
                  onChange={(_, p) => setSubPage(p - 1)}
                />
              </Box>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Campaigns Tab */}
      <TabPanel value={tabValue} index={2}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">SMS Campaigns</Typography>
              <Button startIcon={<AddIcon />} onClick={() => setOpenCampaignDialog(true)}>
                New Campaign
              </Button>
            </Box>
            {campaignsLoading ? <Skeleton height={300} /> : (
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell sx={{ minWidth: 150 }}>
                          <TableSortLabel active={orderBy === 'name'} direction={order} onClick={() => handleSort('name')}>Name</TableSortLabel>
                        </TableCell>
                        <TableCell>List</TableCell>
                        <TableCell>Message</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'sent'} direction={order} onClick={() => handleSort('sent')}>Sent</TableSortLabel>
                        </TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'delivered'} direction={order} onClick={() => handleSort('delivered')}>Delivered</TableSortLabel>
                        </TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'failed'} direction={order} onClick={() => handleSort('failed')}>Failed</TableSortLabel>
                        </TableCell>
                        <TableCell>
                          <TableSortLabel active={orderBy === 'cost'} direction={order} onClick={() => handleSort('cost')}>Cost</TableSortLabel>
                        </TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {campaigns.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={9} align="center">
                            <EmptyState icon={SmsIcon} title="No campaigns found" description="Create your first SMS campaign to reach your subscribers." actionLabel="New Campaign" onAction={() => setOpenCampaignDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        sortedCampaigns.map((campaign: any) => (
                          <TableRow key={campaign.id} hover>
                            <TableCell><Typography fontWeight="medium">{campaign.name}</Typography></TableCell>
                            <TableCell>{campaign.list_name || '-'}</TableCell>
                            <TableCell>
                              <Typography noWrap sx={{ maxWidth: 200 }}>{campaign.message}</Typography>
                            </TableCell>
                            <TableCell>
                              <Chip label={campaign.status} size="small" color={getStatusColor(campaign.status)} />
                            </TableCell>
                            <TableCell>{campaign.sent}</TableCell>
                            <TableCell>{campaign.delivered}</TableCell>
                            <TableCell>{campaign.failed}</TableCell>
                            <TableCell>{formatCurrency(campaign.cost)}</TableCell>
                            <TableCell>
                              {campaign.status === 'draft' && (
                                <IconButton size="small" onClick={() => sendCampaignMutation.mutate(campaign.id)}
                                  title="Send" aria-label="Send campaign">
                                  <SendIcon fontSize="small" />
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

      {/* New List Dialog */}
      <Dialog open={openListDialog} onClose={() => setOpenListDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>New Contact List</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="List Name"
              fullWidth
              value={listForm.name}
              onChange={(e) => setListForm({ ...listForm, name: e.target.value })}
            />
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={2}
              value={listForm.description}
              onChange={(e) => setListForm({ ...listForm, description: e.target.value })}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenListDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => createListMutation.mutate(listForm)}
            disabled={createListMutation.isPending}>
            {createListMutation.isPending ? 'Creating...' : 'Create List'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Subscriber Dialog */}
      <Dialog open={openSubscriberDialog} onClose={() => setOpenSubscriberDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Add Subscriber</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Phone Number"
              fullWidth
              value={subscriberForm.phone}
              onChange={(e) => setSubscriberForm({ ...subscriberForm, phone: e.target.value })}
            />
            <TextField
              label="First Name"
              fullWidth
              value={subscriberForm.first_name}
              onChange={(e) => setSubscriberForm({ ...subscriberForm, first_name: e.target.value })}
            />
            <TextField
              label="Last Name"
              fullWidth
              value={subscriberForm.last_name}
              onChange={(e) => setSubscriberForm({ ...subscriberForm, last_name: e.target.value })}
            />
            <TextField
              label="Add to List"
              fullWidth
              select
              SelectProps={{ native: true }}
              value={subscriberForm.list_id}
              onChange={(e) => setSubscriberForm({ ...subscriberForm, list_id: e.target.value })}
            >
              <option value="">-- Select List --</option>
              {lists.map((list: any) => (
                <option key={list.id} value={list.id}>{list.name}</option>
              ))}
            </TextField>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenSubscriberDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => addSubscriberMutation.mutate(subscriberForm)}
            disabled={!isSubscriberFormValid || addSubscriberMutation.isPending}>
            {addSubscriberMutation.isPending ? 'Adding...' : 'Add Subscriber'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* New Campaign Dialog */}
      <Dialog open={openCampaignDialog} onClose={() => setOpenCampaignDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>New SMS Campaign</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Campaign Name"
              fullWidth
              value={campaignForm.name}
              onChange={(e) => setCampaignForm({ ...campaignForm, name: e.target.value })}
            />
            <TextField
              label="Select List"
              fullWidth
              select
              SelectProps={{ native: true }}
              value={campaignForm.list_id}
              onChange={(e) => setCampaignForm({ ...campaignForm, list_id: e.target.value })}
            >
              <option value="">-- Select List --</option>
              {lists.map((list: any) => (
                <option key={list.id} value={list.id}>{list.name} ({list.subscriber_count} subscribers)</option>
              ))}
            </TextField>
            <TextField
              label="Message"
              fullWidth
              multiline
              rows={4}
              value={campaignForm.message}
              onChange={(e) => setCampaignForm({ ...campaignForm, message: e.target.value })}
              helperText="Supports variables: {{first_name}}, {{last_name}}, {{phone}}"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenCampaignDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => createCampaignMutation.mutate(campaignForm)}
            disabled={!isCampaignFormValid || createCampaignMutation.isPending}>
            {createCampaignMutation.isPending ? 'Creating...' : 'Create Campaign'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
