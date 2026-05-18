---
wave: 5
depends_on: ["06-PLAN.md"]
files_modified:
  - internal/monitoring/prometheus.go
  - configs/grafana/dashboards/system.json
  - configs/grafana/datasources.yaml
  - configs/prometheus/prometheus.yml
  - configs/prometheus/alerts.yml
  - configs/loki/loki-config.yaml
  - configs/promtail/promtail-config.yaml
  - scripts/backup/mysql-backup.sh
  - scripts/backup/restore.sh
  - cmd/system/main.go
  - docker-compose.prod.yml
autonomous: true
requirements_addressed:
  - NFR-06-01
  - NFR-06-02
  - NFR-06-03
  - NFR-06-04
  - NFR-06-05
---

# Plan 08: 可观测性与运维

## Objective

Prometheus + Grafana、Loki 日志、告警规则、MySQL 备份、优雅关闭完善。

## Tasks

### Task 8.1: Prometheus + Grafana

<read_first>
- cmd/system/main.go
- internal/multitenant/pool.go
</read_first>

<action>
internal/monitoring/prometheus.go：Go runtime + HTTP + MySQL + Redis + 租户连接池指标，:9091/metrics。Prometheus/Grafana 配置。docker-compose.prod.yml 添加服务。
</action>

<acceptance_criteria>
- :9091/metrics 暴露 5 类指标
- Grafana dashboard JSON 存在
</acceptance_criteria>

---

### Task 8.2: Zap → Loki

<read_first>
- cmd/system/main.go
</read_first>

<action>
Zap JSON 含 trace_id, tenant_id, request_id。Loki + Promtail 配置。docker-compose.prod.yml 添加服务。
</action>

<acceptance_criteria>
- Zap 日志含 3 个追踪字段
- Loki + Promtail 配置存在
</acceptance_criteria>

---

### Task 8.3: 告警规则

<read_first>
- configs/prometheus/prometheus.yml
</read_first>

<action>
alerts.yml：连接池>80%、Redis>85%、错误率>5%、P99>1s、初始化失败、goroutine>1000。
</action>

<acceptance_criteria>
- 6 条告警规则
</acceptance_criteria>

---

### Task 8.4: MySQL 备份脚本

<read_first>
- internal/config/config.go
</read_first>

<action>
XtraBackup 全量 + binlog → MinIO。7 天保留。按租户恢复。cron 定时执行。
</action>

<acceptance_criteria>
- backup.sh + restore.sh 存在
- 支持按租户恢复
</acceptance_criteria>

---

### Task 8.5: 优雅关闭完善

<read_first>
- cmd/system/main.go
- internal/saga/orchestrator.go
</read_first>

<action>
关闭时持久化 Saga 状态到 MySQL。排空租户连接池。全流程 Zap 日志。
</action>

<acceptance_criteria>
- Saga 状态持久化
- 连接池排空 + 日志
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
docker-compose -f docker-compose.prod.yml config
```

## Must Haves

- [ ] Prometheus 5 类指标
- [ ] Grafana dashboard
- [ ] Loki 日志（3 个追踪字段）
- [ ] 6 条告警规则
- [ ] MySQL 备份脚本
- [ ] 优雅关闭完善
