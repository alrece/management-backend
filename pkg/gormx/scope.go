package gormx

import (
	"fmt"

	"gorm.io/gorm"
)

// DataScopeConfig 数据权限配置
type DataScopeConfig struct {
	Scope   int     // 1=全部 2=自定义 3=本部门 4=本部门及以下 5=仅本人
	DeptIDs []int64 // 自定义/本部门及以下 的部门ID列表
	SelfID  int64   // 当前用户ID（仅本人）
	DeptID  int64   // 当前用户部门ID（本部门）
}

const (
	DataScopeAll       = 1 // 全部数据
	DataScopeCustom    = 2 // 自定义部门
	DataScopeDept      = 3 // 本部门
	DataScopeDeptBelow = 4 // 本部门及以下
	DataScopeSelf      = 5 // 仅本人
)

// DataPermissionScope 根据数据权限配置生成 GORM Scope
func DataPermissionScope(cfg DataScopeConfig, deptField string, userField string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch cfg.Scope {
		case DataScopeAll:
			return db
		case DataScopeCustom, DataScopeDeptBelow:
			if len(cfg.DeptIDs) > 0 {
				return db.Where(fmt.Sprintf("%s IN ?", deptField), cfg.DeptIDs)
			}
			return db.Where("1 = 0")
		case DataScopeDept:
			return db.Where(fmt.Sprintf("%s = ?", deptField), cfg.DeptID)
		case DataScopeSelf:
			return db.Where(fmt.Sprintf("%s = ?", userField), cfg.SelfID)
		default:
			return db.Where("1 = 0")
		}
	}
}
