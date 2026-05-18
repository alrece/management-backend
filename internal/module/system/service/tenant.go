package service

import (
	"context"
	"fmt"
	"time"

	basemodel "management-backend/internal/model"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"

	"gorm.io/gorm"
)

// TenantService 租户业务接口
type TenantService interface {
	Create(ctx context.Context, req *smodel.TenantCreateReq, creator int64) (int64, error)
	Update(ctx context.Context, req *smodel.TenantUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.TenantResp, error)
	Page(ctx context.Context, req *smodel.TenantPageReq) ([]smodel.TenantResp, int64, error)
}

type tenantService struct {
	repo repository.TenantRepo
}

// NewTenantService 创建租户 Service
func NewTenantService(repo repository.TenantRepo) TenantService {
	return &tenantService{repo: repo}
}

func (s *tenantService) Create(ctx context.Context, req *smodel.TenantCreateReq, creator int64) (int64, error) {
	tenant := &smodel.Tenant{
		TenantName:    req.TenantName,
		ContactName:   req.ContactName,
		ContactMobile: req.ContactMobile,
		InitStatus:    smodel.InitStatusPending,
		StatusModel:   basemodel.StatusModel{Status: req.Status},
		RemarkModel:   basemodel.RemarkModel{Remark: req.Remark},
	}
	tenant.ID = snowflake.NextID()
	tenant.Creator = creator
	tenant.Updater = creator

	// 数据库名格式: mb_tenant_{id}
	tenant.DBName = fmt.Sprintf("mb_tenant_%d", tenant.ID)

	if req.ExpireTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.ExpireTime)
		if err == nil {
			tenant.ExpireTime = &t
		}
	}

	if err := s.repo.Create(ctx, tenant); err != nil {
		return 0, err
	}
	return tenant.ID, nil
}

func (s *tenantService) Update(ctx context.Context, req *smodel.TenantUpdateReq, updater int64) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.TenantNotFound)
	}
	if req.TenantName != "" {
		existing.TenantName = req.TenantName
	}
	existing.ContactName = req.ContactName
	existing.ContactMobile = req.ContactMobile
	existing.Status = req.Status
	existing.Remark = req.Remark
	existing.Updater = updater

	if req.ExpireTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.ExpireTime)
		if err == nil {
			existing.ExpireTime = &t
		}
	} else if req.ExpireTime == "" && existing.ExpireTime != nil {
		existing.ExpireTime = nil
	}

	return s.repo.Update(ctx, existing)
}

func (s *tenantService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errcode.Err(errcode.TenantNotFound)
	}
	return s.repo.Delete(ctx, id)
}

func (s *tenantService) GetByID(ctx context.Context, id int64) (*smodel.TenantResp, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.Err(errcode.TenantNotFound)
		}
		return nil, err
	}
	return tenantToResp(tenant), nil
}

func (s *tenantService) Page(ctx context.Context, req *smodel.TenantPageReq) ([]smodel.TenantResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.TenantResp, 0, len(list))
	for i := range list {
		resp = append(resp, *tenantToResp(&list[i]))
	}
	return resp, total, nil
}

func tenantToResp(t *smodel.Tenant) *smodel.TenantResp {
	resp := &smodel.TenantResp{
		ID:            t.ID,
		TenantName:    t.TenantName,
		ContactName:   t.ContactName,
		ContactMobile: t.ContactMobile,
		DBName:        t.DBName,
		Status:        t.Status,
		InitStatus:    t.InitStatus,
		Remark:        t.Remark,
		CreateTime:    t.CreateTime.Format("2006-01-02 15:04:05"),
	}
	if t.ExpireTime != nil {
		resp.ExpireTime = t.ExpireTime.Format("2006-01-02 15:04:05")
	}
	return resp
}
