# TikTok 后端 API

Go + chi + GORM + MySQL。

## 工程结构

见 [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)。

## 环境变量

复制 `.env.example` 为 `.env` 后修改。

## 本地启动

```bash
docker compose up -d mysql
go mod tidy
make run
```

## 构建

```bash
make build
```

前端仓库：[tiktok-frontend](https://github.com/LoveMuZiLi/tiktok-frontend)
