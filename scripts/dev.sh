#!/usr/bin/env bash
# 在 tiktok-backend 目录执行：启动 MySQL，并提示前后端启动命令
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
FRONTEND="$(cd "$ROOT/../tiktok-frontend" 2>/dev/null && pwd || true)"

echo "==> 启动 MySQL (docker compose)"
cd "$ROOT"
docker compose up -d mysql

echo ""
echo "MySQL 已启动。请在两个终端分别执行："
echo ""
echo "  终端1 - 后端:"
echo "    cd $ROOT && make run"
echo ""
if [ -n "$FRONTEND" ]; then
  echo "  终端2 - 前端:"
  echo "    cd $FRONTEND && npm install && npm run dev"
else
  echo "  终端2 - 前端: 请将 tiktok-frontend 克隆到与 tiktok-backend 同级目录"
fi
echo ""
