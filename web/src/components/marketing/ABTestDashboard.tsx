'use client'

import { useQuery } from '@tanstack/react-query'
import {
  Card,
  CardContent,
  Typography,
  Box,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Button,
  LinearProgress,
  IconButton,
} from '@mui/material'
import { Add, PlayArrow, Pause, Stop, Analytics } from '@mui/icons-material'
import { abTestApi } from '@/lib/api-client'
import type { ABTest } from '@bookmarketing/shared-types'

export function ABTestDashboard() {
  const { data: abTestsResponse, isLoading } = useQuery({
    queryKey: ['abTests'],
    queryFn: () => abTestApi.getAll(50, 0),
  })

  const abTests = abTestsResponse?.tests || []

  const getStatusColor = (status: ABTest['status']): 'success' | 'default' | 'info' => {
    const colors: Record<ABTest['status'], 'success' | 'default' | 'info'> = {
      running: 'success',
      completed: 'default',
      draft: 'info',
    }
    return colors[status]
  }

  const getConfidenceColor = (confidence: number): 'success' | 'warning' | 'info' => {
    if (confidence >= 95) return 'success'
    if (confidence >= 80) return 'warning'
    return 'info'
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h6" fontWeight="bold">
          A/B Testing Dashboard
        </Typography>
        <Button variant="contained" startIcon={<Add />}>
          Create New Test
        </Button>
      </Box>

      {/* Active Tests Overview */}
      <Grid container spacing={3} mb={3}>
        {abTests.filter((t) => t.status === 'running').map((test) => (
          <Grid item xs={12} md={6} key={test.id}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                  <Box>
                    <Typography variant="subtitle1" fontWeight="bold">
                      {test.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Started: {test.start_date}
                    </Typography>
                  </Box>
                  <Chip label={`Variant ${test.variant}`} color="primary" size="small" />
                </Box>

                <Grid container spacing={2} mb={2}>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Visitors
                    </Typography>
                    <Typography variant="h6" fontWeight="bold">
                      {test.visitors.toLocaleString()}
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Conversions
                    </Typography>
                    <Typography variant="h6" fontWeight="bold">
                      {test.conversions.toLocaleString()}
                    </Typography>
                  </Grid>
                </Grid>

                <Box mb={2}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                    <Typography variant="body2" color="text.secondary">
                      Statistical Confidence
                    </Typography>
                    <Typography
                      variant="body2"
                      fontWeight="bold"
                      color={getConfidenceColor(test.confidence)}
                    >
                      {test.confidence}%
                    </Typography>
                  </Box>
                  <LinearProgress
                    variant="determinate"
                    value={test.confidence}
                    sx={{ height: 8, borderRadius: 4 }}
                    color={getConfidenceColor(test.confidence)}
                  />
                </Box>

                <Box sx={{ display: 'flex', gap: 1 }}>
                  <IconButton size="small" color="primary">
                    <Pause />
                  </IconButton>
                  <IconButton size="small" color="error">
                    <Stop />
                  </IconButton>
                  <IconButton size="small" color="secondary">
                    <Analytics />
                  </IconButton>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
        {abTests.filter((t) => t.status === 'running').length === 0 && (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" align="center">
                  No running A/B tests. Create a new test to get started!
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>

      {/* All Tests Table */}
      <Card>
        <CardContent>
          <Typography variant="h6" fontWeight="bold" gutterBottom>
            All Tests
          </Typography>
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
                  <TableCell>Duration</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {abTests.length > 0 ? abTests.map((test) => (
                  <TableRow key={test.id}>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {test.name}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip label={test.variant} size="small" variant="outlined" />
                    </TableCell>
                    <TableCell>
                      <Chip label={test.status} size="small" color={getStatusColor(test.status)} />
                    </TableCell>
                    <TableCell>{test.visitors.toLocaleString()}</TableCell>
                    <TableCell>{test.conversions.toLocaleString()}</TableCell>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {test.conversion_rate}%
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <LinearProgress
                          variant="determinate"
                          value={test.confidence}
                          sx={{ width: 50, height: 6, borderRadius: 3 }}
                          color={getConfidenceColor(test.confidence)}
                        />
                        <Typography variant="body2">{test.confidence}%</Typography>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" color="text.secondary">
                        {test.end_date || 'Ongoing'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      {test.status === 'running' && (
                        <IconButton size="small" color="primary">
                          <Pause />
                        </IconButton>
                      )}
                      {test.status === 'draft' && (
                        <IconButton size="small" color="success">
                          <PlayArrow />
                        </IconButton>
                      )}
                    </TableCell>
                  </TableRow>
                )) : (
                  <TableRow>
                    <TableCell colSpan={9} align="center">
                      <Typography color="text.secondary">No A/B tests yet. Create your first test!</Typography>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    </Box>
  )
}
