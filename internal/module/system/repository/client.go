package repository

import (
	"context"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// ClientRepo 客户端数据访问接口
type ClientRepo interface {
	Create(ctx context.Context, client *smodel.Client) error
	Update(ctx context.Context, client *smodel.Client) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Client, error)
	Page(ctx context.Context, req *smodel.ClientPageReq) ([]smodel.Client, int64, error)
}

type clientRepo struct {
	db *gorm.DB
}

func NewClientRepo(db *gorm.DB) ClientRepo {
	return &clientRepo{db: db}
}

func (r *clientRepo) Create(ctx context.Context, client *smodel.Client) error {
	return r.db.WithContext(ctx).Create(client).Error
}

func (r *clientRepo) Update(ctx context.Context, client *smodel.Client) error {
	return r.db.WithContext(ctx).Select("*").Omit("client_id").Updates(client).Error
}

func (r *clientRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&smodel.Client{}, id).Error
}

func (r *clientRepo) GetByID(ctx context.Context, id int64) (*smodel.Client, error) {
	var client smodel.Client
	if err := r.db.WithContext(ctx).First(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepo) Page(ctx context.Context, req *smodel.ClientPageReq) ([]smodel.Client, int64, error) {
	var list []smodel.Client
	var total int64
	db := r.db.WithContext(ctx).Model(&smodel.Client{})

	if req.ClientName != "" {
		db = db.Where("client_name LIKE ?", "%"+req.ClientName+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	db.Count(&total)
	db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list)
	return list, total, nil
}
