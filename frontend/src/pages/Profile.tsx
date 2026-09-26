import { Button, Card, Col, Descriptions, Form, Input, Row, message } from 'antd'
import { useEffect } from 'react'
import { useUserStore } from '@/stores/userStore'

export default function Profile() {
  const { profile, load, save } = useUserStore()
  const [form] = Form.useForm()

  useEffect(() => {
    load().then(() => form.setFieldsValue(profile || {}))
  }, [])

  async function onFinish(values: { real_name: string; phone: string }) {
    await save(values)
    message.success('资料已更新')
  }

  return (
    <Row gutter={16}>
      <Col xs={24} md={10}>
        <Card title="个人资料">
          <Descriptions column={1} size="small">
            <Descriptions.Item label="用户名">{profile?.username}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{profile?.email}</Descriptions.Item>
            <Descriptions.Item label="角色">
              {profile?.role === 'counselor' ? '顾问' : profile?.role === 'admin' ? '管理员' : '学生'}
            </Descriptions.Item>
          </Descriptions>
          <Form form={form} layout="vertical" onFinish={onFinish} style={{ marginTop: 16 }}>
            <Form.Item name="real_name" label="姓名">
              <Input />
            </Form.Item>
            <Form.Item name="phone" label="电话">
              <Input />
            </Form.Item>
            <Button type="primary" htmlType="submit">
              保存修改
            </Button>
          </Form>
        </Card>
      </Col>
    </Row>
  )
}
