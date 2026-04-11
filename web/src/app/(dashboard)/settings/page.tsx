'use client'

import { useState, useCallback, ChangeEvent, FormEvent } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { useThemeMode } from '@/hooks/useTheme'
import { useToast } from '@/hooks/useToast'
import { TabPanel } from '@/components/TabPanel'
import { authApi } from '@/lib/api-client'
import {
  Tabs,
  Tab,
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Avatar,
  Switch,
  FormControlLabel,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Divider,
  Grid,
  IconButton,
  Tooltip,
  Alert,
  Snackbar,
  CircularProgress,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  ListItemSecondaryAction,
  Chip,
  InputAdornment,
  Link,
} from '@mui/material'
import {
  CloudUpload as CloudUploadIcon,
  Visibility as VisibilityIcon,
  VisibilityOff as VisibilityOffIcon,
  Save as SaveIcon,
  Security as SecurityIcon,
  Notifications as NotificationsIcon,
  Extension as ExtensionIcon,
  Person as PersonIcon,
  Palette as PaletteIcon,
  Language as LanguageIcon,
  AccessTime as AccessTimeIcon,
  CalendarToday as CalendarIcon,
  ListAlt as ListAltIcon,
  CheckCircle as CheckCircleIcon,
  Error as ErrorIcon,
  Warning as WarningIcon,
  LaptopWindows as LaptopIcon,
  PhoneAndroid as PhoneIcon,
  Tablet as TabletIcon,
  OpenInNew as OpenInNewIcon,
  Add as AddIcon,
} from '@mui/icons-material'

// --- Tab definitions ---
interface TabItem {
  label: string
  value: number
  icon: React.ReactElement
}

const settingTabs: TabItem[] = [
  { label: 'Profile', value: 0, icon: <PersonIcon /> },
  { label: 'Security', value: 1, icon: <SecurityIcon /> },
  { label: 'Preferences', value: 2, icon: <PaletteIcon /> },
  { label: 'Notifications', value: 3, icon: <NotificationsIcon /> },
  { label: 'Integrations', value: 4, icon: <ExtensionIcon /> },
]

// --- Password strength helper ---
function getPasswordStrength(password: string): { score: number; label: string; color: string } {
  let score = 0
  if (password.length >= 8) score++
  if (password.length >= 12) score++
  if (/[A-Z]/.test(password)) score++
  if (/[a-z]/.test(password)) score++
  if (/[0-9]/.test(password)) score++
  if (/[^A-Za-z0-9]/.test(password)) score++

  if (score <= 2) return { score, label: 'Weak', color: 'error' }
  if (score <= 3) return { score, label: 'Fair', color: 'warning' }
  if (score <= 4) return { score, label: 'Good', color: 'info' }
  return { score, label: 'Strong', color: 'success' }
}

// --- Section header component ---
function SectionHeader({ title, description }: { title: string; description?: string }) {
  return (
    <Box sx={{ mb: 3 }}>
      <Typography variant="h6" fontWeight="bold" gutterBottom>
        {title}
      </Typography>
      {description && (
        <Typography variant="body2" color="text.secondary">
          {description}
        </Typography>
      )}
    </Box>
  )
}

// --- Settings Card wrapper ---
function SettingsCard({
  title,
  description,
  children,
  sx,
}: {
  title: string
  description?: string
  children: React.ReactNode
  sx?: object
}) {
  return (
    <Card sx={{ mb: 3, ...sx }}>
      <CardContent>
        <SectionHeader title={title} description={description} />
        {children}
      </CardContent>
    </Card>
  )
}

// ============================================================
// Tab 1: Profile
// ============================================================
function ProfileTab() {
  const { user } = useAuth()
  const toast = useToast()
  const [fullName, setFullName] = useState(user?.full_name || '')
  const [avatarUrl, setAvatarUrl] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const handleAvatarUpload = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const url = URL.createObjectURL(file)
      setAvatarUrl(url)
      toast.info('Avatar selected. Upload will be saved when you click Save.')
    }
  }, [toast])

  const handleSave = async () => {
    if (!fullName.trim()) {
      toast.error('Full name is required')
      return
    }
    setSaving(true)
    try {
      // Placeholder: call profile update API
      // await authApi.updateProfile({ full_name: fullName })
      await new Promise((r) => setTimeout(r, 500))
      toast.success('Profile updated successfully')
    } catch {
      toast.error('Failed to update profile')
    } finally {
      setSaving(false)
    }
  }

  const initials = fullName
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)

  return (
    <Box>
      {/* Avatar Section */}
      <SettingsCard title="Profile Photo" description="Upload a profile photo to personalize your account.">
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
          <Avatar
            src={avatarUrl || undefined}
            sx={{ width: 96, height: 96, fontSize: '2rem', bgcolor: 'primary.main', color: 'white' }}
          >
            {initials || '?'}
          </Avatar>
          <Box>
            <input
              accept="image/*"
              type="file"
              id="avatar-upload"
              hidden
              onChange={handleAvatarUpload}
            />
            <label htmlFor="avatar-upload">
              <Button
                variant="outlined"
                component="span"
                startIcon={<CloudUploadIcon />}
                size="small"
              >
                Upload Photo
              </Button>
            </label>
            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1 }}>
              JPG, PNG or GIF. Max 2MB.
            </Typography>
          </Box>
        </Box>
      </SettingsCard>

      {/* Personal Info */}
      <SettingsCard title="Personal Information" description="Update your personal details.">
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Full Name"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              placeholder="Enter your full name"
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              label="Email"
              value={user?.email || ''}
              InputProps={{ readOnly: true }}
              helperText={
                <span>
                  To change your email,{' '}
                  <Link href="#" underline="hover" onClick={() => toast.info('Change email flow coming soon!')}>
                    click here
                  </Link>
                </span>
              }
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              label="Role"
              value={user?.role || 'user'}
              InputProps={{ readOnly: true }}
              helperText="Contact an administrator to change your role"
            />
          </Grid>
        </Grid>
        <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            variant="contained"
            startIcon={saving ? <CircularProgress size={20} color="inherit" /> : <SaveIcon />}
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </Button>
        </Box>
      </SettingsCard>
    </Box>
  )
}

// ============================================================
// Tab 2: Security
// ============================================================
function SecurityTab() {
  const toast = useToast()
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showCurrent, setShowCurrent] = useState(false)
  const [showNew, setShowNew] = useState(false)
  const [showConfirm, setShowConfirm] = useState(false)
  const [changing, setChanging] = useState(false)
  const [twoFactorEnabled, setTwoFactorEnabled] = useState(false)

  const strength = getPasswordStrength(newPassword)

  const passwordsMatch = newPassword === confirmPassword && confirmPassword.length > 0
  const canChangePassword =
    currentPassword.length > 0 &&
    newPassword.length >= 8 &&
    passwordsMatch

  const handleChangePassword = async () => {
    if (!canChangePassword) return
    setChanging(true)
    try {
      // Placeholder: await authApi.changePassword({ current_password, new_password })
      await new Promise((r) => setTimeout(r, 800))
      setCurrentPassword('')
      setNewPassword('')
      setConfirmPassword('')
      toast.success('Password changed successfully')
    } catch {
      toast.error('Failed to change password')
    } finally {
      setChanging(false)
    }
  }

  const mockSessions = [
    { device: 'Chrome on Windows', location: 'Bangkok, TH', lastActive: 'Just now', type: 'laptop' as const, current: true },
    { device: 'Safari on iPhone', location: 'Bangkok, TH', lastActive: '2 hours ago', type: 'phone' as const, current: false },
    { device: 'Firefox on macOS', location: 'Chiang Mai, TH', lastActive: '3 days ago', type: 'laptop' as const, current: false },
  ]

  const sessionIconMap = {
    laptop: <LaptopIcon fontSize="small" />,
    phone: <PhoneIcon fontSize="small" />,
    tablet: <TabletIcon fontSize="small" />,
  }

  return (
    <Box>
      {/* Change Password */}
      <SettingsCard title="Change Password" description="Ensure your account stays secure by using a strong password.">
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Current Password"
              type={showCurrent ? 'text' : 'password'}
              value={currentPassword}
              onChange={(e) => setCurrentPassword(e.target.value)}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton onClick={() => setShowCurrent(!showCurrent)} edge="end">
                      {showCurrent ? <VisibilityOffIcon /> : <VisibilityIcon />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="New Password"
              type={showNew ? 'text' : 'password'}
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton onClick={() => setShowNew(!showNew)} edge="end">
                      {showNew ? <VisibilityOffIcon /> : <VisibilityIcon />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
              helperText="Minimum 8 characters"
            />
            {newPassword.length > 0 && (
              <Box sx={{ mt: 1.5 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                  <Box
                    sx={{
                      flex: 1,
                      height: 6,
                      borderRadius: 3,
                      bgcolor: 'grey.200',
                      overflow: 'hidden',
                    }}
                  >
                    <Box
                      sx={{
                        height: '100%',
                        width: `${(strength.score / 5) * 100}%`,
                        bgcolor: `${strength.color}.main`,
                        borderRadius: 3,
                        transition: 'width 0.3s, background-color 0.3s',
                      }}
                    />
                  </Box>
                  <Chip
                    label={strength.label}
                    size="small"
                    color={strength.color as any}
                    variant="outlined"
                    sx={{ minWidth: 60 }}
                  />
                </Box>
              </Box>
            )}
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Confirm New Password"
              type={showConfirm ? 'text' : 'password'}
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              error={confirmPassword.length > 0 && !passwordsMatch}
              helperText={
                confirmPassword.length > 0 && !passwordsMatch
                  ? 'Passwords do not match'
                  : passwordsMatch
                    ? 'Passwords match'
                    : ''
              }
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton onClick={() => setShowConfirm(!showConfirm)} edge="end">
                      {showConfirm ? <VisibilityOffIcon /> : <VisibilityIcon />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
        </Grid>
        <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            variant="contained"
            startIcon={changing ? <CircularProgress size={20} color="inherit" /> : <SecurityIcon />}
            onClick={handleChangePassword}
            disabled={!canChangePassword || changing}
          >
            {changing ? 'Changing...' : 'Change Password'}
          </Button>
        </Box>
      </SettingsCard>

      {/* Two-Factor Authentication */}
      <SettingsCard title="Two-Factor Authentication" description="Add an extra layer of security to your account.">
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box>
            <Typography variant="body1" fontWeight="medium">
              Authenticator App
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {twoFactorEnabled
                ? 'Two-factor authentication is enabled'
                : 'Protect your account with an authenticator app'}
            </Typography>
          </Box>
          <FormControlLabel
            control={
              <Switch
                checked={twoFactorEnabled}
                onChange={(e) => {
                  setTwoFactorEnabled(e.target.checked)
                  toast.info(e.target.checked ? '2FA setup coming soon!' : '2FA disabled')
                }}
              />
            }
            label=""
            labelPlacement="start"
          />
        </Box>
      </SettingsCard>

      {/* Active Sessions */}
      <SettingsCard title="Active Sessions" description="Manage devices where you are currently logged in.">
        <List disablePadding>
          {mockSessions.map((session, idx) => (
            <Box key={idx}>
              <ListItem>
                <ListItemAvatar>
                  <Avatar sx={{ bgcolor: session.current ? 'primary.main' : 'grey.300' }}>
                    {sessionIconMap[session.type]}
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      {session.device}
                      {session.current && (
                        <Chip label="Current" size="small" color="primary" variant="outlined" />
                      )}
                    </Box>
                  }
                  secondary={`${session.location} - Last active: ${session.lastActive}`}
                />
                <ListItemSecondaryAction>
                  {!session.current && (
                    <Tooltip title="Revoke session">
                      <Button size="small" color="error" variant="outlined">
                        Revoke
                      </Button>
                    </Tooltip>
                  )}
                </ListItemSecondaryAction>
              </ListItem>
              {idx < mockSessions.length - 1 && <Divider />}
            </Box>
          ))}
        </List>
      </SettingsCard>
    </Box>
  )
}

// ============================================================
// Tab 3: Preferences
// ============================================================
function PreferencesTab() {
  const themeCtx = useThemeMode()
  const toast = useToast()
  const [theme, setTheme] = useState<'light' | 'dark' | 'system'>(themeCtx.mode)
  const [language, setLanguage] = useState('en')
  const [timezone, setTimezone] = useState('Asia/Bangkok')
  const [dateFormat, setDateFormat] = useState('MM/DD/YYYY')
  const [itemsPerPage, setItemsPerPage] = useState(10)
  const [saving, setSaving] = useState(false)

  const handleThemeChange = (newTheme: 'light' | 'dark' | 'system') => {
    setTheme(newTheme)
    if (newTheme !== 'system') {
      themeCtx.setMode(newTheme)
    } else {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      themeCtx.setMode(prefersDark ? 'dark' : 'light')
    }
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      // Placeholder: await settingsApi.savePreferences({ theme, language, timezone, dateFormat, itemsPerPage })
      await new Promise((r) => setTimeout(r, 500))
      localStorage.setItem('preferences', JSON.stringify({ theme, language, timezone, dateFormat, itemsPerPage }))
      toast.success('Preferences saved successfully')
    } catch {
      toast.error('Failed to save preferences')
    } finally {
      setSaving(false)
    }
  }

  const timezones = [
    'Asia/Bangkok', 'Asia/Ho_Chi_Minh', 'Asia/Jakarta', 'Asia/Manila',
    'Asia/Singapore', 'Asia/Tokyo', 'Asia/Shanghai', 'Asia/Kolkata',
    'America/New_York', 'America/Chicago', 'America/Denver', 'America/Los_Angeles',
    'Europe/London', 'Europe/Paris', 'Europe/Berlin',
    'Australia/Sydney', 'Pacific/Auckland',
    'UTC',
  ]

  return (
    <Box>
      <SettingsCard title="Appearance" description="Customize how the application looks.">
        <Box sx={{ mb: 3 }}>
          <Typography variant="body2" fontWeight="medium" gutterBottom>
            Theme
          </Typography>
          <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
            {(['light', 'dark', 'system'] as const).map((t) => (
              <Button
                key={t}
                variant={theme === t ? 'contained' : 'outlined'}
                onClick={() => handleThemeChange(t)}
                sx={{ textTransform: 'capitalize', minWidth: 120 }}
              >
                {t === 'light' && <PaletteIcon sx={{ mr: 1 }} />}
                {t === 'dark' && <PaletteIcon sx={{ mr: 1 }} />}
                {t === 'system' && <LaptopIcon sx={{ mr: 1 }} />}
                {t.charAt(0).toUpperCase() + t.slice(1)}
              </Button>
            ))}
          </Box>
        </Box>
      </SettingsCard>

      <SettingsCard title="Regional Settings" description="Configure language, timezone, and date preferences.">
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth>
              <InputLabel id="language-label">Language</InputLabel>
              <Select
                labelId="language-label"
                value={language}
                label="Language"
                onChange={(e) => setLanguage(e.target.value)}
                startAdornment={<LanguageIcon sx={{ mr: 1, color: 'text.secondary' }} />}
              >
                <MenuItem value="en">English</MenuItem>
                <MenuItem value="th">Thai (ภาษาไทย)</MenuItem>
                <MenuItem value="ja">Japanese (日本語)</MenuItem>
                <MenuItem value="zh">Chinese (中文)</MenuItem>
                <MenuItem value="vi">Vietnamese (Tiếng Việt)</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth>
              <InputLabel id="timezone-label">Timezone</InputLabel>
              <Select
                labelId="timezone-label"
                value={timezone}
                label="Timezone"
                onChange={(e) => setTimezone(e.target.value)}
                startAdornment={<AccessTimeIcon sx={{ mr: 1, color: 'text.secondary' }} />}
              >
                {timezones.map((tz) => (
                  <MenuItem key={tz} value={tz}>
                    {tz.replace('_', ' ')}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth>
              <InputLabel id="date-format-label">Date Format</InputLabel>
              <Select
                labelId="date-format-label"
                value={dateFormat}
                label="Date Format"
                onChange={(e) => setDateFormat(e.target.value)}
                startAdornment={<CalendarIcon sx={{ mr: 1, color: 'text.secondary' }} />}
              >
                <MenuItem value="MM/DD/YYYY">MM/DD/YYYY</MenuItem>
                <MenuItem value="DD/MM/YYYY">DD/MM/YYYY</MenuItem>
                <MenuItem value="YYYY-MM-DD">YYYY-MM-DD</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth>
              <InputLabel id="ipp-label">Items per Page</InputLabel>
              <Select
                labelId="ipp-label"
                value={itemsPerPage}
                label="Items per Page"
                onChange={(e) => setItemsPerPage(Number(e.target.value))}
                startAdornment={<ListAltIcon sx={{ mr: 1, color: 'text.secondary' }} />}
              >
                <MenuItem value={10}>10</MenuItem>
                <MenuItem value={25}>25</MenuItem>
                <MenuItem value={50}>50</MenuItem>
                <MenuItem value={100}>100</MenuItem>
              </Select>
            </FormControl>
          </Grid>
        </Grid>
        <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            variant="contained"
            startIcon={saving ? <CircularProgress size={20} color="inherit" /> : <SaveIcon />}
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? 'Saving...' : 'Save Preferences'}
          </Button>
        </Box>
      </SettingsCard>
    </Box>
  )
}

// ============================================================
// Tab 4: Notifications
// ============================================================
function NotificationsTab() {
  const toast = useToast()
  const [emailNotifs, setEmailNotifs] = useState({
    campaignStatus: true,
    budgetAlerts: true,
    approvalRequests: true,
    weeklyReports: false,
  })
  const [pushNotifs, setPushNotifs] = useState({
    campaignStatus: true,
    budgetAlerts: false,
    approvalRequests: true,
    weeklyReports: false,
  })
  const [saving, setSaving] = useState(false)

  const handleEmailToggle = (key: keyof typeof emailNotifs) => {
    setEmailNotifs((prev) => ({ ...prev, [key]: !prev[key] }))
  }

  const handlePushToggle = (key: keyof typeof pushNotifs) => {
    setPushNotifs((prev) => ({ ...prev, [key]: !prev[key] }))
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      // Placeholder: await notificationApi.updatePreferences({ email: emailNotifs, push: pushNotifs })
      await new Promise((r) => setTimeout(r, 500))
      localStorage.setItem('notificationPreferences', JSON.stringify({ email: emailNotifs, push: pushNotifs }))
      toast.success('Notification preferences saved')
    } catch {
      toast.error('Failed to save notification preferences')
    } finally {
      setSaving(false)
    }
  }

  const notifItems = [
    {
      key: 'campaignStatus' as const,
      title: 'Campaign Status Changes',
      description: 'Get notified when campaigns are paused, activated, or completed',
    },
    {
      key: 'budgetAlerts' as const,
      title: 'Budget Alerts',
      description: 'Receive alerts when campaign budgets reach critical thresholds',
    },
    {
      key: 'approvalRequests' as const,
      title: 'Approval Requests',
      description: 'Notifications for pending approvals and approval outcomes',
    },
    {
      key: 'weeklyReports' as const,
      title: 'Weekly Reports',
      description: 'Receive a weekly summary of campaign performance',
    },
  ]

  return (
    <Box>
      <SettingsCard title="Notification Channels" description="Choose how you want to receive notifications.">
        {/* Email Notifications */}
        <Box sx={{ mb: 4 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
            <Typography variant="subtitle1" fontWeight="bold">
              Email Notifications
            </Typography>
            <Chip label="Email" size="small" color="primary" variant="outlined" />
          </Box>
          <Grid container spacing={2}>
            {notifItems.map((item) => (
              <Grid item xs={12} key={item.key}>
                <Box
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    p: 2,
                    borderRadius: 2,
                    bgcolor: 'grey.50',
                  }}
                >
                  <Box>
                    <Typography variant="body2" fontWeight="medium">
                      {item.title}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {item.description}
                    </Typography>
                  </Box>
                  <Switch
                    checked={emailNotifs[item.key]}
                    onChange={() => handleEmailToggle(item.key)}
                  />
                </Box>
              </Grid>
            ))}
          </Grid>
        </Box>

        <Divider sx={{ my: 3 }} />

        {/* Push Notifications */}
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
            <Typography variant="subtitle1" fontWeight="bold">
              Push Notifications
            </Typography>
            <Chip label="Push" size="small" color="secondary" variant="outlined" />
          </Box>
          <Grid container spacing={2}>
            {notifItems.map((item) => (
              <Grid item xs={12} key={item.key}>
                <Box
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    p: 2,
                    borderRadius: 2,
                    bgcolor: 'grey.50',
                  }}
                >
                  <Box>
                    <Typography variant="body2" fontWeight="medium">
                      {item.title}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {item.description}
                    </Typography>
                  </Box>
                  <Switch
                    checked={pushNotifs[item.key]}
                    onChange={() => handlePushToggle(item.key)}
                  />
                </Box>
              </Grid>
            ))}
          </Grid>
        </Box>

        <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
          <Button
            variant="contained"
            startIcon={saving ? <CircularProgress size={20} color="inherit" /> : <SaveIcon />}
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? 'Saving...' : 'Save Notification Settings'}
          </Button>
        </Box>
      </SettingsCard>
    </Box>
  )
}

// ============================================================
// Tab 5: Integrations
// ============================================================
interface IntegrationCard {
  name: string
  description: string
  icon: string
  connected: boolean
  category: string
  lastSync?: string
  manageUrl?: string
  connectAction: string
}

function IntegrationsTab() {
  const toast = useToast()

  const integrations: IntegrationCard[] = [
    {
      name: 'Google Ads',
      description: 'Manage Google Ads campaigns and track performance metrics',
      icon: 'G',
      connected: true,
      category: 'Advertising',
      lastSync: '5 minutes ago',
      manageUrl: '#',
      connectAction: 'Manage',
    },
    {
      name: 'Facebook Ads',
      description: 'Run and monitor Facebook advertising campaigns',
      icon: 'f',
      connected: true,
      category: 'Advertising',
      lastSync: '1 hour ago',
      manageUrl: '#',
      connectAction: 'Manage',
    },
    {
      name: 'Google Analytics',
      description: 'Track website traffic and user behavior analytics',
      icon: 'GA',
      connected: true,
      category: 'Analytics',
      lastSync: 'Just now',
      manageUrl: '#',
      connectAction: 'Manage',
    },
    {
      name: 'Mailchimp',
      description: 'Sync email lists and manage email campaigns',
      icon: 'M',
      connected: false,
      category: 'Email',
      connectAction: 'Connect',
    },
    {
      name: 'Shopee',
      description: 'Connect your Shopee store for e-commerce integration',
      icon: 'S',
      connected: false,
      category: 'E-commerce',
      connectAction: 'Connect',
    },
    {
      name: 'TikTok Ads',
      description: 'Create and manage TikTok advertising campaigns',
      icon: 'T',
      connected: false,
      category: 'Advertising',
      connectAction: 'Connect',
    },
  ]

  const statusColor = (connected: boolean, lastSync?: string) => {
    if (!connected) return 'default'
    if (lastSync?.includes('Just now') || lastSync?.includes('minute')) return 'success'
    if (lastSync?.includes('hour')) return 'warning'
    return 'error'
  }

  const statusLabel = (connected: boolean, lastSync?: string) => {
    if (!connected) return 'Not Connected'
    return `Connected - ${lastSync || 'Unknown'}`
  }

  const handleConnect = (name: string) => {
    toast.info(`${name} integration setup coming soon!`)
  }

  const handleManage = (name: string) => {
    toast.info(`Opening ${name} management...`)
  }

  return (
    <Box>
      <Alert severity="info" sx={{ mb: 3 }}>
        Connect your favorite platforms to manage campaigns, analytics, and e-commerce from one place.
      </Alert>

      <Grid container spacing={3}>
        {integrations.map((integration) => (
          <Grid item xs={12} sm={6} md={4} key={integration.name}>
            <Card
              sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                border: integration.connected ? '2px solid' : '1px solid',
                borderColor: integration.connected ? 'success.light' : 'grey.200',
                transition: 'box-shadow 0.2s',
                '&:hover': { boxShadow: 6 },
              }}
            >
              <CardContent sx={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                {/* Header */}
                <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', mb: 2 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                    <Avatar
                      sx={{
                        width: 40,
                        height: 40,
                        bgcolor: integration.connected ? 'primary.main' : 'grey.300',
                        fontSize: '1rem',
                        fontWeight: 'bold',
                      }}
                    >
                      {integration.icon}
                    </Avatar>
                    <Box>
                      <Typography variant="subtitle1" fontWeight="bold">
                        {integration.name}
                      </Typography>
                      <Chip label={integration.category} size="small" variant="outlined" sx={{ mt: 0.5 }} />
                    </Box>
                  </Box>
                  {integration.connected ? (
                    <CheckCircleIcon color="success" fontSize="small" />
                  ) : (
                    <WarningIcon color="disabled" fontSize="small" />
                  )}
                </Box>

                {/* Description */}
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2, flex: 1 }}>
                  {integration.description}
                </Typography>

                {/* Status */}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                  <Box
                    sx={{
                      width: 8,
                      height: 8,
                      borderRadius: '50%',
                      bgcolor: `${statusColor(integration.connected, integration.lastSync)}.main`,
                    }}
                  />
                  <Typography variant="caption" color="text.secondary">
                    {statusLabel(integration.connected, integration.lastSync)}
                  </Typography>
                </Box>

                {/* Action Button */}
                {integration.connected ? (
                  <Button
                    fullWidth
                    variant="outlined"
                    size="small"
                    endIcon={<OpenInNewIcon />}
                    onClick={() => handleManage(integration.name)}
                  >
                    {integration.connectAction}
                  </Button>
                ) : (
                  <Button
                    fullWidth
                    variant="contained"
                    size="small"
                    startIcon={<AddIcon />}
                    onClick={() => handleConnect(integration.name)}
                  >
                    {integration.connectAction}
                  </Button>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  )
}

// ============================================================
// Main Settings Page
// ============================================================
export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState(0)
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' | 'info' | 'warning' }>({
    open: false,
    message: '',
    severity: 'info',
  })

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue)
  }

  const handleCloseSnackbar = () => {
    setSnackbar((prev) => ({ ...prev, open: false }))
  }

  return (
    <Box>
      {/* Page Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight="bold" gutterBottom>
          Settings
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Manage your account settings, preferences, and integrations.
        </Typography>
      </Box>

      {/* Tabs Navigation */}
      <Card sx={{ mb: 3 }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs
            value={activeTab}
            onChange={handleTabChange}
            variant="scrollable"
            scrollButtons="auto"
            aria-label="Settings tabs"
            sx={{
              '& .MuiTab-root': {
                minHeight: 56,
                fontWeight: 500,
                textTransform: 'none',
              },
            }}
          >
            {settingTabs.map((tab) => (
              <Tab
                key={tab.value}
                label={tab.label}
                icon={tab.icon}
                iconPosition="start"
                id={`settings-tab-${tab.value}`}
                aria-controls={`settings-tabpanel-${tab.value}`}
              />
            ))}
          </Tabs>
        </Box>
      </Card>

      {/* Tab Content */}
      <Box sx={{ mt: 2 }}>
        <TabPanel value={activeTab} index={0} id="settings-tabpanel-0">
          <ProfileTab />
        </TabPanel>
        <TabPanel value={activeTab} index={1} id="settings-tabpanel-1">
          <SecurityTab />
        </TabPanel>
        <TabPanel value={activeTab} index={2} id="settings-tabpanel-2">
          <PreferencesTab />
        </TabPanel>
        <TabPanel value={activeTab} index={3} id="settings-tabpanel-3">
          <NotificationsTab />
        </TabPanel>
        <TabPanel value={activeTab} index={4} id="settings-tabpanel-4">
          <IntegrationsTab />
        </TabPanel>
      </Box>

      {/* Snackbar for general feedback */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={4000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={handleCloseSnackbar} severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  )
}
