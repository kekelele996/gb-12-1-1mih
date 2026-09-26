export interface ApiResponse<T = unknown> { code: number; message: string; data: T }
export interface PageData<T> { list: T[]; total: number; page: number; page_size: number }

export interface UserInfo {
  id: number
  username: string
  email: string
  real_name: string
  phone: string
  role: 'student' | 'counselor' | 'admin'
  created_at: string
}

export interface University {
  id: number
  name: string
  country: string
  city: string
  ranking: number
  top_majors: string
  application_deadline: string
  tuition_range: string
  requirements: string
  created_at: string
}

export interface ApplicationProject {
  id: number
  student_id: number
  counselor_id: number
  university_id: number
  major: string
  round: string
  status: string
  created_at: string
  updated_at: string
}

export interface Document {
  id: number
  application_id: number
  doc_type: string
  title: string
  content: string
  current_version: number
  created_at: string
  updated_at: string
}

export interface DocumentVersion {
  id: number
  document_id: number
  content: string
  version_no: number
  change_summary: string
  created_by: number
  created_at: string
}

export interface Annotation {
  id: number
  document_id: number
  counselor_id: number
  content: string
  start_offset: number
  end_offset: number
  created_at: string
}

export interface MaterialItem {
  id: number
  application_id: number
  name: string
  category: string
  is_required: boolean
  status: string
  file_url: string
  uploaded_at: string
  created_at: string
}

export interface TimelineNode {
  id: number
  application_id: number
  title: string
  node_type: string
  due_date: string
  is_done: boolean
  reminder_sent: boolean
  created_at: string
}

export interface Recommendation {
  id: number
  student_id: number
  counselor_id: number
  university_ids: string
  reason: string
  created_at: string
}

export interface Message {
  id: number
  sender_id: number
  receiver_id: number
  content: string
  is_read: boolean
  created_at: string
}
