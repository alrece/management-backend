package handler

import (
	"context"

	pb "management-backend/api/proto/system"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MenuGRPCHandler 菜单 gRPC 服务端实现
type MenuGRPCHandler struct {
	pb.UnimplementedMenuServiceServer
	svc service.MenuService
}

func NewMenuGRPCHandler(svc service.MenuService) *MenuGRPCHandler {
	return &MenuGRPCHandler{svc: svc}
}

func (h *MenuGRPCHandler) GetMenuTree(ctx context.Context, _ *pb.Empty) (*pb.MenuTreeResp, error) {
	tree, err := h.svc.Tree(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询菜单树失败: %v", err)
	}
	return &pb.MenuTreeResp{List: convertMenuTree(tree)}, nil
}

func convertMenuTree(nodes []smodel.MenuTreeResp) []*pb.MenuNode {
	result := make([]*pb.MenuNode, 0, len(nodes))
	for _, n := range nodes {
		node := &pb.MenuNode{
			Id:         n.ID,
			ParentId:   n.ParentID,
			MenuName:   n.MenuName,
			MenuType:   int32(n.MenuType),
			Path:       n.Path,
			Component:  n.Component,
			Permission: n.Permission,
			Icon:       n.Icon,
			Sort:       int32(n.Sort),
			Visible:    int32(n.Visible),
			Status:     int32(n.Status),
		}
		if len(n.Children) > 0 {
			node.Children = convertMenuTree(n.Children)
		}
		result = append(result, node)
	}
	return result
}

func (h *MenuGRPCHandler) CreateMenu(ctx context.Context, req *pb.CreateMenuReq) (*pb.IDResp, error) {
	id, err := h.svc.Create(ctx, &smodel.MenuCreateReq{
		ParentID:   req.ParentId,
		MenuName:   req.MenuName,
		MenuType:   int(req.MenuType),
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Icon:       req.Icon,
		Sort:       int(req.Sort),
		Visible:    int(req.Visible),
		Status:     int(req.Status),
	}, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建菜单失败: %v", err)
	}
	return &pb.IDResp{Id: id}, nil
}

func (h *MenuGRPCHandler) UpdateMenu(ctx context.Context, req *pb.UpdateMenuReq) (*pb.Empty, error) {
	err := h.svc.Update(ctx, &smodel.MenuUpdateReq{
		ID:         req.Id,
		ParentID:   req.ParentId,
		MenuName:   req.MenuName,
		MenuType:   int(req.MenuType),
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Icon:       req.Icon,
		Sort:       int(req.Sort),
		Visible:    int(req.Visible),
		Status:     int(req.Status),
	}, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新菜单失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (h *MenuGRPCHandler) DeleteMenu(ctx context.Context, req *pb.IDReq) (*pb.Empty, error) {
	if err := h.svc.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "删除菜单失败: %v", err)
	}
	return &pb.Empty{}, nil
}
