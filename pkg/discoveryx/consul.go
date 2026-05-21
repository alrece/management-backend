package discoveryx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

// ConsulConfig Consul 服务注册配置
type ConsulConfig struct {
	Enabled     bool
	Address     string
	ServiceAddr string
	TraefikTags []string // Traefik 路由标签
}

// ConsulRegistry Consul 服务注册器
type ConsulRegistry struct {
	client *api.Client
	logger *zap.Logger
	cfg    ConsulConfig
	svcID  string
}

// NewConsulRegistry 创建 Consul 注册器
func NewConsulRegistry(cfg ConsulConfig, logger *zap.Logger) (*ConsulRegistry, error) {
	consulCfg := api.DefaultConfig()
	if cfg.Address != "" {
		consulCfg.Address = cfg.Address
	}

	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, fmt.Errorf("Consul 客户端创建失败: %w", err)
	}

	return &ConsulRegistry{
		client: client,
		logger: logger,
		cfg:    cfg,
	}, nil
}

// Register 注册服务到 Consul（支持 Traefik 自动发现标签）
func (r *ConsulRegistry) Register(serviceName string, port int) error {
	if !r.cfg.Enabled {
		r.logger.Info("Consul 注册已禁用，跳过")
		return nil
	}

	r.svcID = fmt.Sprintf("%s-%d", serviceName, port)
	addr := r.cfg.ServiceAddr
	if addr == "" {
		addr = "127.0.0.1"
	}

	tags := r.cfg.TraefikTags
	if len(tags) == 0 {
		tags = []string{"traefik.enable=true"}
	}

	registration := &api.AgentServiceRegistration{
		ID:      r.svcID,
		Name:    serviceName,
		Port:    port,
		Address: addr,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", addr, port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	if err := r.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("Consul 注册失败: %w", err)
	}

	r.logger.Info("Consul 服务注册成功", zap.String("service", serviceName))
	return nil
}

// Deregister 从 Consul 注销
func (r *ConsulRegistry) Deregister() {
	if r.svcID != "" {
		_ = r.client.Agent().ServiceDeregister(r.svcID)
		r.logger.Info("Consul 服务已注销")
	}
}

// HealthCheck HTTP 健康检查端点
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"UP","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}
