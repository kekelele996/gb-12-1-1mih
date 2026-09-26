import { Avatar, Button, Dropdown, Layout, Menu } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useEffect } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { UserRoleMap } from '@/constants/user'

const { Header, Content } = Layout

const NAV = [
  { key: '/dashboard', label: '申请看板' },
  { key: '/universities', label: '院校库' },
  { key: '/applications', label: '申请项目' },
  { key: '/timeline', label: '时间线' },
  { key: '/messages', label: '站内消息' },
  { key: '/counselor', label: '顾问工作台' },
]

export default function App() {
  const navigate = useNavigate()
  const location = useLocation()
  const { token, user, logout, fetchProfile } = useAuth()

  useEffect(() => {
    if (token && !user) fetchProfile().catch(() => undefined)
  }, [token])

  const selected = NAV.find((n) => location.pathname.startsWith(n.key))?.key || ''

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', gap: 24, background: '#fff' }}>
        <div style={{ fontSize: 20, fontWeight: 700, color: '#2f54eb', cursor: 'pointer', whiteSpace: 'nowrap' }} onClick={() => navigate('/dashboard')}>
          🎓 留学申请管理
        </div>
        <Menu mode="horizontal" selectedKeys={[selected]} items={NAV} onClick={(e) => navigate(e.key)} style={{ flex: 1 }} />
        {token ? (
          <Dropdown
            menu={{
              items: [
                { key: 'profile', label: '个人中心' },
                { key: 'logout', label: '退出登录' },
              ],
              onClick: ({ key }) => {
                if (key === 'profile') navigate('/profile')
                if (key === 'logout') {
                  logout()
                  navigate('/login')
                }
              },
            }}
          >
            <span style={{ cursor: 'pointer' }}>
              <Avatar size="small" icon={<UserOutlined />} style={{ marginRight: 6 }} />
              {user?.real_name || user?.username}
              {user ? `（${UserRoleMap[user.role as keyof typeof UserRoleMap] || user.role}）` : ''}
            </span>
          </Dropdown>
        ) : (
          <Button type="primary" onClick={() => navigate('/login')}>
            登录/注册
          </Button>
        )}
      </Header>
      <Content style={{ padding: 24, maxWidth: 1200, width: '100%', margin: '0 auto' }}>
        <Outlet />
      </Content>
    </Layout>
  )
}
