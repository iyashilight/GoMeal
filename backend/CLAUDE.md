# Ordering System

后端学习项目，从零重新实现外卖订餐系统。

## 目录结构

```
OrderingSystem/
├── backend/                # 旧后端（已废弃）
├── frontend/               # 前端
└── study-backend/          # 学习相关
    ├── CLAUDE.md           # 项目说明
    ├── backend_new/        # 新后端（当前工作目录）
    │   ├── config.yaml
    │   ├── go.mod / go.sum
    │   ├── API.md          # API 文档
    │   ├── cmd/server/     # 入口
    │   └── internal/       # 核心代码
    │       ├── config/     # 配置 + DB/Redis 初始化
    │       ├── cache/      # Redis 缓存封装
    │       ├── model/      # 数据库模型
    │       ├── repository/ # 数据访问层
    │       ├── service/    # 业务逻辑层
    │       ├── handler/    # HTTP 接口层
    │       ├── middleware/ # 鉴权/CORS
    │       └── utils/      # 工具函数
    └── docs/
        ├── learning/       # 学习记录（lesson-*.md）
        └── conclusion/     # 架构总结文档
```

## 后端架构

- 三层架构：Handler → Service → Repository
- 框架：Gin
- 新后端代码在 `study-backend/backend_new/`

## 启动方式

（待补充）

## 学习记录

见 `docs/learning/`

## Go 版本

1.25

## 自动记录

每个阶段/功能完成后自动记录：
1. 更新记忆中的 `learning_progress.md`
2. 更新 `.claude/UPGRADE_PLAN.md`
3. 复制 `UPGRADE_PLAN.md` 到 `docs/`
4. 总结这次学到的内容
