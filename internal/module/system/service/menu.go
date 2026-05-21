package service

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"
)

// MenuService 菜单业务接口
type MenuService interface {
	Create(ctx context.Context, req *smodel.MenuCreateReq, creator int64) (int64, error)
	Update(ctx context.Context, req *smodel.MenuUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	Tree(ctx context.Context) ([]smodel.MenuTreeResp, error)
	GetByIDs(ctx context.Context, ids []int64) ([]smodel.Menu, error)
}

type menuService struct {
	repo repository.MenuRepo
}

// NewMenuService 创建菜单 Service
func NewMenuService(repo repository.MenuRepo) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) Create(ctx context.Context, req *smodel.MenuCreateReq, creator int64) (int64, error) {
	menu := &smodel.Menu{
		ParentID:   req.ParentID,
		MenuName:   req.MenuName,
		MenuType:   req.MenuType,
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Icon:       req.Icon,
		Visible:    req.Visible,
	}
	menu.ID = snowflake.NextID()
	menu.Creator = creator
	menu.Updater = creator
	menu.SortModel.Sort = req.Sort
	menu.StatusModel.Status = req.Status

	if err := s.repo.Create(ctx, menu); err != nil {
		return 0, err
	}
	return menu.ID, nil
}

func (s *menuService) Update(ctx context.Context, req *smodel.MenuUpdateReq, updater int64) error {
	menu, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.MenuNotFound)
	}

	menu.ParentID = req.ParentID
	menu.MenuName = req.MenuName
	menu.MenuType = req.MenuType
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Permission = req.Permission
	menu.Icon = req.Icon
	menu.Visible = req.Visible
	menu.SortModel.Sort = req.Sort
	menu.StatusModel.Status = req.Status
	menu.Updater = updater

	return s.repo.Update(ctx, menu)
}

func (s *menuService) Delete(ctx context.Context, id int64) error {
	hasChildren, _ := s.repo.HasChildren(ctx, id)
	if hasChildren {
		return errcode.Err(errcode.MenuHasChildren)
	}
	return s.repo.Delete(ctx, id)
}

func (s *menuService) Tree(ctx context.Context) ([]smodel.MenuTreeResp, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(list, 0), nil
}

func (s *menuService) GetByIDs(ctx context.Context, ids []int64) ([]smodel.Menu, error) {
	return s.repo.GetByIDs(ctx, ids)
}

func buildMenuTree(menus []smodel.Menu, parentID int64) []smodel.MenuTreeResp {
	var tree []smodel.MenuTreeResp
	for _, m := range menus {
		if m.ParentID == parentID {
			node := smodel.MenuTreeResp{
				ID: m.ID, ParentID: m.ParentID, MenuName: m.MenuName,
				MenuType: m.MenuType, Path: m.Path, Component: m.Component,
				Permission: m.Permission, Icon: m.Icon, Sort: m.Sort,
				Visible: m.Visible, Status: m.Status,
				Children: buildMenuTree(menus, m.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}
