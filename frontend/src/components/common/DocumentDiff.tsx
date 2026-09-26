import { List } from 'antd'
import { diffLines } from '@/utils/diff'

export default function DocumentDiff({ oldText, newText }: { oldText: string; newText: string }) {
  const lines = diffLines(oldText, newText)
  return (
    <List
      size="small"
      dataSource={lines}
      renderItem={(l) => (
        <List.Item style={{ background: l.type === 'add' ? '#f6ffed' : l.type === 'remove' ? '#fff1f0' : 'transparent' }}>
          <span style={{ fontFamily: 'monospace' }}>
            {l.type === 'add' ? '+ ' : l.type === 'remove' ? '- ' : '  '}
            {l.text}
          </span>
        </List.Item>
      )}
    />
  )
}
