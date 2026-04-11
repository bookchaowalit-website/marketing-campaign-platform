'use client'

import { useQuery } from '@tanstack/react-query'
import {
  Card,
  CardContent,
  Typography,
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  IconButton,
  Button,
  LinearProgress,
} from '@mui/material'
import { Add, Edit, Delete, PlayArrow, Pause } from '@mui/icons-material'
import { campaignApi } from '@/lib/api-client'
import type { Campaign } from '@bookmarketing/shared-types'

export function CampaignList() {
  const { data: campaignsResponse, isLoading } = useQuery({
    queryKey: ['campaigns'],
    queryFn: () => campaignApi.getAll(50, 0),
  })

  const campaigns = campaignsResponse?.campaigns || []

  const getStatusColor = (status: Campaign['status']): 'success' | 'warning' | 'default' | 'info' => {
    const colors: Record<Campaign['status'], 'success' | 'warning' | 'default' | 'info'> = {
      active: 'success',
      paused: 'warning',
      completed: 'default',
      draft: 'info',
    }
    return colors[status]
  }

  const getTypeColor = (type: Campaign['type']): 'info' | 'primary' | 'secondary' | 'success' | 'warning' => {
    const colors: Record<Campaign['type'], 'info' | 'primary' | 'secondary' | 'success' | 'warning'> = {
      SEO: 'info',
      SEM: 'primary',
      Social: 'secondary',
      Email: 'success',
      Display: 'warning',
    }
    return colors[type]
  }

  if (isLoading) {
    return <LinearProgress />
  }

  return (
    <Card>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h6" fontWeight="bold">
            Marketing Campaigns
          </Typography>
          <Button variant="contained" startIcon={<Add />}>
            New Campaign
          </Button>
        </Box>

        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Campaign Name</TableCell>
                <TableCell>Type</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Budget</TableCell>
                <TableCell>Spent</TableCell>
                <TableCell>Progress</TableCell>
                <TableCell>Conversions</TableCell>
                <TableCell>CTR</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {campaigns.length > 0 ? campaigns.map((campaign) => {
                const progress = campaign.budget > 0 ? (campaign.spent / campaign.budget) * 100 : 0
                const ctr = campaign.impressions > 0 ? ((campaign.clicks / campaign.impressions) * 100).toFixed(2) : '0.00'
                const conversionRate = campaign.clicks > 0 ? ((campaign.conversions / campaign.clicks) * 100).toFixed(2) : '0.00'

                return (
                  <TableRow key={campaign.id}>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {campaign.name}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={campaign.type}
                        size="small"
                        color={getTypeColor(campaign.type)}
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={campaign.status}
                        size="small"
                        color={getStatusColor(campaign.status)}
                      />
                    </TableCell>
                    <TableCell>${campaign.budget.toLocaleString()}</TableCell>
                    <TableCell>${campaign.spent.toLocaleString()}</TableCell>
                    <TableCell sx={{ minWidth: 150 }}>
                      <LinearProgress
                        variant="determinate"
                        value={Math.min(progress, 100)}
                        sx={{ height: 6, borderRadius: 3 }}
                        color={progress > 90 ? 'error' : progress > 70 ? 'warning' : 'success'}
                      />
                      <Typography variant="caption" color="text.secondary">
                        {progress.toFixed(0)}%
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {campaign.conversions}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        ({conversionRate}%)
                      </Typography>
                    </TableCell>
                    <TableCell>{ctr}%</TableCell>
                    <TableCell>
                      <IconButton size="small" color="primary">
                        {campaign.status === 'active' ? <Pause /> : <PlayArrow />}
                      </IconButton>
                      <IconButton size="small" color="primary">
                        <Edit />
                      </IconButton>
                      <IconButton size="small" color="error">
                        <Delete />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                )
              }) : (
                <TableRow>
                  <TableCell colSpan={9} align="center">
                    <Typography color="text.secondary">No campaigns yet. Create your first campaign!</Typography>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  )
}
