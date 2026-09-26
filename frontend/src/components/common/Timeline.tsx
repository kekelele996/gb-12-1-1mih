import { Button, Tag, Timeline as AntTimeline } from 'antd'
import type { TimelineNode } from '@/types/api'
import { formatDate } from '@/utils/dateFormat'

const typeMap: Record<string, { color: string; text: string }> = {
  language_exam: { color: 'blue', text: '语言/标化' },
  essay_deadline: { color: 'purple', text: '文书截止' },
  application_deadline: { color: 'red', text: '申请截止' },
  interview: { color: 'cyan', text: '面试' },
  result: { color: 'green', text: '结果' },
}

export default function Timeline({
  nodes,
  onDone,
}: {
  nodes: TimelineNode[]
  onDone?: (id: number) => void
}) {
  return (
    <AntTimeline
      items={nodes.map((n) => {
        const t = typeMap[n.node_type] || { color: 'gray', text: n.node_type }
        return {
          color: n.is_done ? 'green' : t.color,
          children: (
            <div>
              <div>
                <Tag color={t.color}>{t.text}</Tag>
                <b>{n.title}</b>
                {n.is_done && <Tag color="green" style={{ marginLeft: 8 }}>已完成</Tag>}
              </div>
              <div style={{ color: '#888', fontSize: 12 }}>截止 {formatDate(n.due_date)}</div>
              {!n.is_done && onDone && (
                <Button size="small" type="link" onClick={() => onDone(n.id)}>
                  标记完成
                </Button>
              )}
            </div>
          ),
        }
      })}
    />
  )
}
