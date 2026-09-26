import { useEffect, useState } from 'react'
import { getDashboardStats, type DashboardStats } from '@/api/dashboard'

export function useApplicationStats() {
  const [data, setData] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(false)

  async function load() {
    setLoading(true)
    try {
      setData(await getDashboardStats())
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  return { data, loading, reload: load }
}
