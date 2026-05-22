package handler

import (
	"context"

	pb "management-backend/api/proto/job"
	"management-backend/services/job-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// JobGRPCHandler 任务 gRPC 服务端实现
type JobGRPCHandler struct {
	pb.UnimplementedJobServiceServer
	svc service.JobService
}

func NewJobGRPCHandler(svc service.JobService) *JobGRPCHandler {
	return &JobGRPCHandler{svc: svc}
}

func (h *JobGRPCHandler) TriggerJob(ctx context.Context, req *pb.TriggerJobRequest) (*pb.TriggerJobResponse, error) {
	if req.TaskId == 0 || req.TenantId == 0 {
		return nil, status.Error(codes.InvalidArgument, "task_id 和 tenant_id 不能为空")
	}

	execID, err := h.svc.Trigger(ctx, req.TaskId, req.TenantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "触发任务失败: %v", err)
	}
	return &pb.TriggerJobResponse{ExecutionId: execID}, nil
}

func (h *JobGRPCHandler) GetJobStatus(ctx context.Context, req *pb.GetJobStatusRequest) (*pb.GetJobStatusResponse, error) {
	if req.TaskId == 0 {
		return nil, status.Error(codes.InvalidArgument, "task_id 不能为空")
	}

	task, err := h.svc.GetByID(ctx, req.TaskId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询任务失败: %v", err)
	}
	return &pb.GetJobStatusResponse{
		TaskId:   task.ID,
		Name:     task.Name,
		Status:   int32(task.Status),
		CronExpr: task.CronExpr,
	}, nil
}
