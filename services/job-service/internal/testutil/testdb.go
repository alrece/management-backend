package testutil

import (
	"context"
	"testing"

	"management-backend/services/job-service/internal/model"
	"management-backend/pkg/snowflake"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewTestDB 创建 SQLite 内存测试库
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	_ = snowflake.Init(1)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&model.JobTask{}, &model.JobExecutionLog{}); err != nil {
		t.Fatalf("迁移失败: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

// SeedTask 创建测试任务
func SeedTask(t *testing.T, db *gorm.DB, name, handler, cronExpr string, tenantID int64) *model.JobTask {
	t.Helper()
	creator := snowflake.NextID()
	task := &model.JobTask{
		ID:       snowflake.NextID(),
		TenantID: tenantID,
		Name:     name,
		Handler:  handler,
		CronExpr: cronExpr,
		Status:   model.TaskStatusNormal,
		Creator:  &creator,
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("创建测试任务失败: %v", err)
	}
	return task
}

// Ctx 返回测试 context
func Ctx() context.Context { return context.Background() }
