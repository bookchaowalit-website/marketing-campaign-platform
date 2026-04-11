'use client'

import { useEffect, useState } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { useSearch } from '@/components/SearchContext'
import GlobalSearch from '@/components/GlobalSearch'
import Sidebar from '@/components/Sidebar'
import NotificationBell from '@/components/NotificationBell'
import { Box, CircularProgress, IconButton, Drawer, AppBar, Toolbar, Typography, Breadcrumbs, Link } from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import HomeIcon from '@mui/icons-material/Home'
import SearchIcon from '@mui/icons-material/Search'
import { useTheme, useMediaQuery } from '@mui/material'

const SIDEBAR_WIDTH = 240

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const { user, isLoading } = useAuth()
  const { setOpen } = useSearch()
  const router = useRouter()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))
  const [mobileOpen, setMobileOpen] = useState(false)
  const pathname = usePathname()

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login')
    }
  }, [user, isLoading, router])

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen)
  }

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '100vh' }}>
        <CircularProgress />
      </Box>
    )
  }

  if (!user) {
    return null
  }

  const sidebarContent = <Sidebar onItemClick={() => isMobile && setMobileOpen(false)} />

  return (
    <Box sx={{ display: 'flex' }}>
      {/* Skip to main content link for accessibility */}
      <a
        href="#main-content"
        className="skip-to-main"
        style={{
          position: 'absolute',
          left: '-9999px',
          top: 'auto',
          width: '1px',
          height: '1px',
          overflow: 'hidden',
          zIndex: 99999,
          padding: '12px 24px',
          background: '#1976d2',
          color: '#fff',
          textDecoration: 'none',
          fontWeight: 'bold',
          borderRadius: '4px',
        }}
        onFocus={(e) => {
          e.currentTarget.style.left = '16px'
          e.currentTarget.style.top = '16px'
          e.currentTarget.style.width = 'auto'
          e.currentTarget.style.height = 'auto'
          e.currentTarget.style.overflow = 'visible'
        }}
        onBlur={(e) => {
          e.currentTarget.style.left = '-9999px'
          e.currentTarget.style.top = 'auto'
          e.currentTarget.style.width = '1px'
          e.currentTarget.style.height = '1px'
          e.currentTarget.style.overflow = 'hidden'
        }}
      >
        Skip to main content
      </a>
      {/* Mobile App Bar */}
      {isMobile && (
        <AppBar position="fixed" sx={{ zIndex: theme.zIndex.drawer + 1 }}>
          <Toolbar>
            <IconButton
              color="inherit"
              edge="start"
              onClick={handleDrawerToggle}
              aria-label="Open navigation menu"
              sx={{ mr: 2 }}
            >
              <MenuIcon />
            </IconButton>
            <Typography variant="h6" noWrap component="div" sx={{ mr: 2 }}>
              MarketPro
            </Typography>
            <IconButton
              color="inherit"
              onClick={() => setOpen(true)}
              aria-label="Open search"
            >
              <SearchIcon />
            </IconButton>
            <NotificationBell />
          </Toolbar>
        </AppBar>
      )}

      {/* Desktop Sidebar + Top Bar */}
      {!isMobile && (
        <Box
          component="nav"
          sx={{ width: { md: SIDEBAR_WIDTH }, flexShrink: { md: 0 } }}
          aria-label="Main navigation"
        >
          <Sidebar />
        </Box>
      )}

      {/* Desktop Top Bar */}
      {!isMobile && (
        <AppBar
          position="fixed"
          sx={{
            width: { md: `calc(100% - ${SIDEBAR_WIDTH}px)` },
            ml: { md: `${SIDEBAR_WIDTH}px` },
            boxShadow: 'none',
            bgcolor: 'background.paper',
            borderBottom: '1px solid',
            borderColor: 'divider',
          }}
        >
          <Toolbar sx={{ justifyContent: 'flex-end' }}>
            <NotificationBell />
          </Toolbar>
        </AppBar>
      )}

      {/* Mobile Drawer */}
      {isMobile && (
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{ keepMounted: true }}
          sx={{
            display: { xs: 'block', md: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: SIDEBAR_WIDTH },
          }}
          aria-label="Mobile navigation menu"
        >
          {sidebarContent}
        </Drawer>
      )}

      {/* Main Content */}
      <Box
        component="main"
        id="main-content"
        sx={{
          flexGrow: 1,
          p: 3,
          minHeight: '100vh',
          bgcolor: 'grey.50',
          width: { md: `calc(100% - ${SIDEBAR_WIDTH}px)` },
          mt: '64px', // Account for fixed top bar height
        }}
      >
        {/* Breadcrumbs */}
        <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
          <Link color="inherit" href="/" underline="hover" sx={{ display: 'flex', alignItems: 'center' }}>
            <HomeIcon sx={{ mr: 0.5 }} fontSize="small" />
            Home
          </Link>
          <Typography color="text.primary" sx={{ textTransform: 'capitalize' }}>
            {pathname === '/' ? 'Dashboard' : pathname.replace('/', '').replace(/-/g, ' ')}
          </Typography>
        </Breadcrumbs>
        {children}
      </Box>

      {/* Global Search Modal */}
      <GlobalSearch />
    </Box>
  )
}
