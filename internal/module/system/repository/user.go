package repository

import (
	"context"

	pkgmiddleware "management-backend/pkg/middleware"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// UserRepo 用户数据访问接口
type UserRepo interface {
	Create(ctx context.Context, user *smodel.User) error
	Update(ctx context.Context, user *smodel.User) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.User, error)
	GetByUsername(ctx context.Context, username string) (*smodel.User, error)
	Page(ctx context.Context, req *smodel.UserPageReq) ([]smodel.User, int64, error)
}

type userRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户 Repository
func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

// getDB 优先使用租户 DB（从 context 获取）
func (r *userRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *userRepo) Create(ctx context.Context, user *smodel.User) error {
	return r.getDB(ctx).WithContext(ctx).Create(user).Error
}

func (r *userRepo) Update(ctx context.Context, user *smodel.User) error {
	return r.getDB(ctx).WithContext(ctx).Select("*").Omit("password").Updates(user).Error
}

func (r *userRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.User{}, id).Error
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*smodel.User, error) {
	var user smodel.User
	if err := r.getDB(ctx).WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*smodel.User, error) {
	var user smodel.User
	if err := r.getDB(ctx).WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Page(ctx context.Context, req *smodel.UserPageReq) ([]smodel.User, int64, error) {
	var list []smodel.User
	var total int64

	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.User{})

	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.Mobile != "" {
		db = db.Where("mobile LIKE ?", "%"+req.Mobile+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
