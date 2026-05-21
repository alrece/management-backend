package multitenant

import (
	"context"
	"database/sql"
	"fmt"

	"management-backend/internal/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Migrator Schema 迁移工具（基于 goose）
type Migrator struct {
	migrationsDir string
}

// NewMigrator 创建迁移器
func NewMigrator() *Migrator {
	return &Migrator{
		migrationsDir: "scripts/sql/migrations",
	}
}

// MigrateDefault 对默认库执行迁移
func (m *Migrator) MigrateDefault(ctx context.Context) error {
	db, err := m.openDefault()
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.UpContext(ctx, db, m.migrationsDir)
}

// MigrateTenant 对指定租户库执行迁移
func (m *Migrator) MigrateTenant(ctx context.Context, tenantID int64) error {
	// 通过默认库获取租户数据库名
	defaultDB, err := gorm.Open(mysql.Open(config.C.MySQL.DSN()), &gorm.Config{})
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, _ := defaultDB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()

	var dbName string
	if err := defaultDB.WithContext(ctx).Table("sys_tenant").
		Where("id = ? AND deleted = 0", tenantID).
		Select("db_name").Scan(&dbName).Error; err != nil {
		return fmt.Errorf("租户不存在: %w", err)
	}

	db, err := m.openDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.UpContext(ctx, db, m.migrationsDir)
}

// MigrateAll 对所有已就绪租户执行迁移
func (m *Migrator) MigrateAll(ctx context.Context) error {
	defaultDB, err := gorm.Open(mysql.Open(config.C.MySQL.DSN()), &gorm.Config{})
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

func (m *Migrator) openDefault() (*sql.DB, error) {
	return m.openDB(config.C.MySQL.Database)
}

func (m *Migrator) openDB(dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.C.MySQL.Username, config.C.MySQL.Password,
		config.C.MySQL.Host, config.C.MySQL.Port,
		dbName, config.C.MySQL.Charset)
	return goose.OpenDBWithDriver("mysql", dsn)
}
