'use client'

import { useState, useRef, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import {
  IconButton,
  Badge,
  Popover,
  Box,
  Typography,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  ListItemButton,
  Avatar,
  Button,
  Divider,
  Chip,
} from '@mui/material'
import NotificationsNoneIcon from '@mui/icons-material/NotificationsNone'
import InfoIcon from '@mui/icons-material/Info'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import WarningIcon from '@mui/icons-material/Warning'
import ErrorIcon from '@mui/icons-material/Error'
import { useNotifications, type Notification } from '@/components/NotificationContext'

function formatRelativeTime(dateString: string): string {
  const now = new Date()
  const date = new Date(dateString)
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)

  if (diffSec < 60) return 'just now'
  if (diffMin < 60) return `${diffMin} minute${diffMin > 1 ? 's' : ''} ago`
  if (diffHour < 24) return `${diffHour} hour${diffHour > 1 ? 's' : ''} ago`
  if (diffDay < 7) return `${diffDay} day${diffDay > 1 ? 's' : ''} ago`
  return date.toLocaleDateString()
}

const typeConfig: Record<Notification['type'], { color: string; bg: string; icon: React.ReactNode }> = {
  info: {
    color: '#1976d2',
    bg: '#e3f2fd',
    icon: <InfoIcon fontSize="small" />,
  },
  success: {
    color: '#2e7d32',
    bg: '#e8f5e9',
    icon: <CheckCircleIcon fontSize="small" />,
  },
  warning: {
    color: '#ed6c02',
    bg: '#fff3e0',
    icon: <WarningIcon fontSize="small" />,
  },
  error: {
    color: '#d32f2f',
    bg: '#ffebee',
    icon: <ErrorIcon fontSize="small" />,
  },
}

function NotificationItem({
  notification,
  onMarkRead,
  onClose,
}: {
  notification: Notification
  onMarkRead: (id: string) => void
  onClose: () => void
}) {
  const router = useRouter()
  const config = typeConfig[notification.type]

  const handleClick = () => {
    onMarkRead(notification.id)
    if (notification.actionUrl) {
      router.push(notification.actionUrl)
      onClose()
    }
  }

  return (
    <ListItem
      disablePadding
      sx={{
        bgcolor: notification.read ? 'transparent' : 'action.hover',
        borderBottom: '1px solid',
        borderColor: 'divider',
      }}
    >
      <ListItemButton onClick={handleClick}>
        <ListItemAvatar>
          <Avatar
            sx={{
              bgcolor: config.bg,
              color: config.color,
              width: 36,
              height: 36,
            }}
          >
            {config.icon}
          </Avatar>
        </ListItemAvatar>
        <ListItemText
          primary={
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap' }}>
              <Typography variant="subtitle2" sx={{ fontWeight: notification.read ? 400 : 600 }}>
                {notification.title}
              </Typography>
              {!notification.read && (
                <Chip
                  label="New"
                  size="small"
                  sx={{
                    height: 18,
                    fontSize: '0.65rem',
                    bgcolor: 'primary.main',
                    color: 'white',
                    fontWeight: 600,
                  }}
                />
              )}
            </Box>
          }
          secondary={
            <>
              <Typography
                variant="body2"
                color="text.secondary"
                sx={{
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  display: '-webkit-box',
                  WebkitLineClamp: 2,
                  WebkitBoxOrient: 'vertical',
                  my: 0.5,
                }}
              >
                {notification.message}
              </Typography>
              <Typography variant="caption" color="text.disabled" sx={{ display: 'block' }}>
                {formatRelativeTime(notification.createdAt)}
              </Typography>
            </>
          }
        />
      </ListItemButton>
      {notification.actionUrl && notification.actionLabel && (
        <Box sx={{ p: 1, pt: 0 }}>
          <Button
            size="small"
            onClick={(e) => {
              e.stopPropagation()
              onMarkRead(notification.id)
              router.push(notification.actionUrl!)
              onClose()
            }}
          >
            {notification.actionLabel}
          </Button>
        </Box>
      )}
    </ListItem>
  )
}

export default function NotificationBell() {
  const { notifications, unreadCount, markRead, markAllRead } = useNotifications()
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null)
  const anchorRef = useRef<HTMLButtonElement>(null)

  const open = Boolean(anchorEl)

  const handleOpen = useCallback((event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget)
  }, [])

  const handleClose = useCallback(() => {
    setAnchorEl(null)
  }, [])

  return (
    <>
      <IconButton
        ref={anchorRef}
        color="inherit"
        onClick={handleOpen}
        aria-label="Open notifications"
        aria-controls={open ? 'notification-popover' : undefined}
        aria-haspopup="true"
        aria-expanded={open ? 'true' : undefined}
      >
        <Badge badgeContent={unreadCount > 0 ? unreadCount : undefined} color="error">
          <NotificationsNoneIcon />
        </Badge>
      </IconButton>

      <Popover
        id="notification-popover"
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        PaperProps={{
          sx: {
            width: { xs: 'calc(100vw - 32px)', sm: 400 },
            maxHeight: 'calc(100vh - 120px)',
          },
        }}
      >
        {/* Header */}
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            px: 2,
            py: 1.5,
          }}
        >
          <Typography variant="h6" sx={{ fontSize: '1rem' }}>
            Notifications
          </Typography>
          {unreadCount > 0 && (
            <Button size="small" onClick={markAllRead} sx={{ textTransform: 'none' }}>
              Mark all as read
            </Button>
          )}
        </Box>
        <Divider />

        {/* Notification List */}
        <List
          sx={{
            maxHeight: 360,
            overflow: 'auto',
            p: 0,
          }}
        >
          {notifications.length === 0 ? (
            <Box sx={{ py: 4, textAlign: 'center' }}>
              <NotificationsNoneIcon sx={{ fontSize: 48, color: 'text.disabled', mb: 1 }} />
              <Typography color="text.secondary">No notifications yet</Typography>
            </Box>
          ) : (
            notifications.map((notification) => (
              <NotificationItem
                key={notification.id}
                notification={notification}
                onMarkRead={markRead}
                onClose={handleClose}
              />
            ))
          )}
        </List>

        {/* Footer */}
        <Divider />
        <Box sx={{ p: 1.5, textAlign: 'center' }}>
          <Button size="small" onClick={handleClose} sx={{ textTransform: 'none' }}>
            View all
          </Button>
        </Box>
      </Popover>
    </>
  )
}
