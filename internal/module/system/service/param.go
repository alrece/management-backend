package service

import (
	"context"
	"fmt"
	"time"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"

	"github.com/redis/go-redis/v9"
)

// ParamService 参数业务接口
type ParamService interface {
	Create(ctx context.Context, req *smodel.ParamCreateReq, creator int64, tenantID int64) (int64, error)
	Update(ctx context.Context, req *smodel.ParamUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.ParamResp, error)
	GetByKey(ctx context.Context, key string) (string, error)
	Page(ctx context.Context, req *smodel.ParamPageReq) ([]smodel.ParamResp, int64, error)
}

type paramService struct {
	repo repository.ParamRepo
	rdb  *redis.Client
}

func NewParamService(repo repository.ParamRepo, rdb *redis.Client) ParamService {
	return &paramService{repo: repo, rdb: rdb}
}

func (s *paramService) Create(ctx context.Context, req *smodel.ParamCreateReq, creator int64, tenantID int64) (int64, error) {
	existing, _ := s.repo.GetByKey(ctx, req.ParamKey)
	if existing != nil {
		return 0, errcode.Err(errcode.ConfigKeyExists)
	}

	param := &smodel.SysParam{
		ParamKey:   req.ParamKey,
		ParamValue: req.ParamValue,
		ParamType:  req.ParamType,
	}
	param.Status = req.Status
	param.Remark = req.Remark
	param.ID = snowflake.NextID()
	param.TenantID = tenantID
	param.Creator = creator
	param.Updater = creator
	if param.ParamType == 0 {
		param.ParamType = 1
	}

	if err := s.repo.Create(ctx, param); err != nil {
		return 0, err
	}
	return param.ID, nil
}

func (s *paramService) Update(ctx context.Context, req *smodel.ParamUpdateReq, updater int64) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.ParamNotFound)
	}
	oldKey := existing.ParamKey
	existing.ParamKey = req.ParamKey
	existing.ParamValue = req.ParamValue
	existing.ParamType = req.ParamType
	existing.Status = req.Status
	existing.Remark = req.Remark
	existing.Updater = updater

	if err := s.repo.Update(ctx, existing); err != nil {
		return err
	}
	s.invalidateCache(ctx, oldKey)
	if oldKey != req.ParamKey {
		s.invalidateCache(ctx, req.ParamKey)
	}
	return nil
}

func (s *paramService) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errcode.Err(errcode.ParamNotFound)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, existing.ParamKey)
	return nil
}

func (s *paramService) GetByID(ctx context.Context, id int64) (*smodel.ParamResp, error) {
	param, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.Err(errcode.ParamNotFound)
	}
	return paramToResp(param), nil
}

func (s *paramService) GetByKey(ctx context.Context, key string) (string, error) {
	cacheKey := fmt.Sprintf("param:%s", key)
	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		return val, nil
	}

	param, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return "", errcode.Err(errcode.ParamNotFound)
	}

	_ = s.rdb.Set(ctx, cacheKey, param.ParamValue, time.Hour).Err()
	return param.ParamValue, nil
}

func (s *paramService) Page(ctx context.Context, req *smodel.ParamPageReq) ([]smodel.ParamResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.ParamResp, 0, len(list))
	for i := range list {
		resp = append(resp, *paramToResp(&list[i]))
	}
	return resp, total, nil
}

func (s *paramService) invalidateCache(ctx context.Context, key string) {
	_ = s.rdb.Del(ctx, fmt.Sprintf("param:%s", key)).Err()
}

func paramToResp(p *smodel.SysParam) *smodel.ParamResp {
	return &smodel.ParamResp{
		ID:         p.ID,
		ParamKey:   p.ParamKey,
		ParamValue: p.ParamValue,
		ParamType:  p.ParamType,
		Status:     p.Status,
		Remark:     p.Remark,
		CreateTime: p.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
