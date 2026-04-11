'use client'

import { useTheme } from '@mui/material/styles'
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts'

interface DataPoint {
  name: string
  value: number
}

interface PieChartProps {
  data: DataPoint[]
  colors?: string[]
  height?: number
  innerRadius?: number
  outerRadius?: number
}

const DEFAULT_COLORS = ['#2563eb', '#7c3aed', '#10b981', '#f59e0b', '#ef4444', '#06b6d4']

export function CustomPieChart({
  data,
  colors = DEFAULT_COLORS,
  height = 300,
  innerRadius = 60,
  outerRadius = 100
}: PieChartProps) {
  const theme = useTheme()

  return (
    <ResponsiveContainer width="100%" height={height}>
      <PieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          innerRadius={innerRadius}
          outerRadius={outerRadius}
          paddingAngle={2}
          dataKey="value"
          label={({ name, percent }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}
          labelLine={false}
        >
          {data.map((_, index) => (
            <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
          ))}
        </Pie>
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
      </PieChart>
    </ResponsiveContainer>
  )
}
