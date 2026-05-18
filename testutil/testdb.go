package testutil

import (
	"context"
	"sync"
	"testing"

	"management-backend/internal/model"
	smodel "management-backend/internal/module/system/model"
	"management-backend/pkg/snowflake"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var snowflakeOnce sync.Once

// initSnowflake 确保 snowflake 节点只初始化一次
func initSnowflake() {
	snowflakeOnce.Do(func() {
		_ = snowflake.Init(1)
	})
}

// NewTestDB 创建 SQLite 内存测试库并自动迁移
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	initSnowflake()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	err = db.AutoMigrate(
		&smodel.User{},
		&smodel.Role{},
		&smodel.RoleMenu{},
		&smodel.Tenant{},
	)
	if err != nil {
		t.Fatalf("迁移失败: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

// NewTestDBWithContext 创建带 context 的测试库
func NewTestDBWithContext(t *testing.T) (*gorm.DB, context.Context) {
	t.Helper()
	db := NewTestDB(t)
	ctx := context.Background()
	return db, ctx
}

// SeedUser 创建测试用户并返回
func SeedUser(t *testing.T, db *gorm.DB, username, password string) *smodel.User {
	t.Helper()
	user := &smodel.User{
		Username:    username,
		Password:    password,
		Nickname:    username,
		StatusModel: model.StatusModel{Status: 0},
	}
	user.ID = snowflake.NextID()
	user.TenantID = 1
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	return user
}

// SeedRole 创建测试角色
func SeedRole(t *testing.T, db *gorm.DB, name, code string) *smodel.Role {
	t.Helper()
	role := &smodel.Role{
		RoleName:    name,
		RoleCode:    code,
		StatusModel: model.StatusModel{Status: 0},
	}
	role.ID = snowflake.NextID()
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("创建测试角色失败: %v", err)
	}
	return role
}

// SeedTenant 创建测试租户
func SeedTenant(t *testing.T, db *gorm.DB, name string) *smodel.Tenant {
	t.Helper()
	tenant := &smodel.Tenant{
		TenantName:    name,
		DBName:        "mb_tenant_test",
		StatusModel:   model.StatusModel{Status: 0},
		InitStatus:    smodel.InitStatusReady,
	}
	tenant.ID = snowflake.NextID()
	if err := db.Create(tenant).Error; err != nil {
		t.Fatalf("创建测试租户失败: %v", err)
	}
	return tenant
}
