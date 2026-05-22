package handler

import (
	"context"
	"testing"

	pb "management-backend/api/proto/system"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- Mock UserService ---

type mockGRPUserService struct {
	users map[int64]*smodel.UserResp
}

func newMockGRPCUserService() *mockGRPUserService {
	return &mockGRPUserService{users: make(map[int64]*smodel.UserResp)}
}

func (m *mockGRPUserService) Create(_ context.Context, _ *smodel.UserCreateReq, _ int64, _ int64) (int64, error) {
	return 9001, nil
}

func (m *mockGRPUserService) Update(_ context.Context, _ *smodel.UserUpdateReq, _ int64) error {
	return nil
}

func (m *mockGRPUserService) Delete(_ context.Context, id int64) error {
	if id == 9999 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *mockGRPUserService) GetByID(_ context.Context, id int64) (*smodel.UserResp, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockGRPUserService) GetByUsername(_ context.Context, _ string) (*smodel.User, error) {
	return nil, nil
}

func (m *mockGRPUserService) GetRawByID(_ context.Context, _ int64) (*smodel.User, error) {
	return nil, nil
}

func (m *mockGRPUserService) UpdateProfile(_ context.Context, _ *service.UserProfileUpdateReq) error {
	return nil
}

func (m *mockGRPUserService) UpdatePassword(_ context.Context, _ int64, _ string) error {
	return nil
}

func (m *mockGRPUserService) Page(_ context.Context, _ *smodel.UserPageReq) ([]smodel.UserResp, int64, error) {
	return []smodel.UserResp{
		{ID: 1, Username: "admin", Nickname: "管理员"},
	}, 1, nil
}

// --- Tests ---

func TestUserGRPC_GetUser(t *testing.T) {
	testutil.NewTestDB(t)
	mockSvc := newMockGRPCUserService()
	mockSvc.users[1] = &smodel.UserResp{
		ID: 1, Username: "admin", Nickname: "管理员", Email: "admin@test.com",
	}

	h := NewUserGRPCHandler(mockSvc)
	resp, err := h.GetUser(context.Background(), &pb.GetUserReq{Id: 1})

	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "admin", resp.Username)
	assert.Equal(t, "管理员", resp.Nickname)
}

func TestUserGRPC_GetUser_NotFound(t *testing.T) {
	testutil.NewTestDB(t)
	mockSvc := newMockGRPCUserService()
	h := NewUserGRPCHandler(mockSvc)

	_, err := h.GetUser(context.Background(), &pb.GetUserReq{Id: 9999})
	assert.Error(t, err)
}

func TestUserGRPC_ListUsers(t *testing.T) {
	testutil.NewTestDB(t)
	mockSvc := newMockGRPCUserService()
	h := NewUserGRPCHandler(mockSvc)

	resp, err := h.ListUsers(context.Background(), &pb.ListUsersReq{
		Username: "admin",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.List, 1)
	assert.Equal(t, "admin", resp.List[0].Username)
}

func TestUserGRPC_CreateUser(t *testing.T) {
	testutil.NewTestDB(t)
	mockSvc := newMockGRPCUserService()
	h := NewUserGRPCHandler(mockSvc)

	resp, err := h.CreateUser(context.Background(), &pb.CreateUserReq{
		Username: "newuser", Password: "pass123", Nickname: "新用户",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(9001), resp.Id)
}

func TestUserGRPC_UpdateUser(t *testing.T) {
	testutil.NewTestDB(t)
	mockSvc := newMockGRPCUserService()
	h := NewUserGRPCHandler(mockSvc)

	_, err := h.UpdateUser(context.Background(), &pb.UpdateUserReq{
		Id: 1, Nickname: "更新后",
	})
	require.NoError(t, err)
}

func TestUserGRPC_DeleteUser(t *testing.T) {
	testutil.NewTestDB(t)
	mockSvc := newMockGRPCUserService()
	h := NewUserGRPCHandler(mockSvc)

	_, err := h.DeleteUser(context.Background(), &pb.IDReq{Id: 1})
	require.NoError(t, err)

	_, err = h.DeleteUser(context.Background(), &pb.IDReq{Id: 9999})
	assert.Error(t, err)
}
