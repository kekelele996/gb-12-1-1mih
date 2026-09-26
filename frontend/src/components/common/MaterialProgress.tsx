import { Progress, Table, Tag } from 'antd'
import type { MaterialItem } from '@/types/api'
import { MaterialStatusMap } from '@/constants/material'

export default function MaterialProgress({ items, progress }: { items: MaterialItem[]; progress: number }) {
  return (
    <div>
      <Progress percent={progress} />
      <Table
        size="small"
        rowKey="id"
        dataSource={items}
        pagination={false}
        columns={[
          { title: '材料', dataIndex: 'name' },
          { title: '分类', dataIndex: 'category' },
          {
            title: '必交',
            dataIndex: 'is_required',
            render: (v: boolean) => (v ? '是' : '否'),
          },
          {
            title: '状态',
            dataIndex: 'status',
            render: (s: string) => {
              const m = MaterialStatusMap[s as keyof typeof MaterialStatusMap]
              return m ? <Tag color={m.color}>{m.text}</Tag> : <Tag>{s}</Tag>
            },
          },
        ]}
      />
    </div>
  )
}
