'use client'

import { useTheme } from '@mui/material/styles'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts'

interface DataPoint {
  name: string
  [key: string]: string | number
}

interface BarChartProps {
  data: DataPoint[]
  bars: { dataKey: string; color: string; name: string }[]
  height?: number
}

export function CustomBarChart({ data, bars, height = 300 }: BarChartProps) {
  const theme = useTheme()

  return (
    <ResponsiveContainer width="100%" height={height}>
      <BarChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
        <XAxis dataKey="name" tick={{ fontSize: 12, fill: theme.palette.text.secondary }} stroke={theme.palette.text.secondary} />
        <YAxis tick={{ fontSize: 12, fill: theme.palette.text.secondary }} stroke={theme.palette.text.secondary} />
        <Tooltip
          contentStyle={{
            backgroundColor: theme.palette.background.paper,
            border: `1px solid ${theme.palette.divider}`,
            borderRadius: '8px',
            color: theme.palette.text.primary,
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
          }}
        />
        <Legend wrapperStyle={{ color: theme.palette.text.secondary }} />
        {bars.map((bar) => (
          <Bar
            key={bar.dataKey}
            dataKey={bar.dataKey}
            name={bar.name}
            fill={bar.color}
            radius={[4, 4, 0, 0]}
          />
        ))}
      </BarChart>
    </ResponsiveContainer>
  )
}
