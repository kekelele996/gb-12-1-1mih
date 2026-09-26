import { Card } from 'antd'
import { useEffect, useState } from 'react'
import { listApplications, listTimeline, markTimelineDone } from '@/api/application'
import Timeline from '@/components/common/Timeline'
import type { TimelineNode } from '@/types/api'

export default function TimelinePage() {
  const [nodes, setNodes] = useState<TimelineNode[]>([])

  async function load() {
    const a = await listApplications()
    const all: TimelineNode[] = []
    for (const app of a) {
      all.push(...(await listTimeline(app.id)))
    }
    setNodes(all)
  }
  useEffect(() => {
    load()
  }, [])

  async function doneNode(id: number) {
    await markTimelineDone(id)
    await load()
  }

  return (
    <div>
      <h1>时间线管理</h1>
      <Card>
        <Timeline nodes={nodes} onDone={doneNode} />
        {!nodes.length && <p style={{ color: '#999' }}>暂无时间线节点</p>}
      </Card>
    </div>
  )
}
