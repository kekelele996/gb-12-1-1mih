import { create } from 'zustand'
import { listVersions } from '@/api/document'
import type { DocumentVersion } from '@/types/api'

interface DocumentState {
  versions: DocumentVersion[]
  loadVersions: (docId: number) => Promise<void>
}

export const useDocumentStore = create<DocumentState>((set) => ({
  versions: [],
  loadVersions: async (docId) => set({ versions: await listVersions(docId) }),
}))
