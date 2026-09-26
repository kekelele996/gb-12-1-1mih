import { Button, Card, Input, List, Segmented, message } from 'antd'
import { useEffect, useState } from 'react'
import { listMessages, sendMessage, markMessageRead } from '@/api/message'
import { useAuth } from '@/hooks/useAuth'
import type { Message } from '@/types/api'
import { formatDateTime } from '@/utils/dateFormat'

export default function Messages() {
  const { user } = useAuth()
  const [items, setItems] = useState<Message[]>([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [filter, setFilter] = useState<'all' | 'unread'>('all')
  const [content, setContent] = useState('')
  const [target, setTarget] = useState('')

  async function load() {
    const res = await listMessages(filter === 'unread')
    setItems(res.items)
    setUnreadCount(res.unread_count)
  }
  useEffect(() => {
    load()
  }, [filter])

  async function send() {
    if (!target || !content.trim()) {
      message.warning('请填写接收者 ID 和内容')
      return
    }
    await sendMessage({ receiver_id: Number(target), content })
    message.success('消息已发送')
    setContent('')
    await load()
  }

  async function read(m: Message) {
    if (!m.is_read) {
      await markMessageRead(m.id)
      await load()
    }
  }

  return (
    <div style={{ maxWidth: 800, margin: '0 auto' }}>
      <h1>站内消息（未读 {unreadCount}）</h1>
      <Segmented
        value={filter}
        onChange={(v) => setFilter(v as 'all' | 'unread')}
        options={[
          { label: '全部', value: 'all' },
          { label: '未读', value: 'unread' },
        ]}
        style={{ marginBottom: 12 }}
      />
      <Card>
        <List
          dataSource={items}
          locale={{ emptyText: '暂无消息' }}
          renderItem={(m) => (
            <List.Item onClick={() => read(m)} style={{ cursor: 'pointer' }}>
              <div>
                <div>
                  {m.sender_id === 0 ? <b>系统</b> : `用户 #${m.sender_id}`} → 我 · {formatDateTime(m.created_at)}
                  {!m.is_read && <span style={{ color: '#f5222d', marginLeft: 8 }}>● 未读</span>}
                </div>
                <div style={{ marginTop: 4 }}>{m.content}</div>
              </div>
            </List.Item>
          )}
        />
      </Card>
      <Card title={`发消息（我是 #${user?.id || '-'}）`} style={{ marginTop: 16 }}>
        <Input
          placeholder="接收者 ID"
          value={target}
          onChange={(e) => setTarget(e.target.value)}
          style={{ marginBottom: 8 }}
        />
        <Input.TextArea rows={3} value={content} onChange={(e) => setContent(e.target.value)} placeholder="消息内容" />
        <Button type="primary" style={{ marginTop: 8 }} onClick={send}>
          发送
        </Button>
      </Card>
    </div>
  )
}
