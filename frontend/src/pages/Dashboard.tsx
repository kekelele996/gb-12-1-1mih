import { Card, Col, Row, Table } from 'antd'
import ReactECharts from 'echarts-for-react'
import { useApplicationStats } from '@/hooks/useApplicationStats'
import ApplicationStatusTag from '@/components/common/ApplicationStatusTag'
import StatCard from '@/components/common/StatCard'
import { ApplicationStatusMap, ApplicationStatusList } from '@/constants/application'

export default function Dashboard() {
  const { data, loading } = useApplicationStats()

  const stats = data?.stats
  const byStatus = stats?.by_status || {}
  const pieOption = {
    title: { text: '申请状态分布', left: 'center' },
    tooltip: { trigger: 'item' },
    series: [
      {
        type: 'pie',
        radius: ['40%', '65%'],
        data: ApplicationStatusList.filter((s) => byStatus[s]).map((s) => ({
          name: ApplicationStatusMap[s].text,
          value: byStatus[s],
        })),
      },
    ],
  }

  return (
    <div>
      <h1>申请数据看板</h1>
      <Row gutter={[16, 16]}>
        <Col xs={12} md={6}><StatCard title="申请项目总数" value={stats?.total || 0} /></Col>
        <Col xs={12} md={6}><StatCard title="已提交申请" value={stats?.applied || 0} /></Col>
        <Col xs={12} md={6}><StatCard title="已录取" value={stats?.admitted || 0} /></Col>
        <Col xs={12} md={6}><StatCard title="材料平均进度" value={stats?.material_avg || 0} suffix="%" /></Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} md={10}>
          <Card loading={loading}>
            <ReactECharts option={pieOption} style={{ height: 320 }} />
          </Card>
        </Col>
        <Col xs={24} md={14}>
          <Card title="申请项目列表" loading={loading}>
            <Table
              rowKey="id"
              size="small"
              dataSource={data?.projects || []}
              pagination={false}
              columns={[
                { title: 'ID', dataIndex: 'id', width: 60 },
                { title: '院校 ID', dataIndex: 'university_id', width: 90 },
                { title: '专业', dataIndex: 'major' },
                { title: '轮次', dataIndex: 'round' },
                { title: '状态', dataIndex: 'status', render: (s: string) => <ApplicationStatusTag status={s} /> },
              ]}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
