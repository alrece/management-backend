#!/bin/bash
# K8s 微服务全栈部署脚本（PostgreSQL 版）
# 用法: bash k8s/deploy.sh [build|deploy|status|clean|all|port-forward]

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

build_system() {
    log "构建 system-service Docker 镜像..."
    docker build -t mb-system:latest "$PROJECT_ROOT"
    log "镜像构建完成: mb-system:latest"
}

build_job() {
    log "构建 job-service Docker 镜像..."
    docker build -t mb-job:latest "$PROJECT_ROOT/services/job-service"
    log "镜像构建完成: mb-job:latest"
}

build_workflow() {
    log "构建 workflow-service Docker 镜像..."
    docker build -t mb-workflow:latest "$PROJECT_ROOT/services/workflow-service"
    log "镜像构建完成: mb-workflow:latest"
}

build_frontend() {
    log "构建前端 Docker 镜像..."
    docker build -t mb-frontend:latest "$PROJECT_ROOT/frontend"
    log "镜像构建完成: mb-frontend:latest"
}

build_all() {
    build_system
    build_job
    build_workflow
    build_frontend
}

deploy_infra() {
    log "创建命名空间..."
    kubectl apply -f "$K8S_DIR/00-namespace.yaml"

    log "部署 PostgreSQL..."
    kubectl apply -f "$K8S_DIR/01-postgresql.yaml"
    log "等待 PostgreSQL 就绪..."
    kubectl rollout status statefulset/postgres -n "$NS" --timeout=120s

    log "部署 Redis..."
    kubectl apply -f "$K8S_DIR/11-redis.yaml"
    log "等待 Redis 就绪..."
    kubectl rollout status deployment/redis -n "$NS" --timeout=60s

    log "部署 MinIO..."
    kubectl apply -f "$K8S_DIR/03-minio.yaml"
    kubectl rollout status deployment/minio -n "$NS" --timeout=60s

    log "部署 Consul..."
    kubectl apply -f "$K8S_DIR/04-consul.yaml"
    kubectl rollout status deployment/consul -n "$NS" --timeout=60s

    log "部署 Jaeger..."
    kubectl apply -f "$K8S_DIR/05-jaeger.yaml"
    kubectl rollout status deployment/jaeger -n "$NS" --timeout=60s
}

deploy_apps() {
    log "部署 system-service..."
    kubectl apply -f "$K8S_DIR/10-system-service.yaml"
    kubectl rollout status deployment/mb-system -n "$NS" --timeout=120s

    log "部署 job-service..."
    kubectl apply -f "$K8S_DIR/11-job-service.yaml"
    kubectl rollout status deployment/mb-job -n "$NS" --timeout=120s

    log "部署 workflow-service..."
    kubectl apply -f "$K8S_DIR/12-workflow-service.yaml"
    kubectl rollout status deployment/mb-workflow -n "$NS" --timeout=120s

    log "部署 frontend..."
    kubectl apply -f "$K8S_DIR/13-frontend.yaml"
    kubectl rollout status deployment/mb-frontend -n "$NS" --timeout=60s
}

deploy() {
    deploy_infra
    deploy_apps
    log "全部部署完成!"
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
    echo "========== Services =========="
    kubectl get svc -n "$NS"
    echo ""
    echo "========== 访问地址 =========="
    echo "  前端: http://localhost:30080"
    echo "  后端 API: kubectl port-forward -n $NS svc/mb-system 8081:8081"
    echo "  Jaeger UI: kubectl port-forward -n $NS svc/jaeger 16686:16686"
    echo "  MinIO Console: kubectl port-forward -n $NS svc/minio 9001:9001"
    echo "  Consul UI: kubectl port-forward -n $NS svc/consul 8500:8500"
}

port_forward() {
    log "端口转发: 前端 http://localhost:30080 (NodePort, 无需转发)"
    log "启动辅助服务端口转发..."
    kubectl port-forward -n "$NS" svc/jaeger 16686:16686 &
    kubectl port-forward -n "$NS" svc/minio 9001:9001 &
    kubectl port-forward -n "$NS" svc/consul 8500:8500 &
    log "按 Ctrl+C 停止"
    wait
}

clean() {
    warn "清理所有 K8s 资源..."
    kubectl delete -f "$K8S_DIR/13-frontend.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/12-workflow-service.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/11-job-service.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/10-system-service.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/05-jaeger.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/04-consul.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/03-minio.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/11-redis.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/01-postgresql.yaml" --ignore-not-found
    kubectl delete -f "$K8S_DIR/00-namespace.yaml" --ignore-not-found
    log "清理完成"
}

logs() {
    local svc="${1:-mb-system}"
    kubectl logs -f deployment/"$svc" -n "$NS" --tail=100
}

case "${1:-all}" in
    build)
        build_all
        ;;
    build-system)
        build_system
        ;;
    build-job)
        build_job
        ;;
    build-workflow)
        build_workflow
        ;;
    build-frontend)
        build_frontend
        ;;
    deploy)
        deploy
        ;;
    deploy-infra)
        deploy_infra
        ;;
    deploy-apps)
        deploy_apps
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
        logs "${2:-mb-system}"
        ;;
    all)
        build_all
        deploy
        ;;
    *)
        echo "用法: bash k8s/deploy.sh [命令]"
        echo ""
        echo "  all               构建并部署全部 (默认)"
        echo "  build             构建所有镜像"
        echo "  build-system      仅构建 system-service"
        echo "  build-job         仅构建 job-service"
        echo "  build-workflow    仅构建 workflow-service"
        echo "  build-frontend    仅构建 frontend"
        echo "  deploy            部署全部 (infra + apps)"
        echo "  deploy-infra      仅部署基础设施"
        echo "  deploy-apps       仅部署应用服务"
        echo "  status            查看部署状态"
        echo "  port-forward      端口转发辅助服务"
        echo "  clean             清理所有资源"
        echo "  logs [service]    查看日志 (默认: mb-system)"
        exit 1
        ;;
esac
