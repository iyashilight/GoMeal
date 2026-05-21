# GoMeal - 外卖订餐系统

后端学习项目，Go + Gin + GORM + MySQL + Redis 实现的外卖订餐系统。

## 目录结构

```
GoMeal/
├── backend/                 # Go 后端
│   ├── cmd/server/          # 入口
│   ├── internal/
│   │   ├── config/          # 配置 + DB/Redis 初始化
│   │   ├── cache/           # Redis 缓存封装
│   │   ├── model/           # 数据库模型
│   │   ├── repository/      # 数据访问层
│   │   ├── service/         # 业务逻辑层
│   │   ├── handler/         # HTTP 接口层
│   │   ├── middleware/      # 鉴权/CORS/限流
│   │   └── docs/            # Swagger 文档
│   ├── Dockerfile
│   ├── docker-entrypoint.sh
│   ├── docker-compose.yml   # Docker 编排（MySQL + Redis + App）
│   └── config.docker.yaml   # Docker 环境配置
└── frontend/                # Vue 前端
```

## 后端架构

- 三层架构：Handler → Service → Repository
- 框架：Gin
- ORM：GORM
- 数据库：MySQL
- 缓存：Redis
- 鉴权：JWT

## 启动方式

### Docker（推荐）

```bash
# 1. 配置环境变量
cp .env.example .env

# 2. 启动所有服务
docker compose up -d

# 3. 查看日志
docker compose logs -f app
```

### 本地开发

```bash
# 1. 确保 MySQL 和 Redis 已启动

# 2. 修改 backend/config.yaml 中的数据库连接信息

# 3. 启动
cd backend && go run ./cmd/server
```

## Go 版本

1.25
