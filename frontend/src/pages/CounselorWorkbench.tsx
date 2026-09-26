import { Button, Card, Col, Row, Select, Table, message } from 'antd'
import { useEffect, useState } from 'react'
import { listStudents } from '@/api/user'
import { listApplications } from '@/api/application'
import { listRecommendations, createRecommendation } from '@/api/recommendation'
import { listUniversities } from '@/api/university'
import ApplicationStatusTag from '@/components/common/ApplicationStatusTag'
import type { ApplicationProject, University, UserInfo } from '@/types/api'

export default function CounselorWorkbench() {
  const [students, setStudents] = useState<UserInfo[]>([])
  const [apps, setApps] = useState<ApplicationProject[]>([])
  const [unis, setUnis] = useState<University[]>([])
  const [selectedStudent, setSelectedStudent] = useState<number | undefined>()
  const [selectedUnis, setSelectedUnis] = useState<number[]>([])
  const [reason, setReason] = useState('')
  const [recommendations, setRecommendations] = useState<{ id: number; reason: string; universities: University[] }[]>([])

  async function loadApps() {
    setApps(await listApplications())
  }
  useEffect(() => {
    listStudents().then(setStudents)
    listUniversities({ page: 1, page_size: 100 }).then((res) => setUnis(res.list))
    loadApps()
  }, [])

  async function loadRecs() {
    if (!selectedStudent) return
    setRecommendations(await listRecommendations(selectedStudent))
  }
  useEffect(() => {
    loadRecs()
  }, [selectedStudent])

  async function submit() {
    if (!selectedStudent || !selectedUnis.length) {
      message.warning('请选择学生和院校')
      return
    }
    await createRecommendation({ student_id: selectedStudent, university_ids: selectedUnis, reason })
    message.success('选校方案已推荐')
    setSelectedUnis([])
    setReason('')
    await loadRecs()
  }

  return (
    <div>
      <h1>顾问工作台</h1>
      <Row gutter={16}>
        <Col xs={24} md={12}>
          <Card title="绑定学生">
            <Select
              style={{ width: '100%' }}
              placeholder="选择学生"
              value={selectedStudent}
              onChange={setSelectedStudent}
              options={students.map((s) => ({ value: s.id, label: `${s.real_name || s.username} (#${s.id})` }))}
            />
          </Card>
          <Card title="推荐选校" style={{ marginTop: 16 }}>
            <Select
              mode="multiple"
              style={{ width: '100%' }}
              placeholder="选择推荐院校"
              value={selectedUnis}
              onChange={setSelectedUnis}
              options={unis.map((u) => ({ value: u.id, label: `${u.name} (${u.country})` }))}
            />
            <input
              style={{ width: '100%', marginTop: 8, padding: 8 }}
              placeholder="推荐理由"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
            />
            <Button type="primary" style={{ marginTop: 8 }} onClick={submit}>
              生成选校方案
            </Button>
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="全部申请项目">
            <Table
              size="small"
              rowKey="id"
              dataSource={apps}
              pagination={false}
              columns={[
                { title: 'ID', dataIndex: 'id', width: 60 },
                { title: '学生', dataIndex: 'student_id', width: 80 },
                { title: '院校', dataIndex: 'university_id', width: 90 },
                { title: '专业', dataIndex: 'major' },
                { title: '状态', dataIndex: 'status', render: (s: string) => <ApplicationStatusTag status={s} /> },
              ]}
            />
          </Card>
          <Card title="已推荐方案" style={{ marginTop: 16 }}>
            {recommendations.map((r) => (
              <Card key={r.id} size="small" style={{ marginBottom: 8 }}>
                <div>{r.universities.map((u) => u.name).join('、')}</div>
                <div style={{ color: '#888' }}>{r.reason}</div>
              </Card>
            ))}
            {!recommendations.length && <p style={{ color: '#999' }}>暂无推荐方案</p>}
          </Card>
        </Col>
      </Row>
    </div>
  )
}
