package multitenant

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"management-backend/internal/config"
	"management-backend/internal/module/system/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Migrator Schema 迁移工具
type Migrator struct {
	tenantRepo    repository.TenantRepo
	migrationsDir string
}

// NewMigrator 创建迁移器
func NewMigrator(tenantRepo repository.TenantRepo) *Migrator {
	return &Migrator{
		tenantRepo:    tenantRepo,
		migrationsDir: "scripts/sql/migrations",
	}
}

// MigrateDefault 对默认库执行迁移
func (m *Migrator) MigrateDefault(ctx context.Context) error {
	cfg := config.C.MySQL
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接默认库失败: %w", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()
	return m.executeMigrations(ctx, db)
}

// MigrateTenant 对指定租户库执行迁移
func (m *Migrator) MigrateTenant(ctx context.Context, tenantID int64) error {
	tenant, err := m.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("租户不存在: %w", err)
	}

	db, err := m.connectTenantDB(tenant.DBName)
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()

	return m.executeMigrations(ctx, db)
}

// MigrateAll 对所有已就绪租户执行迁移
func (m *Migrator) MigrateAll(ctx context.Context) error {
	cfg := config.C.MySQL
	defaultDB, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, _ := defaultDB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()

	var tenantIDs []int64
	defaultDB.WithContext(ctx).Table("sys_tenant").
		Where("init_status = 2 AND deleted = 0").Pluck("id", &tenantIDs)

	var lastErr error
	for _, id := range tenantIDs {
		if err := m.MigrateTenant(ctx, id); err != nil {
			lastErr = fmt.Errorf("租户 %d 迁移失败: %w", id, err)
		}
	}
	return lastErr
}

func (m *Migrator) connectTenantDB(dbName string) (*gorm.DB, error) {
	cfg := config.C.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, dbName, cfg.Charset)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func (m *Migrator) executeMigrations(ctx context.Context, db *gorm.DB) error {
	entries, err := os.ReadDir(m.migrationsDir)
	if err != nil {
		return nil // 目录不存在则跳过
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		content, err := os.ReadFile(filepath.Join(m.migrationsDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("读取 %s 失败: %w", entry.Name(), err)
		}
		if err := db.WithContext(ctx).Exec(string(content)).Error; err != nil {
			return fmt.Errorf("执行 %s 失败: %w", entry.Name(), err)
		}
	}
	return nil
}
