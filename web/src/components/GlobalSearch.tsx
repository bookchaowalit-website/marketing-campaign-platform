'use client'

import { useState, useMemo, useEffect, useCallback, useRef } from 'react'
import { useRouter } from 'next/navigation'
import {
  Dialog,
  DialogContent,
  DialogTitle,
  InputBase,
  Box,
  Typography,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Divider,
  Chip,
  Paper,
  alpha,
} from '@mui/material'
import SearchIcon from '@mui/icons-material/Search'
import DashboardIcon from '@mui/icons-material/Dashboard'
import HubIcon from '@mui/icons-material/Hub'
import AttachMoneyIcon from '@mui/icons-material/AttachMoney'
import CalendarTodayIcon from '@mui/icons-material/CalendarToday'
import TaskAltIcon from '@mui/icons-material/TaskAlt'
import AdsClickIcon from '@mui/icons-material/AdsClick'
import SmsIcon from '@mui/icons-material/Sms'
import EmailIcon from '@mui/icons-material/Email'
import VideoLibraryIcon from '@mui/icons-material/VideoLibrary'
import GroupIcon from '@mui/icons-material/Group'
import PeopleIcon from '@mui/icons-material/People'
import AnalyticsIcon from '@mui/icons-material/Analytics'
import LinkIcon from '@mui/icons-material/Link'
import AddIcon from '@mui/icons-material/Add'
import EventIcon from '@mui/icons-material/Event'
import ScienceIcon from '@mui/icons-material/Science'
import SocialDistanceIcon from '@mui/icons-material/SocialDistance'
import ListAltIcon from '@mui/icons-material/ListAlt'
import { useSearch } from './SearchContext'

// ---------- search index ----------

interface SearchItem {
  id: string
  label: string
  href?: string
  action?: string
  category: string
  icon: React.ReactNode
  keywords?: string[]
}

const searchIndex: SearchItem[] = [
  // Navigation
  { id: 'nav-dashboard', label: 'Dashboard', href: '/', category: 'Navigation', icon: <DashboardIcon />, keywords: ['home', 'overview'] },
  { id: 'nav-campaigns', label: 'All Campaigns', href: '/unified-campaigns', category: 'Navigation', icon: <HubIcon />, keywords: ['unified'] },
  { id: 'nav-classic-campaigns', label: 'Classic Campaigns', href: '/campaigns', category: 'Navigation', icon: <ListAltIcon />, keywords: ['classic'] },
  { id: 'nav-ab-testing', label: 'A/B Testing', href: '/ab-testing', category: 'Navigation', icon: <ScienceIcon />, keywords: ['ab', 'experiment', 'test'] },
  { id: 'nav-marketing-hub', label: 'Marketing Hub', href: '/marketing', category: 'Navigation', icon: <HubIcon /> },
  { id: 'nav-ppc', label: 'PPC / Paid Ads', href: '/ppc', category: 'Navigation', icon: <AdsClickIcon />, keywords: ['pay per click', 'ads', 'paid'] },
  { id: 'nav-sms', label: 'SMS Marketing', href: '/sms', category: 'Navigation', icon: <SmsIcon />, keywords: ['text', 'message'] },
  { id: 'nav-email', label: 'Email Campaigns', href: '/email-campaigns', category: 'Navigation', icon: <EmailIcon />, keywords: ['mail', 'newsletter'] },
  { id: 'nav-video', label: 'Video Marketing', href: '/video', category: 'Navigation', icon: <VideoLibraryIcon />, keywords: ['youtube', 'video'] },
  { id: 'nav-seo', label: 'SEO & SEM', href: '/seo', category: 'Navigation', icon: <SearchIcon />, keywords: ['search engine', 'sem', 'optimization'] },
  { id: 'nav-social', label: 'Social', href: '/social', category: 'Navigation', icon: <SocialDistanceIcon />, keywords: ['social media', 'facebook', 'twitter'] },
  { id: 'nav-competitors', label: 'Competitors', href: '/competitors', category: 'Navigation', icon: <GroupIcon />, keywords: ['competition', 'rival'] },
  { id: 'nav-affiliates', label: 'Affiliates', href: '/affiliates', category: 'Navigation', icon: <PeopleIcon />, keywords: ['affiliate', 'partner'] },
  { id: 'nav-budget', label: 'Budget', href: '/budget', category: 'Navigation', icon: <AttachMoneyIcon />, keywords: ['money', 'spend', 'finance'] },
  { id: 'nav-calendar', label: 'Calendar', href: '/calendar', category: 'Navigation', icon: <CalendarTodayIcon />, keywords: ['schedule', 'event', 'date'] },
  { id: 'nav-approvals', label: 'Approvals', href: '/approvals', category: 'Navigation', icon: <TaskAltIcon />, keywords: ['approve', 'workflow', 'review'] },
  { id: 'nav-analytics', label: 'Analytics', href: '/analytics', category: 'Navigation', icon: <AnalyticsIcon />, keywords: ['report', 'stats', 'metrics'] },
  { id: 'nav-integrations', label: 'Integrations', href: '/integrations', category: 'Navigation', icon: <LinkIcon />, keywords: ['connect', 'api', 'third party'] },
  // Quick Actions
  { id: 'action-new-campaign', label: 'New Campaign', href: '/unified-campaigns', category: 'Actions', icon: <AddIcon />, keywords: ['create', 'add'] },
  { id: 'action-new-budget', label: 'New Budget Period', href: '/budget', category: 'Actions', icon: <AttachMoneyIcon />, keywords: ['create', 'add', 'money'] },
  { id: 'action-new-event', label: 'New Calendar Event', href: '/calendar', category: 'Actions', icon: <EventIcon />, keywords: ['create', 'add', 'schedule'] },
]

const categoryOrder = ['Navigation', 'Actions']

// ---------- fuzzy match ----------

function fuzzyMatch(query: string, text: string): boolean {
  const q = query.toLowerCase()
  const t = text.toLowerCase()
  let qi = 0
  for (let ti = 0; ti < t.length && qi < q.length; ti++) {
    if (t[ti] === q[qi]) qi++
  }
  return qi === q.length
}

function highlightMatch(query: string, text: string): React.ReactNode {
  if (!query) return text
  const lowerText = text.toLowerCase()
  const lowerQuery = query.toLowerCase()
  const segments: { text: string; highlight: boolean }[] = []
  let qi = 0
  let current = ''
  let highlighting = false

  for (let ti = 0; ti < lowerText.length; ti++) {
    if (qi < lowerQuery.length && lowerText[ti] === lowerQuery[qi]) {
      if (!highlighting && current) {
        segments.push({ text: current, highlight: false })
        current = ''
      }
      highlighting = true
      current += text[ti]
      qi++
    } else {
      if (highlighting && current) {
        segments.push({ text: current, highlight: true })
        current = ''
      }
      highlighting = false
      current += text[ti]
    }
  }
  if (current) segments.push({ text: current, highlight: highlighting })

  return segments.map((seg, i) =>
    seg.highlight ? (
      <Box key={i} component="span" sx={{ bgcolor: 'action.hover', borderRadius: 0.5, px: 0.5 }}>
        {seg.text}
      </Box>
    ) : (
      seg.text
    )
  )
}

// ---------- component ----------

export default function GlobalSearch() {
  const { open, setOpen } = useSearch()
  const router = useRouter()
  const [query, setQuery] = useState('')
  const [selectedIndex, setSelectedIndex] = useState(-1)
  const inputRef = useRef<HTMLInputElement>(null)
  const listRef = useRef<HTMLUListElement>(null)

  useEffect(() => {
    if (open) {
      setQuery('')
      setSelectedIndex(-1)
      // Focus input after dialog opens
      const timer = setTimeout(() => inputRef.current?.focus(), 100)
      return () => clearTimeout(timer)
    }
  }, [open])

  const results = useMemo(() => {
    if (!query.trim()) {
      // Show all items when query is empty (recent/popular)
      return searchIndex
    }
    const q = query.toLowerCase()
    return searchIndex.filter(
      (item) =>
        fuzzyMatch(q, item.label) ||
        (item.keywords && item.keywords.some((kw) => kw.includes(q)))
    )
  }, [query])

  const grouped = useMemo(() => {
    const map = new Map<string, SearchItem[]>()
    for (const item of results) {
      const arr = map.get(item.category) || []
      arr.push(item)
      map.set(item.category, arr)
    }
    return categoryOrder
      .filter((cat) => map.has(cat))
      .map((cat) => ({ category: cat, items: map.get(cat)! }))
  }, [results])

  const flatResults = useMemo(() => grouped.flatMap((g) => g.items), [grouped])

  const handleSelect = useCallback(
    (item: SearchItem) => {
      setOpen(false)
      if (item.href) {
        router.push(item.href)
      }
    },
    [setOpen, router]
  )

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === 'ArrowDown') {
        e.preventDefault()
        setSelectedIndex((prev) => Math.min(prev + 1, flatResults.length - 1))
      } else if (e.key === 'ArrowUp') {
        e.preventDefault()
        setSelectedIndex((prev) => Math.max(prev - 1, 0))
      } else if (e.key === 'Enter' && selectedIndex >= 0) {
        e.preventDefault()
        handleSelect(flatResults[selectedIndex])
      }
    },
    [flatResults, selectedIndex, handleSelect]
  )

  // Scroll selected item into view
  useEffect(() => {
    if (selectedIndex >= 0 && listRef.current) {
      const el = listRef.current.querySelector(`[data-index="${selectedIndex}"]`)
      el?.scrollIntoView({ block: 'nearest' })
    }
  }, [selectedIndex])

  return (
    <Dialog
      open={open}
      onClose={() => setOpen(false)}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        sx: {
          mt: '10vh',
          borderRadius: 2,
          boxShadow: (theme) =>
            theme.palette.mode === 'dark'
              ? '0 8px 32px rgba(0,0,0,0.6)'
              : '0 8px 32px rgba(0,0,0,0.15)',
          maxHeight: '80vh',
          overflow: 'hidden',
        },
      }}
      onKeyDown={handleKeyDown}
    >
      <DialogTitle sx={{ pb: 0, pt: 1, px: 2 }}>
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1,
            px: 1,
            py: 0.5,
            borderBottom: 1,
            borderColor: 'divider',
          }}
        >
          <SearchIcon sx={{ color: 'text.secondary' }} />
          <InputBase
            inputRef={inputRef}
            value={query}
            onChange={(e) => {
              setQuery(e.target.value)
              setSelectedIndex(-1)
            }}
            placeholder="Search pages, campaigns, actions..."
            autoFocus
            fullWidth
            sx={{ flex: 1, fontSize: '1rem', py: 1.5 }}
          />
          <Chip
            label="ESC"
            size="small"
            sx={{
              fontSize: '0.65rem',
              height: 20,
              fontWeight: 600,
              bgcolor: 'action.hover',
            }}
          />
        </Box>
      </DialogTitle>

      <DialogContent sx={{ p: 0, overflow: 'auto' }}>
        {flatResults.length === 0 ? (
          <Box sx={{ py: 6, textAlign: 'center' }}>
            <SearchIcon sx={{ fontSize: 48, color: 'text.disabled', mb: 1 }} />
            <Typography color="text.secondary">No results found for &quot;{query}&quot;</Typography>
          </Box>
        ) : (
          <List ref={listRef} sx={{ py: 0.5, px: 1 }}>
            {grouped.map((group) => (
              <Box key={group.category}>
                <Typography
                  variant="caption"
                  sx={{
                    px: 1,
                    py: 0.5,
                    display: 'block',
                    color: 'text.secondary',
                    fontWeight: 600,
                    fontSize: '0.65rem',
                    letterSpacing: 0.5,
                  }}
                >
                  {group.category.toUpperCase()}
                </Typography>
                {group.items.map((item) => {
                  const idx = flatResults.indexOf(item)
                  const selected = idx === selectedIndex
                  return (
                    <ListItemButton
                      key={item.id}
                      data-index={idx}
                      selected={selected}
                      onClick={() => handleSelect(item)}
                      onMouseEnter={() => setSelectedIndex(idx)}
                      sx={{
                        borderRadius: 1,
                        py: 0.75,
                        '&.Mui-selected': {
                          bgcolor: (theme) => alpha(theme.palette.primary.main, 0.12),
                          '&:hover': {
                            bgcolor: (theme) => alpha(theme.palette.primary.main, 0.2),
                          },
                        },
                      }}
                    >
                      <ListItemIcon sx={{ minWidth: 40, color: 'text.secondary' }}>
                        {item.icon}
                      </ListItemIcon>
                      <ListItemText
                        primary={highlightMatch(query, item.label)}
                        primaryTypographyProps={{
                          fontSize: '0.9rem',
                          fontWeight: selected ? 600 : 400,
                        }}
                      />
                      {item.href && (
                        <Typography
                          variant="caption"
                          sx={{ color: 'text.disabled', ml: 1, fontSize: '0.7rem' }}
                        >
                          {item.href}
                        </Typography>
                      )}
                    </ListItemButton>
                  )
                })}
                <Divider sx={{ my: 0.5 }} />
              </Box>
            ))}
          </List>
        )}

        {/* Footer hints */}
        <Paper
          elevation={0}
          sx={{
            px: 2,
            py: 1,
            borderTop: 1,
            borderColor: 'divider',
            display: 'flex',
            gap: 2,
            bgcolor: 'background.paper',
          }}
        >
          <Typography variant="caption" color="text.disabled" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Box component="kbd" sx={{ px: 0.5, py: 0.25, borderRadius: 0.5, bgcolor: 'action.hover', fontSize: '0.6rem' }}>
              ↑↓
            </Box>
            Navigate
          </Typography>
          <Typography variant="caption" color="text.disabled" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Box component="kbd" sx={{ px: 0.5, py: 0.25, borderRadius: 0.5, bgcolor: 'action.hover', fontSize: '0.6rem' }}>
              ↵
            </Box>
            Select
          </Typography>
          <Typography variant="caption" color="text.disabled" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Box component="kbd" sx={{ px: 0.5, py: 0.25, borderRadius: 0.5, bgcolor: 'action.hover', fontSize: '0.6rem' }}>
              ESC
            </Box>
            Close
          </Typography>
        </Paper>
      </DialogContent>
    </Dialog>
  )
}
