import { Card, Descriptions, Tag } from 'antd'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getUniversity } from '@/api/university'
import type { University } from '@/types/api'

export default function UniversityDetail() {
  const { id } = useParams()
  const [uni, setUni] = useState<University | null>(null)

  useEffect(() => {
    getUniversity(id!).then(setUni)
  }, [id])

  if (!uni) return <p>加载中…</p>

  let majors: string[] = []
  try {
    majors = JSON.parse(uni.top_majors || '[]')
  } catch {
    majors = []
  }
  let req: Record<string, string> = {}
  try {
    req = JSON.parse(uni.requirements || '{}')
  } catch {
    req = {}
  }

  return (
    <div style={{ maxWidth: 800, margin: '0 auto' }}>
      <h1>
        {uni.name} <Tag color="blue">QS #{uni.ranking}</Tag>
      </h1>
      <Card title="基本信息">
        <Descriptions column={2} bordered>
          <Descriptions.Item label="国家">{uni.country}</Descriptions.Item>
          <Descriptions.Item label="城市">{uni.city}</Descriptions.Item>
          <Descriptions.Item label="申请截止">{uni.application_deadline || '-'}</Descriptions.Item>
          <Descriptions.Item label="学费范围">{uni.tuition_range || '-'}</Descriptions.Item>
          <Descriptions.Item label="优势专业" span={2}>{majors.join('、')}</Descriptions.Item>
          <Descriptions.Item label="GPA 要求">{req.gpa || '-'}</Descriptions.Item>
          <Descriptions.Item label="语言要求">{req.language || '-'}</Descriptions.Item>
          <Descriptions.Item label="标化考试">{req.test || '-'}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}
