package multitenant

import (
	"context"

	"management-backend/internal/config"
	pkgtenant "management-backend/pkg/tenantx"

	"gorm.io/gorm"
)

// TenantResolver 租户数据库解析器（类型别名）
type TenantResolver = pkgtenant.TenantResolver

// NewTenantResolver 创建租户解析器（从全局配置读取参数）
func NewTenantResolver() *TenantResolver {
	pgCfg := pkgtenant.PostgresConfig{
		Host:     config.C.Postgres.Host,
		Port:     config.C.Postgres.Port,
		Username: config.C.Postgres.Username,
		Password: config.C.Postgres.Password,
		SSLMode:  config.C.Postgres.SSLMode,
	}
	tenantCfg := pkgtenant.TenantConfig{
		MaxPoolSize:       config.C.Tenant.MaxPoolSize,
		IdleTimeout:       config.C.Tenant.IdleTimeout,
		MaxConnsPerTenant: config.C.Tenant.MaxConnsPerTenant,
		DBNamePrefix:      config.C.Tenant.DBNamePrefix,
	}
	return pkgtenant.NewTenantResolver(pgCfg, tenantCfg)
}

// GetDB 获取租户数据库连接（便捷方法，委托到 TenantResolver）
func GetDB(ctx context.Context, r *TenantResolver, tenantID int64) (*gorm.DB, error) {
	return r.GetDB(ctx, tenantID)
}
