'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/useToast'
import { marketingApi } from '@/lib/marketing-features-api'
import { Box, Card, CardContent, Typography, Grid, Chip, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, IconButton, Skeleton, Button, Tabs, Tab, Avatar, Dialog, DialogTitle, DialogContent, DialogActions, TextField, Alert, LinearProgress, Fab, Tooltip, Badge } from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import EmailIcon from '@mui/icons-material/Email'
import ShareIcon from '@mui/icons-material/Share'
import WebIcon from '@mui/icons-material/Web'
import AutorenewIcon from '@mui/icons-material/Autorenew'
import PeopleIcon from '@mui/icons-material/People'
import BarChartIcon from '@mui/icons-material/BarChart'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import TrendingDownIcon from '@mui/icons-material/TrendingDown'
import RefreshIcon from '@mui/icons-material/Refresh'
import ContentCopyIcon from '@mui/icons-material/ContentCopy'
import DeleteIcon from '@mui/icons-material/Delete'
import SendIcon from '@mui/icons-material/Send'
import ScheduleIcon from '@mui/icons-material/Schedule'
import { useState } from 'react'
import { EmailCampaign, SocialPost, EmailList, EmailSubscriber, LandingPage, AutomationWorkflow, Influencer } from '@/types/marketing-features'
import { TabPanel } from '@/components/TabPanel'
import { formatNumber } from '@/lib/utils'

export default function MarketingHubPage() {
  const { user } = useAuth()
  const toast = useToast()
  const queryClient = useQueryClient()
  const [tabValue, setTabValue] = useState(0)

  // Dialog states
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [createType, setCreateType] = useState<string>('')
  const [createForm, setCreateForm] = useState<Record<string, string>>({})

  // Fetch data
  const { data: overview, isLoading: overviewLoading, refetch: refetchOverview } = useQuery({
    queryKey: ['marketingOverview'],
    queryFn: async () => {
      try {
        return await marketingApi.getMarketingOverview()
      } catch {
        return null
      }
    },
  })

  const { data: emailLists, isLoading: listsLoading } = useQuery({
    queryKey: ['emailLists'],
    queryFn: async () => {
      try {
        return await marketingApi.getEmailLists()
      } catch {
        return { lists: [], total: 0 }
      }
    },
  })

  const { data: campaigns, isLoading: campaignsLoading } = useQuery({
    queryKey: ['emailCampaigns'],
    queryFn: async () => {
      try {
        return await marketingApi.getEmailCampaigns(10, 0)
      } catch {
        return { campaigns: [], total: 0 }
      }
    },
  })

  const { data: subscribers, isLoading: subscribersLoading } = useQuery({
    queryKey: ['emailSubscribers'],
    queryFn: async () => {
      try {
        return await marketingApi.getSubscribers(10, 0)
      } catch {
        return { subscribers: [], total: 0 }
      }
    },
  })

  const { data: socialPosts, isLoading: postsLoading } = useQuery({
    queryKey: ['socialPosts'],
    queryFn: async () => {
      try {
        return await marketingApi.getSocialPosts(10, 0)
      } catch {
        return { posts: [], total: 0 }
      }
    },
  })

  const { data: landingPages, isLoading: pagesLoading } = useQuery({
    queryKey: ['landingPages'],
    queryFn: async () => {
      try {
        return await marketingApi.getLandingPages(10, 0)
      } catch {
        return { pages: [], total: 0 }
      }
    },
  })

  const { data: workflows, isLoading: workflowsLoading } = useQuery({
    queryKey: ['automationWorkflows'],
    queryFn: async () => {
      try {
        return await marketingApi.getWorkflows()
      } catch {
        return { workflows: [], total: 0 }
      }
    },
  })

  const { data: influencers, isLoading: influencersLoading } = useQuery({
    queryKey: ['influencers'],
    queryFn: async () => {
      try {
        return await marketingApi.getInfluencers(10, 0)
      } catch {
        return { influencers: [], total: 0 }
      }
    },
  })

  // Create mutations
  const createMutation = useMutation({
    mutationFn: async (data: { type: string; payload: Record<string, unknown> }) => {
      switch (data.type) {
        case 'list':
          return marketingApi.createEmailList(data.payload as { name: string; description?: string })
        case 'campaign':
          return marketingApi.createEmailCampaign(data.payload as { name: string; subject: string; from_name?: string; from_email?: string; html_content?: string; list_ids?: string[] })
        case 'subscriber':
          return marketingApi.addSubscriber(data.payload as { email: string; first_name?: string; last_name?: string; source?: string; list_ids?: string[] })
        case 'post':
          return marketingApi.createSocialPost(data.payload as { content: string; media_urls?: string[]; hashtags?: string[] })
        case 'page':
          return marketingApi.createLandingPage(data.payload as { name: string; slug?: string; title?: string; html_content?: string })
        case 'workflow':
          return marketingApi.createWorkflow(data.payload as { name: string; trigger_type: string; trigger_config?: Record<string, unknown> })
        case 'influencer':
          return marketingApi.createInfluencer(data.payload as { name: string; email?: string; category?: string; niches?: string[] })
        default:
          throw new Error('Unknown type')
      }
    },
    onSuccess: (_, variables) => {
      toast.success(`${variables.type} created successfully!`)
      queryClient.invalidateQueries({ queryKey: ['marketingOverview'] })
      queryClient.invalidateQueries({ queryKey: ['emailLists'] })
      queryClient.invalidateQueries({ queryKey: ['emailCampaigns'] })
      queryClient.invalidateQueries({ queryKey: ['emailSubscribers'] })
      queryClient.invalidateQueries({ queryKey: ['socialPosts'] })
      queryClient.invalidateQueries({ queryKey: ['landingPages'] })
      queryClient.invalidateQueries({ queryKey: ['automationWorkflows'] })
      queryClient.invalidateQueries({ queryKey: ['influencers'] })
      setCreateDialogOpen(false)
      setCreateForm({})
    },
    onError: () => {
      toast.error('Failed to create')
    },
  })

  const handleCreate = () => {
    if (createType && createForm.name) {
      createMutation.mutate({ type: createType, payload: createForm })
    }
  }

  const openCreateDialog = (type: string) => {
    setCreateType(type)
    setCreateForm({})
    setCreateDialogOpen(true)
  }

  const statsCards = [
    { label: 'Email Subscribers', value: overview?.email?.total_subscribers || 0, color: 'primary.main', icon: EmailIcon, tab: 0 },
    { label: 'Social Posts', value: overview?.social?.total_posts || 0, color: 'success.main', icon: ShareIcon, tab: 1 },
    { label: 'Landing Pages', value: overview?.landing_pages?.total_pages || 0, color: 'secondary.main', icon: WebIcon, tab: 2 },
    { label: 'Workflows', value: overview?.influencers?.total_influencers || 0, color: 'warning.main', icon: PeopleIcon, tab: 4 },
  ]

  const getStatusColor = (status: string): 'success' | 'warning' | 'error' | 'default' | 'info' => {
    switch (status?.toLowerCase()) {
      case 'active':
      case 'published':
      case 'sent':
        return 'success'
      case 'draft':
      case 'pending':
      case 'scheduled':
        return 'warning'
      case 'paused':
        return 'default'
      case 'failed':
      case 'error':
        return 'error'
      default:
        return 'info'
    }
  }

  const getCreateDialogTitle = () => {
    switch (createType) {
      case 'list': return 'Create Email List'
      case 'campaign': return 'Create Email Campaign'
      case 'subscriber': return 'Add Subscriber'
      case 'post': return 'Create Social Post'
      case 'page': return 'Create Landing Page'
      case 'workflow': return 'Create Automation Workflow'
      case 'influencer': return 'Add Influencer'
      default: return 'Create'
    }
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Box>
          <Typography variant="h4" fontWeight="bold">
            Marketing Hub
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Manage all your marketing activities in one place
          </Typography>
        </Box>
        <Button startIcon={<RefreshIcon />} onClick={() => refetchOverview()}>
          Refresh
        </Button>
      </Box>

      {/* Stats Grid */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {statsCards.map((stat, i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <Card
              sx={{ height: '100%', cursor: 'pointer', '&:hover': { boxShadow: 4 } }}
              onClick={() => setTabValue(stat.tab)}
            >
              <CardContent>
                {overviewLoading ? (
                  <Skeleton height={80} />
                ) : (
                  <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        {stat.label}
                      </Typography>
                      <Typography variant="h4" fontWeight="bold" sx={{ color: stat.color }}>
                        {formatNumber(stat.value)}
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

      {/* Quick Actions */}
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Quick Actions
          </Typography>
          <Grid container spacing={2}>
            {[
              { label: 'New Email Campaign', icon: <EmailIcon />, type: 'campaign', color: 'primary.main' },
              { label: 'Schedule Post', icon: <ScheduleIcon />, type: 'post', color: 'success.main' },
              { label: 'Create Landing Page', icon: <WebIcon />, type: 'page', color: 'secondary.main' },
              { label: 'Add Subscriber', icon: <PeopleIcon />, type: 'subscriber', color: 'warning.main' },
              { label: 'Create Workflow', icon: <AutorenewIcon />, type: 'workflow', color: 'error.main' },
              { label: 'Add Influencer', icon: <PeopleIcon />, type: 'influencer', color: 'info.main' },
            ].map((action, i) => (
              <Grid item key={i}>
                <Button
                  variant="outlined"
                  startIcon={action.icon}
                  onClick={() => openCreateDialog(action.type)}
                  sx={{ borderColor: action.color, color: action.color, '&:hover': { borderColor: action.color, bgcolor: `${action.color}10` } }}
                >
                  {action.label}
                </Button>
              </Grid>
            ))}
          </Grid>
        </CardContent>
      </Card>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={(e, v) => setTabValue(v)}>
          <Tab icon={<EmailIcon />} iconPosition="start" label="Email Marketing" />
          <Tab icon={<ShareIcon />} iconPosition="start" label="Social Media" />
          <Tab icon={<WebIcon />} iconPosition="start" label="Landing Pages" />
          <Tab icon={<AutorenewIcon />} iconPosition="start" label="Automation" />
          <Tab icon={<PeopleIcon />} iconPosition="start" label="Influencers" />
        </Tabs>
      </Box>

      {/* Email Marketing Tab */}
      <TabPanel value={tabValue} index={0}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6" fontWeight="bold">Email Campaigns</Typography>
                  <Button size="small" startIcon={<AddIcon />} onClick={() => openCreateDialog('campaign')}>
                    New Campaign
                  </Button>
                </Box>
                {campaignsLoading ? <Skeleton height={200} /> : (
                  <TableContainer>
                    <Table size="small">
                      <TableHead>
                        <TableRow>
                          <TableCell>Name</TableCell>
                          <TableCell>Status</TableCell>
                          <TableCell>Open Rate</TableCell>
                          <TableCell>Click Rate</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {campaigns?.campaigns?.slice(0, 5).map((campaign: EmailCampaign) => (
                          <TableRow key={campaign.id} hover>
                            <TableCell>{campaign.name}</TableCell>
                            <TableCell><Chip label={campaign.status} size="small" color={getStatusColor(campaign.status)} /></TableCell>
                            <TableCell>{campaign.open_rate?.toFixed(1)}%</TableCell>
                            <TableCell>{campaign.click_rate?.toFixed(1)}%</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                )}
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6" fontWeight="bold">Email Lists</Typography>
                  <Button size="small" startIcon={<AddIcon />} onClick={() => openCreateDialog('list')}>
                    New List
                  </Button>
                </Box>
                {listsLoading ? <Skeleton height={200} /> : (
                  <TableContainer>
                    <Table size="small">
                      <TableHead>
                        <TableRow>
                          <TableCell>List Name</TableCell>
                          <TableCell>Subscribers</TableCell>
                          <TableCell>Status</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {emailLists?.lists?.slice(0, 5).map((list: EmailList) => (
                          <TableRow key={list.id} hover>
                            <TableCell>{list.name}</TableCell>
                            <TableCell>{list.subscriber_count}</TableCell>
                            <TableCell><Chip label={list.status} size="small" color={getStatusColor(list.status)} /></TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                )}
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6" fontWeight="bold">Recent Subscribers</Typography>
                  <Button size="small" startIcon={<AddIcon />} onClick={() => openCreateDialog('subscriber')}>
                    Add Subscriber
                  </Button>
                </Box>
                {subscribersLoading ? <Skeleton height={150} /> : (
                  <TableContainer>
                    <Table size="small">
                      <TableHead>
                        <TableRow>
                          <TableCell>Email</TableCell>
                          <TableCell>Name</TableCell>
                          <TableCell>Status</TableCell>
                          <TableCell>Source</TableCell>
                          <TableCell>Added</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {subscribers?.subscribers?.slice(0, 5).map((sub: EmailSubscriber) => (
                          <TableRow key={sub.id} hover>
                            <TableCell>{sub.email}</TableCell>
                            <TableCell>{sub.first_name} {sub.last_name}</TableCell>
                            <TableCell><Chip label={sub.status} size="small" color={getStatusColor(sub.status)} /></TableCell>
                            <TableCell>{sub.source || '-'}</TableCell>
                            <TableCell>{new Date(sub.created_at).toLocaleDateString()}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                )}
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>

      {/* Social Media Tab */}
      <TabPanel value={tabValue} index={1}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">Social Media Posts</Typography>
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => openCreateDialog('post')}>
                Create Post
              </Button>
            </Box>
            {postsLoading ? <Skeleton height={300} /> : (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Content</TableCell>
                      <TableCell>Hashtags</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Scheduled</TableCell>
                      <TableCell>Created</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {socialPosts?.posts?.slice(0, 10).map((post: SocialPost) => (
                      <TableRow key={post.id} hover>
                        <TableCell sx={{ maxWidth: 300 }}>
                          <Typography noWrap>{post.content}</Typography>
                        </TableCell>
                        <TableCell>
                          {post.hashtags?.slice(0, 3).map((tag, i) => (
                            <Chip key={i} label={`#${tag}`} size="small" sx={{ mr: 0.5 }} />
                          ))}
                        </TableCell>
                        <TableCell><Chip label={post.status} size="small" color={getStatusColor(post.status)} /></TableCell>
                        <TableCell>
                          {post.schedules?.[0]?.scheduled_time
                            ? new Date(post.schedules[0].scheduled_time).toLocaleString()
                            : '-'}
                        </TableCell>
                        <TableCell>{new Date(post.created_at).toLocaleDateString()}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Landing Pages Tab */}
      <TabPanel value={tabValue} index={2}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">Landing Pages</Typography>
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => openCreateDialog('page')}>
                Create Page
              </Button>
            </Box>
            {pagesLoading ? <Skeleton height={300} /> : (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Slug</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Views</TableCell>
                      <TableCell>Conversions</TableCell>
                      <TableCell>Created</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {landingPages?.pages?.map((page: LandingPage) => (
                      <TableRow key={page.id} hover>
                        <TableCell>{page.name}</TableCell>
                        <TableCell>
                          <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                            /{page.slug}
                          </Typography>
                        </TableCell>
                        <TableCell><Chip label={page.status} size="small" color={getStatusColor(page.status)} /></TableCell>
                        <TableCell>{formatNumber(page.view_count)}</TableCell>
                        <TableCell>{page.conversion_count}</TableCell>
                        <TableCell>{new Date(page.created_at).toLocaleDateString()}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Automation Tab */}
      <TabPanel value={tabValue} index={3}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">Automation Workflows</Typography>
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => openCreateDialog('workflow')}>
                Create Workflow
              </Button>
            </Box>
            {workflowsLoading ? <Skeleton height={300} /> : (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Trigger</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Active Contacts</TableCell>
                      <TableCell>Total Contacts</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {workflows?.workflows?.map((workflow: AutomationWorkflow) => (
                      <TableRow key={workflow.id} hover>
                        <TableCell>{workflow.name}</TableCell>
                        <TableCell>
                          <Chip label={workflow.trigger_type} size="small" variant="outlined" />
                        </TableCell>
                        <TableCell><Chip label={workflow.status} size="small" color={getStatusColor(workflow.status)} /></TableCell>
                        <TableCell>{workflow.active_contacts}</TableCell>
                        <TableCell>{workflow.total_contacts}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Influencers Tab */}
      <TabPanel value={tabValue} index={4}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">Influencers</Typography>
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => openCreateDialog('influencer')}>
                Add Influencer
              </Button>
            </Box>
            {influencersLoading ? <Skeleton height={300} /> : (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Email</TableCell>
                      <TableCell>Category</TableCell>
                      <TableCell>Niches</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Rating</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {influencers?.influencers?.map((inf: Influencer) => (
                      <TableRow key={inf.id} hover>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <Avatar sx={{ width: 32, height: 32 }}>{inf.name?.[0]}</Avatar>
                            {inf.name}
                          </Box>
                        </TableCell>
                        <TableCell>{inf.email || '-'}</TableCell>
                        <TableCell>{inf.category || '-'}</TableCell>
                        <TableCell>
                          {inf.niches?.slice(0, 2).map((niche, i) => (
                            <Chip key={i} label={niche} size="small" sx={{ mr: 0.5 }} />
                          ))}
                        </TableCell>
                        <TableCell><Chip label={inf.status} size="small" color={getStatusColor(inf.status)} /></TableCell>
                        <TableCell>
                          {inf.rating > 0 ? `${inf.rating.toFixed(1)} ⭐` : '-'}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      {/* Create Dialog */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{getCreateDialogTitle()}</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Name"
            fullWidth
            value={createForm.name || ''}
            onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
            sx={{ mb: 2 }}
          />

          {createType === 'campaign' && (
            <>
              <TextField
                margin="dense"
                label="Subject"
                fullWidth
                value={createForm.subject || ''}
                onChange={(e) => setCreateForm({ ...createForm, subject: e.target.value })}
                sx={{ mb: 2 }}
              />
              <TextField
                margin="dense"
                label="From Name"
                fullWidth
                value={createForm.from_name || ''}
                onChange={(e) => setCreateForm({ ...createForm, from_name: e.target.value })}
              />
            </>
          )}

          {createType === 'subscriber' && (
            <>
              <TextField
                margin="dense"
                label="Email"
                type="email"
                fullWidth
                value={createForm.email || ''}
                onChange={(e) => setCreateForm({ ...createForm, email: e.target.value })}
                sx={{ mb: 2 }}
              />
              <TextField
                margin="dense"
                label="First Name"
                fullWidth
                value={createForm.first_name || ''}
                onChange={(e) => setCreateForm({ ...createForm, first_name: e.target.value })}
                sx={{ mb: 2 }}
              />
              <TextField
                margin="dense"
                label="Last Name"
                fullWidth
                value={createForm.last_name || ''}
                onChange={(e) => setCreateForm({ ...createForm, last_name: e.target.value })}
              />
            </>
          )}

          {createType === 'post' && (
            <TextField
              margin="dense"
              label="Content"
              fullWidth
              multiline
              rows={4}
              value={createForm.content || ''}
              onChange={(e) => setCreateForm({ ...createForm, content: e.target.value })}
            />
          )}

          {createType === 'page' && (
            <>
              <TextField
                margin="dense"
                label="Slug"
                fullWidth
                value={createForm.slug || ''}
                onChange={(e) => setCreateForm({ ...createForm, slug: e.target.value })}
                sx={{ mb: 2 }}
              />
              <TextField
                margin="dense"
                label="Title"
                fullWidth
                value={createForm.title || ''}
                onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
              />
            </>
          )}

          {createType === 'workflow' && (
            <>
              <TextField
                margin="dense"
                label="Trigger Type"
                fullWidth
                select
                value={createForm.trigger_type || 'list_subscription'}
                onChange={(e) => setCreateForm({ ...createForm, trigger_type: e.target.value })}
                sx={{ mb: 2 }}
              >
                <option value="list_subscription">List Subscription</option>
                <option value="tag_added">Tag Added</option>
                <option value="form_submission">Form Submission</option>
                <option value="page_visit">Page Visit</option>
                <option value="purchase">Purchase</option>
              </TextField>
            </>
          )}

          {createType === 'influencer' && (
            <>
              <TextField
                margin="dense"
                label="Email"
                type="email"
                fullWidth
                value={createForm.email || ''}
                onChange={(e) => setCreateForm({ ...createForm, email: e.target.value })}
                sx={{ mb: 2 }}
              />
              <TextField
                margin="dense"
                label="Category"
                fullWidth
                value={createForm.category || ''}
                onChange={(e) => setCreateForm({ ...createForm, category: e.target.value })}
              />
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleCreate} disabled={createMutation.isPending || !createForm.name}>
            {createMutation.isPending ? 'Creating...' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
