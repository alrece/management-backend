package repository

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	pkgmiddleware "management-backend/pkg/middleware"
	"management-backend/pkg/snowflake"

	"gorm.io/gorm"
)

// RoleRepo 角色数据访问接口
type RoleRepo interface {
	Create(ctx context.Context, role *smodel.Role) error
	Update(ctx context.Context, role *smodel.Role) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Role, error)
	Page(ctx context.Context, req *smodel.RolePageReq) ([]smodel.Role, int64, error)
	GetByUserID(ctx context.Context, userID int64) ([]smodel.Role, error)
	GetByRoleCode(ctx context.Context, roleCode string) (*smodel.Role, error)
	AssignMenus(ctx context.Context, roleID int64, menuIDs []int64, creator int64) error
	GetMenuIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error)
	GetMenuIDsByRoleIDs(ctx context.Context, roleIDs []int64) ([]int64, error)
}

type roleRepo struct {
	db *gorm.DB
}

// NewRoleRepo 创建角色 Repository
func NewRoleRepo(db *gorm.DB) RoleRepo {
	return &roleRepo{db: db}
}

func (r *roleRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *roleRepo) Create(ctx context.Context, role *smodel.Role) error {
	return r.getDB(ctx).WithContext(ctx).Create(role).Error
}

func (r *roleRepo) Update(ctx context.Context, role *smodel.Role) error {
	return r.getDB(ctx).WithContext(ctx).Save(role).Error
}

func (r *roleRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.Role{}, id).Error
}

func (r *roleRepo) GetByID(ctx context.Context, id int64) (*smodel.Role, error) {
	var role smodel.Role
	err := r.getDB(ctx).WithContext(ctx).First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) Page(ctx context.Context, req *smodel.RolePageReq) ([]smodel.Role, int64, error) {
	var list []smodel.Role
	var total int64
	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.Role{})

	if req.RoleName != "" {
		db = db.Where("role_name LIKE ?", "%"+req.RoleName+"%")
	}
	if req.RoleCode != "" {
		db = db.Where("role_code LIKE ?", "%"+req.RoleCode+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	db.Count(&total)
	db.Offset(req.Offset()).Limit(req.Limit()).Order("sort ASC, id DESC").Find(&list)
	return list, total, nil
}

func (r *roleRepo) GetByUserID(ctx context.Context, userID int64) ([]smodel.Role, error) {
	var roles []smodel.Role
	err := r.getDB(ctx).WithContext(ctx).
		Table("sys_role").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *roleRepo) GetByRoleCode(ctx context.Context, roleCode string) (*smodel.Role, error) {
	var role smodel.Role
	err := r.getDB(ctx).WithContext(ctx).Where("role_code = ?", roleCode).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) AssignMenus(ctx context.Context, roleID int64, menuIDs []int64, creator int64) error {
	db := r.getDB(ctx)
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("role_id = ?", roleID).Delete(&smodel.RoleMenu{})

		for _, menuID := range menuIDs {
			rm := &smodel.RoleMenu{
				ID:      snowflake.NextID(),
				RoleID:  roleID,
				MenuID:  menuID,
				Creator: creator,
				Updater: creator,
			}
			if err := tx.Create(rm).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepo) GetMenuIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var ids []int64
	err := r.getDB(ctx).WithContext(ctx).
		Model(&smodel.RoleMenu{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &ids).Error
	return ids, err
}

func (r *roleRepo) GetMenuIDsByRoleIDs(ctx context.Context, roleIDs []int64) ([]int64, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var ids []int64
	err := r.getDB(ctx).WithContext(ctx).
		Model(&smodel.RoleMenu{}).
		Where("role_id IN ?", roleIDs).
		Pluck("menu_id", &ids).Error
	return ids, err
}
