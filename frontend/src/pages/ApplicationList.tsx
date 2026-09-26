import { Button, Card, Select, Table, message } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createApplication, listApplications } from '@/api/application'
import { listUniversities } from '@/api/university'
import ApplicationStatusTag from '@/components/common/ApplicationStatusTag'
import { ApplicationStatusList } from '@/constants/application'
import type { ApplicationProject, University } from '@/types/api'

export default function ApplicationList() {
  const navigate = useNavigate()
  const [apps, setApps] = useState<ApplicationProject[]>([])
  const [unis, setUnis] = useState<University[]>([])
  const [status, setStatus] = useState('')
  const [createOpen, setCreateOpen] = useState(false)
  const [form, setForm] = useState({ university_id: 0, major: '', round: '' })

  async function load() {
    const all = await listApplications()
    setApps(status ? all.filter((a) => a.status === status) : all)
  }
  useEffect(() => {
    load()
    listUniversities({ page: 1, page_size: 100 }).then((res) => setUnis(res.list))
  }, [status])

  async function create() {
    if (!form.university_id || !form.major) {
      message.warning('请选择院校并填写专业')
      return
    }
    await createApplication(form)
    message.success('申请项目已创建')
    setCreateOpen(false)
    setForm({ university_id: 0, major: '', round: '' })
    await load()
  }

  return (
    <div>
      <h1>申请项目</h1>
      <div style={{ marginBottom: 12, display: 'flex', gap: 12 }}>
        <Select
          style={{ width: 180 }}
          placeholder="按状态筛选"
          allowClear
          value={status || undefined}
          onChange={(v) => setStatus(v || '')}
          options={ApplicationStatusList.map((s) => ({ value: s, label: s }))}
        />
        <Button type="primary" onClick={() => setCreateOpen(true)}>
          创建申请项目
        </Button>
      </div>
      <Table
        rowKey="id"
        dataSource={apps}
        pagination={false}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '院校 ID', dataIndex: 'university_id', width: 90 },
          { title: '专业', dataIndex: 'major' },
          { title: '轮次', dataIndex: 'round' },
          { title: '状态', dataIndex: 'status', render: (s: string) => <ApplicationStatusTag status={s} /> },
          {
            title: '操作',
            render: (_, r) => (
              <Button size="small" onClick={() => navigate(`/applications/${r.id}`)}>
                详情
              </Button>
            ),
          },
        ]}
      />
      {createOpen && (
        <Card title="创建申请项目" style={{ marginTop: 16, maxWidth: 480 }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            <Select
              placeholder="选择院校"
              value={form.university_id || undefined}
              onChange={(v) => setForm({ ...form, university_id: v })}
              options={unis.map((u) => ({ value: u.id, label: `${u.name} (${u.country})` }))}
            />
            <input
              style={{ padding: 8 }}
              placeholder="专业"
              value={form.major}
              onChange={(e) => setForm({ ...form, major: e.target.value })}
            />
            <input
              style={{ padding: 8 }}
              placeholder="轮次（如：秋季）"
              value={form.round}
              onChange={(e) => setForm({ ...form, round: e.target.value })}
            />
            <div style={{ display: 'flex', gap: 8 }}>
              <Button type="primary" onClick={create}>
                创建
              </Button>
              <Button onClick={() => setCreateOpen(false)}>取消</Button>
            </div>
          </div>
        </Card>
      )}
    </div>
  )
}
