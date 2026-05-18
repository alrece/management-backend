package service

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户业务接口
type UserService interface {
	Create(ctx context.Context, req *smodel.UserCreateReq, creator int64, tenantID int64) (int64, error)
	Update(ctx context.Context, req *smodel.UserUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.UserResp, error)
	GetByUsername(ctx context.Context, username string) (*smodel.User, error)
	Page(ctx context.Context, req *smodel.UserPageReq, tenantScope func(*gorm.DB) *gorm.DB) ([]smodel.UserResp, int64, error)
}

type userService struct {
	repo repository.UserRepo
}

// NewUserService 创建用户 Service（构造函数注入 Repository 接口）
func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, req *smodel.UserCreateReq, creator int64, tenantID int64) (int64, error) {
	// 校验用户名唯一
	existing, _ := s.repo.GetByUsername(ctx, req.Username)
	if existing != nil {
		return 0, errcode.Err(errcode.UserExists)
	}

	// 密码加密
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &smodel.User{
		Username: req.Username,
		Password: string(hashed),
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
		DeptID:   req.DeptID,
	}
	user.ID = snowflake.NextID()
	user.TenantID = tenantID
	user.Creator = creator
	user.Updater = creator
	user.Status = req.Status
	user.Remark = req.Remark

	if err := s.repo.Create(ctx, user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *userService) Update(ctx context.Context, req *smodel.UserUpdateReq, updater int64) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.UserNotFound)
	}
	existing.Nickname = req.Nickname
	existing.Email = req.Email
	existing.Mobile = req.Mobile
	existing.Sex = req.Sex
	existing.Status = req.Status
	existing.DeptID = req.DeptID
	existing.Remark = req.Remark
	existing.Updater = updater
	return s.repo.Update(ctx, existing)
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errcode.Err(errcode.UserNotFound)
	}
	return s.repo.Delete(ctx, id)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*smodel.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *userService) GetByID(ctx context.Context, id int64) (*smodel.UserResp, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.Err(errcode.UserNotFound)
	}
	return toResp(user), nil
}

func (s *userService) Page(ctx context.Context, req *smodel.UserPageReq, tenantScope func(*gorm.DB) *gorm.DB) ([]smodel.UserResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req, tenantScope)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.UserResp, 0, len(list))
	for i := range list {
		resp = append(resp, *toResp(&list[i]))
	}
	return resp, total, nil
}

// toResp Entity → Resp 转换
func toResp(u *smodel.User) *smodel.UserResp {
	return &smodel.UserResp{
		ID:         u.ID,
		Username:   u.Username,
		Nickname:   u.Nickname,
		Email:      u.Email,
		Mobile:     u.Mobile,
		Sex:        u.Sex,
		Avatar:     u.Avatar,
		Status:     u.Status,
		DeptID:     u.DeptID,
		TenantID:   u.TenantID,
		CreateTime: u.CreateTime.Format("2006-01-02 15:04:05"),
		Remark:     u.Remark,
	}
}
