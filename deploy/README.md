# TikTok 线上部署

**勿将服务器 IP、root 密码、SSH 账号写入仓库。** 凭据仅保存在本机 `deploy/deploy.local.env`（已 gitignore）。

## 架构

```
浏览器 :8088 → Nginx → 静态前端
                  └→ /api/* → Go :8080 → MySQL
```

## 本机配置（一次性）

```bash
cp deploy/env/deploy.local.env.example deploy/deploy.local.env
# 编辑 deploy.local.env，填写：
# DEPLOY_HOST=your-user@your-server-ip
# PUBLIC_PORT=8088
```

## 服务器初始化（仅需一次）

```bash
scp deploy/scripts/server-init.sh "${DEPLOY_HOST}:/tmp/"
ssh "${DEPLOY_HOST}" 'WEB_PORT=8088 MYSQL_APP_PASS=你的数据库密码 bash /tmp/server-init.sh'
```

在服务器创建 `/opt/tiktok/backend/.env`（参考 `deploy/env/production.env.example`），设置 `DB_PASSWORD` 与 `ALLOWED_ORIGINS`（如 `http://你的IP:8088`）。

## 发布

```bash
./deploy/scripts/deploy.sh
```

## 验证

- 首页：`http://<服务器IP>:8088/`
- 健康检查：`http://<服务器IP>:8088/api/v1/health`

云厂商安全组需放行 **8088**（及 SSH 22）。
