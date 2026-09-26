import request from '@/utils/request'
import type { Message } from '@/types/api'

export function listMessages(unreadOnly = false) {
  return request.get<never, { items: Message[]; unread_count: number }>('/messages', { params: { unread: unreadOnly ? 1 : 0 } })
}
export function sendMessage(payload: { receiver_id: number; content: string }) {
  return request.post<never, Message>('/messages', payload)
}
export function markMessageRead(id: number) { return request.put<never, Message>(`/messages/${id}/read`) }
