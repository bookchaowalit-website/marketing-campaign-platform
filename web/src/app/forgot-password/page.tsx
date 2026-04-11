'use client'

import { useState } from 'react'
import {
  Box,
  Typography,
  TextField,
  Button,
  Paper,
  Alert,
  Link as MuiLink,
  CircularProgress,
} from '@mui/material'
import Link from 'next/link'
import LockResetIcon from '@mui/icons-material/LockReset'
import { authApi } from '@/lib/api-client'

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setSuccess(false)

    if (!email.trim()) {
      setError('Please enter your email address.')
      return
    }

    setIsLoading(true)
    try {
      await authApi.forgotPassword(email)
      setSuccess(true)
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to send reset email. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

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
          p: 4,
          maxWidth: 420,
          width: '100%',
          display: 'flex',
          flexDirection: 'column',
          gap: 2,
        }}
      >
        <Box sx={{ display: 'flex', justifyContent: 'center', mb: 1 }}>
          <LockResetIcon sx={{ fontSize: 48, color: 'primary.main' }} />
        </Box>
        <Typography variant="h5" fontWeight="bold" textAlign="center">
          Forgot Password
        </Typography>
        <Typography variant="body2" color="text.secondary" textAlign="center">
          Enter your email address and we will send you a link to reset your password.
        </Typography>

        {error && (
          <Alert severity="error" onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success">
            Password reset link has been sent to your email. Please check your inbox.
          </Alert>
        )}

        {!success && (
          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="Email Address"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              sx={{ mb: 2 }}
            />
            <Button
              type="submit"
              variant="contained"
              fullWidth
              size="large"
              disabled={isLoading}
              sx={{ mb: 2 }}
            >
              {isLoading ? <CircularProgress size={24} /> : 'Send Reset Link'}
            </Button>
          </form>
        )}

        <Typography variant="body2" textAlign="center">
          Remember your password?{' '}
          <MuiLink component={Link} href="/login" sx={{ fontWeight: 'medium' }}>
            Back to Login
          </MuiLink>
        </Typography>
      </Paper>
    </Box>
  )
}
