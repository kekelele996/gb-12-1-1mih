import request from '@/utils/request'
import type { ApplicationProject, Document, MaterialItem, TimelineNode } from '@/types/api'

export function listApplications() { return request.get<never, ApplicationProject[]>('/applications') }
export function getApplication(id: number | string) { return request.get<never, ApplicationProject>(`/applications/${id}`) }
export function createApplication(payload: { university_id: number; major: string; round?: string }) {
  return request.post<never, ApplicationProject>('/applications', payload)
}
export function updateApplicationStatus(id: number, status: string) {
  return request.put<never, ApplicationProject>(`/applications/${id}/status`, { status })
}
export function listDocuments(appId: number) { return request.get<never, Document[]>(`/applications/${appId}/documents`) }
export function createDocument(appId: number, payload: { doc_type: string; title: string; content?: string }) {
  return request.post<never, Document>(`/applications/${appId}/documents`, payload)
}
export function listMaterials(appId: number) {
  return request.get<never, { items: MaterialItem[]; progress: number }>(`/applications/${appId}/materials`)
}
export function createMaterial(appId: number, payload: { name: string; category?: string; is_required?: boolean }) {
  return request.post<never, MaterialItem>(`/applications/${appId}/materials`, payload)
}
export function updateMaterialStatus(id: number, payload: { status: string; file_url?: string }) {
  return request.put<never, MaterialItem>(`/materials/${id}/status`, payload)
}
export function listTimeline(appId: number) { return request.get<never, TimelineNode[]>(`/applications/${appId}/timeline`) }
export function createTimeline(appId: number, payload: { title: string; node_type?: string; due_date: string }) {
  return request.post<never, TimelineNode>(`/applications/${appId}/timeline`, payload)
}
export function markTimelineDone(id: number) { return request.put<never, TimelineNode>(`/timeline/${id}/done`) }
