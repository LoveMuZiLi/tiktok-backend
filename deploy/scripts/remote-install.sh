#!/usr/bin/env bash
# 服务器端：解压发布包、配置 Nginx 与 systemd
set -euo pipefail

[ "$(id -u)" -eq 0 ] || { echo "需要 root"; exit 1; }

mkdir -p /opt/tiktok/backend/bin /var/www/tiktok/frontend
tar -xzf /tmp/tiktok-release.tar.gz -C /tmp

rsync -a --delete /tmp/frontend/ /var/www/tiktok/frontend/
install -m 755 /tmp/backend/bin/tiktok-api /opt/tiktok/backend/bin/tiktok-api 2>/dev/null || \
  cp /tmp/backend/bin/tiktok-api /opt/tiktok/backend/bin/tiktok-api
chmod +x /opt/tiktok/backend/bin/tiktok-api

if [ ! -f /opt/tiktok/backend/.env ]; then
  cp /tmp/backend/.env.production.example /opt/tiktok/backend/.env
  echo "已生成 /opt/tiktok/backend/.env ，请修改 DB_PASSWORD 后重启服务"
fi

install -m 644 /tmp/tiktok.conf /etc/nginx/conf.d/tiktok.conf
sed -i 's/YOUR_DOMAIN/_/' /etc/nginx/conf.d/tiktok.conf 2>/dev/null || \
  sed -i '' 's/YOUR_DOMAIN/_/' /etc/nginx/conf.d/tiktok.conf 2>/dev/null || true

install -m 644 /tmp/tiktok-api.service /etc/systemd/system/tiktok-api.service
systemctl daemon-reload
systemctl enable tiktok-api
systemctl restart tiktok-api

nginx -t
systemctl reload nginx || systemctl restart nginx

echo "部署完成"
systemctl status tiktok-api --no-pager | head -15
