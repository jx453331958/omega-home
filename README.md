# Omega Home

个人导航门户，参考 [gethomepage/homepage](https://github.com/gethomepage/homepage) 但**零 YAML 配置**，所有设置通过 Web 管理后台完成。

## 特性

- 🏠 **美观门户首页** — 暗色科技风主题，支持多种配色
- 📦 **服务分组管理** — 自定义分组和服务卡片
- 🔍 **状态检测** — 自动检测服务在线状态（HTTP/TCP）
- 🔖 **书签管理** — 快捷链接收藏
- 🎨 **主题切换** — 多种暗色主题可选
- 🔐 **管理后台** — JWT 认证，bcrypt 密码加密，支持修改密码
- 📦 **单二进制部署** — 模板嵌入编译产物，开箱即用

## 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Tailwind CSS + Alpine.js
- **认证**: JWT + bcrypt
- **部署**: 单二进制 / Docker Compose

## 快速开始

### Docker 部署（推荐）

预构建的多平台镜像（amd64/arm64）托管在 GitHub Container Registry，每次推送到 main 分支时通过 GitHub Actions 自动构建。

```bash
git clone https://github.com/jx453331958/omega-home.git
cd omega-home
./deploy.sh deploy
```

或者不克隆仓库，直接使用 Docker：

```bash
docker pull ghcr.io/jx453331958/omega-home:latest

docker run -d \
  --name omega-home \
  -p 3000:3000 \
  -v omega-home-data:/app/data \
  ghcr.io/jx453331958/omega-home:latest
```

可用的镜像标签：
- `latest` — main 分支最新构建
- `<commit-sha>` — 指定提交（如 `ghcr.io/jx453331958/omega-home:abc1234`）
- `<version>` — 语义化版本号（如 `ghcr.io/jx453331958/omega-home:1.0.0`）

首次启动会自动生成随机管理密码，通过日志查看：

```bash
docker compose logs | grep "Initial admin password"
```

输出示例：
```
========================================
Initial admin password: mOl3UyW0zoym
========================================
```

访问 `http://localhost:3000`，管理后台 `http://localhost:3000/admin`。

### 本地编译运行

```bash
cp .env.example .env
go build -o omega-home .
./omega-home
```

初始密码会直接打印在控制台。

## 管理密码

- **首次启动**：自动生成 12 位随机密码，打印到日志
- **指定初始密码**：设置环境变量 `ADMIN_PASSWORD=your-password`（仅首次生效）
- **修改密码**：登录管理后台 → 设置 → 修改密码
- **密码存储**：bcrypt 哈希存储在数据库中

## 配置

通过环境变量或 `.env` 文件配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `3000` | 监听端口 |
| `ADMIN_PASSWORD` | *随机生成* | 初始管理密码（仅首次启动生效） |
| `DATABASE_URL` | `sqlite:///data/omega.db` | 数据库连接 |
| `SECRET_KEY` | `change-me-to-random` | JWT 签名密钥 |
| `CHECK_INTERVAL` | `60` | 状态检测间隔（秒） |

## 部署脚本

```bash
./deploy.sh deploy   # 首次部署
./deploy.sh update   # 拉取最新镜像并重启
./deploy.sh start    # 启动服务
./deploy.sh stop     # 停止服务
./deploy.sh restart  # 重启服务
./deploy.sh status   # 查看状态
./deploy.sh logs     # 查看日志
./deploy.sh backup   # 备份数据库
./deploy.sh clean    # 删除容器和镜像
```

## 项目结构

```
omega-home/
├── main.go              # 入口，路由注册
├── config/              # 环境变量配置
├── models/              # 数据模型（Group/Service/Setting/Bookmark）
├── handlers/            # HTTP 处理器（含密码管理）
├── middleware/           # JWT 认证中间件
├── services/            # 状态检测服务
├── templates/           # HTML 模板（编译时嵌入）
├── static/              # 静态资源（Tailwind CSS）
├── Dockerfile           # Docker 多阶段构建（BuildKit 缓存优化）
├── docker-compose.yml
├── deploy.sh            # 部署脚本
└── update.sh            # 一键更新脚本
```

## License

MIT
