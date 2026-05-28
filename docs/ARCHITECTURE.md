# 后端工程结构

```
tiktok-backend/
├── cmd/
│   └── api/                    # 程序入口 main
├── internal/
│   ├── bootstrap/              # 应用启动与依赖组装
│   ├── config/                 # 配置加载
│   ├── entity/                 # 领域实体 / 数据模型
│   ├── repository/             # 数据访问层
│   ├── service/                # 业务逻辑层
│   ├── transport/
│   │   └── http/               # HTTP 传输层
│   │       ├── handler/        # 控制器
│   │       └── router.go       # 路由注册
│   └── infra/
│       └── persistence/        # MySQL / GORM
├── migrations/                 # SQL 迁移脚本
├── deploy/                     # 部署配置（nginx/systemd/脚本）
└── docs/
```

分层约定：`handler` → `service` → `repository` → `entity`，禁止跨层跳跃访问数据库。

## 数据表

| 表 | 说明 |
|----|------|
| `users` | 用户资料 |
| `videos` | 短视频（关联 `user_id`） |
| `follows` | 关注关系 |
| `conversations` / `conversation_members` | 私信会话 |
| `messages` | 聊天消息 |
| `notifications` | 收件箱通知 |

## API 概览（`/api/v1`）

- `users` — CRUD + `GET /{id}/profile`
- `videos` — CRUD + `POST /{id}/like`；列表支持 `?feed=following|friends|user&target_id=`
- `follows` — 关注/取关/状态
- `conversations` — 会话列表、消息列表、发消息
- `notifications` — 通知 CRUD

无登录阶段默认 `?user_id=1` 表示当前用户。
