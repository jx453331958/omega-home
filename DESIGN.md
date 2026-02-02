# Omega Home - 设计文档（Go 版）

## 定位
简洁美观的个人导航门户，参考 gethomepage/homepage 但**零 YAML 配置**，所有设置通过 Web 管理后台完成。

## 技术栈
- **后端**: Go + Gin + GORM + SQLite（支持外部 PostgreSQL）
- **前端**: 原生 HTML + Tailwind CSS CDN + Alpine.js CDN（embed 进二进制）
- **部署**: 单二进制 / Docker Compose

## 文件结构

```
omega-home/
├── main.go                  # 入口
├── go.mod / go.sum
├── config/
│   └── config.go            # 环境变量配置
├── models/
│   ├── database.go          # GORM 初始化 + 迁移
│   ├── group.go             # 分组模型
│   ├── service.go           # 服务模型
│   ├── setting.go           # 设置模型
│   └── bookmark.go          # 书签模型
├── handlers/
│   ├── portal.go            # 门户页面 + API
│   ├── admin.go             # 管理后台页面 + CRUD API
│   ├── auth.go              # 管理员认证（JWT）
│   ├── status.go            # 服务状态 API
│   └── upload.go            # 图片上传
├── services/
│   └── checker.go           # 后台状态检测 goroutine
├── middleware/
│   └── auth.go              # JWT 认证中间件
├── templates/
│   ├── index.html           # 门户首页
│   └── admin.html           # 管理后台（SPA 风格）
├── static/
│   └── uploads/             # 用户上传的图片（运行时）
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── .dockerignore
```

## 数据库模型

```go
type Group struct {
    ID        uint   `gorm:"primarykey" json:"id"`
    Name      string `gorm:"not null" json:"name"`
    Icon      string `json:"icon"`       // emoji or url
    SortOrder int    `json:"sort_order"`
    Columns   int    `gorm:"default:3" json:"columns"` // 1-4
    Services  []Service `json:"services,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Service struct {
    ID          uint   `gorm:"primarykey" json:"id"`
    Name        string `gorm:"not null" json:"name"`
    URL         string `gorm:"not null" json:"url"`
    Icon        string `json:"icon"`        // emoji, url, or icon name
    Description string `json:"description"`
    GroupID     uint   `json:"group_id"`
    SortOrder   int    `json:"sort_order"`
    Target      string `gorm:"default:_blank" json:"target"`
    StatusCheck bool   `gorm:"default:true" json:"status_check"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Setting struct {
    Key   string `gorm:"primarykey" json:"key"`
    Value string `json:"value"`
}

type Bookmark struct {
    ID        uint   `gorm:"primarykey" json:"id"`
    Name      string `gorm:"not null" json:"name"`
    URL       string `gorm:"not null" json:"url"`
    Icon      string `json:"icon"`
    SortOrder int    `json:"sort_order"`
}
```

## API 设计

```
# 公开 API
GET  /                        # 门户首页（HTML）
GET  /api/config              # 门户配置 JSON（分组+服务+设置+书签）
GET  /api/status              # 所有服务状态

# 管理后台
GET  /admin                   # 管理后台页面（HTML）
POST /api/admin/login         # 登录，返回 JWT

# 需要 JWT 认证
GET    /api/admin/services       # 服务列表
POST   /api/admin/services       # 添加
PUT    /api/admin/services/:id   # 编辑
DELETE /api/admin/services/:id   # 删除
PUT    /api/admin/services/reorder  # 批量排序

GET    /api/admin/groups
POST   /api/admin/groups
PUT    /api/admin/groups/:id
DELETE /api/admin/groups/:id
PUT    /api/admin/groups/reorder

GET    /api/admin/bookmarks
POST   /api/admin/bookmarks
PUT    /api/admin/bookmarks/:id
DELETE /api/admin/bookmarks/:id

GET    /api/admin/settings
PUT    /api/admin/settings       # 批量更新

POST   /api/admin/upload         # 上传图片
```

## 门户首页设计

### 视觉要求
- **毛玻璃卡片**: `backdrop-blur-md bg-white/10 border border-white/20`
- **渐变背景**: 5 种预设主题
  - 靛蓝紫: `from-indigo-600 via-purple-600 to-pink-500`
  - 深空黑: `from-gray-900 via-slate-900 to-zinc-900`
  - 翠绿: `from-emerald-600 via-teal-600 to-cyan-500`
  - 暖橙: `from-orange-500 via-red-500 to-pink-500`
  - 极光: `from-green-400 via-cyan-500 to-blue-600`
- **卡片 hover**: `hover:-translate-y-1 hover:shadow-xl transition-all duration-300`
- **状态指示灯**: 在线绿色呼吸灯动画，离线红色静态
- **图标**: 支持 emoji、外部图片 URL
- **响应式**: mobile 1列, sm 2列, md 3列, lg 4列
- **暗色模式**: 跟随系统 + 手动切换

### 页面结构
```
Header: Logo + 站点标题 + 主题切换 + 时钟
搜索栏: 实时过滤服务
书签栏: 快捷链接（可选）
分组1:
  [卡片] [卡片] [卡片]
分组2:
  [卡片] [卡片]
Footer: 自定义文字 + ⚙️ 管理入口
```

## 管理后台设计

单页面 + Tab 切换:
- **服务管理 Tab**: 表格列表 + 添加/编辑模态框
- **分组管理 Tab**: 同上
- **书签管理 Tab**: 同上
- **外观设置 Tab**: 表单（标题/副标题/主题/背景/布局/问候语/页脚）
- **组件设置 Tab**: 搜索栏开关/天气开关/时钟开关等

## 状态检测

- 后台 goroutine，每 60 秒（可配置）检测所有开启了 status_check 的服务
- HTTP GET，超时 5 秒
- 结果缓存在内存 map 中（sync.Map）
- 首页通过 /api/status 获取，轮询间隔 30 秒

## Docker 部署

```dockerfile
# 多阶段构建
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o omega-home .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=build /app/omega-home /usr/local/bin/
COPY --from=build /app/templates /app/templates
COPY --from=build /app/static /app/static
WORKDIR /app
EXPOSE 3000
CMD ["omega-home"]
```

```yaml
# docker-compose.yml
services:
  omega-home:
    build: .
    ports:
      - "${PORT:-3000}:3000"
    environment:
      - ADMIN_PASSWORD=${ADMIN_PASSWORD:-admin}
      - DATABASE_URL=${DATABASE_URL:-sqlite:///data/omega.db}
      - SECRET_KEY=${SECRET_KEY:-change-me-to-random}
      - CHECK_INTERVAL=${CHECK_INTERVAL:-60}
    volumes:
      - omega-data:/app/data
    restart: unless-stopped

volumes:
  omega-data:
```

## 配置（环境变量）

```
PORT=3000
ADMIN_PASSWORD=admin          # 管理后台密码
DATABASE_URL=sqlite:///data/omega.db  # 或 postgres://...
SECRET_KEY=random-string      # JWT 签名
CHECK_INTERVAL=60             # 状态检测间隔（秒）
```

## 初始数据

首次启动自动创建：
- 默认设置（标题 "Omega Home"、靛蓝紫主题、搜索栏开启）
- 示例分组 "常用工具"
- 示例服务（Google、GitHub、Wikipedia）

## 实现顺序

1. go mod init + 基础骨架（main.go + config + models + 数据库迁移 + 初始数据）
2. 管理员认证（auth handler + middleware）
3. CRUD API（services + groups + bookmarks + settings + upload）
4. 状态检测 goroutine
5. 门户首页 HTML（Tailwind + Alpine.js，要非常好看！）
6. 管理后台 HTML（Tailwind + Alpine.js，完整 CRUD 界面）
7. Docker 部署文件
8. git push

每完成一个核心模块就 git commit。
