package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"

	"github.com/redis/go-redis/v9"
)

// DictService 字典业务接口
type DictService interface {
	CreateType(ctx context.Context, req *smodel.DictTypeCreateReq, creator int64, tenantID int64) (int64, error)
	UpdateType(ctx context.Context, req *smodel.DictTypeUpdateReq, updater int64) error
	DeleteType(ctx context.Context, id int64) error
	GetTypeByID(ctx context.Context, id int64) (*smodel.DictTypeResp, error)
	PageType(ctx context.Context, req *smodel.DictTypePageReq) ([]smodel.DictTypeResp, int64, error)

	CreateData(ctx context.Context, req *smodel.DictDataCreateReq, creator int64, tenantID int64) (int64, error)
	UpdateData(ctx context.Context, req *smodel.DictDataUpdateReq, updater int64) error
	DeleteData(ctx context.Context, id int64) error
	PageData(ctx context.Context, req *smodel.DictDataPageReq) ([]smodel.DictDataResp, int64, error)
	ListDataByType(ctx context.Context, dictType string) ([]smodel.DictDataResp, error)
}

type dictService struct {
	typeRepo repository.DictTypeRepo
	dataRepo repository.DictDataRepo
	rdb      *redis.Client
}

func NewDictService(typeRepo repository.DictTypeRepo, dataRepo repository.DictDataRepo, rdb *redis.Client) DictService {
	return &dictService{typeRepo: typeRepo, dataRepo: dataRepo, rdb: rdb}
}

func (s *dictService) CreateType(ctx context.Context, req *smodel.DictTypeCreateReq, creator int64, tenantID int64) (int64, error) {
	existing, _ := s.typeRepo.GetByDictType(ctx, req.DictType)
	if existing != nil {
		return 0, errcode.Err(errcode.DictNotFound)
	}

	dict := &smodel.DictType{
		DictName: req.DictName,
		DictType: req.DictType,
	}
	dict.Status = req.Status
	dict.Remark = req.Remark
	dict.ID = snowflake.NextID()
	dict.TenantID = tenantID
	dict.Creator = creator
	dict.Updater = creator

	if err := s.typeRepo.Create(ctx, dict); err != nil {
		return 0, err
	}
	return dict.ID, nil
}

func (s *dictService) UpdateType(ctx context.Context, req *smodel.DictTypeUpdateReq, updater int64) error {
	existing, err := s.typeRepo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.DictNotFound)
	}
	existing.DictName = req.DictName
	existing.DictType = req.DictType
	existing.Status = req.Status
	existing.Remark = req.Remark
	existing.Updater = updater
	return s.typeRepo.Update(ctx, existing)
}

func (s *dictService) DeleteType(ctx context.Context, id int64) error {
	if _, err := s.typeRepo.GetByID(ctx, id); err != nil {
		return errcode.Err(errcode.DictNotFound)
	}
	return s.typeRepo.Delete(ctx, id)
}

func (s *dictService) GetTypeByID(ctx context.Context, id int64) (*smodel.DictTypeResp, error) {
	dict, err := s.typeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.Err(errcode.DictNotFound)
	}
	return dictTypeToResp(dict), nil
}

func (s *dictService) PageType(ctx context.Context, req *smodel.DictTypePageReq) ([]smodel.DictTypeResp, int64, error) {
	list, total, err := s.typeRepo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.DictTypeResp, 0, len(list))
	for i := range list {
		resp = append(resp, *dictTypeToResp(&list[i]))
	}
	return resp, total, nil
}

func (s *dictService) CreateData(ctx context.Context, req *smodel.DictDataCreateReq, creator int64, tenantID int64) (int64, error) {
	data := &smodel.DictData{
		DictType: req.DictType,
		Label:    req.Label,
		Value:    req.Value,
	}
	data.Sort = req.Sort
	data.Status = req.Status
	data.Remark = req.Remark
	data.ID = snowflake.NextID()
	data.TenantID = tenantID
	data.Creator = creator
	data.Updater = creator

	if err := s.dataRepo.Create(ctx, data); err != nil {
		return 0, err
	}
	s.invalidateCache(ctx, req.DictType)
	return data.ID, nil
}

func (s *dictService) UpdateData(ctx context.Context, req *smodel.DictDataUpdateReq, updater int64) error {
	existing, err := s.dataRepo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.DictNotFound)
	}
	oldType := existing.DictType
	existing.DictType = req.DictType
	existing.Label = req.Label
	existing.Value = req.Value
	existing.Sort = req.Sort
	existing.Status = req.Status
	existing.Remark = req.Remark
	existing.Updater = updater

	if err := s.dataRepo.Update(ctx, existing); err != nil {
		return err
	}
	s.invalidateCache(ctx, oldType)
	if oldType != req.DictType {
		s.invalidateCache(ctx, req.DictType)
	}
	return nil
}

func (s *dictService) DeleteData(ctx context.Context, id int64) error {
	existing, err := s.dataRepo.GetByID(ctx, id)
	if err != nil {
		return errcode.Err(errcode.DictNotFound)
	}
	if err := s.dataRepo.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, existing.DictType)
	return nil
}

func (s *dictService) PageData(ctx context.Context, req *smodel.DictDataPageReq) ([]smodel.DictDataResp, int64, error) {
	list, total, err := s.dataRepo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.DictDataResp, 0, len(list))
	for i := range list {
		resp = append(resp, *dictDataToResp(&list[i]))
	}
	return resp, total, nil
}

func (s *dictService) ListDataByType(ctx context.Context, dictType string) ([]smodel.DictDataResp, error) {
	cacheKey := fmt.Sprintf("dict:%s", dictType)

	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var cached []smodel.DictDataResp
		if json.Unmarshal([]byte(val), &cached) == nil {
			return cached, nil
		}
	}

	list, err := s.dataRepo.ListByDictType(ctx, dictType)
	if err != nil {
		return nil, err
	}

	resp := make([]smodel.DictDataResp, 0, len(list))
	for i := range list {
		resp = append(resp, *dictDataToResp(&list[i]))
	}

	if data, e := json.Marshal(resp); e == nil {
		_ = s.rdb.Set(ctx, cacheKey, data, time.Hour).Err()
	}
	return resp, nil
}

func (s *dictService) invalidateCache(ctx context.Context, dictType string) {
	_ = s.rdb.Del(ctx, fmt.Sprintf("dict:%s", dictType)).Err()
}

func dictTypeToResp(d *smodel.DictType) *smodel.DictTypeResp {
	return &smodel.DictTypeResp{
		ID:         d.ID,
		DictName:   d.DictName,
		DictType:   d.DictType,
		Status:     d.Status,
		Remark:     d.Remark,
		CreateTime: d.CreateTime.Format("2006-01-02 15:04:05"),
	}
}

func dictDataToResp(d *smodel.DictData) *smodel.DictDataResp {
	return &smodel.DictDataResp{
		ID:         d.ID,
		DictType:   d.DictType,
		Label:      d.Label,
		Value:      d.Value,
		Sort:       d.Sort,
		Status:     d.Status,
		Remark:     d.Remark,
		CreateTime: d.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
