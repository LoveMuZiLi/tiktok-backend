#!/usr/bin/env bash
# 本地执行：构建并上传到服务器（不在仓库中保存 SSH 地址与密码）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
DEPLOY_ENV="${DEPLOY_ENV:-$ROOT/deploy/deploy.local.env}"

if [ -f "$DEPLOY_ENV" ]; then
  # shellcheck source=/dev/null
  . "$DEPLOY_ENV"
fi

SERVER="${1:-${DEPLOY_HOST:-}}"
if [ -z "$SERVER" ]; then
  echo "用法: DEPLOY_HOST=user@host ./deploy/scripts/deploy.sh"
  echo "或:   cp deploy/env/deploy.local.env.example deploy/deploy.local.env 后填写 DEPLOY_HOST"
  exit 1
fi

PUBLIC_PORT="${PUBLIC_PORT:-8088}"
PUBLIC_SCHEME="${PUBLIC_SCHEME:-http}"

FRONTEND="${ROOT}/../tiktok-frontend"
BACKEND="${ROOT}"
STAGE="${ROOT}/.deploy-stage"

if [ ! -d "$FRONTEND" ]; then
  FRONTEND="$(cd "$ROOT/../../tiktok-frontend" 2>/dev/null && pwd)" || true
fi
if [ ! -d "$FRONTEND" ]; then
  echo "找不到 tiktok-frontend 目录"
  exit 1
fi

echo "==> 构建前端"
cd "$FRONTEND"
npm install
npm run build

echo "==> 交叉编译后端 (linux/amd64)"
cd "$BACKEND"
export GOOS=linux GOARCH=amd64 CGO_ENABLED=0
mkdir -p "$STAGE/backend/bin"
go build -o "$STAGE/backend/bin/tiktok-api" ./cmd/server

echo "==> 打包"
rm -rf "$STAGE/frontend"
mkdir -p "$STAGE/frontend" "$STAGE/backend/bin"
cp -R "$FRONTEND/dist/"* "$STAGE/frontend/"
cp "$BACKEND/.env.example" "$STAGE/backend/.env.example"
cp "$ROOT/deploy/env/production.env.example" "$STAGE/backend/.env.production.example"

tar -czf "${ROOT}/tiktok-release.tar.gz" -C "$STAGE" frontend backend

echo "==> 上传到服务器"
scp -o StrictHostKeyChecking=no "${ROOT}/tiktok-release.tar.gz" "${SERVER}:/tmp/"
scp -o StrictHostKeyChecking=no "${ROOT}/deploy/nginx/tiktok.conf" "${SERVER}:/tmp/tiktok.conf"
scp -o StrictHostKeyChecking=no "${ROOT}/deploy/systemd/tiktok-api.service" "${SERVER}:/tmp/tiktok-api.service"
scp -o StrictHostKeyChecking=no "${ROOT}/deploy/scripts/remote-install.sh" "${SERVER}:/tmp/remote-install.sh"

echo "==> 远程安装"
ssh -o StrictHostKeyChecking=no "$SERVER" "bash /tmp/remote-install.sh"

HOST_PART="${SERVER#*@}"
echo "完成。请在安全组放行 TCP ${PUBLIC_PORT}，访问: ${PUBLIC_SCHEME}://${HOST_PART}:${PUBLIC_PORT}/"
