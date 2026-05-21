package service

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"

	basemodel "management-backend/internal/model"
)

// PostService 岗位业务接口
type PostService interface {
	Create(ctx context.Context, req *smodel.PostCreateReq, creator int64) (int64, error)
	Update(ctx context.Context, req *smodel.PostUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	Page(ctx context.Context, req *smodel.PostPageReq) ([]smodel.PostResp, int64, error)
}

type postService struct {
	repo repository.PostRepo
}

// NewPostService 创建岗位 Service
func NewPostService(repo repository.PostRepo) PostService {
	return &postService{repo: repo}
}

func (s *postService) Create(ctx context.Context, req *smodel.PostCreateReq, creator int64) (int64, error) {
	existing, _ := s.repo.GetByCode(ctx, req.PostCode)
	if existing != nil {
		return 0, errcode.Err(errcode.PostExists)
	}

	post := &smodel.Post{
		PostCode:    req.PostCode,
		PostName:    req.PostName,
		RemarkModel: basemodel.RemarkModel{Remark: req.Remark},
	}
	post.ID = snowflake.NextID()
	post.Creator = creator
	post.Updater = creator
	post.SortModel.Sort = req.Sort
	post.StatusModel.Status = req.Status

	if err := s.repo.Create(ctx, post); err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (s *postService) Update(ctx context.Context, req *smodel.PostUpdateReq, updater int64) error {
	post, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.PostNotFound)
	}

	post.PostCode = req.PostCode
	post.PostName = req.PostName
	post.SortModel.Sort = req.Sort
	post.StatusModel.Status = req.Status
	post.RemarkModel.Remark = req.Remark
	post.Updater = updater

	return s.repo.Update(ctx, post)
}

func (s *postService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errcode.Err(errcode.PostNotFound)
	}
	return s.repo.Delete(ctx, id)
}

func (s *postService) Page(ctx context.Context, req *smodel.PostPageReq) ([]smodel.PostResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.PostResp, 0, len(list))
	for i := range list {
		resp = append(resp, *toPostResp(&list[i]))
	}
	return resp, total, nil
}

func toPostResp(p *smodel.Post) *smodel.PostResp {
	return &smodel.PostResp{
		ID: p.ID, PostCode: p.PostCode, PostName: p.PostName,
		Sort: p.Sort, Status: p.Status, Remark: p.Remark,
		CreateTime: p.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
