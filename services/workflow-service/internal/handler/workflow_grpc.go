package handler

import (
	"context"

	pb "management-backend/api/proto/workflow"
	"management-backend/services/workflow-service/internal/model"
	"management-backend/services/workflow-service/internal/repository"
	"management-backend/services/workflow-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WorkflowGRPCHandler 工作流 gRPC 服务端实现
type WorkflowGRPCHandler struct {
	pb.UnimplementedWorkflowServiceServer
	svc        service.WorkflowService
	instanceRp repository.InstanceRepo
}

func NewWorkflowGRPCHandler(svc service.WorkflowService, instanceRp repository.InstanceRepo) *WorkflowGRPCHandler {
	return &WorkflowGRPCHandler{svc: svc, instanceRp: instanceRp}
}

func (h *WorkflowGRPCHandler) CreateWorkflow(ctx context.Context, req *pb.CreateWorkflowReq) (*pb.CreateWorkflowResp, error) {
	if req.Name == "" || req.TenantId == 0 {
		return nil, status.Error(codes.InvalidArgument, "name 和 tenant_id 不能为空")
	}

	id, err := h.svc.Create(ctx, &model.WorkflowCreateReq{Name: req.Name, CategoryID: req.CategoryId}, req.Creator, req.TenantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建工作流失败: %v", err)
	}
	return &pb.CreateWorkflowResp{Id: id}, nil
}

func (h *WorkflowGRPCHandler) ActivateWorkflow(ctx context.Context, req *pb.WorkflowIDReq) (*pb.Empty, error) {
	if err := h.svc.Activate(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "激活工作流失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (h *WorkflowGRPCHandler) DeactivateWorkflow(ctx context.Context, req *pb.WorkflowIDReq) (*pb.Empty, error) {
	if err := h.svc.Deactivate(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "停用工作流失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (h *WorkflowGRPCHandler) ExecuteWorkflow(ctx context.Context, req *pb.ExecuteWorkflowReq) (*pb.ExecuteWorkflowResp, error) {
	if req.WorkflowId == 0 || req.TenantId == 0 {
		return nil, status.Error(codes.InvalidArgument, "workflow_id 和 tenant_id 不能为空")
	}

	instID, err := h.svc.Execute(ctx, req.WorkflowId, &model.ExecuteReq{Variables: req.Variables}, req.Creator, req.TenantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "执行工作流失败: %v", err)
	}
	return &pb.ExecuteWorkflowResp{InstanceId: instID}, nil
}

func (h *WorkflowGRPCHandler) GetInstance(ctx context.Context, req *pb.GetInstanceReq) (*pb.InstanceResp, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id 不能为空")
	}

	inst, err := h.instanceRp.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询实例失败: %v", err)
	}
	if inst.TenantID != req.TenantId {
		return nil, status.Error(codes.PermissionDenied, "无权访问该实例")
	}

	resp := &pb.InstanceResp{
		Id:          inst.ID,
		WorkflowId:  inst.WorkflowID,
		Status:      inst.Status,
		Variables:   inst.Variables,
		Result:      inst.Result,
		DurationMs:  inst.DurationMs,
		ErrorMsg:    inst.ErrorMsg,
		StartedAt:   inst.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	return resp, nil
}
