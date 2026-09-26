export type MaterialStatus = 'pending' | 'uploaded' | 'approved'

export const MaterialStatusMap: Record<MaterialStatus, { text: string; color: string }> = {
  pending: { text: '待上传', color: 'gold' },
  uploaded: { text: '已上传', color: 'blue' },
  approved: { text: '已审核', color: 'green' },
}
