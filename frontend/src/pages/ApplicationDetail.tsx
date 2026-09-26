import { Alert, Button, Card, Col, Popconfirm, Row, Select, Space, Tag, message } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import {
  getApplication,
  listDocuments,
  listMaterials,
  listTimeline,
  updateApplicationStatus,
  updateMaterialStatus,
  markTimelineDone,
} from '@/api/application'
import ApplicationStatusTag from '@/components/common/ApplicationStatusTag'
import MaterialProgress from '@/components/common/MaterialProgress'
import Timeline from '@/components/common/Timeline'
import UploadButton from '@/components/common/UploadButton'
import { ApplicationStatusMap } from '@/constants/application'
import { useAuth } from '@/hooks/useAuth'
import type { ApplicationProject, Document, MaterialItem, TimelineNode } from '@/types/api'

// 状态机允许的前向流转（与后端一致）
const NextStatuses: Record<string, string[]> = {
  planning: ['preparing'],
  preparing: ['submitted'],
  submitted: ['waiting'],
  waiting: ['admitted', 'rejected', 'waitlisted'],
}

export default function ApplicationDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user, role } = useAuth()
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

  async function uploadMaterial(item: MaterialItem, fileUrl: string) {
    await updateMaterialStatus(item.id, { status: 'uploaded', file_url: fileUrl })
    message.success(`材料「${item.name}」已上传，等待顾问审核`)
    await load()
  }

  async function approveMaterial(item: MaterialItem) {
    await updateMaterialStatus(item.id, { status: 'approved' })
    message.success(`材料「${item.name}」已审核通过`)
    await load()
  }

  async function doneNode(nodeId: number) {
    await markTimelineDone(nodeId)
    await load()
  }

  if (!app) return <p>加载中…</p>

  const editable = app.status === 'planning' || app.status === 'preparing'
  const isStudent = role === 'student'
  const canReview = role === 'admin' || (role === 'counselor' && user?.id === app.counselor_id)
  // 阻断提交的必交材料（未审核通过）
  const missingRequired = materials.items.filter((m) => m.is_required && m.status !== 'approved')

  function renderMaterialActions(item: MaterialItem) {
    // 学生：准备阶段可上传/重新上传
    if (isStudent && editable && item.status !== 'approved') {
      return <UploadButton onUploaded={(url) => uploadMaterial(item, url)} />
    }
    // 负责顾问/管理员：审核已上传的材料
    if (canReview && item.status === 'uploaded') {
      return (
        <Button size="small" type="primary" onClick={() => approveMaterial(item)}>
          审核通过
        </Button>
      )
    }
    return null
  }

  return (
    <div>
      <h1>
        申请项目 #{app.id} <ApplicationStatusTag status={app.status} />
      </h1>
      {isStudent && !editable && (
        <Alert
          style={{ marginBottom: 16 }}
          type="info"
          showIcon
          message="申请已提交，材料与文书已锁定"
          description="如需修改，请联系负责顾问将申请退回「材料准备中」。"
        />
      )}
      {isStudent && editable && app.status === 'preparing' && missingRequired.length > 0 && (
        <Alert
          style={{ marginBottom: 16 }}
          type="warning"
          showIcon
          message="暂不能提交申请"
          description={`以下必交材料尚未审核通过：${missingRequired.map((m) => m.name).join('、')}`}
        />
      )}
      <Row gutter={16}>
        <Col xs={24} md={8}>
          <Card title="项目信息">
            <p>院校 ID：{app.university_id}</p>
            <p>专业：{app.major}</p>
            <p>轮次：{app.round || '-'}</p>
            <p>
              状态：<ApplicationStatusTag status={app.status} />
            </p>
            <Space direction="vertical" style={{ width: '100%' }}>
              {isStudent && app.status === 'planning' && (
                <Button type="primary" block onClick={() => changeStatus('preparing')}>
                  开始准备材料
                </Button>
              )}
              {isStudent && app.status === 'preparing' && (
                <Popconfirm
                  title="确认提交申请？"
                  description="提交后材料与文书将锁定，无法修改。"
                  onConfirm={() => changeStatus('submitted')}
                >
                  <Button type="primary" block disabled={missingRequired.length > 0}>
                    提交申请
                  </Button>
                </Popconfirm>
              )}
              {canReview && app.status === 'submitted' && (
                <Popconfirm
                  title="退回准备中？"
                  description="退回后学生可重新修改材料与文书。"
                  onConfirm={() => changeStatus('preparing')}
                >
                  <Button danger block>
                    退回准备中
                  </Button>
                </Popconfirm>
              )}
              {!isStudent && (NextStatuses[app.status] || []).length > 0 && (
                <Select
                  style={{ width: '100%' }}
                  placeholder="更新状态"
                  value={undefined}
                  onChange={changeStatus}
                  options={(NextStatuses[app.status] || []).map((s) => ({
                    value: s,
                    label: ApplicationStatusMap[s as keyof typeof ApplicationStatusMap]?.text || s,
                  }))}
                />
              )}
            </Space>
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
            <MaterialProgress items={materials.items} progress={materials.progress} renderActions={renderMaterialActions} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
