import { List, Tag } from 'antd'
import type { Annotation } from '@/types/api'
import { formatDateTime } from '@/utils/dateFormat'

export default function AnnotationList({ annotations }: { annotations: Annotation[] }) {
  return (
    <List
      dataSource={annotations}
      locale={{ emptyText: '暂无批注' }}
      renderItem={(a) => (
        <List.Item>
          <div>
            <div>
              <Tag color="gold">顾问 #{a.counselor_id}</Tag>
              <span style={{ color: '#999', fontSize: 12 }}>{formatDateTime(a.created_at)}</span>
            </div>
            <div>{a.content}</div>
            <div style={{ color: '#bbb', fontSize: 12 }}>
              位置 {a.start_offset}-{a.end_offset}
            </div>
          </div>
        </List.Item>
      )}
    />
  )
}
