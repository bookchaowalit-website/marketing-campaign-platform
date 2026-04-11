'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent,
  DialogActions, TextField, MenuItem, Skeleton, Grid,
  FormControl, InputLabel, Select, IconButton,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import EventIcon from '@mui/icons-material/Event'
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline'
import { calendarApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { ConfirmDialog, useConfirmDialog } from '@/components/ConfirmDialog'
import { EmptyState } from '@/components/EmptyState'
import { formatDate } from '@/lib/utils'

const EVENT_TYPES = ['campaign', 'content', 'launch', 'deadline', 'meeting', 'review', 'other']

const EVENT_TYPE_COLORS: Record<string, 'primary' | 'secondary' | 'success' | 'warning' | 'info' | 'error'> = {
  campaign: 'primary',
  content: 'secondary',
  launch: 'success',
  deadline: 'error',
  meeting: 'info',
  review: 'warning',
  other: 'primary',
}

const CHANNELS = ['all', 'Google Ads', 'Facebook', 'Instagram', 'TikTok', 'Twitter/X', 'LinkedIn', 'YouTube', 'Email', 'SEO', 'Offline']

const defaultEvents: any[] = []

export default function CalendarPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const confirmDialog = useConfirmDialog()

  const [openDialog, setOpenDialog] = useState(false)
  const [editingEvent, setEditingEvent] = useState<any>(null)
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    start_date: new Date().toISOString().split('T')[0],
    end_date: '',
    start_time: '09:00',
    end_time: '10:00',
    event_type: 'campaign',
    channel: '',
    campaign_id: '',
    color: '',
  })

  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')
  const [channelFilter, setChannelFilter] = useState('all')
  const [eventTypeFilter, setEventTypeFilter] = useState('all')

  const { data: eventsData, isLoading, error, refetch } = useQuery<{ events: any[] }>({
    queryKey: ['calendar-events', startDate, endDate, channelFilter, eventTypeFilter],
    queryFn: () => {
      const params: any = {}
      if (startDate) params.start = startDate
      if (endDate) params.end = endDate
      if (channelFilter !== 'all') params.channel = channelFilter
      if (eventTypeFilter !== 'all') params.event_type = eventTypeFilter
      return calendarApi.getEvents(params)
    },
    initialData: { events: defaultEvents },
  })

  const events = eventsData?.events || defaultEvents

  const createMutation = useMutation({
    mutationFn: async (data: typeof formData) => {
      await calendarApi.createEvent(data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
      setOpenDialog(false)
      setFormData({ title: '', description: '', start_date: new Date().toISOString().split('T')[0], end_date: '', start_time: '09:00', end_time: '10:00', event_type: 'campaign', channel: '', campaign_id: '', color: '' })
      setEditingEvent(null)
      toast.success('Event created successfully')
    },
    onError: () => {
      toast.error('Failed to create event')
    },
  })

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: typeof formData }) => {
      await calendarApi.updateEvent(id, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
      setOpenDialog(false)
      setEditingEvent(null)
      toast.success('Event updated successfully')
    },
    onError: () => {
      toast.error('Failed to update event')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await calendarApi.deleteEvent(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
      toast.success('Event deleted successfully')
    },
    onError: () => {
      toast.error('Failed to delete event')
    },
  })

  const handleOpenDialog = () => {
    setEditingEvent(null)
    setFormData({ title: '', description: '', start_date: new Date().toISOString().split('T')[0], end_date: '', start_time: '09:00', end_time: '10:00', event_type: 'campaign', channel: '', campaign_id: '', color: '' })
    setOpenDialog(true)
  }

  const handleSubmit = () => {
    if (editingEvent) {
      updateMutation.mutate({ id: editingEvent.id, data: formData })
    } else {
      createMutation.mutate(formData)
    }
  }

  const handleDelete = (event: any) => {
    confirmDialog.showConfirm(
      'Delete Event',
      `Are you sure you want to delete "${event.title}"?`,
      () => deleteMutation.mutate(event.id)
    )
  }

  // Validation: title required, end_time after start_time if both set
  const isFormValid =
    formData.title.trim() !== '' &&
    (formData.start_time === '' || formData.end_time === '' || formData.end_time > formData.start_time)

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Marketing Calendar
        </Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={handleOpenDialog}>
          New Event
        </Button>
      </Box>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={6} md={3}>
              <TextField label="Start Date" fullWidth type="date" value={startDate}
                onChange={(e) => setStartDate(e.target.value)} InputLabelProps={{ shrink: true }} size="small" />
            </Grid>
            <Grid item xs={6} md={3}>
              <TextField label="End Date" fullWidth type="date" value={endDate}
                onChange={(e) => setEndDate(e.target.value)} InputLabelProps={{ shrink: true }} size="small" />
            </Grid>
            <Grid item xs={6} md={2}>
              <FormControl fullWidth size="small">
                <InputLabel>Channel</InputLabel>
                <Select value={channelFilter} label="Channel" onChange={(e) => setChannelFilter(e.target.value)}>
                  {CHANNELS.map((ch) => (
                    <MenuItem key={ch} value={ch}>{ch === 'all' ? 'All Channels' : ch}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={6} md={2}>
              <FormControl fullWidth size="small">
                <InputLabel>Event Type</InputLabel>
                <Select value={eventTypeFilter} label="Event Type" onChange={(e) => setEventTypeFilter(e.target.value)}>
                  <MenuItem value="all">All Types</MenuItem>
                  {EVENT_TYPES.map((t) => (
                    <MenuItem key={t} value={t}>{t.charAt(0).toUpperCase() + t.slice(1)}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={2}>
              <Typography variant="body2" color="text.secondary">
                {events.length} event(s)
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Events Table */}
      <Card>
        <CardContent>
          {isLoading ? <Skeleton height={300} /> : error ? (
            <Box sx={{ py: 8, textAlign: 'center' }}>
              <ErrorOutlineIcon sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
              <Typography variant="h6" color="error" gutterBottom>Failed to load calendar events</Typography>
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
                      <TableCell>Date</TableCell>
                      <TableCell>Title</TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Channel</TableCell>
                      <TableCell>Time</TableCell>
                      <TableCell>Description</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {events.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} align="center">
                          <EmptyState icon={EventIcon} title="No events scheduled" description="Create your first calendar event to start planning." actionLabel="New Event" onAction={() => setOpenDialog(true)} />
                        </TableCell>
                      </TableRow>
                    ) : (
                      events.map((event: any) => (
                        <TableRow key={event.id} hover>
                          <TableCell>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                              <EventIcon fontSize="small" color="action" />
                              <Box>
                                <Typography variant="body2" fontWeight="medium">{formatDate(event.start_date)}</Typography>
                                {event.end_date && event.end_date !== event.start_date && (
                                  <Typography variant="caption" color="text.secondary">to {formatDate(event.end_date)}</Typography>
                                )}
                              </Box>
                            </Box>
                          </TableCell>
                          <TableCell>
                            <Typography fontWeight="medium">{event.title}</Typography>
                          </TableCell>
                          <TableCell>
                            <Chip
                              label={event.event_type || 'other'}
                              size="small"
                              color={EVENT_TYPE_COLORS[event.event_type || 'other'] || 'default'}
                            />
                          </TableCell>
                          <TableCell>
                            {event.channel ? (
                              <Chip label={event.channel} size="small" variant="outlined" />
                            ) : '-'}
                          </TableCell>
                          <TableCell>
                            {event.start_time ? `${event.start_time} - ${event.end_time || ''}` : 'All day'}
                          </TableCell>
                          <TableCell>
                            <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                              {event.description || '-'}
                            </Typography>
                          </TableCell>
                          <TableCell>
                            <Box sx={{ display: 'flex', gap: 0.5 }}>
                              <IconButton size="small" onClick={() => {
                                setEditingEvent(event)
                                setFormData({
                                  title: event.title || '',
                                  description: event.description || '',
                                  start_date: event.start_date || new Date().toISOString().split('T')[0],
                                  end_date: event.end_date || '',
                                  start_time: event.start_time || '09:00',
                                  end_time: event.end_time || '10:00',
                                  event_type: event.event_type || 'campaign',
                                  channel: event.channel || '',
                                  campaign_id: event.campaign_id || '',
                                  color: event.color || '',
                                })
                                setOpenDialog(true)
                              }} title="Edit" aria-label="Edit event">
                                <EditIcon fontSize="small" />
                              </IconButton>
                              <IconButton size="small" color="error" onClick={() => handleDelete(event)} title="Delete" aria-label="Delete event">
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

      {/* Create/Edit Event Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingEvent ? 'Edit Event' : 'Create New Event'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField label="Title" fullWidth value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })} required />
            <FormControl fullWidth>
              <InputLabel>Event Type</InputLabel>
              <Select value={formData.event_type} label="Event Type"
                onChange={(e) => setFormData({ ...formData, event_type: e.target.value })}>
                {EVENT_TYPES.map((t) => (
                  <MenuItem key={t} value={t}>{t.charAt(0).toUpperCase() + t.slice(1)}</MenuItem>
                ))}
              </Select>
            </FormControl>
            <TextField label="Channel" fullWidth value={formData.channel}
              onChange={(e) => setFormData({ ...formData, channel: e.target.value })} />
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField label="Start Date" fullWidth type="date" value={formData.start_date}
                  onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                  InputLabelProps={{ shrink: true }} required />
              </Grid>
              <Grid item xs={6}>
                <TextField label="End Date" fullWidth type="date" value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  InputLabelProps={{ shrink: true }} />
              </Grid>
            </Grid>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField label="Start Time" fullWidth type="time" value={formData.start_time}
                  onChange={(e) => setFormData({ ...formData, start_time: e.target.value })}
                  InputLabelProps={{ shrink: true }} />
              </Grid>
              <Grid item xs={6}>
                <TextField label="End Time" fullWidth type="time" value={formData.end_time}
                  onChange={(e) => setFormData({ ...formData, end_time: e.target.value })}
                  InputLabelProps={{ shrink: true }} />
              </Grid>
            </Grid>
            <TextField label="Campaign ID" fullWidth value={formData.campaign_id}
              onChange={(e) => setFormData({ ...formData, campaign_id: e.target.value })} />
            <TextField label="Description" fullWidth multiline rows={3} value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })} />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleSubmit}
            disabled={!isFormValid || createMutation.isPending || updateMutation.isPending}>
            {createMutation.isPending || updateMutation.isPending
              ? (editingEvent ? 'Updating...' : 'Creating...')
              : (editingEvent ? 'Update' : 'Create')}
          </Button>
        </DialogActions>
      </Dialog>

      <ConfirmDialog {...confirmDialog.dialogProps} />
    </Box>
  )
}
