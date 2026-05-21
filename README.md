# GoMeal

外卖订餐系统后端 — Go 学习项目。

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **鉴权**: JWT
- **部署**: Docker Compose

## 快速启动

```bash
# 配置密码
cp .env.example .env

# 启动（MySQL + Redis + App）
docker compose up -d

# API 在 http://localhost:8080
```

## 项目结构

```
backend/             # Go 后端
├── cmd/server/      # 入口
├── internal/        # 核心代码
│   ├── handler/     # HTTP 接口
│   ├── service/     # 业务逻辑
│   ├── repository/  # 数据访问
│   ├── model/       # 数据模型
│   ├── middleware/   # 鉴权/限流
│   ├── cache/       # Redis 缓存
│   └── config/      # 配置管理
├── docs/            # Swagger 文档
├── Dockerfile
└── docker-compose.yml

frontend/            # Vue 前端
```

## API 文档

详见 [API.md](backend/API.md) 或启动后访问 Swagger UI。

## License

MIT
