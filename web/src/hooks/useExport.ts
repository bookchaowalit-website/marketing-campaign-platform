'use client'

import { useToast } from './useToast'
import { exportToCSV } from '@/lib/export'

export function useExport() {
  const toast = useToast()

  const handleExport = async (data: any[], filename: string, columns?: string[]) => {
    if (!data || data.length === 0) {
      toast.warning('No data to export')
      return
    }

    try {
      const total = data.length
      exportToCSV(data, filename, columns, (progress) => {
        if (total > 100 && progress % 25 === 0) {
          toast.info(`Exporting... ${Math.round((progress / total) * 100)}%`)
        }
      })

      toast.success(`Exported ${total} records to ${filename}.csv`)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to export data')
    }
  }

  return { handleExport }
}
