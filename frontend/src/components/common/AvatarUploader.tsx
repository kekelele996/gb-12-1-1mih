import { Avatar, Upload, message } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import axios from 'axios'

export default function AvatarUploader({ value, onChange }: { value: string; onChange: (url: string) => void }) {
  return (
    <Upload
      showUploadList={false}
      customRequest={async (opt) => {
        const token = localStorage.getItem('gbstudyapply_token')
        const form = new FormData()
        form.append('file', opt.file as File)
        try {
          const res = await axios.post('/api/v1/uploads', form, {
            headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'multipart/form-data' },
          })
          message.success('头像已上传')
          onChange(res.data.data.url)
        } catch {
          message.error('上传失败')
        }
      }}
    >
      <Avatar size={64} src={value} icon={<UserOutlined />} style={{ cursor: 'pointer' }} />
    </Upload>
  )
}
