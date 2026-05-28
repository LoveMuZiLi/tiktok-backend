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
