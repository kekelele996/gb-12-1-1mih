import { create } from 'zustand'
import { listApplications } from '@/api/application'
import type { ApplicationProject } from '@/types/api'

interface ApplicationState {
  applications: ApplicationProject[]
  load: () => Promise<void>
}

export const useApplicationStore = create<ApplicationState>((set) => ({
  applications: [],
  load: async () => set({ applications: await listApplications() }),
}))
