package service

import (
	"context"

	"management-backend/internal/authz"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"

	basemodel "management-backend/internal/model"
)

// RoleService 角色业务接口
type RoleService interface {
	Create(ctx context.Context, req *smodel.RoleCreateReq, creator int64) (int64, error)
	Update(ctx context.Context, req *smodel.RoleUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.RoleResp, error)
	Page(ctx context.Context, req *smodel.RolePageReq) ([]smodel.RoleResp, int64, error)
	AssignMenus(ctx context.Context, req *smodel.RoleMenuAssignReq, creator int64, tenantID int64) error
	GetRolesByUserID(ctx context.Context, userID int64) ([]smodel.Role, error)
}

type roleService struct {
	repo     repository.RoleRepo
	menuRepo repository.MenuRepo
	authzMgr *authz.EnforcerManager
}

// NewRoleService 创建角色 Service
func NewRoleService(repo repository.RoleRepo, menuRepo repository.MenuRepo, authzMgr *authz.EnforcerManager) RoleService {
	return &roleService{repo: repo, menuRepo: menuRepo, authzMgr: authzMgr}
}

func (s *roleService) Create(ctx context.Context, req *smodel.RoleCreateReq, creator int64) (int64, error) {
	existing, _ := s.repo.GetByRoleCode(ctx, req.RoleCode)
	if existing != nil {
		return 0, errcode.Err(errcode.RoleExists)
	}

	role := &smodel.Role{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		DataScope:   req.DataScope,
		RemarkModel: basemodel.RemarkModel{Remark: req.Remark},
		StatusModel: basemodel.StatusModel{Status: req.Status},
		SortModel:   basemodel.SortModel{Sort: req.Sort},
	}
	role.ID = snowflake.NextID()
	role.Creator = creator
	role.Updater = creator

	if err := s.repo.Create(ctx, role); err != nil {
		return 0, err
	}
	return role.ID, nil
}

func (s *roleService) Update(ctx context.Context, req *smodel.RoleUpdateReq, updater int64) error {
	role, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.RoleNotFound)
	}

	role.RoleName = req.RoleName
	role.RoleCode = req.RoleCode
	role.DataScope = req.DataScope
	role.DataScopeDeptIDs = req.DataScopeDeptIDs
	role.SortModel = basemodel.SortModel{Sort: req.Sort}
	role.StatusModel = basemodel.StatusModel{Status: req.Status}
	role.RemarkModel = basemodel.RemarkModel{Remark: req.Remark}
	role.Updater = updater

	return s.repo.Update(ctx, role)
}

func (s *roleService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errcode.Err(errcode.RoleNotFound)
	}
	return s.repo.Delete(ctx, id)
}

func (s *roleService) GetByID(ctx context.Context, id int64) (*smodel.RoleResp, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.Err(errcode.RoleNotFound)
	}
	return toRoleResp(role), nil
}

func (s *roleService) Page(ctx context.Context, req *smodel.RolePageReq) ([]smodel.RoleResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.RoleResp, 0, len(list))
	for i := range list {
		resp = append(resp, *toRoleResp(&list[i]))
	}
	return resp, total, nil
}

func (s *roleService) AssignMenus(ctx context.Context, req *smodel.RoleMenuAssignReq, creator int64, tenantID int64) error {
	role, err := s.repo.GetByID(ctx, req.RoleID)
	if err != nil {
		return errcode.Err(errcode.RoleNotFound)
	}

	if err := s.repo.AssignMenus(ctx, req.RoleID, req.MenuIDs, creator); err != nil {
		return err
	}

	// 同步 Casbin 策略
	if s.authzMgr != nil {
		menus, _ := s.menuRepo.GetByIDs(ctx, req.MenuIDs)
		var perms []string
		for _, m := range menus {
			if m.Permission != "" {
				perms = append(perms, m.Permission)
			}
		}
		_ = s.authzMgr.SyncRolePermissions(ctx, tenantID, role.RoleCode, perms)
	}

	return nil
}

func (s *roleService) GetRolesByUserID(ctx context.Context, userID int64) ([]smodel.Role, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func toRoleResp(r *smodel.Role) *smodel.RoleResp {
	return &smodel.RoleResp{
		ID:               r.ID,
		RoleName:         r.RoleName,
		RoleCode:         r.RoleCode,
		Sort:             r.Sort,
		DataScope:        r.DataScope,
		DataScopeDeptIDs: r.DataScopeDeptIDs,
		Status:           r.Status,
		Remark:           r.Remark,
		CreateTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
