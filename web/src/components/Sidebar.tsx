'use client'

import { usePathname } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { useThemeMode } from '@/hooks/useTheme'
import { useSearch } from '@/components/SearchContext'
import Link from 'next/link'
import { Box, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Typography, Avatar, Divider, Button, IconButton, Tooltip } from '@mui/material'
import DashboardIcon from '@mui/icons-material/Dashboard'
import SearchIcon from '@mui/icons-material/Search'
import GroupIcon from '@mui/icons-material/Group'
import ScienceIcon from '@mui/icons-material/Science'
import AnalyticsIcon from '@mui/icons-material/Analytics'
import PeopleIcon from '@mui/icons-material/People'
import LinkIcon from '@mui/icons-material/Link'
import HubIcon from '@mui/icons-material/Hub'
import LogoutIcon from '@mui/icons-material/Logout'
import HelpOutlineIcon from '@mui/icons-material/HelpOutline'
import DarkModeIcon from '@mui/icons-material/DarkMode'
import LightModeIcon from '@mui/icons-material/LightMode'
import AttachMoneyIcon from '@mui/icons-material/AttachMoney'
import CalendarTodayIcon from '@mui/icons-material/CalendarToday'
import TaskAltIcon from '@mui/icons-material/TaskAlt'
import AdsClickIcon from '@mui/icons-material/AdsClick'
import SmsIcon from '@mui/icons-material/Sms'
import EmailIcon from '@mui/icons-material/Email'
import VideoLibraryIcon from '@mui/icons-material/VideoLibrary'
import SocialDistanceIcon from '@mui/icons-material/SocialDistance'
import ListAltIcon from '@mui/icons-material/ListAlt'
import SettingsIcon from '@mui/icons-material/Settings'
import ListSubheader from '@mui/material/ListSubheader'

const navGroups = [
  {
    subheader: 'Overview',
    items: [
      { href: '/', icon: <DashboardIcon />, label: 'Dashboard' },
    ],
  },
  {
    subheader: 'Campaigns',
    items: [
      { href: '/unified-campaigns', icon: <HubIcon />, label: 'All Campaigns' },
      { href: '/campaigns', icon: <ListAltIcon />, label: 'Classic Campaigns' },
      { href: '/ab-testing', icon: <ScienceIcon />, label: 'A/B Testing' },
    ],
  },
  {
    subheader: 'Channels',
    items: [
      { href: '/marketing', icon: <HubIcon />, label: 'Marketing Hub' },
      { href: '/ppc', icon: <AdsClickIcon />, label: 'PPC / Paid Ads' },
      { href: '/sms', icon: <SmsIcon />, label: 'SMS Marketing' },
      { href: '/email-campaigns', icon: <EmailIcon />, label: 'Email Campaigns' },
      { href: '/video', icon: <VideoLibraryIcon />, label: 'Video Marketing' },
      { href: '/seo', icon: <SearchIcon />, label: 'SEO & SEM' },
      { href: '/social', icon: <SocialDistanceIcon />, label: 'Social' },
    ],
  },
  {
    subheader: 'Audience',
    items: [
      { href: '/competitors', icon: <GroupIcon />, label: 'Competitors' },
      { href: '/affiliates', icon: <PeopleIcon />, label: 'Affiliates' },
    ],
  },
  {
    subheader: 'Planning',
    items: [
      { href: '/budget', icon: <AttachMoneyIcon />, label: 'Budget' },
      { href: '/calendar', icon: <CalendarTodayIcon />, label: 'Calendar' },
      { href: '/approvals', icon: <TaskAltIcon />, label: 'Approvals' },
    ],
  },
  {
    subheader: 'Analytics',
    items: [
      { href: '/analytics', icon: <AnalyticsIcon />, label: 'Analytics' },
      { href: '/integrations', icon: <LinkIcon />, label: 'Integrations' },
    ],
  },
  {
    subheader: 'Settings',
    items: [
      { href: '/settings', icon: <SettingsIcon />, label: 'Settings' },
    ],
  },
]

const drawerWidth = 240

interface SidebarProps {
  onItemClick?: () => void
}

export default function Sidebar({ onItemClick }: SidebarProps) {
  const pathname = usePathname()
  const { user, logout } = useAuth()
  const { mode, toggleTheme } = useThemeMode()
  const { setOpen } = useSearch()

  // Don't show sidebar on auth pages
  if (pathname === '/login' || pathname === '/signup') {
    return null
  }

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: drawerWidth,
          boxSizing: 'border-box',
        },
      }}
    >
      <Box sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
        <Typography variant="h6" color="primary" fontWeight="bold">
          MarketPro
        </Typography>
      </Box>

      {/* Search Button */}
      <Box sx={{ px: 2, pb: 1 }}>
        <Button
          fullWidth
          variant="outlined"
          size="small"
          startIcon={<SearchIcon />}
          onClick={() => setOpen(true)}
          sx={{
            justifyContent: 'flex-start',
            color: 'text.secondary',
            borderColor: 'divider',
            px: 1,
            py: 0.5,
          }}
        >
          <Typography variant="body2" sx={{ flex: 1, textAlign: 'left' }}>
            Search...
          </Typography>
          <Box
            component="kbd"
            sx={{
              px: 0.5,
              py: 0.25,
              borderRadius: 0.5,
              bgcolor: 'action.hover',
              fontSize: '0.65rem',
              fontWeight: 600,
              lineHeight: 1,
            }}
          >
            ⌘K
          </Box>
        </Button>
      </Box>

      <Divider />

      <List sx={{ px: 1 }}>
        {navGroups.map((group) => (
          <Box key={group.subheader}>
            <ListSubheader sx={{ lineHeight: 2, bgcolor: 'transparent', fontSize: '0.7rem' }} disableSticky>
              {group.subheader.toUpperCase()}
            </ListSubheader>
            {group.items.map((item) => (
              <ListItem key={item.href} disablePadding>
                <ListItemButton
                  component={Link}
                  href={item.href}
                  selected={pathname === item.href}
                  onClick={onItemClick}
                  aria-label={`Navigate to ${item.label}`}
                  sx={{
                    borderRadius: 1,
                    mb: 0.5,
                    '&.Mui-selected': {
                      bgcolor: 'primary.main',
                      color: 'white',
                      '& .MuiListItemIcon-root': { color: 'white' },
                      '&:hover': { bgcolor: 'primary.dark' },
                    },
                  }}
                >
                  <ListItemIcon sx={{ minWidth: 40 }}>
                    {item.icon}
                  </ListItemIcon>
                  <ListItemText primary={item.label} />
                </ListItemButton>
              </ListItem>
            ))}
          </Box>
        ))}
      </List>

      <Divider sx={{ mt: 'auto' }} />

      {/* Theme Toggle */}
      <Box sx={{ p: 1, display: 'flex', justifyContent: 'center' }}>
        <Tooltip title={mode === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}>
          <IconButton onClick={toggleTheme} color="inherit">
            {mode === 'light' ? <DarkModeIcon /> : <LightModeIcon />}
          </IconButton>
        </Tooltip>
      </Box>

      <Divider />

      <Box sx={{ p: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
          <Avatar sx={{ bgcolor: 'primary.main', width: 32, height: 32 }}>
            {user?.full_name?.[0] || user?.email?.[0] || 'U'}
          </Avatar>
          <Box>
            <Typography variant="body2" fontWeight="medium">
              {user?.full_name || user?.email?.split('@')[0] || 'User'}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {user?.role || 'user'}
            </Typography>
          </Box>
        </Box>
        <Button
          fullWidth
          variant="outlined"
          size="small"
          startIcon={<LogoutIcon />}
          onClick={logout}
        >
          Logout
        </Button>
        <Button
          fullWidth
          variant="text"
          size="small"
          startIcon={<HelpOutlineIcon />}
          href="/docs"
          component="a"
          target="_blank"
          rel="noopener noreferrer"
          sx={{ mt: 1 }}
        >
          Help
        </Button>
      </Box>
    </Drawer>
  )
}
