package service

import (
	"testing"

	"management-backend/services/job-service/internal/model"
	"management-backend/services/job-service/internal/repository"
	"management-backend/services/job-service/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestJobService(t *testing.T) JobService {
	t.Helper()
	db := testutil.NewTestDB(t)
	return NewJobService(repository.NewJobTaskRepo(db), repository.NewExecLogRepo(db))
}

func TestJobService_Create(t *testing.T) {
	svc := newTestJobService(t)
	ctx := testutil.Ctx()

	req := &model.JobTaskCreateReq{
		Name:     "测试任务",
		Handler:  "TestHandler",
		CronExpr: "*/5 * * * *",
		Params:   `{"key":"val"}`,
	}
	id, err := svc.Create(ctx, req, 1, 100)
	require.NoError(t, err)
	assert.NotZero(t, id)
}

func TestJobService_Page(t *testing.T) {
	svc := newTestJobService(t)
	ctx := testutil.Ctx()

	// 创建 3 个任务
	for i := 0; i < 3; i++ {
		_, _ = svc.Create(ctx, &model.JobTaskCreateReq{
			Name:     "任务" + string(rune('A'+i)),
			Handler:  "Handler",
			CronExpr: "0 * * * *",
		}, 1, 100)
	}

	req := &model.JobTaskPageReq{Page: 1, Size: 10}
	list, total, err := svc.Page(ctx, req, 100)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 3)
}

func TestJobService_Page_TenantIsolation(t *testing.T) {
	svc := newTestJobService(t)
	ctx := testutil.Ctx()

	_, _ = svc.Create(ctx, &model.JobTaskCreateReq{Name: "A", Handler: "H", CronExpr: "* * * * *"}, 1, 100)
	_, _ = svc.Create(ctx, &model.JobTaskCreateReq{Name: "B", Handler: "H", CronExpr: "* * * * *"}, 1, 200)

	list100, total100, _ := svc.Page(ctx, &model.JobTaskPageReq{Page: 1, Size: 10}, 100)
	list200, total200, _ := svc.Page(ctx, &model.JobTaskPageReq{Page: 1, Size: 10}, 200)

	assert.Equal(t, int64(1), total100)
	assert.Equal(t, int64(1), total200)
	assert.Equal(t, "A", list100[0].Name)
	assert.Equal(t, "B", list200[0].Name)
}

func TestJobService_Update(t *testing.T) {
	svc := newTestJobService(t)
	ctx := testutil.Ctx()

	id, _ := svc.Create(ctx, &model.JobTaskCreateReq{Name: "原始", Handler: "H", CronExpr: "* * * * *"}, 1, 100)

	newStatus := int8(model.TaskStatusPaused)
	err := svc.Update(ctx, id, &model.JobTaskUpdateReq{Name: "更新", Status: &newStatus})
	require.NoError(t, err)

	list, _, _ := svc.Page(ctx, &model.JobTaskPageReq{Page: 1, Size: 10}, 100)
	assert.Equal(t, "更新", list[0].Name)
	assert.Equal(t, int8(model.TaskStatusPaused), list[0].Status)
}

func TestJobService_Delete(t *testing.T) {
	svc := newTestJobService(t)
	ctx := testutil.Ctx()

	id, _ := svc.Create(ctx, &model.JobTaskCreateReq{Name: "删除测试", Handler: "H", CronExpr: "* * * * *"}, 1, 100)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)

	_, total, _ := svc.Page(ctx, &model.JobTaskPageReq{Page: 1, Size: 10}, 100)
	assert.Equal(t, int64(0), total)
}

func TestJobService_Trigger(t *testing.T) {
	svc := newTestJobService(t)
	ctx := testutil.Ctx()

	id, _ := svc.Create(ctx, &model.JobTaskCreateReq{Name: "触发测试", Handler: "H", CronExpr: "* * * * *"}, 1, 100)

	logID, err := svc.Trigger(ctx, id, 100)
	require.NoError(t, err)
	assert.NotZero(t, logID)
}

func TestJobService_Trigger_WrongTenant(t *testing.T) {
	svc := newTestJobService(t)
	ctx := testutil.Ctx()

	id, _ := svc.Create(ctx, &model.JobTaskCreateReq{Name: "隔离测试", Handler: "H", CronExpr: "* * * * *"}, 1, 100)

	_, err := svc.Trigger(ctx, id, 999)
	assert.Error(t, err)
}
