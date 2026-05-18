package monitoring

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "HTTP 请求总数"},
		[]string{"method", "path", "status"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds", Help: "HTTP 请求耗时",
			Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)
	MysqlQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "mysql_query_duration_seconds", Help: "MySQL 查询耗时",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"operation"},
	)
	MysqlConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "mysql_connections", Help: "MySQL 连接数"},
		[]string{"state"},
	)
	RedisCommandsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "redis_commands_total", Help: "Redis 命令总数"},
		[]string{"command"},
	)
	RedisCommandDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "redis_command_duration_seconds", Help: "Redis 命令耗时",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
		},
		[]string{"command"},
	)
	TenantPoolSize = prometheus.NewGauge(
		prometheus.GaugeOpts{Name: "tenant_pool_connections", Help: "租户连接池总连接数"},
	)
	TenantPoolCount = prometheus.NewGauge(
		prometheus.GaugeOpts{Name: "tenant_pool_tenants", Help: "活跃租户连接池数量"},
	)
)

func init() {
	prometheus.MustRegister(
		HttpRequestsTotal, HttpRequestDuration,
		MysqlQueryDuration, MysqlConnections,
		RedisCommandsTotal, RedisCommandDuration,
		TenantPoolSize, TenantPoolCount,
	)
}

// StartMetricsServer 启动 Prometheus 指标 HTTP 服务
func StartMetricsServer(addr string, logger *zap.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP"}`))
	})

	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	go func() {
		logger.Info("Prometheus 指标服务启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("指标服务异常", zap.Error(err))
		}
	}()
}

// RecordHTTPRequest 记录 HTTP 请求指标
func RecordHTTPRequest(method, path, status string, duration time.Duration) {
	HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
	HttpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// RecordRedisCommand 记录 Redis 命令指标
func RecordRedisCommand(cmd string, duration time.Duration) {
	RedisCommandsTotal.WithLabelValues(cmd).Inc()
	RedisCommandDuration.WithLabelValues(cmd).Observe(duration.Seconds())
}

// ShutdownMetrics 优雅关闭指标服务
func ShutdownMetrics() {}
