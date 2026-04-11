'use client'

import { Box, Typography, Button, Paper } from '@mui/material'
import Link from 'next/link'
import SearchOffIcon from '@mui/icons-material/SearchOff'

export default function NotFound() {
  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'background.default',
        p: 2,
      }}
    >
      <Paper
        elevation={3}
        sx={{
          p: 6,
          maxWidth: 480,
          width: '100%',
          textAlign: 'center',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: 2,
        }}
      >
        <SearchOffIcon sx={{ fontSize: 80, color: 'text.secondary' }} />
        <Typography variant="h1" sx={{ fontSize: 72, fontWeight: 'bold', lineHeight: 1 }}>
          404
        </Typography>
        <Typography variant="h5" fontWeight="bold">
          Page Not Found
        </Typography>
        <Typography variant="body1" color="text.secondary">
          The page you are looking for does not exist or has been moved.
        </Typography>
        <Button
          component={Link}
          href="/"
          variant="contained"
          size="large"
          sx={{ mt: 2 }}
        >
          Back to Dashboard
        </Button>
      </Paper>
    </Box>
  )
}
