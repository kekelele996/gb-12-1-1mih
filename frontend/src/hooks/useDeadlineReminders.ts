import { useEffect, useState } from 'react'
import { listApplications } from '@/api/application'
import { listTimeline } from '@/api/application'
import type { TimelineNode } from '@/types/api'
import dayjs from 'dayjs'

export function useDeadlineReminders() {
  const [reminders, setReminders] = useState<TimelineNode[]>([])

  useEffect(() => {
    ;(async () => {
      const apps = await listApplications()
      const all: TimelineNode[] = []
      for (const a of apps.slice(0, 10)) {
        try {
          all.push(...(await listTimeline(a.id)))
        } catch {
          // skip
        }
      }
      const soon = all.filter((n) => !n.is_done && dayjs(n.due_date).diff(dayjs(), 'day') <= 14)
      setReminders(soon)
    })()
  }, [])

  return { reminders }
}
