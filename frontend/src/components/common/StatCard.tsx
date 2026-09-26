import { Card, Statistic } from 'antd'

export default function StatCard({ title, value, suffix }: { title: string; value: number | string; suffix?: string }) {
  return (
    <Card>
      <Statistic title={title} value={value as number} suffix={suffix} />
    </Card>
  )
}
