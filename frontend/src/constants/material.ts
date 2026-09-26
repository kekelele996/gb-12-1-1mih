export type MaterialStatus = 'pending' | 'uploaded' | 'approved' | 'rejected'

export const MaterialStatusMap: Record<MaterialStatus, { text: string; color: string }> = {
  pending: { text: '待上传', color: 'gold' },
  uploaded: { text: '已上传·待审核', color: 'blue' },
  approved: { text: '已审核', color: 'green' },
  rejected: { text: '审核不通过', color: 'red' },
}
