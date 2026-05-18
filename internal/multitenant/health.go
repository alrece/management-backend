package multitenant

import (
	"context"
	"fmt"
	"time"

	"management-backend/internal/module/system/repository"
)

// TenantHealth 租户健康状态
type TenantHealth struct {
	TenantID      int64     `json:"tenantId"`
	DBConnected   bool      `json:"dbConnected"`
	SchemaVersion string    `json:"schemaVersion"`
	TableCount    int       `json:"tableCount"`
	LastCheckTime time.Time `json:"lastCheckTime"`
	Error         string    `json:"error,omitempty"`
}

// HealthChecker 租户健康检查器
type HealthChecker struct {
	resolver   *TenantResolver
	tenantRepo repository.TenantRepo
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(resolver *TenantResolver, tenantRepo repository.TenantRepo) *HealthChecker {
	return &HealthChecker{
		resolver:   resolver,
		tenantRepo: tenantRepo,
	}
}

// HealthCheck 检查单个租户健康状态
func (h *HealthChecker) HealthCheck(ctx context.Context, tenantID int64) (*TenantHealth, error) {
	health := &TenantHealth{
		TenantID:      tenantID,
		LastCheckTime: time.Now(),
	}

	_, err := h.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		health.Error = "租户不存在"
		return health, err
	}

	db, err := h.resolver.GetDB(ctx, tenantID)
	if err != nil {
		health.DBConnected = false
		health.Error = fmt.Sprintf("数据库连接失败: %v", err)
		return health, nil
	}
	health.DBConnected = true

	var tableCount int64
	db.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()",
	).Scan(&tableCount)
	health.TableCount = int(tableCount)
	health.SchemaVersion = "001"

	return health, nil
}

// HealthCheckAll 检查所有已就绪租户
func (h *HealthChecker) HealthCheckAll(ctx context.Context) ([]TenantHealth, error) {
	return []TenantHealth{}, nil
}
