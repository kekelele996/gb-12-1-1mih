# StudyPath（留学申请管理系统）

面向留学中介机构和自助留学学生：院校库维护、申请项目与状态跟踪、文书在线编辑与版本对比、顾问批注、材料清单进度、时间线与截止提醒、站内消息与选校推荐。

## Docker Compose 一键启动（推荐）

```bash
cp .env.example .env
docker compose up -d --build
```

访问地址：

- 前端：http://localhost:8012
- 后端 API：http://localhost:3012
- MinIO 控制台：http://localhost:9010（minioadmin / minioadmin）
- 健康检查：http://localhost:3012/healthz

默认账号：`admin / admin123`（管理员）、`counselor1 / counselor123`（顾问）、`student1 / student123`（学生）。

关闭并清理：

```bash
docker compose down -v --remove-orphans
```

## 本地开发（备选）

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
go build ./...
go test ./...
```

前端：

```bash
cd frontend
npm install
npm run dev
npm run build
```

## 技术栈

| 分层 | 技术 |
| --- | --- |
| 前端 | React 18 + TypeScript + Ant Design + ECharts + Vite + Zustand + React Router |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 16 + MinIO（文书/材料附件，minio-go/v7） |
| 认证 | JWT（github.com/golang-jwt/jwt/v5）+ RBAC（student/counselor/admin） |
| 其他 | validator/v10、log/slog、Nginx |

## 项目目录结构

```
gb-12-1/
├── docker-compose.yml
├── .env.example
├── README.md
├── database/init.sql
├── backend/
│   ├── cmd/server/            # main.go + migrate/seed
│   └── internal/
│       ├── config/            # DB/MinIO/JWT/限流配置
│       ├── model/             # 10 个实体
│       ├── repository/        # 按实体分文件
│       ├── service/           # 按实体分文件 + notification/stats
│       ├── handler/           # 按实体分文件 + upload/dashboard
│       ├── router/            # router.go + 按实体分文件
│       ├── middleware/        # auth/rbac/rate_limiter/error_handler/logger/cors
│       ├── dto/
│       ├── constants/         # application/document/material/user/error_codes/log_templates/messages
│       └── util/              # jwt/logger/formatters/minio_client
└── frontend/
    ├── nginx.conf
    └── src/
        ├── api/               # user/university/application/document/recommendation/message/dashboard
        ├── stores/            # authStore/userStore/applicationStore/documentStore/universityStore
        ├── components/common/ # ApplicationStatusTag/Timeline/UniversityCard/AnnotationList/...
        ├── hooks/             # useAuth/useApplicationStats/useDocumentVersions/useDeadlineReminders
        ├── pages/             # Dashboard/UniversityList/.../Login
        ├── router/            # index.tsx
        ├── utils/             # request/dateFormat/diff
        └── constants/         # application/document/material/user/errorCodes
```

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | gbstudyapply | Compose 项目名/容器前缀 |
| DB_NAME / DB_USER / DB_PASSWORD | gbstudyapply_db / gbstudyapply_user / gbstudyapply_pwd | PostgreSQL |
| MINIO_ROOT_USER / MINIO_ROOT_PASSWORD | minioadmin / minioadmin | MinIO 账号 |
| MINIO_PORT / MINIO_CONSOLE_PORT | 9009 / 9010 | MinIO API/控制台 |
| JWT_SECRET | change_me_to_a_long_random_string | JWT 密钥（生产必改） |
| FRONTEND_PORT / BACKEND_PORT / DB_PORT | 8012 / 3012 / 5505 | 端口映射 |

## Docker 部署说明

- 端口：前端 `8012:80`、后端 `${BACKEND_PORT:-3012}:8080`、DB `${DB_PORT:-5505}:5432`、MinIO `${MINIO_PORT:-9009}:9000` + `${MINIO_CONSOLE_PORT:-9010}:9001`
- 数据卷：`db_data`、`minio_data`（命名卷持久化）
- 依赖顺序：db/minio healthcheck → backend `depends_on: service_healthy` → frontend
- 常见问题：端口冲突改 `.env`；重置数据 `docker compose down -v`；文书附件通过 `/api/v1/files/:key` 代理 MinIO 读取

## API 接口清单

> 后端统一前缀 `/api/v1`，响应统一为 `{ "code": 0, "message": "ok", "data": ... }`。标注「登录」的接口需携带 `Authorization: Bearer <JWT>`。

| 方法 | 路径 | 权限 | 说明 |
| --- | --- | --- | --- |
| GET | /healthz | 公开 | 健康检查 |
| GET | /api/v1/files/:key | 公开 | 读取 MinIO 上传文件 |
| POST | /api/v1/uploads | 登录（限流） | 上传文件到 MinIO |
| GET | /api/v1/dashboard/stats | 登录 | 工作台统计 |
| POST | /api/v1/users/register | 公开（限流） | 注册并返回 JWT |
| POST | /api/v1/users/login | 公开（限流） | 登录并返回 JWT |
| GET | /api/v1/users/me | 登录 | 获取当前用户 |
| PUT | /api/v1/users/me | 登录 | 更新当前用户资料 |
| GET | /api/v1/users/students | counselor/admin | 学生列表 |
| GET | /api/v1/universities | 公开 | 院校分页列表 |
| GET | /api/v1/universities/:id | 公开 | 院校详情 |
| POST | /api/v1/universities | admin（限流） | 新增院校 |
| PUT | /api/v1/universities/:id | admin | 更新院校 |
| DELETE | /api/v1/universities/:id | admin | 删除院校 |
| GET | /api/v1/applications | 登录 | 申请项目列表 |
| GET | /api/v1/applications/:id | 登录 | 申请项目详情 |
| POST | /api/v1/applications | student（限流） | 创建申请项目 |
| PUT | /api/v1/applications/:id/status | 登录 | 申请状态流转 |
| GET | /api/v1/applications/:id/documents | 登录 | 项目文档列表 |
| POST | /api/v1/applications/:id/documents | 登录（限流） | 创建文档（事务：文档+初始版本） |
| GET | /api/v1/applications/:id/materials | 登录 | 材料清单 |
| POST | /api/v1/applications/:id/materials | 登录 | 添加材料 |
| PUT | /api/v1/materials/:id/status | 登录 | 材料状态流转 |
| GET | /api/v1/applications/:id/timeline | 登录 | 项目时间线 |
| POST | /api/v1/applications/:id/timeline | 登录（限流） | 添加时间线节点 |
| GET | /api/v1/documents/:id | 登录 | 文档详情 |
| PUT | /api/v1/documents/:id | 登录（限流） | 保存文档（事务：更新文档+写入新版本） |
| GET | /api/v1/documents/:id/versions | 登录 | 版本历史 |
| POST | /api/v1/documents/:id/rollback | 登录 | 回滚到历史版本 |
| GET | /api/v1/documents/:id/annotations | 登录 | 文档批注列表 |
| POST | /api/v1/documents/:id/annotations | counselor/admin（限流） | 添加批注 |
| PUT | /api/v1/timeline/:id/done | 登录 | 标记时间线节点完成 |
| GET | /api/v1/recommendations/student/:studentId | 登录 | 学生选校方案 |
| POST | /api/v1/recommendations | counselor（限流） | 生成选校推荐 |
| GET | /api/v1/messages | 登录 | 消息列表 |
| POST | /api/v1/messages | 登录（限流） | 发送站内消息 |
| PUT | /api/v1/messages/:id/read | 登录 | 标记消息已读 |

## 枚举出现位置清单

### ApplicationStatus（planning/preparing/submitted/waiting/admitted/rejected/waitlisted）

- 后端：`internal/constants/application.go`（定义+状态机）、`internal/model/application_project.go`（模型）、`internal/service/application_service.go`（流转校验）、`internal/util/formatters.go`（AppStatusText）、`internal/constants/log_templates.go`、`internal/service/stats_service.go`（看板统计）、`database/init.sql`
- 前端：`src/constants/application.ts`（定义）、`src/components/common/ApplicationStatusTag.tsx`、`src/pages/ApplicationList.tsx`（筛选）、`src/pages/Dashboard.tsx`（看板饼图）、`src/pages/CounselorWorkbench.tsx`

### DocumentType（ps/rl/cv/essay）

- 后端：`internal/constants/document.go`（定义）、`internal/model/document.go`、`internal/service/document_service.go`（校验）、`internal/util/formatters.go`（DocTypeText）、`internal/constants/log_templates.go`、`database/init.sql`
- 前端：`src/constants/document.ts`（定义）、`src/pages/ApplicationDetail.tsx`（文书入口）、`src/pages/DocumentEditor.tsx`

### MaterialStatus（pending/uploaded/approved）

- 后端：`internal/constants/material.go`（定义）、`internal/model/material_item.go`、`internal/service/material_service.go`（状态流转+进度）、`internal/util/formatters.go`、`database/init.sql`
- 前端：`src/constants/material.ts`（定义）、`src/components/common/MaterialProgress.tsx`、`src/pages/ApplicationDetail.tsx`

### UserRole（student/counselor/admin）

- 后端：`internal/constants/user.go`（定义）、`internal/model/user.go`、`internal/middleware/rbac.go`、`internal/router/*.go`（角色路由）、`internal/util/formatters.go`（RoleText）、`database/init.sql`
- 前端：`src/constants/user.ts`（定义）、`src/router/index.tsx`（守卫）、`src/components/common/RoleGuard.tsx`、`src/App.tsx`（角色显示）、`src/pages/CounselorWorkbench.tsx`

## 横切关注点

- 认证授权（JWT + RBAC）：`middleware/auth.go`、`rbac.go`、`util/jwt.go`、学生/顾问/管理员路由、前端 `RoleGuard.tsx`
- 接口限流：`middleware/rate_limiter.go`（登录/创建/保存/上传启用）
- 全局错误处理：`middleware/error_handler.go`、`util/app_error.go`、`constants/error_codes.go`、前端 `utils/request.ts`
- 文件上传（MinIO）：`handler/upload_handler.go`、`util/minio_client.go`、`components/common/UploadButton.tsx`
- 截止日期提醒：`service/notification_service.go`、`hooks/useDeadlineReminders.ts`、站内消息

## License

MIT
