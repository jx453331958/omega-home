# Omega Home

个人导航门户，参考 [gethomepage/homepage](https://github.com/gethomepage/homepage) 但**零 YAML 配置**，所有设置通过 Web 管理后台完成。

## 特性

- 🏠 **美观门户首页** — 暗色科技风主题，支持多种配色
- 📦 **服务分组管理** — 自定义分组和服务卡片
- 🔍 **状态检测** — 自动检测服务在线状态（HTTP/TCP）
- 🔖 **书签管理** — 快捷链接收藏
- 🎨 **主题切换** — 多种暗色主题可选
- 🔐 **管理后台** — JWT 认证保护的 Web 管理界面
- 📦 **单二进制部署** — 模板嵌入编译产物，开箱即用

## 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Tailwind CSS CDN + Alpine.js CDN
- **认证**: JWT
- **部署**: 单二进制 / Docker Compose

## 快速开始

### 本地运行

```bash
# 克隆项目
git clone https://github.com/jx453331958/omega-home.git
cd omega-home

# 复制配置
cp .env.example .env

# 编译运行
go build -o omega-home .
./omega-home
```

访问 `http://localhost:3000`，管理后台 `http://localhost:3000/admin`（默认密码 `admin`）。

### Docker

```bash
docker compose up -d
```

## 配置

通过环境变量或 `.env` 文件配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `3000` | 监听端口 |
| `ADMIN_PASSWORD` | `admin` | 管理后台密码 |
| `DATABASE_URL` | `sqlite:///data/omega.db` | 数据库连接 |
| `SECRET_KEY` | `change-me-to-random` | JWT 签名密钥 |
| `CHECK_INTERVAL` | `60` | 状态检测间隔（秒） |

## 项目结构

```
omega-home/
├── main.go              # 入口，路由注册
├── config/              # 环境变量配置
├── models/              # 数据模型（Group/Service/Setting/Bookmark）
├── handlers/            # HTTP 处理器
├── middleware/           # JWT 认证中间件
├── services/            # 状态检测服务
├── templates/           # HTML 模板（编译时嵌入）
├── static/              # 静态资源
├── Dockerfile           # 多阶段构建
└── docker-compose.yml
```

## License

MIT
