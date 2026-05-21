package discovery

import (
	"management-backend/internal/config"
	pkgdiscovery "management-backend/pkg/discoveryx"

	"go.uber.org/zap"
)

// ConsulRegistry Consul 服务注册器（类型别名）
type ConsulRegistry = pkgdiscovery.ConsulRegistry

// NewConsulRegistry 创建 Consul 注册器
func NewConsulRegistry(logger *zap.Logger) (*ConsulRegistry, error) {
	cfg := pkgdiscovery.ConsulConfig{
		Enabled:     config.C.Consul.Enabled,
		Address:     config.C.Consul.Address,
		ServiceAddr: config.C.Consul.ServiceAddr,
	}
	return pkgdiscovery.NewConsulRegistry(cfg, logger)
}

// HealthCheck HTTP 健康检查端点（委托到 pkg）
var HealthCheck = pkgdiscovery.HealthCheck
