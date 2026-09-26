import { Button, Card, Col, Row, Select, Tag, message } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import {
  getApplication,
  listDocuments,
  listMaterials,
  listTimeline,
  updateApplicationStatus,
  markTimelineDone,
} from '@/api/application'
import ApplicationStatusTag from '@/components/common/ApplicationStatusTag'
import MaterialProgress from '@/components/common/MaterialProgress'
import Timeline from '@/components/common/Timeline'
import { ApplicationStatusList } from '@/constants/application'
import type { ApplicationProject, Document, MaterialItem, TimelineNode } from '@/types/api'

export default function ApplicationDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [app, setApp] = useState<ApplicationProject | null>(null)
  const [docs, setDocs] = useState<Document[]>([])
  const [materials, setMaterials] = useState<{ items: MaterialItem[]; progress: number }>({ items: [], progress: 0 })
  const [nodes, setNodes] = useState<TimelineNode[]>([])

  async function load() {
    const a = await getApplication(id!)
    setApp(a)
    setDocs(await listDocuments(a.id))
    setMaterials(await listMaterials(a.id))
    setNodes(await listTimeline(a.id))
  }
  useEffect(() => {
    load()
  }, [id])

  async function changeStatus(status: string) {
    await updateApplicationStatus(app!.id, status)
    message.success('状态已更新')
    await load()
  }

  async function doneNode(nodeId: number) {
    await markTimelineDone(nodeId)
    await load()
  }

  if (!app) return <p>加载中…</p>
  return (
    <div>
      <h1>
        申请项目 #{app.id} <ApplicationStatusTag status={app.status} />
      </h1>
      <Row gutter={16}>
        <Col xs={24} md={8}>
          <Card title="项目信息">
            <p>院校 ID：{app.university_id}</p>
            <p>专业：{app.major}</p>
            <p>轮次：{app.round || '-'}</p>
            <p>状态：<ApplicationStatusTag status={app.status} /></p>
            <Select
              style={{ width: '100%' }}
              placeholder="更新状态"
              value={undefined}
              onChange={changeStatus}
              options={ApplicationStatusList.map((s) => ({ value: s, label: s }))}
            />
          </Card>
          <Card title="时间线" style={{ marginTop: 16 }}>
            <Timeline nodes={nodes} onDone={doneNode} />
          </Card>
        </Col>
        <Col xs={24} md={16}>
          <Card title="文书" style={{ marginBottom: 16 }}>
            {docs.map((d) => (
              <Button key={d.id} style={{ margin: 4 }} onClick={() => navigate(`/documents/${d.id}`)}>
                <Tag>{d.doc_type.toUpperCase()}</Tag> {d.title} (v{d.current_version})
              </Button>
            ))}
            {!docs.length && <p style={{ color: '#999' }}>暂无文书</p>}
          </Card>
          <Card title="材料清单">
            <MaterialProgress items={materials.items} progress={materials.progress} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
