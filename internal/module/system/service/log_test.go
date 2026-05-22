package service

import (
	"context"
	"errors"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOperLogRepo struct {
	logs map[int64]*smodel.SysOperLog
}

type mockLoginLogRepo struct {
	logs map[int64]*smodel.SysLoginLog
}

func newMockOperLogRepo() *mockOperLogRepo {
	return &mockOperLogRepo{logs: make(map[int64]*smodel.SysOperLog)}
}

func newMockLoginLogRepo() *mockLoginLogRepo {
	return &mockLoginLogRepo{logs: make(map[int64]*smodel.SysLoginLog)}
}

func (m *mockOperLogRepo) Create(_ context.Context, log *smodel.SysOperLog) error {
	m.logs[log.ID] = log
	return nil
}

func (m *mockOperLogRepo) Page(_ context.Context, _ *smodel.OperLogPageReq) ([]smodel.SysOperLog, int64, error) {
	var list []smodel.SysOperLog
	for _, l := range m.logs {
		list = append(list, *l)
	}
	return list, int64(len(list)), nil
}

func (m *mockOperLogRepo) DeleteByID(_ context.Context, id int64) error {
	if _, ok := m.logs[id]; !ok {
		return errors.New("not found")
	}
	delete(m.logs, id)
	return nil
}

func (m *mockOperLogRepo) Clean(_ context.Context) error {
	m.logs = make(map[int64]*smodel.SysOperLog)
	return nil
}

func (m *mockLoginLogRepo) Create(_ context.Context, log *smodel.SysLoginLog) error {
	m.logs[log.ID] = log
	return nil
}

func (m *mockLoginLogRepo) Page(_ context.Context, _ *smodel.LoginLogPageReq) ([]smodel.SysLoginLog, int64, error) {
	var list []smodel.SysLoginLog
	for _, l := range m.logs {
		list = append(list, *l)
	}
	return list, int64(len(list)), nil
}

func (m *mockLoginLogRepo) DeleteByID(_ context.Context, id int64) error {
	if _, ok := m.logs[id]; !ok {
		return errors.New("not found")
	}
	delete(m.logs, id)
	return nil
}

func (m *mockLoginLogRepo) Clean(_ context.Context) error {
	m.logs = make(map[int64]*smodel.SysLoginLog)
	return nil
}

func TestLogService_OperLog(t *testing.T) {
	testutil.NewTestDB(t)
	operRepo := newMockOperLogRepo()
	svc := NewLogService(operRepo, newMockLoginLogRepo())
	ctx := context.Background()

	log := &smodel.SysOperLog{Title: "用户登录", Method: "POST"}
	log.ID = 20001
	err := svc.CreateOperLog(ctx, log)
	require.NoError(t, err)

	list, total, err := svc.PageOperLog(ctx, &smodel.OperLogPageReq{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
	assert.Equal(t, "用户登录", list[0].Title)

	err = svc.DeleteOperLog(ctx, 20001)
	require.NoError(t, err)

	err = svc.DeleteOperLog(ctx, 9999)
	assert.Error(t, err)

	err = svc.CleanOperLog(ctx)
	require.NoError(t, err)
	assert.Len(t, operRepo.logs, 0)
}

func TestLogService_LoginLog(t *testing.T) {
	testutil.NewTestDB(t)
	loginRepo := newMockLoginLogRepo()
	svc := NewLogService(newMockOperLogRepo(), loginRepo)
	ctx := context.Background()

	log := &smodel.SysLoginLog{Username: "admin", LoginIP: "127.0.0.1"}
	log.ID = 20002
	err := svc.CreateLoginLog(ctx, log)
	require.NoError(t, err)

	list, total, err := svc.PageLoginLog(ctx, &smodel.LoginLogPageReq{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "admin", list[0].Username)

	err = svc.DeleteLoginLog(ctx, 20002)
	require.NoError(t, err)

	err = svc.CleanLoginLog(ctx)
	require.NoError(t, err)
}
