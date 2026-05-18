package multitenant

import (
	"context"
	"fmt"
	"sync"
	"time"

	"management-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TenantResolver 租户数据库解析器
type TenantResolver struct {
	pool *TenantPool
	mu   sync.RWMutex
	dbs  map[int64]*gorm.DB
}

// NewTenantResolver 创建租户解析器
func NewTenantResolver() *TenantResolver {
	cfg := config.C.Tenant
	maxSize := cfg.MaxPoolSize
	if maxSize <= 0 {
		maxSize = 50
	}
	idleTimeout := time.Duration(cfg.IdleTimeout) * time.Minute
	if idleTimeout <= 0 {
		idleTimeout = 30 * time.Minute
	}

	r := &TenantResolver{
		dbs: make(map[int64]*gorm.DB),
	}
	r.pool = NewTenantPool(maxSize, idleTimeout)
	go r.pool.cleanupLoop(context.Background())
	return r
}

// GetDB 获取租户专属数据库连接
func (r *TenantResolver) GetDB(ctx context.Context, tenantID int64) (*gorm.DB, error) {
	r.mu.RLock()
	db, ok := r.dbs[tenantID]
	r.mu.RUnlock()
	if ok {
		r.pool.Touch(tenantID)
		return db, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if db, ok = r.dbs[tenantID]; ok {
		r.pool.Touch(tenantID)
		return db, nil
	}

	// LRU 驱逐
	if len(r.dbs) >= r.pool.maxSize {
		evictedID := r.pool.EvictOldest()
		if evictedID > 0 {
			if oldDB, exists := r.dbs[evictedID]; exists {
				sqlDB, _ := oldDB.DB()
				if sqlDB != nil {
					_ = sqlDB.Close()
				}
				delete(r.dbs, evictedID)
			}
		}
	}

	newDB, err := r.createTenantDB(tenantID)
	if err != nil {
		return nil, fmt.Errorf("创建租户数据库连接失败: %w", err)
	}

	r.dbs[tenantID] = newDB
	r.pool.Touch(tenantID)
	return newDB, nil
}

func (r *TenantResolver) createTenantDB(tenantID int64) (*gorm.DB, error) {
	cfg := config.C.Tenant
	prefix := cfg.DBNamePrefix
	if prefix == "" {
		prefix = "mb_tenant_"
	}
	dbName := fmt.Sprintf("%s%d", prefix, tenantID)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.C.MySQL.Username,
		config.C.MySQL.Password,
		config.C.MySQL.Host,
		config.C.MySQL.Port,
		dbName,
		config.C.MySQL.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		maxConns := cfg.MaxConnsPerTenant
		if maxConns <= 0 {
			maxConns = 10
		}
		sqlDB.SetMaxIdleConns(maxConns / 2)
		sqlDB.SetMaxOpenConns(maxConns)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db, nil
}

// PoolStats 返回连接池统计
func (r *TenantResolver) PoolStats() PoolMetrics {
	return r.pool.GetMetrics()
}
