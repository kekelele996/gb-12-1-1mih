import request from '@/utils/request'
import type { University } from '@/types/api'

export interface RecommendationView {
  id: number
  reason: string
  created_at: string
  universities: University[]
}

export function listRecommendations(studentId: number) {
  return request.get<never, RecommendationView[]>(`/recommendations/student/${studentId}`)
}
export function createRecommendation(payload: { student_id: number; university_ids: number[]; reason?: string }) {
  return request.post<never, { id: number }>('/recommendations', payload)
}
