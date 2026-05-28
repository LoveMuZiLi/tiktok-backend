# TikTok 线上部署（39.96.28.124）

## 架构

```
浏览器 → Nginx:80 → 静态前端 (/var/www/tiktok/frontend)
                 → /api/* → Go API :8080 → MySQL
```

## 一、服务器初始化（仅需一次）

SSH 登录后执行：

```bash
curl -fsSL https://raw.githubusercontent.com/LoveMuZiLi/tiktok-backend/main/../deploy/scripts/server-init.sh
# 或从本地上传：
scp deploy/scripts/server-init.sh root@39.96.28.124:/tmp/
ssh root@39.96.28.124 'bash /tmp/server-init.sh'
```

编辑 `/opt/tiktok/backend/.env` 中的数据库密码。

## 二、本地一键发布

```bash
chmod +x deploy/scripts/deploy.sh deploy/scripts/server-init.sh
./deploy/scripts/deploy.sh root@39.96.28.124
```

## 三、验证

- 首页：http://39.96.28.124/
- 健康检查：http://39.96.28.124/api/v1/health

## SSH 说明

若密码登录失败，请在阿里云控制台：

1. 确认 root 密码已重置
2. 安全组放行 **22、80、443**
3. 或配置 SSH 公钥后使用：`./deploy/scripts/deploy.sh root@39.96.28.124`

**切勿将 root 密码提交到 Git。**
