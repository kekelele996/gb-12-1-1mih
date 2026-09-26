import request from '@/utils/request'
import type { Annotation, Document, DocumentVersion } from '@/types/api'

export function getDocument(id: number | string) { return request.get<never, Document>(`/documents/${id}`) }
export function saveDocument(id: number, payload: { content: string; change_summary?: string }) {
  return request.put<never, Document>(`/documents/${id}`, payload)
}
export function listVersions(id: number) { return request.get<never, DocumentVersion[]>(`/documents/${id}/versions`) }
export function rollbackDocument(id: number, version_no: number) {
  return request.post<never, Document>(`/documents/${id}/rollback`, { version_no })
}
export function listAnnotations(id: number) { return request.get<never, Annotation[]>(`/documents/${id}/annotations`) }
export function addAnnotation(id: number, payload: { content: string; start_offset?: number; end_offset?: number }) {
  return request.post<never, Annotation>(`/documents/${id}/annotations`, payload)
}
