'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  Box, Card, CardContent, Typography, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Chip, Button, Dialog, DialogTitle, DialogContent,
  DialogActions, TextField, MenuItem, Skeleton, Tabs, Tab, Grid,
  FormControl, InputLabel, Select, Alert, AlertTitle, Pagination,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import AttachMoneyIcon from '@mui/icons-material/AttachMoney'
import WarningIcon from '@mui/icons-material/Warning'
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline'
import AccountBalanceIcon from '@mui/icons-material/AccountBalance'
import ReceiptLongIcon from '@mui/icons-material/ReceiptLong'
import NotificationsIcon from '@mui/icons-material/Notifications'
import { budgetApi } from '@/lib/api-client'
import { useToast } from '@/hooks/useToast'
import { ConfirmDialog, useConfirmDialog } from '@/components/ConfirmDialog'
import { EmptyState } from '@/components/EmptyState'
import { formatCurrency, formatDate } from '@/lib/utils'

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

const defaultPeriods: any[] = []
const defaultTransactions: any[] = []
const defaultAlerts: any[] = []

export default function BudgetPage() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const confirmDialog = useConfirmDialog()

  const [tabValue, setTabValue] = useState(0)
  const [txnPage, setTxnPage] = useState(0)
  const txnPageSize = 10
  const [openPeriodDialog, setOpenPeriodDialog] = useState(false)
  const [openTransactionDialog, setOpenTransactionDialog] = useState(false)
  const [openAllocationDialog, setOpenAllocationDialog] = useState(false)
  const [selectedPeriod, setSelectedPeriod] = useState<any>(null)

  const [periodForm, setPeriodForm] = useState({
    name: '',
    start_date: new Date().toISOString().split('T')[0],
    end_date: '',
    total_budget: 0,
  })

  const [transactionForm, setTransactionForm] = useState({
    period_id: '',
    campaign_id: '',
    amount: 0,
    type: 'spend',
    description: '',
    date: new Date().toISOString().split('T')[0],
  })

  const [allocationForm, setAllocationForm] = useState({
    campaign_id: '',
    campaign_name: '',
    amount: 0,
    channel: '',
  })

  const { data: periodsData, isLoading: periodsLoading, error: periodsError, refetch: refetchPeriods } = useQuery<{ periods: any[] }>({
    queryKey: ['budget-periods'],
    queryFn: () => budgetApi.getPeriods(),
    initialData: { periods: defaultPeriods },
  })

  const { data: transactionsData, isLoading: transactionsLoading, error: transactionsError, refetch: refetchTransactions } = useQuery<{ transactions: any[] }>({
    queryKey: ['budget-transactions'],
    queryFn: () => budgetApi.getTransactions(),
    initialData: { transactions: defaultTransactions },
  })

  const { data: alertsData, isLoading: alertsLoading, error: alertsError, refetch: refetchAlerts } = useQuery<{ alerts: any[] }>({
    queryKey: ['budget-alerts'],
    queryFn: () => budgetApi.getAlerts(),
    initialData: { alerts: defaultAlerts },
  })

  const periods = periodsData?.periods || defaultPeriods
  const transactions = transactionsData?.transactions || defaultTransactions
  const alerts = alertsData?.alerts || defaultAlerts

  const createPeriodMutation = useMutation({
    mutationFn: async (data: typeof periodForm) => {
      await budgetApi.createPeriod({ ...data, total_budget: Number(data.total_budget) })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budget-periods'] })
      setOpenPeriodDialog(false)
      setPeriodForm({ name: '', start_date: new Date().toISOString().split('T')[0], end_date: '', total_budget: 0 })
      toast.success('Budget period created successfully')
    },
    onError: () => {
      toast.error('Failed to create budget period')
    },
  })

  const createAllocationMutation = useMutation({
    mutationFn: async ({ periodId, data }: { periodId: string; data: typeof allocationForm }) => {
      await budgetApi.createAllocation(periodId, { ...data, amount: Number(data.amount) })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budget-periods'] })
      setOpenAllocationDialog(false)
      setAllocationForm({ campaign_id: '', campaign_name: '', amount: 0, channel: '' })
      toast.success('Allocation created successfully')
    },
    onError: () => {
      toast.error('Failed to create allocation')
    },
  })

  const recordTransactionMutation = useMutation({
    mutationFn: async (data: typeof transactionForm) => {
      await budgetApi.recordTransaction({ ...data, amount: Number(data.amount) })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budget-transactions'] })
      setOpenTransactionDialog(false)
      setTransactionForm({ period_id: '', campaign_id: '', amount: 0, type: 'spend', description: '', date: new Date().toISOString().split('T')[0] })
      toast.success('Transaction recorded successfully')
    },
    onError: () => {
      toast.error('Failed to record transaction')
    },
  })

  const getSeverityColor = (severity: string): 'error' | 'warning' | 'info' | 'success' => {
    switch (severity) {
      case 'critical': return 'error'
      case 'warning': return 'warning'
      case 'info': return 'info'
      default: return 'info'
    }
  }

  // Validation for Create Budget Period
  const isPeriodFormValid =
    periodForm.name.trim() !== '' &&
    periodForm.total_budget > 0 &&
    periodForm.start_date !== '' &&
    periodForm.end_date !== '' &&
    periodForm.end_date > periodForm.start_date

  // Validation for Add Allocation
  const isAllocationFormValid =
    allocationForm.campaign_id.trim() !== '' &&
    allocationForm.campaign_name.trim() !== '' &&
    allocationForm.amount > 0

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          Budget Management
        </Typography>
      </Box>

      {/* Tabs */}
      <Card>
        <Tabs value={tabValue} onChange={(_, v) => setTabValue(v)} sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tab label="Budget Periods" />
          <Tab label="Transactions" />
          <Tab label="Alerts" />
        </Tabs>

        {/* Periods Tab */}
        <TabPanel value={tabValue} index={0}>
          <Box sx={{ mb: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">Budget Periods</Typography>
            <Box sx={{ display: 'flex', gap: 1 }}>
              {selectedPeriod && (
                <Button variant="outlined" startIcon={<AttachMoneyIcon />} onClick={() => setOpenAllocationDialog(true)}>
                  Add Allocation
                </Button>
              )}
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpenPeriodDialog(true)}>
                New Period
              </Button>
            </Box>
          </Box>

          {periodsLoading ? <Skeleton height={200} /> : periodsError ? (
            <Box sx={{ py: 8, textAlign: 'center' }}>
              <ErrorOutlineIcon sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
              <Typography variant="h6" color="error" gutterBottom>Failed to load budget periods</Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                {periodsError instanceof Error ? periodsError.message : 'An unexpected error occurred'}
              </Typography>
              <Button variant="contained" onClick={() => refetchPeriods()}>Retry</Button>
            </Box>
          ) : (
            <>
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell>Start Date</TableCell>
                        <TableCell>End Date</TableCell>
                        <TableCell>Total Budget</TableCell>
                        <TableCell>Allocated</TableCell>
                        <TableCell>Remaining</TableCell>
                        <TableCell>Status</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {periods.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={7} align="center">
                            <EmptyState icon={AccountBalanceIcon} title="No budget periods yet" description="Create a budget period to start planning your marketing spend." actionLabel="New Budget Period" onAction={() => setOpenPeriodDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        periods.map((period: any) => (
                          <TableRow
                            key={period.id}
                            hover
                            selected={selectedPeriod?.id === period.id}
                            onClick={() => setSelectedPeriod(period)}
                            sx={{ cursor: 'pointer' }}
                          >
                            <TableCell><Typography fontWeight="medium">{period.name}</Typography></TableCell>
                            <TableCell>{formatDate(period.start_date)}</TableCell>
                            <TableCell>{period.end_date ? formatDate(period.end_date) : 'Ongoing'}</TableCell>
                            <TableCell>{formatCurrency(period.total_budget || 0)}</TableCell>
                            <TableCell>{formatCurrency(period.allocated || 0)}</TableCell>
                            <TableCell>{formatCurrency((period.total_budget || 0) - (period.allocated || 0))}</TableCell>
                            <TableCell>
                              <Chip label={period.status || 'active'} size="small"
                                color={period.status === 'active' ? 'success' : period.status === 'closed' ? 'default' : 'warning'}
                              />
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
            </>
          )}
        </TabPanel>

        {/* Transactions Tab */}
        <TabPanel value={tabValue} index={1}>
          <Box sx={{ mb: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">Transactions</Typography>
            <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpenTransactionDialog(true)}>
              Record Spend
            </Button>
          </Box>

          {transactionsLoading ? <Skeleton height={200} /> : transactionsError ? (
            <Box sx={{ py: 8, textAlign: 'center' }}>
              <ErrorOutlineIcon sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
              <Typography variant="h6" color="error" gutterBottom>Failed to load transactions</Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                {transactionsError instanceof Error ? transactionsError.message : 'An unexpected error occurred'}
              </Typography>
              <Button variant="contained" onClick={() => refetchTransactions()}>Retry</Button>
            </Box>
          ) : (
            <>
              <Box sx={{ overflowX: 'auto' }}>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Date</TableCell>
                        <TableCell>Campaign</TableCell>
                        <TableCell>Period</TableCell>
                        <TableCell>Type</TableCell>
                        <TableCell>Amount</TableCell>
                        <TableCell>Description</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {transactions.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={6} align="center">
                            <EmptyState icon={ReceiptLongIcon} title="No transactions recorded yet" description="Record your first spend or refund to track budget usage." actionLabel="Record Spend" onAction={() => setOpenTransactionDialog(true)} />
                          </TableCell>
                        </TableRow>
                      ) : (
                        transactions.slice(txnPage * txnPageSize, (txnPage + 1) * txnPageSize).map((tx: any) => (
                          <TableRow key={tx.id} hover>
                            <TableCell>{tx.date ? formatDate(tx.date) : formatDate(tx.created_at)}</TableCell>
                            <TableCell>{tx.campaign_name || tx.campaign_id || '-'}</TableCell>
                            <TableCell>{tx.period_name || tx.period_id || '-'}</TableCell>
                            <TableCell>
                              <Chip label={tx.type || 'spend'} size="small"
                                color={tx.type === 'refund' ? 'success' : 'warning'}
                              />
                            </TableCell>
                            <TableCell>
                              <Typography fontWeight="medium" color={tx.type === 'refund' ? 'success.main' : 'error.main'}>
                                {tx.type === 'refund' ? '+' : '-'}{formatCurrency(tx.amount || 0)}
                              </Typography>
                            </TableCell>
                            <TableCell>{tx.description || '-'}</TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
              {transactions.length > txnPageSize && (
                <Box sx={{ display: 'flex', justifyContent: 'center', pt: 2 }}>
                  <Pagination
                    count={Math.ceil(transactions.length / txnPageSize)}
                    page={txnPage + 1}
                    onChange={(_, p) => setTxnPage(p - 1)}
                  />
                </Box>
              )}
            </>
          )}
        </TabPanel>

        {/* Alerts Tab */}
        <TabPanel value={tabValue} index={2}>
          <Typography variant="h6" sx={{ mb: 2 }}>Budget Alerts</Typography>

          {alertsLoading ? <Skeleton height={200} /> : alertsError ? (
            <Box sx={{ py: 8, textAlign: 'center' }}>
              <ErrorOutlineIcon sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
              <Typography variant="h6" color="error" gutterBottom>Failed to load alerts</Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                {alertsError instanceof Error ? alertsError.message : 'An unexpected error occurred'}
              </Typography>
              <Button variant="contained" onClick={() => refetchAlerts()}>Retry</Button>
            </Box>
          ) : (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              {alerts.length === 0 ? (
                <EmptyState icon={NotificationsIcon} title="No budget alerts triggered" description="Alerts will appear here when budget thresholds are exceeded." />
              ) : (
                alerts.map((alert: any) => (
                  <Alert key={alert.id} severity={getSeverityColor(alert.severity || 'info')} sx={{ p: 2 }}>
                    <AlertTitle>{alert.title || 'Budget Alert'}</AlertTitle>
                    {alert.message}
                    {alert.period_name && <Typography variant="caption" sx={{ display: 'block', mt: 1 }}>Period: {alert.period_name}</Typography>}
                    {alert.created_at && <Typography variant="caption">{new Date(alert.created_at).toLocaleString()}</Typography>}
                  </Alert>
                ))
              )}
            </Box>
          )}
        </TabPanel>
      </Card>

      {/* Create Period Dialog */}
      <Dialog open={openPeriodDialog} onClose={() => setOpenPeriodDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create Budget Period</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Period Name"
              fullWidth
              value={periodForm.name}
              onChange={(e) => setPeriodForm({ ...periodForm, name: e.target.value })}
              required
              placeholder="e.g., Q2 2025"
            />
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField label="Start Date" fullWidth type="date" value={periodForm.start_date}
                  onChange={(e) => setPeriodForm({ ...periodForm, start_date: e.target.value })}
                  InputLabelProps={{ shrink: true }} required />
              </Grid>
              <Grid item xs={6}>
                <TextField label="End Date" fullWidth type="date" value={periodForm.end_date}
                  onChange={(e) => setPeriodForm({ ...periodForm, end_date: e.target.value })}
                  InputLabelProps={{ shrink: true }} />
              </Grid>
            </Grid>
            <TextField label="Total Budget" fullWidth type="number" value={periodForm.total_budget}
              onChange={(e) => setPeriodForm({ ...periodForm, total_budget: Number(e.target.value) })} required />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenPeriodDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => createPeriodMutation.mutate(periodForm)}
            disabled={!isPeriodFormValid || createPeriodMutation.isPending}>
            {createPeriodMutation.isPending ? 'Creating...' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Create Allocation Dialog */}
      <Dialog open={openAllocationDialog} onClose={() => setOpenAllocationDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Add Allocation {selectedPeriod ? `for ${selectedPeriod.name}` : ''}</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField label="Campaign ID" fullWidth value={allocationForm.campaign_id}
              onChange={(e) => setAllocationForm({ ...allocationForm, campaign_id: e.target.value })} required />
            <TextField label="Campaign Name" fullWidth value={allocationForm.campaign_name}
              onChange={(e) => setAllocationForm({ ...allocationForm, campaign_name: e.target.value })} required />
            <TextField label="Amount" fullWidth type="number" value={allocationForm.amount}
              onChange={(e) => setAllocationForm({ ...allocationForm, amount: Number(e.target.value) })} required />
            <TextField label="Channel" fullWidth value={allocationForm.channel}
              onChange={(e) => setAllocationForm({ ...allocationForm, channel: e.target.value })} />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAllocationDialog(false)}>Cancel</Button>
          <Button variant="contained"
            onClick={() => selectedPeriod && createAllocationMutation.mutate({ periodId: selectedPeriod.id, data: allocationForm })}
            disabled={!isAllocationFormValid || createAllocationMutation.isPending || !selectedPeriod}>
            {createAllocationMutation.isPending ? 'Creating...' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Record Transaction Dialog */}
      <Dialog open={openTransactionDialog} onClose={() => setOpenTransactionDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Record Transaction</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <FormControl fullWidth>
              <InputLabel>Period</InputLabel>
              <Select value={transactionForm.period_id} label="Period"
                onChange={(e) => setTransactionForm({ ...transactionForm, period_id: e.target.value })}>
                {periods.map((p: any) => (
                  <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>
                ))}
              </Select>
            </FormControl>
            <TextField label="Campaign ID" fullWidth value={transactionForm.campaign_id}
              onChange={(e) => setTransactionForm({ ...transactionForm, campaign_id: e.target.value })} />
            <TextField label="Amount" fullWidth type="number" value={transactionForm.amount}
              onChange={(e) => setTransactionForm({ ...transactionForm, amount: Number(e.target.value) })} required />
            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select value={transactionForm.type} label="Type"
                onChange={(e) => setTransactionForm({ ...transactionForm, type: e.target.value })}>
                <MenuItem value="spend">Spend</MenuItem>
                <MenuItem value="refund">Refund</MenuItem>
                <MenuItem value="adjustment">Adjustment</MenuItem>
              </Select>
            </FormControl>
            <TextField label="Description" fullWidth value={transactionForm.description}
              onChange={(e) => setTransactionForm({ ...transactionForm, description: e.target.value })} multiline rows={2} />
            <TextField label="Date" fullWidth type="date" value={transactionForm.date}
              onChange={(e) => setTransactionForm({ ...transactionForm, date: e.target.value })}
              InputLabelProps={{ shrink: true }} required />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenTransactionDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => recordTransactionMutation.mutate(transactionForm)}
            disabled={recordTransactionMutation.isPending}>
            {recordTransactionMutation.isPending ? 'Recording...' : 'Record'}
          </Button>
        </DialogActions>
      </Dialog>

      <ConfirmDialog {...confirmDialog.dialogProps} />
    </Box>
  )
}
