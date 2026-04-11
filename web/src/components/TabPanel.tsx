import Box from '@mui/material/Box'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
  id?: string
}

export function TabPanel(props: TabPanelProps) {
  const { children, value, index, id, ...other } = props
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={id || `tabpanel-${index}`}
      aria-labelledby={id ? `${id}-tab` : `tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 2 }}>{children}</Box>}
    </div>
  )
}
