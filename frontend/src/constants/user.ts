export type UserRole = 'student' | 'counselor' | 'admin'

export const UserRoleMap: Record<UserRole, string> = {
  student: '学生',
  counselor: '顾问',
  admin: '管理员',
}
