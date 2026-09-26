export type DocumentType = 'ps' | 'rl' | 'cv' | 'essay'

export const DocumentTypeMap: Record<DocumentType, string> = {
  ps: '个人陈述',
  rl: '推荐信',
  cv: '简历',
  essay: 'Essay',
}
