package repository

import (
	"context"

	pkgmiddleware "management-backend/pkg/middleware"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// NoticeRepo 通知公告数据访问接口
type NoticeRepo interface {
	Create(ctx context.Context, notice *smodel.Notice) error
	Update(ctx context.Context, notice *smodel.Notice) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Notice, error)
	Page(ctx context.Context, req *smodel.NoticePageReq) ([]smodel.Notice, int64, error)
}

type noticeRepo struct {
	db *gorm.DB
}

func NewNoticeRepo(db *gorm.DB) NoticeRepo {
	return &noticeRepo{db: db}
}

func (r *noticeRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *noticeRepo) Create(ctx context.Context, notice *smodel.Notice) error {
	return r.getDB(ctx).WithContext(ctx).Create(notice).Error
}

func (r *noticeRepo) Update(ctx context.Context, notice *smodel.Notice) error {
	return r.getDB(ctx).WithContext(ctx).Save(notice).Error
}

func (r *noticeRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.Notice{}, id).Error
}

func (r *noticeRepo) GetByID(ctx context.Context, id int64) (*smodel.Notice, error) {
	var notice smodel.Notice
	if err := r.getDB(ctx).WithContext(ctx).First(&notice, id).Error; err != nil {
		return nil, err
	}
	return &notice, nil
}

func (r *noticeRepo) Page(ctx context.Context, req *smodel.NoticePageReq) ([]smodel.Notice, int64, error) {
	var list []smodel.Notice
	var total int64

	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.Notice{})

	if req.NoticeTitle != "" {
		db = db.Where("notice_title LIKE ?", "%"+req.NoticeTitle+"%")
	}
	if req.NoticeType != nil {
		db = db.Where("notice_type = ?", *req.NoticeType)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
