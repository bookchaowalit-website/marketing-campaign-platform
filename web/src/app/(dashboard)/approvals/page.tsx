'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent,
  DialogActions, TextField, Skeleton, Tabs, Tab, Avatar,
} from '@mui/material'
import CheckIcon from '@mui/icons-material/Check'
import CloseIcon from '@mui/icons-material/Close'
import CommentIcon from '@mui/icons-material/Comment'
import PersonIcon from '@mui/icons-material/Person'
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn'
import { approvalApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { ConfirmDialog, useConfirmDialog } from '@/components/ConfirmDialog'
import { EmptyState } from '@/components/EmptyState'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props
  return (
    <div role="tabpanel" hidden={value !== index} {...other}>
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  )
}

const defaultRequests: any[] = []

export default function ApprovalsPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const confirmDialog = useConfirmDialog()

  const [tabValue, setTabValue] = useState(0)
  const [openRespondDialog, setOpenRespondDialog] = useState(false)
  const [selectedRequest, setSelectedRequest] = useState<any>(null)
  const [responseAction, setResponseAction] = useState<'approve' | 'reject' | 'comment'>('approve')
  const [comments, setComments] = useState('')

  const { data: myRequestsData, isLoading: myRequestsLoading } = useQuery<{ requests: any[] }>({
    queryKey: ['approval-requests', 'my'],
    queryFn: () => approvalApi.getRequests(),
    initialData: { requests: defaultRequests },
  })

  const { data: pendingData, isLoading: pendingLoading } = useQuery<{ requests: any[] }>({
    queryKey: ['approval-requests', 'pending'],
    queryFn: () => approvalApi.getRequests('pending'),
    initialData: { requests: defaultRequests },
  })

  const myRequests = myRequestsData?.requests || defaultRequests
  const pendingRequests = pendingData?.requests || defaultRequests

  const respondMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: { action: string; comments: string } }) => {
      await approvalApi.respond(id, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['approval-requests'] })
      setOpenRespondDialog(false)
      setComments('')
      setSelectedRequest(null)
      toast.success(`Request ${responseAction === 'approve' ? 'approved' : responseAction === 'reject' ? 'rejected' : 'commented on'} successfully`)
    },
    onError: () => {
      toast.error('Failed to respond to request')
    },
  })

  const getStatusColor = (status: string): 'success' | 'warning' | 'default' | 'info' | 'error' => {
    switch (status) {
      case 'approved': return 'success'
      case 'rejected': return 'error'
      case 'pending': return 'warning'
      case 'submitted': return 'info'
      case 'commented': return 'info'
      default: return 'default'
    }
  }

  const getEntityIcon = (entityType: string): string => {
    switch (entityType) {
      case 'campaign': return 'C'
      case 'budget': return 'B'
      case 'content': return 'C'
      case 'creative': return 'A'
      default: return 'E'
    }
  }

  const openRespond = (request: any, action: 'approve' | 'reject' | 'comment') => {
    setSelectedRequest(request)
    setResponseAction(action)
    setComments('')
    setOpenRespondDialog(true)
  }

  const handleRespond = () => {
    if (selectedRequest) {
      respondMutation.mutate({
        id: selectedRequest.id,
        data: { action: responseAction, comments },
      })
    }
  }

  const renderRequestsTable = (requests: any[], showActions: boolean) => (
    <Box sx={{ overflowX: 'auto' }}>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Entity</TableCell>
              <TableCell>Title</TableCell>
              <TableCell>Submitted By</TableCell>
              <TableCell>Date</TableCell>
              <TableCell>Status</TableCell>
              {showActions && <TableCell>Actions</TableCell>}
            </TableRow>
          </TableHead>
          <TableBody>
            {requests.length === 0 ? (
              <TableRow>
                <TableCell colSpan={showActions ? 6 : 5} align="center">
                  <EmptyState icon={AssignmentTurnedInIcon} title="No requests found" description="Approval requests will appear here when campaigns or budgets are submitted for review." />
                </TableCell>
              </TableRow>
            ) : (
              requests.map((req: any) => (
                <TableRow key={req.id} hover>
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Avatar sx={{ bgcolor: 'primary.main', width: 32, height: 32, fontSize: 14 }}>
                        {getEntityIcon(req.entity_type || 'entity')}
                      </Avatar>
                      <Typography variant="body2" textTransform="capitalize">
                        {req.entity_type || 'entity'}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Typography fontWeight="medium">{req.title || req.entity_title || '-'}</Typography>
                  </TableCell>
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <PersonIcon fontSize="small" color="action" />
                      {req.submitted_by || req.submitter_name || '-'}
                    </Box>
                  </TableCell>
                  <TableCell>
                    {req.submitted_at ? new Date(req.submitted_at).toLocaleDateString() : '-'}
                  </TableCell>
                  <TableCell>
                    <Chip label={req.status || 'pending'} size="small" color={getStatusColor(req.status || 'pending')} />
                  </TableCell>
                  {showActions && (
                    <TableCell>
                      {req.status === 'pending' ? (
                        <Box sx={{ display: 'flex', gap: 0.5 }}>
                          <Button size="small" variant="outlined" color="success" startIcon={<CheckIcon />}
                            onClick={() => openRespond(req, 'approve')}>
                            Approve
                          </Button>
                          <Button size="small" variant="outlined" color="error" startIcon={<CloseIcon />}
                            onClick={() => openRespond(req, 'reject')}>
                            Reject
                          </Button>
                          <Button size="small" variant="outlined" startIcon={<CommentIcon />}
                            onClick={() => openRespond(req, 'comment')}>
                            Comment
                          </Button>
                        </Box>
                      ) : (
                        <Typography variant="caption" color="text.secondary">
                          {req.status === 'approved' ? 'Approved' : req.status === 'rejected' ? 'Rejected' : 'Reviewed'}
                        </Typography>
                      )}
                    </TableCell>
                  )}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  )

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight="bold">
          Approvals
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Manage approval requests for campaigns, budgets, and content
        </Typography>
      </Box>

      {/* Tabs */}
      <Card>
        <Tabs value={tabValue} onChange={(_, v) => setTabValue(v)} sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tab label={`My Requests (${myRequests.length})`} />
          <Tab label={`Pending Approvals (${pendingRequests.length})`} />
        </Tabs>

        {/* My Requests Tab */}
        <TabPanel value={tabValue} index={0}>
          <Typography variant="h6" sx={{ mb: 2 }}>My Approval Requests</Typography>
          {myRequestsLoading ? <Skeleton height={200} /> : renderRequestsTable(myRequests, false)}
        </TabPanel>

        {/* Pending Approvals Tab */}
        <TabPanel value={tabValue} index={1}>
          <Typography variant="h6" sx={{ mb: 2 }}>Pending Approvals</Typography>
          {pendingLoading ? <Skeleton height={200} /> : renderRequestsTable(pendingRequests, true)}
        </TabPanel>
      </Card>

      {/* Respond Dialog */}
      <Dialog open={openRespondDialog} onClose={() => setOpenRespondDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {responseAction === 'approve' ? 'Approve Request' :
            responseAction === 'reject' ? 'Reject Request' : 'Add Comment'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2 }}>
            {selectedRequest && (
              <Box sx={{ mb: 2, p: 2, bgcolor: 'action.hover', borderRadius: 1 }}>
                <Typography variant="body2" fontWeight="medium">{selectedRequest.title}</Typography>
                <Typography variant="caption" color="text.secondary">
                  {selectedRequest.entity_type} | Submitted by {selectedRequest.submitted_by}
                </Typography>
              </Box>
            )}
            <TextField
              label={responseAction === 'approve' ? 'Approval Comments (optional)' :
                responseAction === 'reject' ? 'Rejection Reason (required)' : 'Comment'}
              fullWidth
              multiline
              rows={4}
              value={comments}
              onChange={(e) => setComments(e.target.value)}
              required={responseAction === 'reject'}
              placeholder={responseAction === 'approve' ? 'Add any approval notes...' :
                responseAction === 'reject' ? 'Explain why this is being rejected...' : 'Add your comment...'}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenRespondDialog(false)}>Cancel</Button>
          <Button
            variant="contained"
            color={responseAction === 'approve' ? 'success' : responseAction === 'reject' ? 'error' : 'primary'}
            onClick={handleRespond}
            disabled={respondMutation.isPending || (responseAction === 'reject' && !comments)}
          >
            {respondMutation.isPending ? 'Submitting...' :
              responseAction === 'approve' ? 'Approve' :
                responseAction === 'reject' ? 'Reject' : 'Submit Comment'}
          </Button>
        </DialogActions>
      </Dialog>

      <ConfirmDialog {...confirmDialog.dialogProps} />
    </Box>
  )
}
