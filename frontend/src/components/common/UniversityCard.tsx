import { Card, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import type { University } from '@/types/api'

export default function UniversityCard({ university }: { university: University }) {
  const navigate = useNavigate()
  return (
    <Card hoverable onClick={() => navigate(`/universities/${university.id}`)} style={{ marginBottom: 12 }}>
      <Card.Meta
        title={
          <span>
            {university.name} <Tag color="blue">#{university.ranking}</Tag>
          </span>
        }
        description={`${university.country} · ${university.city} · ${university.tuition_range || '-'}`}
      />
    </Card>
  )
}
