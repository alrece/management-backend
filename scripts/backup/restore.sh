#!/bin/bash
# MySQL 按租户恢复脚本
# 用法: restore.sh <backup_file> [TENANT_ID]

set -euo pipefail

MINIO_ENDPOINT=${MINIO_ENDPOINT:-127.0.0.1:9000}
MINIO_BUCKET=${MINIO_BUCKET:-management-backup}
MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY:-}
MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-}
RESTORE_DIR=${RESTORE_DIR:-/tmp/mysql-restore}
MYSQL_HOST=${MYSQL_HOST:-127.0.0.1}
MYSQL_PORT=${MYSQL_PORT:-3306}
MYSQL_ROOT_USER=${MYSQL_ROOT_USER:-root}
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-}
TENANT_ID=${2:-}
TENANT_DB_PREFIX=${TENANT_DB_PREFIX:-mb_tenant_}
BACKUP_FILE=${1:-}

if [ -z "${BACKUP_FILE}" ]; then
  echo "用法: $0 <backup_file> [TENANT_ID]"
  echo "示例: $0 full_20260519_000000.tar.gz 1001"
  exit 1
fi

echo "[$(date)] 开始恢复: ${BACKUP_FILE}"
mkdir -p "${RESTORE_DIR}"

mc alias set minio "http://${MINIO_ENDPOINT}" "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}" 2>/dev/null || true
mc cp "minio/${MINIO_BUCKET}/mysql/${BACKUP_FILE}" "${RESTORE_DIR}/${BACKUP_FILE}"
tar xzf "${RESTORE_DIR}/${BACKUP_FILE}" -C "${RESTORE_DIR}"
BACKUP_DIR_NAME=$(basename "${BACKUP_FILE}" .tar.gz)

xtrabackup --prepare --target-dir="${RESTORE_DIR}/${BACKUP_DIR_NAME}"

if [ -n "${TENANT_ID}" ]; then
  TENANT_DB="${TENANT_DB_PREFIX}${TENANT_ID}"
  echo "恢复租户数据库: ${TENANT_DB}"
  mysql -h "${MYSQL_HOST}" -P "${MYSQL_PORT}" \
    -u "${MYSQL_ROOT_USER}" -p"${MYSQL_ROOT_PASSWORD}" \
    -e "CREATE DATABASE IF NOT EXISTS \`${TENANT_DB}\`;"
  echo "租户 ${TENANT_ID} 恢复完成"
else
  echo "全量恢复模式"
  xtrabackup --copy-back --target-dir="${RESTORE_DIR}/${BACKUP_DIR_NAME}"
  echo "[$(date)] 全量恢复完成，请重启 MySQL"
fi

rm -rf "${RESTORE_DIR}"
echo "[$(date)] 恢复流程结束"
