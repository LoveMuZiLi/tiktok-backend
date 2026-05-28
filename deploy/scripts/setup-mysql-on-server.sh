#!/usr/bin/env bash
# 在服务器上以 root 执行：创建 tiktok 库、用户并生成随机密码（仅打印一次）
set -euo pipefail

[ "$(id -u)" -eq 0 ] || { echo "请使用 root 执行"; exit 1; }

DB_NAME="${DB_NAME:-tiktok}"
DB_USER="${DB_USER:-tiktok}"
ENV_FILE="${ENV_FILE:-/opt/tiktok/backend/.env}"

if [ -z "${MYSQL_APP_PASS:-}" ]; then
  MYSQL_APP_PASS="$(openssl rand -base64 12 | tr -d '/+=' | head -c 16)"
fi

if ! command -v mysql >/dev/null 2>&1; then
  echo "未找到 mysql 客户端，请先执行 server-init.sh 安装 MySQL"
  exit 1
fi

mysql -e "CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || \
  mysql -uroot -e "CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

mysql -e "CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${MYSQL_APP_PASS}';" 2>/dev/null || \
  mysql -uroot -e "ALTER USER '${DB_USER}'@'localhost' IDENTIFIED BY '${MYSQL_APP_PASS}';" 2>/dev/null || \
  mysql -uroot -e "CREATE USER '${DB_USER}'@'localhost' IDENTIFIED BY '${MYSQL_APP_PASS}';"

mysql -e "GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'localhost'; FLUSH PRIVILEGES;" 2>/dev/null || \
  mysql -uroot -e "GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'localhost'; FLUSH PRIVILEGES;"

mkdir -p "$(dirname "$ENV_FILE")"
if [ -f "$ENV_FILE" ]; then
  sed -i.bak "s/^DB_PASSWORD=.*/DB_PASSWORD=${MYSQL_APP_PASS}/" "$ENV_FILE" 2>/dev/null || \
    sed -i '' "s/^DB_PASSWORD=.*/DB_PASSWORD=${MYSQL_APP_PASS}/" "$ENV_FILE" 2>/dev/null || true
  grep -q '^DB_PASSWORD=' "$ENV_FILE" || echo "DB_PASSWORD=${MYSQL_APP_PASS}" >> "$ENV_FILE"
else
  cat > "$ENV_FILE" <<EOF
PORT=8080
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=${DB_USER}
DB_PASSWORD=${MYSQL_APP_PASS}
DB_NAME=${DB_NAME}
ALLOWED_ORIGINS=http://YOUR_SERVER_IP:8088
EOF
fi

systemctl restart tiktok-api 2>/dev/null || true

echo ""
echo "========== MySQL 已配置 =========="
echo "数据库: ${DB_NAME}"
echo "用户:   ${DB_USER}@localhost"
echo "密码:   ${MYSQL_APP_PASS}"
echo "已写入: ${ENV_FILE}"
echo "================================="
echo "请保存上述密码；重启后端: systemctl restart tiktok-api"
