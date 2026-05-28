#!/usr/bin/env bash
# 在服务器上以 root 执行：安装 Nginx、MySQL（不包含任何账号密码）
set -euo pipefail

if [ "$(id -u)" -ne 0 ]; then
  echo "请使用 root 执行"
  exit 1
fi

WEB_PORT="${WEB_PORT:-8088}"

if [ -f /etc/os-release ]; then
  # shellcheck source=/dev/null
  . /etc/os-release
else
  echo "无法识别系统"
  exit 1
fi

install_centos() {
  yum install -y epel-release || true
  yum install -y nginx mysql-server git curl firewalld || yum install -y nginx mariadb-server git curl firewalld
  systemctl enable --now mysqld 2>/dev/null || systemctl enable --now mariadb
  systemctl enable --now nginx
  systemctl enable --now firewalld 2>/dev/null || true
  firewall-cmd --permanent --add-port="${WEB_PORT}/tcp" 2>/dev/null || true
  firewall-cmd --reload 2>/dev/null || true
}

install_debian() {
  export DEBIAN_FRONTEND=noninteractive
  apt-get update -y
  apt-get install -y nginx mysql-server git curl ufw
  systemctl enable --now mysql
  systemctl enable --now nginx
  ufw allow "${WEB_PORT}/tcp" 2>/dev/null || true
}

case "${ID:-}" in
  centos|rhel|rocky|almalinux|fedora|alinux|anolis)
    install_centos
    ;;
  ubuntu|debian)
    install_debian
    ;;
  *)
    command -v yum >/dev/null && install_centos && exit 0
    command -v apt-get >/dev/null && install_debian && exit 0
    echo "不支持的发行版: ${ID}"
    exit 1
    ;;
esac

if [ -z "${MYSQL_APP_PASS:-}" ]; then
  echo "请设置环境变量 MYSQL_APP_PASS 后重新执行数据库初始化，例如："
  echo "  MYSQL_APP_PASS='强密码' bash server-init.sh"
  exit 0
fi

mysql -e "CREATE DATABASE IF NOT EXISTS tiktok CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || \
  mysql -uroot -e "CREATE DATABASE IF NOT EXISTS tiktok CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -e "CREATE USER IF NOT EXISTS 'tiktok'@'localhost' IDENTIFIED BY '${MYSQL_APP_PASS}';" 2>/dev/null || true
mysql -e "GRANT ALL PRIVILEGES ON tiktok.* TO 'tiktok'@'localhost'; FLUSH PRIVILEGES;" 2>/dev/null || true

mkdir -p /opt/tiktok/backend /var/www/tiktok/frontend /etc/nginx/conf.d
echo "基础环境安装完成（Nginx 将监听 ${WEB_PORT}）。请运行 deploy.sh 发布应用。"
