import request from '@/utils/request'
import type { ApplicationProject } from '@/types/api'

export interface DashboardStats {
  stats: {
    total: number
    by_status: Record<string, number>
    admitted: number
    applied: number
    material_avg: number
  }
  projects: ApplicationProject[]
}

export function getDashboardStats() { return request.get<never, DashboardStats>('/dashboard/stats') }
