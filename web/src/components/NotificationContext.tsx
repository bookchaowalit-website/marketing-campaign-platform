'use client'
import { createContext, useContext, useState, useCallback, ReactNode } from 'react'

export interface Notification {
  id: string
  type: 'info' | 'success' | 'warning' | 'error'
  title: string
  message: string
  read: boolean
  createdAt: string
  actionUrl?: string
  actionLabel?: string
}

interface NotificationContextType {
  notifications: Notification[]
  unreadCount: number
  markRead: (id: string) => void
  markAllRead: () => void
  addNotification: (n: Omit<Notification, 'id' | 'read' | 'createdAt'>) => void
}

const NotificationContext = createContext<NotificationContextType>({
  notifications: [],
  unreadCount: 0,
  markRead: () => {},
  markAllRead: () => {},
  addNotification: () => {},
})

// Seed with demo notifications
const demoNotifications: Notification[] = [
  {
    id: '1',
    type: 'info',
    title: 'Campaign Approved',
    message: 'Your campaign "Summer Sale 2024" has been approved and is now active.',
    read: false,
    createdAt: new Date(Date.now() - 1000 * 60 * 30).toISOString(), // 30 min ago
    actionUrl: '/unified-campaigns',
    actionLabel: 'View Campaign',
  },
  {
    id: '2',
    type: 'warning',
    title: 'Budget Alert',
    message: 'Google Ads budget is at 85% of allocated amount for this month.',
    read: false,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(), // 2 hours ago
    actionUrl: '/budget',
    actionLabel: 'View Budget',
  },
  {
    id: '3',
    type: 'success',
    title: 'Email Campaign Sent',
    message: 'Newsletter Q2 was successfully sent to 2,450 subscribers.',
    read: false,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 5).toISOString(),
    actionUrl: '/email-campaigns',
    actionLabel: 'View Results',
  },
  {
    id: '4',
    type: 'error',
    title: 'SMS Delivery Failed',
    message: '12 messages failed to deliver in the "Flash Sale" SMS campaign.',
    read: true,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
    actionUrl: '/sms',
    actionLabel: 'View Details',
  },
  {
    id: '5',
    type: 'info',
    title: 'Approval Request',
    message: 'New campaign "Holiday Promo" is waiting for your approval.',
    read: true,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 48).toISOString(),
    actionUrl: '/approvals',
    actionLabel: 'Review',
  },
]

export function NotificationProvider({ children }: { children: ReactNode }) {
  const [notifications, setNotifications] = useState<Notification[]>(demoNotifications)

  const unreadCount = notifications.filter(n => !n.read).length
  const markRead = useCallback((id: string) => {
    setNotifications(prev => prev.map(n => n.id === id ? { ...n, read: true } : n))
  }, [])
  const markAllRead = useCallback(() => {
    setNotifications(prev => prev.map(n => ({ ...n, read: true })))
  }, [])
  const addNotification = useCallback((n: Omit<Notification, 'id' | 'read' | 'createdAt'>) => {
    const newN: Notification = {
      ...n,
      id: Date.now().toString(),
      read: false,
      createdAt: new Date().toISOString(),
    }
    setNotifications(prev => [newN, ...prev])
  }, [])

  return (
    <NotificationContext.Provider value={{ notifications, unreadCount, markRead, markAllRead, addNotification }}>
      {children}
    </NotificationContext.Provider>
  )
}

export function useNotifications() {
  return useContext(NotificationContext)
}
