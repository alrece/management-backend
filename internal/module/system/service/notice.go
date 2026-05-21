package service

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"
)

// NoticeService 通知公告业务接口
type NoticeService interface {
	Create(ctx context.Context, req *smodel.NoticeCreateReq, creator int64, tenantID int64) (int64, error)
	Update(ctx context.Context, req *smodel.NoticeUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.NoticeResp, error)
	Page(ctx context.Context, req *smodel.NoticePageReq) ([]smodel.NoticeResp, int64, error)
}

type noticeService struct {
	repo repository.NoticeRepo
}

func NewNoticeService(repo repository.NoticeRepo) NoticeService {
	return &noticeService{repo: repo}
}

func (s *noticeService) Create(ctx context.Context, req *smodel.NoticeCreateReq, creator int64, tenantID int64) (int64, error) {
	notice := &smodel.Notice{
		NoticeTitle: req.NoticeTitle,
		NoticeType:  req.NoticeType,
		Content:     req.Content,
	}
	notice.Status = req.Status
	notice.Remark = req.Remark
	notice.ID = snowflake.NextID()
	notice.TenantID = tenantID
	notice.Creator = creator
	notice.Updater = creator

	if err := s.repo.Create(ctx, notice); err != nil {
		return 0, err
	}
	return notice.ID, nil
}

func (s *noticeService) Update(ctx context.Context, req *smodel.NoticeUpdateReq, updater int64) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.NotFound)
	}
	existing.NoticeTitle = req.NoticeTitle
	existing.NoticeType = req.NoticeType
	existing.Content = req.Content
	existing.Status = req.Status
	existing.Remark = req.Remark
	existing.Updater = updater
	return s.repo.Update(ctx, existing)
}

func (s *noticeService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errcode.Err(errcode.NotFound)
	}
	return s.repo.Delete(ctx, id)
}

func (s *noticeService) GetByID(ctx context.Context, id int64) (*smodel.NoticeResp, error) {
	notice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.Err(errcode.NotFound)
	}
	return noticeToResp(notice), nil
}

func (s *noticeService) Page(ctx context.Context, req *smodel.NoticePageReq) ([]smodel.NoticeResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.NoticeResp, 0, len(list))
	for i := range list {
		resp = append(resp, *noticeToResp(&list[i]))
	}
	return resp, total, nil
}

func noticeToResp(n *smodel.Notice) *smodel.NoticeResp {
	return &smodel.NoticeResp{
		ID:          n.ID,
		NoticeTitle: n.NoticeTitle,
		NoticeType:  n.NoticeType,
		Content:     n.Content,
		Status:      n.Status,
		Remark:      n.Remark,
		Creator:     n.Creator,
		CreateTime:  n.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
