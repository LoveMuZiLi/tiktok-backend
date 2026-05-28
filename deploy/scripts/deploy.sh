#!/usr/bin/env bash
# 本地执行：构建并上传到服务器
# 用法: ./deploy.sh [user@host]
set -euo pipefail

SERVER="${1:-root@39.96.28.124}"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
FRONTEND="${ROOT}/tiktok-frontend"
BACKEND="${ROOT}/tiktok-backend"
STAGE="${ROOT}/.deploy-stage"

echo "==> 构建前端"
cd "$FRONTEND"
npm install
npm run build

echo "==> 交叉编译后端 (linux/amd64)"
cd "$BACKEND"
export GOOS=linux GOARCH=amd64 CGO_ENABLED=0
go build -o "${STAGE}/tiktok-api" ./cmd/server

echo "==> 打包"
rm -rf "$STAGE/frontend" "$STAGE/backend"
mkdir -p "$STAGE/frontend" "$STAGE/backend/bin"
cp -R "$FRONTEND/dist/"* "$STAGE/frontend/"
cp "${STAGE}/tiktok-api" "$STAGE/backend/bin/"
cp "$BACKEND/.env.example" "$STAGE/backend/.env.example"
cp "$ROOT/deploy/env/production.env.example" "$STAGE/backend/.env.production.example"

tar -czf "${ROOT}/tiktok-release.tar.gz" -C "$STAGE" frontend backend

echo "==> 上传到 ${SERVER}"
scp -o StrictHostKeyChecking=no "${ROOT}/tiktok-release.tar.gz" "${SERVER}:/tmp/"
scp -o StrictHostKeyChecking=no "${ROOT}/deploy/nginx/tiktok.conf" "${SERVER}:/tmp/tiktok.conf"
scp -o StrictHostKeyChecking=no "${ROOT}/deploy/systemd/tiktok-api.service" "${SERVER}:/tmp/tiktok-api.service"
scp -o StrictHostKeyChecking=no "${ROOT}/deploy/scripts/remote-install.sh" "${SERVER}:/tmp/remote-install.sh"

echo "==> 远程安装"
ssh -o StrictHostKeyChecking=no "$SERVER" "bash /tmp/remote-install.sh"

echo "完成。访问: http://39.96.28.124/"
