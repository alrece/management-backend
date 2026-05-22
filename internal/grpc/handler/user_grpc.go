package handler

import (
	"context"

	pb "management-backend/api/proto/system"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserGRPCHandler 用户 gRPC 服务端实现
type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserGRPCHandler(svc service.UserService) *UserGRPCHandler {
	return &UserGRPCHandler{svc: svc}
}

func (h *UserGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.UserResp, error) {
	user, err := h.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询用户失败: %v", err)
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return &pb.UserResp{
		Id:         user.ID,
		Username:   user.Username,
		Nickname:   user.Nickname,
		Email:      user.Email,
		Mobile:     user.Mobile,
		Sex:        int32(user.Sex),
		Status:     int32(user.Status),
		DeptId:     user.DeptID,
		CreateTime: user.CreateTime,
	}, nil
}

func (h *UserGRPCHandler) ListUsers(ctx context.Context, req *pb.ListUsersReq) (*pb.UserListResp, error) {
	st := int(0)
	statusPtr := &st
	if req.Status > 0 {
		*statusPtr = int(req.Status)
	} else {
		statusPtr = nil
	}

	list, total, err := h.svc.Page(ctx, &smodel.UserPageReq{
		Username: req.Username,
		Status:   statusPtr,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询用户列表失败: %v", err)
	}

	resp := &pb.UserListResp{Total: total}
	for _, u := range list {
		resp.List = append(resp.List, &pb.UserResp{
			Id:         u.ID,
			Username:   u.Username,
			Nickname:   u.Nickname,
			Email:      u.Email,
			Mobile:     u.Mobile,
			Sex:        int32(u.Sex),
			Status:     int32(u.Status),
			DeptId:     u.DeptID,
			CreateTime: u.CreateTime,
		})
	}
	return resp, nil
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.UserIDResp, error) {
	id, err := h.svc.Create(ctx, &smodel.UserCreateReq{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Sex:      int(req.Sex),
		DeptID:   req.DeptId,
		Status:   int(req.Status),
		Remark:   req.Remark,
	}, 0, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建用户失败: %v", err)
	}
	return &pb.UserIDResp{Id: id}, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.Empty, error) {
	err := h.svc.Update(ctx, &smodel.UserUpdateReq{
		ID:       req.Id,
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Sex:      int(req.Sex),
		DeptID:   req.DeptId,
		Status:   int(req.Status),
		Remark:   req.Remark,
	}, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新用户失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (h *UserGRPCHandler) DeleteUser(ctx context.Context, req *pb.IDReq) (*pb.Empty, error) {
	if err := h.svc.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "删除用户失败: %v", err)
	}
	return &pb.Empty{}, nil
}
