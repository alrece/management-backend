package service

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
)

// LogService 日志业务接口
type LogService interface {
	CreateOperLog(ctx context.Context, log *smodel.SysOperLog) error
	PageOperLog(ctx context.Context, req *smodel.OperLogPageReq) ([]smodel.OperLogResp, int64, error)
	DeleteOperLog(ctx context.Context, id int64) error
	CleanOperLog(ctx context.Context) error

	CreateLoginLog(ctx context.Context, log *smodel.SysLoginLog) error
	PageLoginLog(ctx context.Context, req *smodel.LoginLogPageReq) ([]smodel.LoginLogResp, int64, error)
	DeleteLoginLog(ctx context.Context, id int64) error
	CleanLoginLog(ctx context.Context) error
}

type logService struct {
	operRepo  repository.OperLogRepo
	loginRepo repository.LoginLogRepo
}

func NewLogService(operRepo repository.OperLogRepo, loginRepo repository.LoginLogRepo) LogService {
	return &logService{operRepo: operRepo, loginRepo: loginRepo}
}

func (s *logService) CreateOperLog(ctx context.Context, log *smodel.SysOperLog) error {
	return s.operRepo.Create(ctx, log)
}

func (s *logService) PageOperLog(ctx context.Context, req *smodel.OperLogPageReq) ([]smodel.OperLogResp, int64, error) {
	list, total, err := s.operRepo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.OperLogResp, 0, len(list))
	for i := range list {
		resp = append(resp, operLogToResp(&list[i]))
	}
	return resp, total, nil
}

func (s *logService) DeleteOperLog(ctx context.Context, id int64) error {
	if err := s.operRepo.DeleteByID(ctx, id); err != nil {
		return errcode.Err(errcode.NotFound)
	}
	return nil
}

func (s *logService) CleanOperLog(ctx context.Context) error {
	return s.operRepo.Clean(ctx)
}

func (s *logService) CreateLoginLog(ctx context.Context, log *smodel.SysLoginLog) error {
	return s.loginRepo.Create(ctx, log)
}

func (s *logService) PageLoginLog(ctx context.Context, req *smodel.LoginLogPageReq) ([]smodel.LoginLogResp, int64, error) {
	list, total, err := s.loginRepo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]smodel.LoginLogResp, 0, len(list))
	for i := range list {
		resp = append(resp, loginLogToResp(&list[i]))
	}
	return resp, total, nil
}

func (s *logService) DeleteLoginLog(ctx context.Context, id int64) error {
	if err := s.loginRepo.DeleteByID(ctx, id); err != nil {
		return errcode.Err(errcode.NotFound)
	}
	return nil
}

func (s *logService) CleanLoginLog(ctx context.Context) error {
	return s.loginRepo.Clean(ctx)
}

func operLogToResp(l *smodel.SysOperLog) smodel.OperLogResp {
	return smodel.OperLogResp{
		ID:           l.ID,
		TenantID:     l.TenantID,
		Title:        l.Title,
		BusinessType: l.BusinessType,
		Method:       l.Method,
		RequestURL:   l.RequestURL,
		OperIP:       l.OperIP,
		OperUserID:   l.OperUserID,
		OperName:     l.OperName,
		RequestID:    l.RequestID,
		Status:       l.Status,
		ErrorMsg:     l.ErrorMsg,
		OperTime:     l.OperTime.Format("2006-01-02 15:04:05"),
	}
}

func loginLogToResp(l *smodel.SysLoginLog) smodel.LoginLogResp {
	return smodel.LoginLogResp{
		ID:            l.ID,
		Username:      l.Username,
		LoginIP:       l.LoginIP,
		LoginLocation: l.LoginLocation,
		Browser:       l.Browser,
		OS:            l.OS,
		Status:        l.Status,
		Msg:           l.Msg,
		LoginTime:     l.LoginTime.Format("2006-01-02 15:04:05"),
	}
}
