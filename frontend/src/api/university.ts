import request from '@/utils/request'
import type { University } from '@/types/api'
import type { PageData } from '@/types/api'

export function listUniversities(params: { page?: number; page_size?: number; country?: string; rank_min?: number; rank_max?: number; keyword?: string }) {
  return request.get<never, PageData<University>>('/universities', { params })
}
export function getUniversity(id: number | string) { return request.get<never, University>(`/universities/${id}`) }
export function createUniversity(payload: Partial<University>) { return request.post<never, University>('/universities', payload) }
export function updateUniversity(id: number, payload: Partial<University>) { return request.put<never, University>(`/universities/${id}`, payload) }
export function deleteUniversity(id: number) { return request.delete<never, { deleted: boolean }>(`/universities/${id}`) }
