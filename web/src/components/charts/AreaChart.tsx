'use client'

import { useTheme } from '@mui/material/styles'
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts'

interface DataPoint {
  name: string
  [key: string]: string | number
}

interface AreaChartProps {
  data: DataPoint[]
  areas: { dataKey: string; color: string; name: string }[]
  height?: number
}

export function CustomAreaChart({ data, areas, height = 300 }: AreaChartProps) {
  const theme = useTheme()

  return (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
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
        {areas.map((area) => (
          <Area
            key={area.dataKey}
            type="monotone"
            dataKey={area.dataKey}
            name={area.name}
            stroke={area.color}
            fill={area.color}
            fillOpacity={0.3}
          />
        ))}
      </AreaChart>
    </ResponsiveContainer>
  )
}
