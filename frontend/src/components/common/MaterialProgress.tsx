import { Button, Input, Modal, Popconfirm, Progress, Space, Table, Tag, message } from 'antd'
import { UploadOutlined } from '@ant-design/icons'
import { useState } from 'react'
import type { ColumnsType } from 'antd/es/table'
import UploadButton from '@/components/common/UploadButton'
import type { MaterialItem } from '@/types/api'
import { MaterialStatusMap } from '@/constants/material'

interface Props {
  items: MaterialItem[]
  progress: number
  canUpload: boolean
  canReview: boolean
  onChange: () => Promise<void> | void
  onUpdateStatus: (id: number, payload: { status: string; file_url?: string; review_remark?: string }) => Promise<unknown>
}

export default function MaterialProgress({ items, progress, canUpload, canReview, onChange, onUpdateStatus }: Props) {
  const [rejecting, setRejecting] = useState<MaterialItem | null>(null)
  const [remark, setRemark] = useState('')

  async function handleUpload(m: MaterialItem, url: string) {
    await onUpdateStatus(m.id, { status: 'uploaded', file_url: url })
    message.success('材料已上传，等待顾问审核')
    await onChange()
  }

  async function approve(m: MaterialItem) {
    await onUpdateStatus(m.id, { status: 'approved' })
    message.success('已审核通过')
    await onChange()
  }

  async function confirmReject() {
    if (!rejecting) return
    await onUpdateStatus(rejecting.id, { status: 'rejected', review_remark: remark })
    message.success('已退回学生重新上传')
    setRejecting(null)
    setRemark('')
    await onChange()
  }

  const columns: ColumnsType<MaterialItem> = [
    { title: '材料', dataIndex: 'name' },
    { title: '分类', dataIndex: 'category' },
    {
      title: '必交',
      dataIndex: 'is_required',
      width: 70,
      render: (v: boolean) => (v ? <Tag color="red">必交</Tag> : '选交'),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 130,
      render: (s: string) => {
        const m = MaterialStatusMap[s as keyof typeof MaterialStatusMap]
        return m ? <Tag color={m.color}>{m.text}</Tag> : <Tag>{s}</Tag>
      },
    },
    {
      title: '文件 / 审核意见',
      width: 220,
      render: (_, m) => (
        <Space direction="vertical" size={0}>
          {m.file_url
            ? <a href={m.file_url} target="_blank" rel="noreferrer">查看附件</a>
            : <span style={{ color: '#999' }}>未上传</span>}
          {m.review_remark && <span style={{ color: '#cf1322' }}>意见：{m.review_remark}</span>}
        </Space>
      ),
    },
    {
      title: '操作',
      width: 220,
      render: (_, m) => {
        const studentActions = canUpload && (m.status === 'pending' || m.status === 'uploaded' || m.status === 'rejected') && (
          <UploadButton onUploaded={(url) => handleUpload(m, url)}>
            <Button size="small" type={m.status === 'rejected' ? 'primary' : 'default'} icon={<UploadOutlined />}>
              {m.status === 'uploaded' ? '重新上传' : m.status === 'rejected' ? '按意见重传' : '上传'}
            </Button>
          </UploadButton>
        )
        const reviewActions = canReview && m.status === 'uploaded' && (
          <Space size={4}>
            <Popconfirm title="确认审核通过该材料？" onConfirm={() => approve(m)}>
              <Button size="small" type="primary">通过</Button>
            </Popconfirm>
            <Button size="small" danger onClick={() => { setRejecting(m); setRemark(m.review_remark || '') }}>
              退回
            </Button>
          </Space>
        )
        const locked = m.status === 'approved'
        return (
          <Space>
            {studentActions}
            {reviewActions}
            {locked && <span style={{ color: '#999' }}>已锁定</span>}
          </Space>
        )
      },
    },
  ]

  return (
    <div>
      <Progress percent={progress} format={(p) => `必交材料审核通过 ${p}%`} />
      <Table size="small" rowKey="id" dataSource={items} pagination={false} columns={columns} />
      <Modal
        title={`退回材料：${rejecting?.name || ''}`}
        open={!!rejecting}
        onOk={confirmReject}
        onCancel={() => setRejecting(null)}
        okText="确认退回"
        cancelText="取消"
      >
        <Input.TextArea
          rows={3}
          value={remark}
          onChange={(e) => setRemark(e.target.value)}
          placeholder="请填写退回原因，学生将据此重新上传"
        />
      </Modal>
    </div>
  )
}
