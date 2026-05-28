# TikTok 后端 API

Go + chi + GORM + MySQL 短视频 Feed API。

前端仓库：[tiktok-frontend](https://github.com/LoveMuZiLi/tiktok-frontend)

## 本地启动

```bash
docker compose up -d mysql
cp .env.example .env
go mod tidy
make run
```

默认 http://localhost:8080

## API

| 方法 | 路径 |
|------|------|
| GET | `/api/v1/health` |
| GET | `/api/v1/videos` |
| GET | `/api/v1/videos/{id}` |
| POST | `/api/v1/videos/{id}/like` |

## CI

推送至 `main` 时 GitHub Actions 启动 MySQL 服务容器，执行 `go vet`、`go test`、`go build`。
