export type ApplicationStatus =
  | 'planning'
  | 'preparing'
  | 'submitted'
  | 'waiting'
  | 'admitted'
  | 'rejected'
  | 'waitlisted'

export const ApplicationStatusMap: Record<ApplicationStatus, { text: string; color: string }> = {
  planning: { text: '规划中', color: 'default' },
  preparing: { text: '材料准备中', color: 'blue' },
  submitted: { text: '已提交', color: 'cyan' },
  waiting: { text: '等待结果', color: 'gold' },
  admitted: { text: '已录取', color: 'green' },
  rejected: { text: '已拒绝', color: 'red' },
  waitlisted: { text: '候补名单', color: 'orange' },
}

export const ApplicationStatusList = Object.keys(ApplicationStatusMap) as ApplicationStatus[]
