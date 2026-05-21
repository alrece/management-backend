package service

import (
	"context"
	"fmt"

	smodel "management-backend/internal/module/system/model"

	"github.com/xuri/excelize/v2"
)

const maxExportRows = 100000

// UserExportService 用户导出接口
type UserExportService interface {
	Export(ctx context.Context, req *smodel.UserPageReq) (*excelize.File, error)
}

type userExportService struct {
	userSvc UserService
}

// NewUserExportService 创建导出 Service
func NewUserExportService(userSvc UserService) UserExportService {
	return &userExportService{userSvc: userSvc}
}

// Export 导出用户列表为 Excel（StreamWriter，限 10 万行）
func (s *userExportService) Export(ctx context.Context, req *smodel.UserPageReq) (*excelize.File, error) {
	req.Page = 1
	req.PageSize = maxExportRows

	list, _, err := s.userSvc.Page(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	f := excelize.NewFile()
	sw, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("创建 StreamWriter 失败: %w", err)
	}

	headers := []string{"ID", "用户名", "昵称", "邮箱", "手机号", "性别", "状态", "创建时间"}
	if err := sw.SetRow("A1", anySlice(headers)); err != nil {
		return nil, err
	}

	for i, u := range list {
		row := i + 2
		cell, _ := excelize.CoordinatesToCellName(1, row)
		sexStr := "未知"
		switch u.Sex {
		case 1:
			sexStr = "男"
		case 2:
			sexStr = "女"
		}
		statusStr := "正常"
		if u.Status == 1 {
			statusStr = "停用"
		}
		if err := sw.SetRow(cell, anySlice([]string{
			fmt.Sprintf("%d", u.ID), u.Username, u.Nickname, u.Email, u.Mobile,
			sexStr, statusStr, u.CreateTime,
		})); err != nil {
			return nil, err
		}
	}

	if err := sw.Flush(); err != nil {
		return nil, fmt.Errorf("刷新 StreamWriter 失败: %w", err)
	}

	return f, nil
}

func anySlice(s []string) []interface{} {
	r := make([]interface{}, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}
