import { Alert, Button, Card, Col, Input, Row, Select, Space, Tag, message } from 'antd'
import { LockOutlined } from '@ant-design/icons'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  addAnnotation,
  getDocument,
  listAnnotations,
  listVersions,
  rollbackDocument,
  saveDocument,
} from '@/api/document'
import { getApplication, getApplicationAbilities } from '@/api/application'
import AnnotationList from '@/components/common/AnnotationList'
import DocumentDiff from '@/components/common/DocumentDiff'
import { useAuth } from '@/hooks/useAuth'
import type { Annotation, ApplicationAbilities, Document, DocumentVersion } from '@/types/api'

export default function DocumentEditor() {
  const { id } = useParams()
  const { role } = useAuth()
  const [doc, setDoc] = useState<Document | null>(null)
  const [abilities, setAbilities] = useState<ApplicationAbilities | null>(null)
  const [content, setContent] = useState('')
  const [versions, setVersions] = useState<DocumentVersion[]>([])
  const [annotations, setAnnotations] = useState<Annotation[]>([])
  const [compareVersion, setCompareVersion] = useState<number | undefined>()
  const [compareText, setCompareText] = useState('')
  const [showDiff, setShowDiff] = useState(false)
  const [annContent, setAnnContent] = useState('')

  async function load() {
    const d = await getDocument(id!)
    setDoc(d)
    setContent(d.content)
    setVersions(await listVersions(d.id))
    setAnnotations(await listAnnotations(d.id))
    const app = await getApplication(d.application_id)
    setAbilities(await getApplicationAbilities(app.id))
  }
  useEffect(() => {
    load()
  }, [id])

  const canEdit = !!abilities?.can_edit_documents
  const canAnnotate = !!abilities?.is_managing_staff && (role === 'counselor' || role === 'admin')

  async function save() {
    await saveDocument(doc!.id, { content, change_summary: `v${(doc!.current_version + 1)} 编辑` })
    message.success('文书已保存')
    await load()
  }

  async function rollback(versionNo: number) {
    await rollbackDocument(doc!.id, versionNo)
    message.success('已回滚')
    await load()
  }

  async function submitAnnotation() {
    if (!annContent.trim()) return
    await addAnnotation(doc!.id, { content: annContent })
    setAnnContent('')
    await load()
  }

  if (!doc || !abilities) return <p>加载中…</p>
  return (
    <div>
      <h1>
        文书编辑：{doc.title} {!canEdit && <LockOutlined style={{ color: '#999' }} />}
      </h1>
      {role === 'student' && !canEdit && (
        <Alert
          style={{ marginBottom: 16 }}
          type="info"
          showIcon
          message="申请已提交，文书已锁定（只读）"
          description="负责顾问将申请退回「材料准备中」后，才能继续编辑或回滚版本。"
        />
      )}
      <Row gutter={16}>
        <Col xs={24} md={16}>
          <Card
            title={`当前版本 v${doc.current_version}`}
            extra={canEdit ? <Button type="primary" onClick={save}>保存新版本</Button> : <Tag>只读</Tag>}
          >
            <Input.TextArea rows={16} value={content} disabled={!canEdit} onChange={(e) => setContent(e.target.value)} />
          </Card>
          <Card title="版本对比" style={{ marginTop: 16 }}>
            <Space>
              <Select
                placeholder="选择历史版本"
                style={{ width: 200 }}
                value={compareVersion}
                onChange={(v) => {
                  setCompareVersion(v)
                  const ver = versions.find((x) => x.version_no === v)
                  setCompareText(ver?.content || '')
                }}
                options={versions.map((v) => ({ value: v.version_no, label: `v${v.version_no} ${v.change_summary}` }))}
              />
              <Button onClick={() => setShowDiff(!showDiff)}>对比 Diff</Button>
              {canEdit && compareVersion ? (
                <Button danger onClick={() => rollback(compareVersion)}>
                  回滚到此版本
                </Button>
              ) : null}
            </Space>
            {showDiff && compareText !== undefined && <DocumentDiff oldText={compareText} newText={content} />}
          </Card>
        </Col>
        <Col xs={24} md={8}>
          <Card title="版本历史" style={{ marginBottom: 16 }}>
            {versions.map((v) => (
              <div key={v.id} style={{ marginBottom: 8 }}>
                <b>v{v.version_no}</b> {v.change_summary}
              </div>
            ))}
          </Card>
          <Card title="顾问批注">
            <AnnotationList annotations={annotations} />
            {canAnnotate && (
              <>
                <Input.TextArea rows={3} value={annContent} onChange={(e) => setAnnContent(e.target.value)} placeholder="添加批注…" style={{ marginTop: 8 }} />
                <Button type="primary" size="small" style={{ marginTop: 8 }} onClick={submitAnnotation}>
                  添加批注
                </Button>
              </>
            )}
          </Card>
        </Col>
      </Row>
    </div>
  )
}
