package service

import (
	"context"

	"management-backend/internal/model"
	smodel "management-backend/internal/module/system/model"
	"management-backend/pkg/snowflake"
)

// ClientService 客户端管理接口
type ClientService interface {
	Create(ctx context.Context, req *smodel.ClientCreateReq, creator int64) (int64, error)
	Update(ctx context.Context, req *smodel.ClientUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*smodel.ClientResp, error)
	Page(ctx context.Context, req *smodel.ClientPageReq) ([]smodel.ClientResp, int64, error)
}

type clientRepo interface {
	Create(ctx context.Context, client *smodel.Client) error
	Update(ctx context.Context, client *smodel.Client) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Client, error)
	Page(ctx context.Context, req *smodel.ClientPageReq) ([]smodel.Client, int64, error)
}

type clientSvc struct {
	repo clientRepo
}

func NewClientService(repo clientRepo) ClientService {
	return &clientSvc{repo: repo}
}

func (s *clientSvc) Create(ctx context.Context, req *smodel.ClientCreateReq, creator int64) (int64, error) {
	client := &smodel.Client{
		BaseModel:             model.BaseModel{ID: snowflake.NextID()},
		ClientID:              req.ClientID,
		ClientKey:             req.ClientKey,
		ClientName:            req.ClientName,
		GrantType:             req.GrantType,
		RedirectURI:           req.RedirectURI,
		AccessTokenValidity:   req.AccessTokenValidity,
		RefreshTokenValidity:  req.RefreshTokenValidity,
	}
	client.StatusModel.Status = req.Status
	client.RemarkModel.Remark = req.Remark
	client.Creator = creator
	if client.GrantType == "" {
		client.GrantType = "authorization_code"
	}
	if client.AccessTokenValidity == 0 {
		client.AccessTokenValidity = 1800
	}
	if client.RefreshTokenValidity == 0 {
		client.RefreshTokenValidity = 604800
	}
	if err := s.repo.Create(ctx, client); err != nil {
		return 0, err
	}
	return client.ID, nil
}

func (s *clientSvc) Update(ctx context.Context, req *smodel.ClientUpdateReq, updater int64) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	existing.ClientKey = req.ClientKey
	existing.ClientName = req.ClientName
	existing.GrantType = req.GrantType
	existing.RedirectURI = req.RedirectURI
	existing.AccessTokenValidity = req.AccessTokenValidity
	existing.RefreshTokenValidity = req.RefreshTokenValidity
	existing.StatusModel.Status = req.Status
	existing.RemarkModel.Remark = req.Remark
	existing.Updater = updater
	return s.repo.Update(ctx, existing)
}

func (s *clientSvc) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *clientSvc) Get(ctx context.Context, id int64) (*smodel.ClientResp, error) {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toClientResp(client), nil
}

func (s *clientSvc) Page(ctx context.Context, req *smodel.ClientPageReq) ([]smodel.ClientResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]smodel.ClientResp, 0, len(list))
	for i := range list {
		resps = append(resps, *toClientResp(&list[i]))
	}
	return resps, total, nil
}

func toClientResp(c *smodel.Client) *smodel.ClientResp {
	return &smodel.ClientResp{
		ID:                   c.ID,
		ClientID:             c.ClientID,
		ClientKey:            c.ClientKey,
		ClientName:           c.ClientName,
		GrantType:            c.GrantType,
		RedirectURI:          c.RedirectURI,
		AccessTokenValidity:  c.AccessTokenValidity,
		RefreshTokenValidity: c.RefreshTokenValidity,
		Status:               c.Status,
		Remark:               c.Remark,
		CreateTime:           c.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
