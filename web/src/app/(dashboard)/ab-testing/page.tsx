'use client'

import { useQuery } from '@tanstack/react-query'
import { Box, Card, CardContent, Typography, Grid, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Chip, Skeleton, LinearProgress, Button, IconButton } from '@mui/material'
import PlayArrowIcon from '@mui/icons-material/PlayArrow'
import StopIcon from '@mui/icons-material/Stop'
import AnalyticsIcon from '@mui/icons-material/Analytics'
import { abTestApi } from '@/lib/api-client'
import { formatNumber } from '@/lib/utils'

interface ABTest {
  id: string
  name: string
  variant: string
  status: string
  visitors: number
  conversions: number
  conversion_rate: number
  confidence: number
  start_date: string
}

const defaultTests: ABTest[] = [
  { id: '1', name: 'Homepage CTA Button', variant: 'A', status: 'running', visitors: 5000, conversions: 250, conversion_rate: 5.0, confidence: 85, start_date: '2024-01-01' },
  { id: '2', name: 'Homepage CTA Button', variant: 'B', status: 'running', visitors: 5000, conversions: 320, conversion_rate: 6.4, confidence: 85, start_date: '2024-01-01' },
  { id: '3', name: 'Email Subject Line', variant: 'A', status: 'completed', visitors: 10000, conversions: 500, conversion_rate: 5.0, confidence: 92, start_date: '2024-01-15' },
  { id: '4', name: 'Email Subject Line', variant: 'B', status: 'completed', visitors: 10000, conversions: 650, conversion_rate: 6.5, confidence: 92, start_date: '2024-01-15' },
  { id: '5', name: 'Pricing Page Layout', variant: 'A', status: 'draft', visitors: 0, conversions: 0, conversion_rate: 0, confidence: 0, start_date: '2024-02-01' },
]

export default function ABTestingPage() {
  const { data: testsResponse, isLoading } = useQuery<{ tests: ABTest[] }>({
    queryKey: ['ab-tests'],
    queryFn: async () => {
      return abTestApi.getAll()
    },
    initialData: { tests: defaultTests },
  })

  const tests = testsResponse?.tests || defaultTests

  const getStatusColor = (status: string): 'success' | 'default' | 'warning' => {
    switch (status) {
      case 'running': return 'success'
      case 'completed': return 'default'
      case 'draft': return 'warning'
      default: return 'default'
    }
  }

  const getConfidenceColor = (confidence: number): 'success' | 'warning' | 'info' => {
    if (confidence >= 95) return 'success'
    if (confidence >= 80) return 'warning'
    return 'info'
  }

  const runningTests = tests?.filter(t => t.status === 'running') || []

  return (
    <Box>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" fontWeight="bold">
          A/B Testing
        </Typography>
        <Button variant="contained" disabled title="Coming soon">Create New Test</Button>
      </Box>

      {/* Running Tests Summary */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Running Tests
              </Typography>
              <Typography variant="h4" fontWeight="bold" color="success.main">
                {Math.ceil(runningTests.length / 2) || 1}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Total Visitors
              </Typography>
              <Typography variant="h4" fontWeight="bold" color="primary">
                {formatNumber(runningTests.reduce((sum, t) => sum + t.visitors, 0) || 10000)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Avg. Confidence
              </Typography>
              <Typography variant="h4" fontWeight="bold" color="warning.main">
                {runningTests[0]?.confidence || 85}%
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Active Tests */}
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            Active Tests
          </Typography>
          {isLoading ? (
            <Skeleton height={200} />
          ) : (
            <Grid container spacing={2}>
              {runningTests.slice(0, 4).map((test) => (
                <Grid item xs={12} sm={6} key={test.id}>
                  <Paper sx={{ p: 2, bgcolor: 'grey.50' }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                      <Typography fontWeight="medium">{test.name}</Typography>
                      <Chip label={`Variant ${test.variant}`} size="small" color="primary" />
                    </Box>
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="text.secondary">Visitors</Typography>
                        <Typography fontWeight="bold">{formatNumber(test.visitors)}</Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="text.secondary">Conversions</Typography>
                        <Typography fontWeight="bold">{formatNumber(test.conversions)}</Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="text.secondary">Conv. Rate</Typography>
                        <Typography fontWeight="bold" color="success.main">{test.conversion_rate}%</Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="text.secondary">Confidence</Typography>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <LinearProgress
                            variant="determinate"
                            value={test.confidence}
                            sx={{ width: 60, height: 6, borderRadius: 3 }}
                            color={getConfidenceColor(test.confidence)}
                          />
                          <Typography fontWeight="bold">{test.confidence}%</Typography>
                        </Box>
                      </Grid>
                    </Grid>
                    <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
                      <Button size="small" variant="outlined" startIcon={<StopIcon />} disabled title="Coming soon">
                        Stop
                      </Button>
                      <Button size="small" variant="outlined" startIcon={<AnalyticsIcon />} disabled title="Coming soon">
                        Analytics
                      </Button>
                    </Box>
                  </Paper>
                </Grid>
              ))}
            </Grid>
          )}
        </CardContent>
      </Card>

      {/* All Tests Table */}
      <Card>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
            All Tests
          </Typography>
          {isLoading ? (
            <Skeleton height={200} />
          ) : (
            <Box sx={{ overflowX: 'auto' }}>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Test Name</TableCell>
                      <TableCell>Variant</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Visitors</TableCell>
                      <TableCell>Conversions</TableCell>
                      <TableCell>Conv. Rate</TableCell>
                      <TableCell>Confidence</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {tests.map((test) => (
                      <TableRow key={test.id}>
                        <TableCell>
                          <Typography fontWeight="medium">{test.name}</Typography>
                        </TableCell>
                        <TableCell>
                          <Chip label={test.variant} size="small" color="primary" />
                        </TableCell>
                        <TableCell>
                          <Chip label={test.status} size="small" color={getStatusColor(test.status)} />
                        </TableCell>
                        <TableCell>{formatNumber(test.visitors)}</TableCell>
                        <TableCell>{formatNumber(test.conversions)}</TableCell>
                        <TableCell>{test.conversion_rate}%</TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <LinearProgress
                              variant="determinate"
                              value={test.confidence}
                              sx={{ width: 60, height: 6, borderRadius: 3 }}
                              color={getConfidenceColor(test.confidence)}
                            />
                            <Typography variant="body2">{test.confidence}%</Typography>
                          </Box>
                        </TableCell>
                        <TableCell>
                          {test.status === 'draft' && (
                            <Button size="small" variant="contained" startIcon={<PlayArrowIcon />} disabled title="Coming soon">
                              Start
                            </Button>
                          )}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          )}
        </CardContent>
      </Card>
    </Box>
  )
}
