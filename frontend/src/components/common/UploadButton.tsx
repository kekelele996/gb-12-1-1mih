import { Button, Upload, message } from 'antd'
import { UploadOutlined } from '@ant-design/icons'
import axios from 'axios'

export default function UploadButton({ onUploaded }: { onUploaded: (url: string) => void }) {
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
          message.success('上传成功')
          onUploaded(res.data.data.url)
        } catch {
          message.error('上传失败')
        }
      }}
    >
      <Button icon={<UploadOutlined />}>上传附件</Button>
    </Upload>
  )
}
