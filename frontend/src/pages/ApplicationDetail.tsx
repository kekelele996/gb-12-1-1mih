import {
  Alert,
  Button,
  Card,
  Col,
  Form,
  Input,
  Modal,
  Popconfirm,
  Row,
  Select,
  Space,
  Switch,
  Tag,
  message,
} from 'antd'
import { LockOutlined } from '@ant-design/icons'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import {
  createDocument,
  createMaterial,
  getApplication,
  getApplicationAbilities,
  listDocuments,
  listMaterials,
  listTimeline,
  markTimelineDone,
  updateApplicationStatus,
  updateMaterialStatus,
} from '@/api/application'
import ApplicationStatusTag from '@/components/common/ApplicationStatusTag'
import MaterialProgress from '@/components/common/MaterialProgress'
import Timeline from '@/components/common/Timeline'
import { useAuth } from '@/hooks/useAuth'
import type {
  ApplicationAbilities,
  ApplicationProject,
  Document,
  MaterialItem,
  TimelineNode,
} from '@/types/api'

const DOC_TYPE_OPTIONS = [
  { value: 'ps', label: '个人陈述 (PS)' },
  { value: 'rl', label: '推荐信 (RL)' },
  { value: 'cv', label: '简历 (CV)' },
  { value: 'essay', label: 'Essay' },
]

const STAFF_NEXT: Record<string, { value: string; label: string }[]> = {
  submitted: [{ value: 'waiting', label: '标记为等待结果' }],
  waiting: [
    { value: 'admitted', label: '标记为已录取' },
    { value: 'rejected', label: '标记为已拒绝' },
    { value: 'waitlisted', label: '标记为候补' },
  ],
}

export default function ApplicationDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { role } = useAuth()
  const [app, setApp] = useState<ApplicationProject | null>(null)
  const [abilities, setAbilities] = useState<ApplicationAbilities | null>(null)
  const [docs, setDocs] = useState<Document[]>([])
  const [materials, setMaterials] = useState<{ items: MaterialItem[]; progress: number }>({ items: [], progress: 0 })
  const [nodes, setNodes] = useState<TimelineNode[]>([])
  const [docModalOpen, setDocModalOpen] = useState(false)
  const [matModalOpen, setMatModalOpen] = useState(false)
  const [docForm] = Form.useForm()
  const [matForm] = Form.useForm()

  async function load() {
    const a = await getApplication(id!)
    setApp(a)
    setAbilities(await getApplicationAbilities(a.id))
    setDocs(await listDocuments(a.id))
    setMaterials(await listMaterials(a.id))
    setNodes(await listTimeline(a.id))
  }
  useEffect(() => {
    load()
  }, [id])

  async function changeStatus(status: string, okText?: string) {
    await updateApplicationStatus(app!.id, status)
    message.success(okText ?? '状态已更新')
    await load()
  }

  async function submitApp() {
    await changeStatus('submitted', '申请已提交，材料与文书已锁定')
  }

  async function returnForRevision() {
    await changeStatus('preparing', '已退回准备中，学生可继续修改')
  }

  async function doneNode(nodeId: number) {
    await markTimelineDone(nodeId)
    await load()
  }

  async function handleCreateDoc() {
    const v = await docForm.validateFields()
    await createDocument(app!.id, v)
    message.success('文书已创建')
    setDocModalOpen(false)
    docForm.resetFields()
    await load()
  }

  async function handleCreateMaterial() {
    const v = await matForm.validateFields()
    await createMaterial(app!.id, { ...v, is_required: v.is_required ?? true })
    message.success('材料项已添加')
    setMatModalOpen(false)
    matForm.resetFields()
    await load()
  }

  if (!app || !abilities) return <p>加载中…</p>

  const staffNext = STAFF_NEXT[app.status] || []
  const missingNames = abilities.missing_required_materials.map((m) => m.name).join('、')

  return (
    <div>
      <h1>
        申请项目 #{app.id} <ApplicationStatusTag status={app.status} />
      </h1>

      {/* 提交后锁定提示 / 退回后恢复提示 */}
      {abilities.is_student_owner && !abilities.editable && (
        <Alert
          style={{ marginBottom: 16 }}
          type="info"
          showIcon
          icon={<LockOutlined />}
          message="申请已提交，材料和文书已锁定"
          description="如需修改，请等待负责顾问将申请退回「材料准备中」。"
        />
      )}
      {abilities.is_student_owner && abilities.editable && app.status === 'preparing' && (
        <Alert
          style={{ marginBottom: 16 }}
          type="warning"
          showIcon
          message="材料准备中：提交前所有必交材料须由学生上传并经顾问/管理员审核通过"
          description={
            missingNames
              ? `尚未通过审核的必交材料：${missingNames}`
              : '全部必交材料已审核通过，可以提交申请。'
          }
        />
      )}

      {/* 阻断原因列表（按角色） */}
      {abilities.blocked.length > 0 && (
        <Alert
          style={{ marginBottom: 16 }}
          type="error"
          showIcon
          message="当前不可执行的操作"
          description={
            <ul style={{ margin: 0, paddingLeft: 18 }}>
              {abilities.blocked.map((b, i) => (
                <li key={i}>{b.reason}</li>
              ))}
            </ul>
          }
        />
      )}

      <Row gutter={16}>
        <Col xs={24} md={8}>
          <Card title="项目信息">
            <p>院校 ID：{app.university_id}</p>
            <p>专业：{app.major}</p>
            <p>轮次：{app.round || '-'}</p>
            <p>状态：<ApplicationStatusTag status={app.status} /></p>

            {/* 学生操作 */}
            {role === 'student' && (
              <Space direction="vertical" style={{ width: '100%' }}>
                {app.status === 'planning' && (
                  <Button type="primary" block onClick={() => changeStatus('preparing')}>
                    开始准备材料
                  </Button>
                )}
                {abilities.can_submit ? (
                  <Popconfirm
                    title="提交后材料和文书将锁定，不能再修改。确认提交？"
                    onConfirm={submitApp}
                    okText="确认提交"
                    cancelText="再检查一下"
                  >
                    <Button type="primary" block danger>提交申请</Button>
                  </Popconfirm>
                ) : app.status === 'preparing' ? (
                  <Button type="primary" block danger disabled title={missingNames ? `缺少：${missingNames}` : ''}>
                    提交申请（必交材料未齐）
                  </Button>
                ) : null}
              </Space>
            )}

            {/* 顾问/管理员操作 */}
            {abilities.is_managing_staff && (role === 'counselor' || role === 'admin') && (
              <Space direction="vertical" style={{ width: '100%' }}>
                {abilities.can_return && (
                  <Popconfirm
                    title="退回后学生可重新修改材料和文书，确认退回？"
                    onConfirm={returnForRevision}
                    okText="退回准备中"
                    cancelText="取消"
                  >
                    <Button block>退回「材料准备中」</Button>
                  </Popconfirm>
                )}
                {staffNext.length > 0 && (
                  <Select
                    style={{ width: '100%' }}
                    placeholder="更新申请状态"
                    value={undefined}
                    onChange={(v) => changeStatus(v ?? '')}
                    options={staffNext}
                  />
                )}
                {!abilities.can_return && staffNext.length === 0 && (
                  <p style={{ color: '#999' }}>当前状态暂无可执行流转</p>
                )}
              </Space>
            )}
          </Card>
          <Card title="时间线" style={{ marginTop: 16 }}>
            <Timeline nodes={nodes} onDone={doneNode} />
          </Card>
        </Col>
        <Col xs={24} md={16}>
          <Card
            title="文书"
            style={{ marginBottom: 16 }}
            extra={
              abilities.can_edit_documents ? (
                <Button type="primary" size="small" onClick={() => setDocModalOpen(true)}>新建文书</Button>
              ) : <Tag icon={<LockOutlined />}>{role === 'student' ? '已锁定·只读' : '只读'}</Tag>
            }
          >
            {docs.map((d) => (
              <Button key={d.id} style={{ margin: 4 }} onClick={() => navigate(`/documents/${d.id}`)}>
                <Tag>{d.doc_type.toUpperCase()}</Tag> {d.title} (v{d.current_version})
              </Button>
            ))}
            {!docs.length && <p style={{ color: '#999' }}>暂无文书</p>}
          </Card>
          <Card
            title="材料清单"
            extra={
              abilities.can_manage_materials ? (
                <Button type="primary" size="small" onClick={() => setMatModalOpen(true)}>添加材料</Button>
              ) : null
            }
          >
            <MaterialProgress
              items={materials.items}
              progress={materials.progress}
              canUpload={abilities.editable}
              canReview={abilities.can_review_materials}
              onChange={load}
              onUpdateStatus={(mid, payload) => updateMaterialStatus(mid, payload)}
            />
          </Card>
        </Col>
      </Row>

      <Modal title="新建文书" open={docModalOpen} onOk={handleCreateDoc} onCancel={() => setDocModalOpen(false)} okText="创建" cancelText="取消">
        <Form form={docForm} layout="vertical" initialValues={{ doc_type: 'ps' }}>
          <Form.Item name="doc_type" label="文书类型" rules={[{ required: true }]}>
            <Select options={DOC_TYPE_OPTIONS} />
          </Form.Item>
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input maxLength={255} placeholder="如：个人陈述" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="添加材料项" open={matModalOpen} onOk={handleCreateMaterial} onCancel={() => setMatModalOpen(false)} okText="添加" cancelText="取消">
        <Form form={matForm} layout="vertical" initialValues={{ is_required: true }}>
          <Form.Item name="name" label="材料名称" rules={[{ required: true, message: '请输入材料名称' }]}>
            <Input maxLength={128} placeholder="如：成绩单、语言成绩" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Input maxLength={64} placeholder="如：学术材料 / 标化 / 文书" />
          </Form.Item>
          <Form.Item name="is_required" label="是否必交" valuePropName="checked">
            <Switch checkedChildren="必交" unCheckedChildren="选交" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
