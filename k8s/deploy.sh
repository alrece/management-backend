#!/bin/bash
# K8s 验证性部署脚本
# 用法: bash k8s/deploy.sh [build|deploy|status|clean|all]

set -e

NS="mb"
K8S_DIR="$(dirname "$0")"
PROJECT_ROOT="$(cd "$K8S_DIR/.." && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERROR]${NC} $1"; }

build_backend() {
    log "构建后端 Docker 镜像..."
    docker build -t mb-backend:latest "$PROJECT_ROOT"
    log "后端镜像构建完成: mb-backend:latest"
}

build_frontend() {
    log "构建前端 Docker 镜像..."
    docker build -t mb-frontend:latest "$PROJECT_ROOT/frontend"
    log "前端镜像构建完成: mb-frontend:latest"
}

build_all() {
    build_backend
    build_frontend
}

deploy() {
    log "创建命名空间..."
    kubectl apply -f "$K8S_DIR/00-namespace.yaml"

    log "部署 MySQL..."
    kubectl apply -f "$K8S_DIR/10-mysql.yaml"

    log "等待 MySQL 就绪..."
    kubectl rollout status deployment/mysql -n "$NS" --timeout=120s

    log "部署 Redis..."
    kubectl apply -f "$K8S_DIR/11-redis.yaml"

    log "等待 Redis 就绪..."
    kubectl rollout status deployment/redis -n "$NS" --timeout=60s

    log "部署后端..."
    kubectl apply -f "$K8S_DIR/20-backend.yaml"

    log "等待后端就绪..."
    kubectl rollout status deployment/mb-backend -n "$NS" --timeout=120s

    log "部署前端..."
    kubectl apply -f "$K8S_DIR/21-frontend.yaml"

    log "等待前端就绪..."
    kubectl rollout status deployment/mb-frontend -n "$NS" --timeout=60s

    log "部署完成!"
    status
}

status() {
    echo ""
    echo "========== K8s 部署状态 =========="
    kubectl get all -n "$NS" -o wide
    echo ""
    echo "========== Pods 详情 =========="
    kubectl get pods -n "$NS" -o custom-columns="NAME:.metadata.name,STATUS:.status.phase,READY:.status.containerStatuses[0].ready,RESTARTS:.status.containerStatuses[0].restartCount,IP:.status.podIP"
    echo ""

    # 检查 Ingress
    if kubectl get ingress -n "$NS" 2>/dev/null | grep -q mb; then
        echo "========== Ingress =========="
        kubectl get ingress -n "$NS"
    fi
}

port_forward() {
    log "端口转发: 前端 http://localhost:8080 -> mb-frontend:80"
    log "端口转发: 后端 http://localhost:8081 -> mb-backend:8081"
    kubectl port-forward -n "$NS" svc/mb-frontend 8080:80 &
    kubectl port-forward -n "$NS" svc/mb-backend 8081:8081 &
    wait
}

clean() {
    warn "清理所有 K8s 资源..."
    kubectl delete -f "$K8S_DIR/21-frontend.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/20-backend.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/11-redis.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/10-mysql.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/00-namespace.yaml" --ignore-not-found
    log "清理完成"
}

logs() {
    local svc="${1:-mb-backend}"
    kubectl logs -f deployment/"$svc" -n "$NS" --tail=100
}

case "${1:-all}" in
    build)
        build_all
        ;;
    build-backend)
        build_backend
        ;;
    build-frontend)
        build_frontend
        ;;
    deploy)
        deploy
        ;;
    status)
        status
        ;;
    port-forward|pf)
        port_forward
        ;;
    clean)
        clean
        ;;
    logs)
        logs "${2:-mb-backend}"
        ;;
    all)
        build_all
        deploy
        ;;
    *)
        echo "用法: bash k8s/deploy.sh [build|build-backend|build-frontend|deploy|status|port-forward|clean|logs [service]|all]"
        echo ""
        echo "  all              构建并部署全部 (默认)"
        echo "  build            构建前后端镜像"
        echo "  build-backend    仅构建后端镜像"
        echo "  build-frontend   仅构建前端镜像"
        echo "  deploy           部署到 K8s (需先 build)"
        echo "  status           查看部署状态"
        echo "  port-forward     端口转发到本地"
        echo "  clean            清理所有资源"
        echo "  logs [service]   查看日志 (默认: mb-backend)"
        exit 1
        ;;
esac
