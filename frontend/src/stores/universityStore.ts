import { create } from 'zustand'
import { listUniversities } from '@/api/university'
import type { University } from '@/types/api'

interface UniversityState {
  universities: University[]
  total: number
  load: (params?: { page?: number; page_size?: number; country?: string; keyword?: string }) => Promise<void>
}

export const useUniversityStore = create<UniversityState>((set) => ({
  universities: [],
  total: 0,
  load: async (params = {}) => {
    const res = await listUniversities(params)
    set({ universities: res.list, total: res.total })
  },
}))
