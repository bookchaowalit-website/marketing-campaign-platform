'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, Skeleton, MenuItem, Grid
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import PlayArrowIcon from '@mui/icons-material/PlayArrow'
import VisibilityIcon from '@mui/icons-material/Visibility'
import ThumbUpIcon from '@mui/icons-material/ThumbUp'
import VideoLibraryIcon from '@mui/icons-material/VideoLibrary'
import { unifiedCampaignApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { EmptyState } from '@/components/EmptyState'
import { formatNumber } from '@/lib/utils'

const VIDEO_CHANNELS = ['video_youtube', 'video_vimeo', 'webinar']

const defaultCampaigns = [
  { id: '1', name: 'Product Demo - YouTube', channel: 'video_youtube', platform: 'YouTube', status: 'active', title: 'How to Use Our Product', views: 15000, engagement: 4.5, likes: 320, comments: 45 },
  { id: '2', name: 'Tutorial Series Ep1', channel: 'video_youtube', platform: 'YouTube', status: 'active', title: 'Getting Started Guide', views: 8500, engagement: 3.8, likes: 180, comments: 22 },
  { id: '3', name: 'Webinar Q1 Results', channel: 'webinar', platform: 'Zoom', status: 'completed', title: 'Q1 Business Review', views: 500, engagement: 7.2, likes: 45, comments: 30 },
  { id: '4', name: 'Brand Story - Vimeo', channel: 'video_vimeo', platform: 'Vimeo', status: 'draft', title: 'Our Story', views: 0, engagement: 0, likes: 0, comments: 0 },
]

export default function VideoMarketingPage() {
  const queryClient = useQueryClient()
  const toast = useToast()

  const [openDialog, setOpenDialog] = useState(false)
  const [campaignForm, setCampaignForm] = useState({
    name: '', title: '', channel: 'video_youtube', platform: 'YouTube', description: '',
  })

  const { data: allCampaigns, isLoading } = useQuery({
    queryKey: ['video-campaigns'],
    queryFn: async () => {
      try {
        // Fetch from each video channel using the API's channel filter parameter
        const [youtubeRes, vimeoRes, webinarRes] = await Promise.all([
          unifiedCampaignApi.getAll('video_youtube', undefined, 50, 0),
          unifiedCampaignApi.getAll('video_vimeo', undefined, 50, 0),
          unifiedCampaignApi.getAll('webinar', undefined, 50, 0),
        ])
        const combined = [
          ...(youtubeRes.campaigns || []),
          ...(vimeoRes.campaigns || []),
          ...(webinarRes.campaigns || []),
        ]
        return { campaigns: combined.length > 0 ? combined : defaultCampaigns }
      } catch {
        return { campaigns: defaultCampaigns }
      }
    },
  })

  const campaigns = allCampaigns?.campaigns || defaultCampaigns

  const createMutation = useMutation({
    mutationFn: async (data: typeof campaignForm) => {
      return unifiedCampaignApi.create({
        ...data,
        type: data.channel,
        channel: data.channel,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['video-campaigns'] })
      setOpenDialog(false)
      setCampaignForm({ name: '', title: '', channel: 'video_youtube', platform: 'YouTube', description: '' })
      toast.success('Video campaign created successfully')
    },
    onError: () => toast.error('Failed to create video campaign'),
  })

  const getStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' => {
    switch (status) {
      case 'active': case 'published': case 'live': return 'success'
      case 'draft': case 'pending': return 'warning'
      case 'completed': return 'default'
      case 'scheduled': return 'info'
      default: return 'default'
    }
  }

  const getChannelColor = (channel: string): 'primary' | 'secondary' | 'success' | 'warning' | 'info' => {
    switch (channel) {
      case 'video_youtube': return 'primary'
      case 'video_vimeo': return 'secondary'
      case 'webinar': return 'success'
      default: return 'info'
    }
  }

  const getChannelLabel = (channel: string) => {
    switch (channel) {
      case 'video_youtube': return 'YouTube'
      case 'video_vimeo': return 'Vimeo'
      case 'webinar': return 'Webinar'
      default: return channel
    }
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Video Marketing
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setOpenDialog(true)}
        >
          New Video Campaign
        </Button>
      </Box>

      {/* Campaigns Table */}
      <Card>
        <CardContent>
          {isLoading ? <Skeleton height={300} /> : (
            <Box sx={{ overflowX: 'auto' }}>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Title</TableCell>
                      <TableCell>Platform</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Views</TableCell>
                      <TableCell>Likes</TableCell>
                      <TableCell>Comments</TableCell>
                      <TableCell>Engagement</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {campaigns.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={8} align="center">
                          <EmptyState icon={VideoLibraryIcon} title="No video campaigns found" description="Create your first video campaign to start reaching viewers." actionLabel="New Video Campaign" onAction={() => setOpenDialog(true)} />
                        </TableCell>
                      </TableRow>
                    ) : (
                      campaigns.map((campaign: any) => (
                        <TableRow key={campaign.id} hover>
                          <TableCell>
                            <Typography fontWeight="medium">{campaign.name}</Typography>
                          </TableCell>
                          <TableCell>{campaign.title || '-'}</TableCell>
                          <TableCell>
                            <Chip
                              label={campaign.platform || getChannelLabel(campaign.channel)}
                              size="small"
                              color={getChannelColor(campaign.channel)}
                            />
                          </TableCell>
                          <TableCell>
                            <Chip label={campaign.status} size="small" color={getStatusColor(campaign.status)} />
                          </TableCell>
                          <TableCell>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                              <VisibilityIcon fontSize="small" color="action" />
                              {formatNumber(campaign.views || 0)}
                            </Box>
                          </TableCell>
                          <TableCell>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                              <ThumbUpIcon fontSize="small" color="action" />
                              {campaign.likes || 0}
                            </Box>
                          </TableCell>
                          <TableCell>{campaign.comments || 0}</TableCell>
                          <TableCell>
                            {campaign.engagement ? `${campaign.engagement.toFixed(1)}%` : '-'}
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

      {/* New Video Campaign Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>New Video Campaign</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Campaign Name"
              fullWidth
              value={campaignForm.name}
              onChange={(e) => setCampaignForm({ ...campaignForm, name: e.target.value })}
            />
            <TextField
              label="Video Title"
              fullWidth
              value={campaignForm.title}
              onChange={(e) => setCampaignForm({ ...campaignForm, title: e.target.value })}
            />
            <TextField
              label="Platform"
              fullWidth
              select
              value={campaignForm.channel}
              onChange={(e) => setCampaignForm({
                ...campaignForm,
                channel: e.target.value,
                platform: e.target.value === 'video_youtube' ? 'YouTube' : e.target.value === 'video_vimeo' ? 'Vimeo' : 'Zoom',
              })}
            >
              <MenuItem value="video_youtube">YouTube</MenuItem>
              <MenuItem value="video_vimeo">Vimeo</MenuItem>
              <MenuItem value="webinar">Webinar</MenuItem>
            </TextField>
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={3}
              value={campaignForm.description}
              onChange={(e) => setCampaignForm({ ...campaignForm, description: e.target.value })}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => createMutation.mutate(campaignForm)}
            disabled={createMutation.isPending}>
            {createMutation.isPending ? 'Creating...' : 'Create Campaign'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
