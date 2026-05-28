#!/usr/bin/env bash
# 服务器端：部署前端静态资源 + Go 后端服务（systemd），由 Nginx :8088 转发 /api
set -euo pipefail

[ "$(id -u)" -eq 0 ] || { echo "需要 root"; exit 1; }

WEB_PORT="${WEB_PORT:-8088}"
API_PORT="${API_PORT:-8080}"

mkdir -p /opt/tiktok/backend/bin /var/www/tiktok/frontend
tar -xzf /tmp/tiktok-release.tar.gz -C /tmp

if [ ! -f /tmp/backend/bin/tiktok-api ]; then
  echo "错误: 发布包中缺少 backend/bin/tiktok-api"
  exit 1
fi

echo "==> 部署前端静态文件"
mkdir -p /var/www/tiktok/frontend
if command -v rsync >/dev/null 2>&1; then
  rsync -a --delete /tmp/frontend/ /var/www/tiktok/frontend/
else
  find /var/www/tiktok/frontend -mindepth 1 -delete 2>/dev/null || rm -rf /var/www/tiktok/frontend/*
  cp -a /tmp/frontend/. /var/www/tiktok/frontend/
fi

echo "==> 部署后端二进制"
install -m 755 /tmp/backend/bin/tiktok-api /opt/tiktok/backend/bin/tiktok-api

if [ ! -f /opt/tiktok/backend/.env ]; then
  cp /tmp/backend/.env.production.example /opt/tiktok/backend/.env
  echo "已生成 /opt/tiktok/backend/.env"
  echo "请编辑 DB_PASSWORD、ALLOWED_ORIGINS 后执行: systemctl restart tiktok-api"
fi

echo "==> 配置 Nginx (${WEB_PORT}) 与 systemd (API :${API_PORT})"
install -m 644 /tmp/tiktok.conf /etc/nginx/conf.d/tiktok.conf
install -m 644 /tmp/tiktok-api.service /etc/systemd/system/tiktok-api.service

systemctl daemon-reload
systemctl enable tiktok-api nginx
systemctl restart tiktok-api

nginx -t
systemctl reload nginx || systemctl restart nginx

echo "==> 等待后端启动"
for i in $(seq 1 30); do
  if curl -fsS "http://127.0.0.1:${API_PORT}/api/v1/health" >/dev/null 2>&1; then
    echo "后端健康检查通过 (127.0.0.1:${API_PORT})"
    break
  fi
  if [ "$i" -eq 30 ]; then
    echo "警告: 后端未响应，请检查: journalctl -u tiktok-api -n 50"
    journalctl -u tiktok-api -n 20 --no-pager || true
    exit 1
  fi
  sleep 1
done

if curl -fsS "http://127.0.0.1:${WEB_PORT}/api/v1/health" >/dev/null 2>&1; then
  echo "Nginx 反代 /api 正常 (127.0.0.1:${WEB_PORT})"
else
  echo "警告: Nginx :${WEB_PORT} 无法访问 /api，请检查防火墙与安全组"
fi

echo "部署完成。前端与 API 均由 Nginx :${WEB_PORT} 对外提供（/api -> :${API_PORT}）。"
systemctl status tiktok-api --no-pager | head -12
