package multitenant

import (
	"context"
	"fmt"
	"time"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/internal/config"
	"management-backend/pkg/errcode"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initializer 租户初始化器（幂等状态机）
type Initializer struct {
	tenantRepo repository.TenantRepo
	rdb        *redis.Client
}

// NewInitializer 创建租户初始化器
func NewInitializer(tenantRepo repository.TenantRepo, rdb *redis.Client) *Initializer {
	return &Initializer{
		tenantRepo: tenantRepo,
		rdb:        rdb,
	}
}

// Initialize 租户初始化（幂等方法）
// 状态机: 0(待初始化)→1(初始化中)→2(已就绪)，1→3(失败)，3→0(可重试)
func (in *Initializer) Initialize(ctx context.Context, tenantID int64) error {
	tenant, err := in.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return errcode.Err(errcode.TenantNotFound)
	}

	if tenant.InitStatus == smodel.InitStatusReady {
		return nil
	}

	if tenant.InitStatus != smodel.InitStatusPending && tenant.InitStatus != smodel.InitStatusFailed {
		return fmt.Errorf("租户当前状态 %d 不允许初始化", tenant.InitStatus)
	}

	// Redis 分布式锁
	lockKey := fmt.Sprintf("tenant:init:%d", tenantID)
	locked, err := in.rdb.SetNX(ctx, lockKey, 1, 10*time.Minute).Result()
	if err != nil {
		return fmt.Errorf("获取分布式锁失败: %w", err)
	}
	if !locked {
		return fmt.Errorf("租户初始化正在进行中")
	}
	defer in.rdb.Del(ctx, lockKey)

	// CAS: 0/3 → 1（初始化中）
	oldStatus := tenant.InitStatus
	ok, err := in.tenantRepo.UpdateInitStatus(ctx, tenantID, oldStatus, smodel.InitStatusIniting)
	if err != nil || !ok {
		return fmt.Errorf("状态转换失败")
	}

	if err := in.doInit(ctx, tenant); err != nil {
		_, _ = in.tenantRepo.UpdateInitStatus(ctx, tenantID, smodel.InitStatusIniting, smodel.InitStatusFailed)
		return fmt.Errorf("初始化失败: %w", err)
	}

	_, err = in.tenantRepo.UpdateInitStatus(ctx, tenantID, smodel.InitStatusIniting, smodel.InitStatusReady)
	if err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}
	return nil
}

func (in *Initializer) doInit(ctx context.Context, tenant *smodel.Tenant) error {
	if err := in.createDatabase(ctx, tenant.DBName); err != nil {
		return fmt.Errorf("创建数据库失败: %w", err)
	}
	return nil
}

func (in *Initializer) createDatabase(ctx context.Context, dbName string) error {
	cfg := config.C.Postgres
	// 连接默认的 postgres 库来创建新数据库
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()

	// PG 创建数据库（不能用事务，需直接 Exec）
	return db.WithContext(ctx).Exec(
		fmt.Sprintf(`CREATE DATABASE "%s"`, dbName),
	).Error
}

// RetryInit 重试初始化
func (in *Initializer) RetryInit(ctx context.Context, tenantID int64) error {
	tenant, err := in.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return errcode.Err(errcode.TenantNotFound)
	}
	if tenant.InitStatus != smodel.InitStatusFailed {
		return fmt.Errorf("仅在失败状态可重试，当前: %d", tenant.InitStatus)
	}
	return in.Initialize(ctx, tenantID)
}
