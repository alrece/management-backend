#!/bin/bash
# MySQL 全量备份脚本
# XtraBackup 全量 + binlog → MinIO，保留 7 天

set -euo pipefail

MYSQL_HOST=${MYSQL_HOST:-127.0.0.1}
MYSQL_PORT=${MYSQL_PORT:-3306}
MYSQL_USER=${MYSQL_USER:-backup}
MYSQL_PASSWORD=${MYSQL_PASSWORD:-}
BACKUP_DIR=${BACKUP_DIR:-/tmp/mysql-backup}
MINIO_ENDPOINT=${MINIO_ENDPOINT:-127.0.0.1:9000}
MINIO_BUCKET=${MINIO_BUCKET:-management-backup}
MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY:-}
MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-}
RETENTION_DAYS=${RETENTION_DAYS:-7}
DATE=$(date +%Y%m%d_%H%M%S)

echo "[$(date)] 开始 MySQL 备份..."
mkdir -p "${BACKUP_DIR}"
BACKUP_PATH="${BACKUP_DIR}/full_${DATE}"

# XtraBackup 全量备份
xtrabackup --backup \
  --host="${MYSQL_HOST}" \
  --port="${MYSQL_PORT}" \
  --user="${MYSQL_USER}" \
  --password="${MYSQL_PASSWORD}" \
  --target-dir="${BACKUP_PATH}" \
  --parallel=4

# 压缩
tar czf "${BACKUP_PATH}.tar.gz" -C "${BACKUP_DIR}" "full_${DATE}"
rm -rf "${BACKUP_PATH}"

# 上传 MinIO
mc alias set minio "http://${MINIO_ENDPOINT}" "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}" 2>/dev/null || true
mc cp "${BACKUP_PATH}.tar.gz" "minio/${MINIO_BUCKET}/mysql/full_${DATE}.tar.gz"
rm -f "${BACKUP_PATH}.tar.gz"

# 清理过期备份
mc rm --recursive --force "minio/${MINIO_BUCKET}/mysql/" \
  --older-than "${RETENTION_DAYS}d" 2>/dev/null || true

echo "[$(date)] MySQL 备份完成: full_${DATE}.tar.gz"
