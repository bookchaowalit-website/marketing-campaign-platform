'use client'

import { useState } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  IconButton,
  CircularProgress,
} from '@mui/material'
import DeleteIcon from '@mui/icons-material/Delete'
import WarningIcon from '@mui/icons-material/Warning'

interface ConfirmDialogProps {
  open: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  onConfirm: () => void
  onCancel: () => void
  isLoading?: boolean
  severity?: 'error' | 'warning' | 'info'
}

export function ConfirmDialog({
  open,
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  onConfirm,
  onCancel,
  isLoading = false,
  severity = 'warning',
}: ConfirmDialogProps) {
  const colors = {
    error: { bg: '#fef2f2', color: '#dc2626', button: 'error' as const },
    warning: { bg: '#fffbeb', color: '#d97706', button: 'warning' as const },
    info: { bg: '#eff6ff', color: '#2563eb', button: 'info' as const },
  }

  return (
    <Dialog open={open} onClose={onCancel} maxWidth="sm" fullWidth>
      <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <WarningIcon sx={{ color: colors[severity].color }} />
        {title}
      </DialogTitle>
      <DialogContent>
        <DialogContentText>{message}</DialogContentText>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button onClick={onCancel} disabled={isLoading}>
          {cancelText}
        </Button>
        <Button
          variant="contained"
          color={colors[severity].button}
          onClick={onConfirm}
          disabled={isLoading}
          startIcon={isLoading ? <CircularProgress size={16} /> : null}
        >
          {confirmText}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

// Hook for easier usage
export function useConfirmDialog() {
  const [state, setState] = useState<{
    open: boolean
    title: string
    message: string
    onConfirm: (() => void) | null
    isLoading: boolean
  }>({
    open: false,
    title: '',
    message: '',
    onConfirm: null,
    isLoading: false,
  })

  const showConfirm = (title: string, message: string, onConfirm: () => void) => {
    setState({ open: true, title, message, onConfirm, isLoading: false })
  }

  const handleConfirm = async () => {
    if (state.onConfirm) {
      setState((prev) => ({ ...prev, isLoading: true }))
      try {
        await state.onConfirm()
      } finally {
        setState((prev) => ({ ...prev, open: false, isLoading: false }))
      }
    }
  }

  const handleCancel = () => {
    setState((prev) => ({ ...prev, open: false, onConfirm: null }))
  }

  return {
    showConfirm,
    dialogProps: {
      open: state.open,
      title: state.title,
      message: state.message,
      onConfirm: handleConfirm,
      onCancel: handleCancel,
      isLoading: state.isLoading,
    },
  }
}