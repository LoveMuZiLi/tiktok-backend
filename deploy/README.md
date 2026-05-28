# TikTok 线上部署（前端 + 后端）

**勿将服务器 IP、root、密码写入仓库。** 凭据仅保存在本机 `deploy/deploy.local.env`（已 gitignore）。

## 架构（均在同一台服务器）

```
浏览器
  │
  ▼ :8088  Nginx
  ├─ /          → /var/www/tiktok/frontend   （React 静态）
  └─ /api/*     → 127.0.0.1:8080            （Go 后端 systemd: tiktok-api）
                      │
                      ▼
                   MySQL :3306
```

前端通过相对路径 `/api/v1/...` 调用后端，**无需**在构建时写死后端地址。

## 本机配置

```bash
cp deploy/env/deploy.local.env.example deploy/deploy.local.env
# DEPLOY_HOST=your-user@your-server-ip
# PUBLIC_PORT=8088
```

## 1. 服务器初始化（一次）

```bash
scp deploy/scripts/server-init.sh "${DEPLOY_HOST}:/tmp/"
ssh "${DEPLOY_HOST}" 'WEB_PORT=8088 MYSQL_APP_PASS=你的数据库密码 bash /tmp/server-init.sh'
```

## 2. 配置后端环境（一次）

在服务器编辑 `/opt/tiktok/backend/.env`（参考 `deploy/env/production.env.example`）：

```env
PORT=8080
DB_HOST=127.0.0.1
DB_USER=tiktok
DB_PASSWORD=你的数据库密码
DB_NAME=tiktok
ALLOWED_ORIGINS=http://你的公网IP:8088
```

## 3. 一键发布（前端 + 后端）

```bash
./deploy/scripts/deploy.sh
```

脚本会：

1. 构建前端 `dist/`
2. 交叉编译 Linux 后端 `tiktok-api`
3. 上传并安装 systemd 服务 + Nginx
4. 自检 `http://127.0.0.1:8080/api/v1/health` 与 Nginx `:8088/api`

## 验证

| 检查项 | 地址 |
|--------|------|
| 前端页面 | `http://<IP>:8088/` |
| API 健康 | `http://<IP>:8088/api/v1/health` |
| 视频列表 | `http://<IP>:8088/api/v1/videos` |

安全组放行：**8088**（对外）、**22**（SSH）；**不要**对公网开放 8080、3306。

## 运维命令（服务器上）

```bash
systemctl status tiktok-api
journalctl -u tiktok-api -f
systemctl restart tiktok-api
```
