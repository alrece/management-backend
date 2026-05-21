package testutil

import (
	"context"
	"testing"

	"management-backend/services/workflow-service/internal/model"
	"management-backend/pkg/snowflake"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	_ = snowflake.Init(3)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&model.WfCategory{}, &model.WfWorkflow{}, &model.WfInstance{}); err != nil {
		t.Fatalf("迁移失败: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

func SeedCategory(t *testing.T, db *gorm.DB, name string, tenantID int64) *model.WfCategory {
	t.Helper()
	c := &model.WfCategory{ID: snowflake.NextID(), TenantID: tenantID, Name: name, Status: 0}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("创建测试分类失败: %v", err)
	}
	return c
}

func SeedWorkflow(t *testing.T, db *gorm.DB, name, n8nID string, tenantID int64) *model.WfWorkflow {
	t.Helper()
	wf := &model.WfWorkflow{
		ID: snowflake.NextID(), TenantID: tenantID, Name: name,
		N8nWorkflowID: n8nID, Status: "ACTIVE",
	}
	if err := db.Create(wf).Error; err != nil {
		t.Fatalf("创建测试工作流失败: %v", err)
	}
	return wf
}

func Ctx() context.Context { return context.Background() }
