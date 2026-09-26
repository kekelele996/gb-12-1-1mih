import { useEffect, useState } from 'react'
import { listVersions } from '@/api/document'
import type { DocumentVersion } from '@/types/api'

export function useDocumentVersions(docId: number) {
  const [versions, setVersions] = useState<DocumentVersion[]>([])
  async function load() {
    setVersions(await listVersions(docId))
  }
  useEffect(() => {
    if (docId) load()
  }, [docId])
  return { versions, reload: load }
}
