import { Col, Row, Select } from 'antd'
import { useEffect, useState } from 'react'
import UniversityCard from '@/components/common/UniversityCard'
import SearchFilter from '@/components/common/SearchFilter'
import { useUniversityStore } from '@/stores/universityStore'

export default function UniversityList() {
  const { universities, total, load } = useUniversityStore()
  const [country, setCountry] = useState('')

  useEffect(() => {
    load({ page: 1, page_size: 20, country })
  }, [country])

  function onSearch(kw: string) {
    load({ page: 1, page_size: 20, country, keyword: kw })
  }

  return (
    <div>
      <h1>院校库</h1>
      <SearchFilter onSearch={onSearch}>
        <Select
          placeholder="国家"
          allowClear
          style={{ width: 140 }}
          value={country || undefined}
          onChange={(v) => setCountry(v || '')}
          options={['美国', '英国', '新加坡', '澳大利亚', '加拿大'].map((c) => ({ value: c, label: c }))}
        />
      </SearchFilter>
      <Row gutter={[16, 16]}>
        {universities.map((u) => (
          <Col xs={24} sm={12} md={8} key={u.id}>
            <UniversityCard university={u} />
          </Col>
        ))}
      </Row>
      {!universities.length && <p style={{ color: '#999' }}>暂无院校（共 {total} 条）</p>}
    </div>
  )
}
