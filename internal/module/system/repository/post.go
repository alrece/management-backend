package repository

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	pkgmiddleware "management-backend/pkg/middleware"

	"gorm.io/gorm"
)

// PostRepo 岗位数据访问接口
type PostRepo interface {
	Create(ctx context.Context, post *smodel.Post) error
	Update(ctx context.Context, post *smodel.Post) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Post, error)
	GetByCode(ctx context.Context, code string) (*smodel.Post, error)
	Page(ctx context.Context, req *smodel.PostPageReq) ([]smodel.Post, int64, error)
}

type postRepo struct {
	db *gorm.DB
}

// NewPostRepo 创建岗位 Repository
func NewPostRepo(db *gorm.DB) PostRepo {
	return &postRepo{db: db}
}

func (r *postRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *postRepo) Create(ctx context.Context, post *smodel.Post) error {
	return r.getDB(ctx).WithContext(ctx).Create(post).Error
}

func (r *postRepo) Update(ctx context.Context, post *smodel.Post) error {
	return r.getDB(ctx).WithContext(ctx).Save(post).Error
}

func (r *postRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.Post{}, id).Error
}

func (r *postRepo) GetByID(ctx context.Context, id int64) (*smodel.Post, error) {
	var post smodel.Post
	err := r.getDB(ctx).WithContext(ctx).First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) GetByCode(ctx context.Context, code string) (*smodel.Post, error) {
	var post smodel.Post
	err := r.getDB(ctx).WithContext(ctx).Where("post_code = ?", code).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) Page(ctx context.Context, req *smodel.PostPageReq) ([]smodel.Post, int64, error) {
	var list []smodel.Post
	var total int64
	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.Post{})

	if req.PostCode != "" {
		db = db.Where("post_code LIKE ?", "%"+req.PostCode+"%")
	}
	if req.PostName != "" {
		db = db.Where("post_name LIKE ?", "%"+req.PostName+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	db.Count(&total)
	db.Offset(req.Offset()).Limit(req.Limit()).Order("sort ASC, id DESC").Find(&list)
	return list, total, nil
}
