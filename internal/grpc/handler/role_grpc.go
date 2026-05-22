package handler

import (
	"context"

	pb "management-backend/api/proto/system"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RoleGRPCHandler 角色 gRPC 服务端实现
type RoleGRPCHandler struct {
	pb.UnimplementedRoleServiceServer
	svc service.RoleService
}

func NewRoleGRPCHandler(svc service.RoleService) *RoleGRPCHandler {
	return &RoleGRPCHandler{svc: svc}
}

func (h *RoleGRPCHandler) GetRole(ctx context.Context, req *pb.IDReq) (*pb.RoleResp, error) {
	role, err := h.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询角色失败: %v", err)
	}
	if role == nil {
		return nil, status.Error(codes.NotFound, "角色不存在")
	}
	return &pb.RoleResp{
		Id:         role.ID,
		RoleName:   role.RoleName,
		RoleCode:   role.RoleCode,
		DataScope:  int32(role.DataScope),
		Sort:       int32(role.Sort),
		Status:     int32(role.Status),
		CreateTime: role.CreateTime,
	}, nil
}

func (h *RoleGRPCHandler) ListRoles(ctx context.Context, req *pb.ListRolesReq) (*pb.RoleListResp, error) {
	list, total, err := h.svc.Page(ctx, &smodel.RolePageReq{
		RoleName: req.RoleName,
		RoleCode: req.RoleCode,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询角色列表失败: %v", err)
	}

	resp := &pb.RoleListResp{Total: total}
	for _, r := range list {
		resp.List = append(resp.List, &pb.RoleResp{
			Id:         r.ID,
			RoleName:   r.RoleName,
			RoleCode:   r.RoleCode,
			DataScope:  int32(r.DataScope),
			Sort:       int32(r.Sort),
			Status:     int32(r.Status),
			CreateTime: r.CreateTime,
		})
	}
	return resp, nil
}

func (h *RoleGRPCHandler) CreateRole(ctx context.Context, req *pb.CreateRoleReq) (*pb.IDResp, error) {
	id, err := h.svc.Create(ctx, &smodel.RoleCreateReq{
		RoleName:  req.RoleName,
		RoleCode:  req.RoleCode,
		DataScope: int(req.DataScope),
		Sort:      int(req.Sort),
		Status:    int(req.Status),
		Remark:    req.Remark,
	}, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建角色失败: %v", err)
	}
	return &pb.IDResp{Id: id}, nil
}

func (h *RoleGRPCHandler) UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (*pb.Empty, error) {
	err := h.svc.Update(ctx, &smodel.RoleUpdateReq{
		ID:        req.Id,
		RoleName:  req.RoleName,
		RoleCode:  req.RoleCode,
		DataScope: int(req.DataScope),
		Sort:      int(req.Sort),
		Status:    int(req.Status),
		Remark:    req.Remark,
	}, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新角色失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (h *RoleGRPCHandler) DeleteRole(ctx context.Context, req *pb.IDReq) (*pb.Empty, error) {
	if err := h.svc.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "删除角色失败: %v", err)
	}
	return &pb.Empty{}, nil
}
